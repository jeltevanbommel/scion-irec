package main

import (
	"context"
	"fmt"
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/scionproto/scion/pkg/addr"
	libgrpc "github.com/scionproto/scion/pkg/grpc"
	"github.com/scionproto/scion/pkg/log"
	cppb "github.com/scionproto/scion/pkg/proto/control_plane"
	racpb "github.com/scionproto/scion/pkg/proto/rac"
	"github.com/scionproto/scion/private/app"
	"github.com/scionproto/scion/private/topology"
	"github.com/scionproto/scion/rac"
	"github.com/scionproto/scion/rac/env/wasm"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"net"
	_ "net/http/pprof"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/scionproto/scion/pkg/private/serrors"
	"github.com/scionproto/scion/private/app/launcher"
)

var globalCfg Config

func main() {
	application := launcher.Application{
		TOMLConfig: &globalCfg,
		ShortName:  "SCION RAC",
		Main:       realMain,
	}
	application.Run()
}

func realMain(ctx context.Context) error {

	topo, err := topology.NewLoader(topology.LoaderCfg{
		File:      globalCfg.General.Topology(),
		Reload:    app.SIGHUPChannel(ctx),
		Validator: &topology.DefaultValidator{},
	})
	log.Info("config", "is", globalCfg.RAC)
	if err != nil {
		return serrors.WrapStr("creating topology loader", err)
	}
	g, errCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		defer log.HandlePanic()
		return topo.Run(errCtx)
	})

	dialer := &libgrpc.TCPDialer{
		SvcResolver: func(dst addr.HostSVC) []resolver.Address {
			if base := dst.Base(); base != addr.SvcCS {
				panic("Unsupported address type, implementation error?")
			}
			targets := []resolver.Address{}
			for _, entry := range topo.ControlServiceAddresses() {
				targets = append(targets, resolver.Address{Addr: entry.String()})
			}
			return targets
		},
	}
	listen := globalCfg.RAC.RACAddr
	host, port, err := net.SplitHostPort(listen)
	switch {
	case err != nil:
		listen = net.JoinHostPort(listen, strconv.Itoa(31033))
	case port == "0", port == "":
		listen = net.JoinHostPort(host, strconv.Itoa(31033))
	}

	listener, err := net.Listen("tcp", listen)
	if err != nil {
		return serrors.WrapStr("listening", err)
	}

	server := grpc.NewServer(grpc.MaxRecvMsgSize(1024 * 1024 * 20))

	env := wasm.WasmEnv{}
	env.Initialize()
	algCache := make(map[string]*wasmtime.Module)
	for _, alg := range globalCfg.RAC.LocalAlgorithms {
		module, err := wasmtime.NewModuleFromFile(env.Engine, alg.FilePath)
		if err != nil {
			log.Error("Algorithm not loaded due to error", alg.HexHash, err)
			continue
		}
		log.Info(fmt.Sprintf("Loaded algorithm %s", alg.HexHash))
		algCache[alg.HexHash] = module
	}

	conn, err := dialer.Dial(ctx, addr.SvcCS)
	if err != nil {
		return err
	}
	defer conn.Close()
	rs := &rac.Server{
		Writer: rac.Writer{
			Conn: conn,
		},
		ModuleStore: algCache,
		Environment: env,
	}
	//racpb.RegisterRACServiceServer(server, rs)

	var cleanup app.Cleanup
	g.Go(func() error {
		defer log.HandlePanic()
		if err := server.Serve(listener); err != nil {
			return serrors.WrapStr("serving gRPC API", err, "addr", listen)
		}
		return nil
	})
	cleanup.Add(func() error { server.GracefulStop(); return nil })
	ctr := atomic.Uint64{}
	// main loop
	for true {
		client := cppb.NewIngressServiceClient(conn)
		// TODO(jvb): create a shortcut client.GetMeAnything
		// First get possible sources from the ingress gateway (source=originas, algorithmid, alghash combo)
		sources, err1 := client.BeaconSources(ctx, &cppb.RACBeaconSourcesRequest{
			IgnoreIntfGroup: false,
		}, libgrpc.RetryProfile...)
		if err1 != nil {
			log.Error("Error when retrieving beacon sources", "err", err1)
			time.Sleep(1 * time.Second)
			continue
		}
		if len(sources.Sources) > 0 {
			// If there are PCB sources to process, get the job. This will mark the PCB's as taken such that other
			// RACS do not reprocess them.
			exec, err2 := client.GetJob(ctx, &cppb.RACBeaconRequest{
				Maximum:         0,
				AlgorithmHash:   sources.Sources[0].AlgorithmHash,
				AlgorithmID:     sources.Sources[0].AlgorithmID,
				OriginAS:        sources.Sources[0].OriginAS,
				OriginIntfGroup: sources.Sources[0].OriginIntfGroup,
				IgnoreIntfGroup: false,
			})
			if err2 != nil {
				log.Error("Error when retrieving job for sources", "err", err2)
				time.Sleep(1 * time.Second)
				continue
			}
			_, err3 := rs.Execute(ctx, &racpb.ExecutionRequest{
				Beacons:       exec.Beacons,
				AlgorithmHash: exec.AlgorithmHash,
				Intfs:         nil,
				PropIntfs:     exec.PropIntfs,
			}, int32(ctr.Load()))
			ctr.Add(1)
			// TODO(jvb): Acknowledge to the ingress gateway that we have completed the job.
			// Ingress gateway needs to keep a map of outstanding jobs, if one is not completed within a specific
			// time it returns it to the pool.
			if err3 != nil {
				log.Error("Error when executing rac for sources", "err", err3)
				time.Sleep(1 * time.Second)
				continue
			}
		} else {
			time.Sleep(1 * time.Second)
		}
	}

	g.Go(func() error {
		defer log.HandlePanic()
		<-errCtx.Done()
		return cleanup.Do()
	})
	return nil
	//return g.Wait()
}
