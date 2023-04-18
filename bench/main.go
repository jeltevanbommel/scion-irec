package main

import (
	"context"
	"fmt"
	"github.com/scionproto/scion/control/beacon"
	"github.com/scionproto/scion/control/irec/ingress"
	"github.com/scionproto/scion/pkg/addr"
	libgrpc "github.com/scionproto/scion/pkg/grpc"
	"github.com/scionproto/scion/pkg/log"
	"github.com/scionproto/scion/pkg/private/common"
	"github.com/scionproto/scion/pkg/private/serrors"
	"github.com/scionproto/scion/pkg/private/xtest/graph"
	cppb "github.com/scionproto/scion/pkg/proto/control_plane"
	seg "github.com/scionproto/scion/pkg/segment"
	"github.com/scionproto/scion/pkg/segment/extensions/staticinfo"
	"github.com/scionproto/scion/pkg/slayers/path"
	"github.com/scionproto/scion/private/app"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var (
	IA311 = addr.MustIAFrom(1, 0xff0000000311)
	IA330 = addr.MustIAFrom(1, 0xff0000000330)
	IA331 = addr.MustIAFrom(1, 0xff0000000331)
	IA332 = addr.MustIAFrom(1, 0xff0000000332)
	IA333 = addr.MustIAFrom(1, 0xff0000000333)

	IA334 = addr.MustIAFrom(1, 0xff0000000334)
	IA335 = addr.MustIAFrom(1, 0xff0000000335)
	IA336 = addr.MustIAFrom(1, 0xff0000000336)
	IA337 = addr.MustIAFrom(1, 0xff0000000337)
	IA338 = addr.MustIAFrom(1, 0xff0000000338)
	IA339 = addr.MustIAFrom(1, 0xff0000000339)
	IA340 = addr.MustIAFrom(1, 0xff0000000340)
	IA341 = addr.MustIAFrom(1, 0xff0000000341)

	Info1 = []IfInfo{
		{
			IA:     IA330,
			Next:   IA331,
			Egress: 5,
		},
		{
			IA:      IA331,
			Next:    IA332,
			Ingress: 2,
			Egress:  3,
			Peers:   []PeerEntry{{IA: IA311, Ingress: 6}},
		},
		{
			IA:      IA332,
			Next:    IA333,
			Ingress: 1,
			Egress:  7,
		},
	}
	Info2 = []IfInfo{
		{
			IA:     IA330,
			Next:   IA334,
			Egress: 4,
		},
		{
			IA:     IA334,
			Next:   IA335,
			Egress: 5,
		},
		{
			IA:     IA335,
			Next:   IA336,
			Egress: 5,
		},

		{
			IA:     IA336,
			Next:   IA337,
			Egress: 5,
		},
		{
			IA:     IA337,
			Next:   IA338,
			Egress: 5,
		},
		{
			IA:     IA338,
			Next:   IA339,
			Egress: 5,
		},
		{
			IA:     IA339,
			Next:   IA340,
			Egress: 5,
		},
		{
			IA:     IA340,
			Next:   IA341,
			Egress: 5,
		},
		{
			IA:      IA341,
			Next:    IA332,
			Ingress: 2,
			Egress:  3,
			Peers:   []PeerEntry{{IA: IA311, Ingress: 6}},
		},
		{
			IA:      IA332,
			Next:    IA333,
			Ingress: 1,
			Egress:  7,
		},
	}

	Info3 = []IfInfo{
		{
			IA:     IA330,
			Next:   IA334,
			Egress: 5,
		},
		{
			IA:     IA334,
			Next:   IA335,
			Egress: 5,
		},
		{
			IA:     IA335,
			Next:   IA336,
			Egress: 5,
		},

		{
			IA:     IA336,
			Next:   IA337,
			Egress: 5,
		},
		{
			IA:     IA337,
			Next:   IA338,
			Egress: 5,
		},
		{
			IA:     IA338,
			Next:   IA339,
			Egress: 5,
		},
		{
			IA:      IA339,
			Next:    IA332,
			Ingress: 2,
			Egress:  3,
			Peers:   []PeerEntry{{IA: IA311, Ingress: 6}},
		},
		{
			IA:      IA332,
			Next:    IA333,
			Ingress: 1,
			Egress:  7,
		},
	}

	Info4 = []IfInfo{
		{
			IA:     IA330,
			Next:   IA333,
			Egress: 8,
		},
		{
			IA:      IA332,
			Next:    IA333,
			Ingress: 1,
			Egress:  7,
		},
	}

	Info5 = []IfInfo{
		{
			IA:     IA330,
			Next:   IA334,
			Egress: 5,
		},
		{
			IA:     IA334,
			Next:   IA335,
			Egress: 5,
		},
		{
			IA:     IA335,
			Next:   IA336,
			Egress: 5,
		},

		{
			IA:     IA336,
			Next:   IA337,
			Egress: 5,
		},
		{
			IA:     IA337,
			Next:   IA338,
			Egress: 5,
		},
		{
			IA:     IA338,
			Next:   IA339,
			Egress: 5,
		},
		{
			IA:     IA339,
			Next:   IA340,
			Egress: 5,
		},
		{
			IA:     IA340,
			Next:   IA338,
			Egress: 5,
		},
		{
			IA:     IA338,
			Next:   IA339,
			Egress: 5,
		},
		{
			IA:     IA339,
			Next:   IA340,
			Egress: 5,
		},
		{
			IA:     IA340,
			Next:   IA341,
			Egress: 5,
		},
		{
			IA:      IA341,
			Next:    IA332,
			Ingress: 2,
			Egress:  3,
			Peers:   []PeerEntry{{IA: IA311, Ingress: 6}},
		},
		{
			IA:      IA332,
			Next:    IA333,
			Ingress: 1,
			Egress:  7,
		},
	}
)

type PeerEntry struct {
	IA      addr.IA
	Ingress common.IFIDType
}

type IfInfo struct {
	IA      addr.IA
	Next    addr.IA
	Ingress common.IFIDType
	Egress  common.IFIDType
	Peers   []PeerEntry
}

func MockBeacon(ases []IfInfo, id uint16, inIfId uint16, infoTS int64) (beacon.Beacon, *cppb.IRECPathSegment, []byte) {

	entries := make([]seg.ASEntry, len(ases))
	for i, as := range ases {
		var mtu int
		if i != 0 {
			mtu = 1500
		}
		var peers []seg.PeerEntry
		for _, peer := range as.Peers {
			peers = append(peers, seg.PeerEntry{
				Peer:          peer.IA,
				PeerInterface: 1337,
				PeerMTU:       1500,
				HopField: seg.HopField{
					ExpTime:     63,
					ConsIngress: uint16(peer.Ingress),
					ConsEgress:  uint16(as.Egress),
					MAC:         [path.MacLen]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
			})
		}
		entries[i] = seg.ASEntry{
			Local: as.IA,
			Next:  as.Next,
			MTU:   1500,
			HopEntry: seg.HopEntry{
				IngressMTU: mtu,
				HopField: seg.HopField{
					ExpTime:     63,
					ConsIngress: uint16(as.Ingress),
					ConsEgress:  uint16(as.Egress),
					MAC:         [path.MacLen]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
				},
			},
			PeerEntries: peers,
			Extensions: seg.Extensions{
				HiddenPath: seg.HiddenPathExtension{},
				StaticInfo: &staticinfo.Extension{
					Latency: staticinfo.LatencyInfo{
						Intra: nil,
						Inter: nil,
					},
					Bandwidth:    staticinfo.BandwidthInfo{},
					Geo:          nil,
					LinkType:     nil,
					InternalHops: nil,
					Note:         "",
				},
				Digests: nil,
				Irec:    nil,
			},
		}
	}

	// XXX(roosd): deterministic beacon needed.
	pseg, err := seg.CreateSegment(time.Unix(int64(infoTS), 0), id)
	if err != nil {
		//log.Error("err", "err", err)
		return beacon.Beacon{}, &cppb.IRECPathSegment{}, nil
	}

	asEntries := make([]*cppb.IRECASEntry, 0, len(entries))
	for _, entry := range entries {
		signer := graph.NewSigner()
		// for testing purposes set the signer timestamp equal to infoTS
		signer.Timestamp = time.Unix(int64(infoTS), 0)
		err := pseg.AddASEntry(context.Background(), entry, signer)
		if err != nil {
			return beacon.Beacon{}, &cppb.IRECPathSegment{}, nil
		}
	}
	// above loop adds to the pseg asentries, now transform them into irec as entries.
	// bit hacky

	for _, entry := range pseg.ASEntries {
		asEntries = append(asEntries, &cppb.IRECASEntry{
			Signed:     entry.Signed,
			SignedBody: entry.SignedBody,
			Unsigned:   seg.UnsignedExtensionsToPB(entry.UnsignedExtensions),
		})
	}
	return beacon.Beacon{Segment: pseg, InIfId: inIfId}, &cppb.IRECPathSegment{
		SegmentInfo: pseg.Info.Raw,
		AsEntries:   asEntries,
	}, pseg.ID()
}

var globalCfg Config

func main() {
	if err := log.Setup(log.Config{
		Console: log.ConsoleConfig{Format: "human", Level: "debug"},
	}); err != nil {
		fmt.Println(err)
		return
	}
	defer log.Flush()
	defer log.HandlePanic()

	amountOfBeacons, _ := strconv.Atoi(os.Args[1])
	beacons := make([]beacon.Beacon, amountOfBeacons)
	irecbeacons := make([]*cppb.IRECBeacon, amountOfBeacons)
	if os.Args[2] == "rand" {
		s1 := rand.NewSource(42)
		r1 := rand.New(s1)
		for i := 0; i < amountOfBeacons; i++ {
			var bcn beacon.Beacon
			var xpb *cppb.IRECPathSegment
			switch r1.Intn(5) {
			case 0:
				bcn, xpb, _ = MockBeacon(Info1, uint16(i), 0, time.Now().Unix())
				break
			case 1:
				bcn, xpb, _ = MockBeacon(Info2, uint16(i), 0, time.Now().Unix())
				break
			case 2:
				bcn, xpb, _ = MockBeacon(Info3, uint16(i), 0, time.Now().Unix())
				break
			case 3:
				bcn, xpb, _ = MockBeacon(Info4, uint16(i), 0, time.Now().Unix())
				break
			case 4:
				bcn, xpb, _ = MockBeacon(Info5, uint16(i), 0, time.Now().Unix())
				break
			}

			beacons[i] = bcn
			irecbeacons[i] = &cppb.IRECBeacon{
				PathSeg: xpb,
				InIfId:  0,
				Id:      uint32(i),
			}
		}
	} else {
		for i := 0; i < amountOfBeacons; i++ {
			bcn, xpb, _ := MockBeacon(Info1, uint16(i), 0, time.Now().Unix())

			beacons[i] = bcn
			irecbeacons[i] = &cppb.IRECBeacon{
				PathSeg: xpb,
				InIfId:  0,
				Id:      uint32(i),
			}
		}
	}
	env := BenchEnv{
		beacons:     beacons,
		irecbeacons: irecbeacons,
		RACManager: ingress.RACManager{
			RACCount:      &atomic.Uint64{},
			RACList:       &sync.Map{},
			RACListUnsafe: make(map[string]*ingress.RACInfo),
			Dialer: &libgrpc.TCPDialer{
				SvcResolver: func(dst addr.HostSVC) []resolver.Address {
					return []resolver.Address{}
				},
			},
		},
	}
	//ip, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:20000")
	//env.RACManager.RACList["rac1"] = ingress.RACInfo{State: ingress.Free, Addr: ip}

	//t := make(map[uint32]*int64)
	if os.Args[3] == "stress" {
		// Used in stress tests; measure the PCB/s of incoming requests.
		t := atomic.Uint64{}
		env.RegisterServers(&SimpleEgressCounter{Results: &t})
		time.Sleep(5 * time.Second)
		amtIters, _ := strconv.Atoi(os.Args[4])
		for i := 0; i < amtIters; i++ {
			t.Store(0)
			start := time.Now()
			time.Sleep(10 * time.Second)
			count := t.Load()
			duration := time.Since(start)
			fmt.Printf("%d, %d, %f pcb/s\n", duration.Nanoseconds(), count, float64(count*uint64(amountOfBeacons))/duration.Seconds())
		}
	} else if os.Args[3] == "rac" {
		// Spawns a server, which just serves requests and logs the time at which the egress prop request came in
		// Use this to measure on the RACs itself
		env.RegisterServers(&EgressLogger{})
		g, errCtx := errgroup.WithContext(context.Background())
		g.Go(func() error {
			defer log.HandlePanic()
			<-errCtx.Done()
			return nil
		})

		g.Wait()
	} else if os.Args[3] == "scion" {
		// Measure SCION Latency for the whole algorithm.
		time.Sleep(5 * time.Second)
		amtIters, _ := strconv.Atoi(os.Args[4])
		//env.RegisterIngressServer(&EgressLogger{})
		for i := 0; i < amtIters; i++ {
			env.BenchSCION()
		}
	} else if os.Args[3] == "scionthr" {
		// Measure SCION throughput for the algorithm over 10 seconds.
		amtIters, _ := strconv.Atoi(os.Args[4])
		for rep := 0; rep < amtIters; rep += 1 {
			i := 0
			start := time.Now()
			for time.Since(start) < 10*time.Second {
				i += 1
				baseAlgo{}.SelectBeacons(env.beacons, 20)

			}
			duration := time.Since(start)
			fmt.Printf("%d, %d, %f pcb/s\n", duration.Nanoseconds(), i, float64(uint64(i)*uint64(amountOfBeacons))/duration.Seconds())

		}
	}

	//testing.Benchmark(BenchmarkRac)

}

type BenchEnv struct {
	beacons     []beacon.Beacon
	irecbeacons []*cppb.IRECBeacon
	RACManager  ingress.RACManager
	tcpServer   *grpc.Server
}

type EgressLogger struct {
}

func (e *EgressLogger) RequestPropagation(ctx context.Context, request *cppb.PropagationRequest) (*cppb.PropagationRequestResponse, error) {
	completed := time.Now().UnixNano()
	fmt.Println("\t[RAC]=e", request.Beacon[0].Egress[0], ",", completed)
	return &cppb.PropagationRequestResponse{}, nil
}

type SimpleEgressCounter struct {
	Results *atomic.Uint64
}

func (e *SimpleEgressCounter) RequestPropagation(ctx context.Context, request *cppb.PropagationRequest) (*cppb.PropagationRequestResponse, error) {
	e.Results.Add(1)
	return &cppb.PropagationRequestResponse{}, nil
}

type EgressCounter struct {
	Results map[uint32]*int64
}

func (e *EgressCounter) RequestPropagation(ctx context.Context, request *cppb.PropagationRequest) (*cppb.PropagationRequestResponse, error) {
	atomic.AddInt64(e.Results[request.Beacon[0].Egress[0]], 1)
	return &cppb.PropagationRequestResponse{}, nil
}

type IngressDummy struct {
	env *BenchEnv
}

func (i *IngressDummy) GetJob(ctx context.Context, request *cppb.RACBeaconRequest) (*cppb.RACJob, error) {
	return &cppb.RACJob{
		Beacons:       i.env.irecbeacons,
		AlgorithmHash: []byte{0x00, 0xde, 0xad, 0xbe, 0xef},
		PropIntfs:     []uint32{request.AlgorithmID, 1, 2},
	}, nil
}

func (i *IngressDummy) Handle(ctx context.Context, req *cppb.IncomingBeacon) (*cppb.IncomingBeaconResponse, error) {
	return &cppb.IncomingBeaconResponse{}, nil
}

func (i *IngressDummy) RegisterRAC(ctx context.Context, req *cppb.RegisterRACRequest) (*cppb.RegisterRACResponse, error) {
	return &cppb.RegisterRACResponse{}, nil
}

func (i *IngressDummy) BeaconSources(ctx context.Context, request *cppb.RACBeaconSourcesRequest) (*cppb.RACBeaconSources, error) {
	return &cppb.RACBeaconSources{
		Sources: []*cppb.RACBeaconSource{&cppb.RACBeaconSource{
			AlgorithmHash:   []byte{0x00, 0xde, 0xad, 0xbe, 0xef},
			AlgorithmID:     0,
			OriginAS:        0,
			OriginIntfGroup: 0,
		}},
	}, nil
}

func (i *IngressDummy) RequestBeacons(ctx context.Context, req *cppb.RACBeaconRequest) (*cppb.RACBeaconResponse, error) {
	return &cppb.RACBeaconResponse{}, nil
}

func (env *BenchEnv) RegisterServers(server cppb.EgressServiceServer) {
	tcpServer := grpc.NewServer()
	//is := &ingress.IngressServer{
	//	IncomingHandler: ingress.Handler{},
	//	RACManager:      env.RACManager,
	//}
	cppb.RegisterEgressServiceServer(tcpServer, server)
	cppb.RegisterIngressServiceServer(tcpServer, &IngressDummy{env: env})
	var cleanup app.Cleanup
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		ip, _ := net.ResolveTCPAddr("tcp", "127.0.0.1:33333")
		listener, _ := net.ListenTCP("tcp", ip)
		if err := tcpServer.Serve(listener); err != nil {
			return serrors.WrapStr("serving gRPC/TCP API", err)
		}
		return nil
	})
	cleanup.Add(func() error { tcpServer.GracefulStop(); return nil })
	env.tcpServer = tcpServer
}

func (env *BenchEnv) GetBeaconsSCION() []beacon.Beacon {
	return env.beacons
}
func (env *BenchEnv) BenchSCION() {
	start := time.Now()
	baseAlgo{}.SelectBeacons(env.GetBeaconsSCION(), 20)
	duration := time.Since(start)
	fmt.Println("[SCION] done in", duration.Nanoseconds())
}

//func (env *BenchEnv) GetBeaconsRAC() (*cppb.RACBeaconResponse, error) {
//	//storedBeacon := make([]*cppb.StoreadBeacon, 0, len(env.beacons))
//	//for i, b := range env.beacons {
//	//	// TODO: get rid of this overhead: raw-bytes in sql -> pathseg -> beacon -> pathseg proto -> raw bytes.
//	//	storedBeacon = append(storedBeacon, &cppb.IRECBeacon{
//	//		PathSeg: seg.PathSegmentToPB(b.Segment),
//	//		InIfId:  uint32(b.InIfId),
//	//		Id:      uint32(i),
//	//	})
//	//}
//	return &cppb.RACBeaconResponse{Beacons: env.irecbeacons}, nil
//}
//
//var wg sync.WaitGroup
//
//func (env *BenchEnv) BenchRac(iteration uint32) {
//	started := time.Now().UnixNano()
//	bcns, _ := env.GetBeaconsRAC()
//	//start := time.Now()
//	//fmt.Println("start at", time.Now().UnixNano())
//	env.RACManager.RunImm(context.Background(), &racservice.ExecutionRequest{
//		Beacons:       bcns.Beacons,
//		AlgorithmHash: []byte{0x00, 0xde, 0xad, 0xbe, 0xef},
//		Intfs:         nil,
//		PropIntfs:     []uint32{iteration, 1, 2},
//	}, &wg)
//	fmt.Println("\t[RAC]=s", iteration, ",", started)
//
//	wg.Wait()
//}
