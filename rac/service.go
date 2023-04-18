package rac

import (
	"context"
	"encoding/hex"
	"fmt"
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/scionproto/scion/pkg/log"
	"github.com/scionproto/scion/pkg/private/serrors"
	cppb "github.com/scionproto/scion/pkg/proto/control_plane"
	"github.com/scionproto/scion/pkg/proto/rac"
	"github.com/scionproto/scion/rac/env/wasm"
	"path"
	"time"
)

type Server struct {
	Writer
	ModuleStore map[string]*wasmtime.Module
	Environment wasm.WasmEnv
}

func (s *Server) Execute(ctx context.Context, request *rac.ExecutionRequest, our_num int32) (*rac.ExecutionResponse, error) {
	hexHash := hex.EncodeToString(request.AlgorithmHash)

	algorithm, ok := s.ModuleStore[hexHash]
	if hexHash == "00010203040506070809" { // fallback rac
		algorithm = s.ModuleStore["default"]
	} else if !ok {
		var err error
		algorithm, err = wasmtime.NewModuleFromFile(s.Environment.Engine, path.Join("algorithms", path.Clean(fmt.Sprintf("%s.wasm", hexHash))))
		if err != nil {
			log.Error("Algorithm not loaded due to error", hexHash, err)
			//	todo: fetch the algorithm from ingress gateway algorithm storage.
			return &rac.ExecutionResponse{}, serrors.New("Algorithm Fetching")
		}
	}
	execCtx, cancelF := context.WithTimeout(ctx, 1*time.Second) // Only allow an execution of max. 1 second
	defer cancelF()
	s.Environment.Execute(execCtx, int32(our_num), request.Beacons, algorithm, wasm.AdditionalInfo{PropagationInterfaces: request.PropIntfs},
		func(response *rac.RACResponse) {
			selectedBeacons := make([]*cppb.StoredBeaconAndIntf, 0)
			// The selected beacons by the RAC come in the format beacon_id, egress_intfs[]
			// We do not return copies of the beacons back from the RAC, but rather find the beacons
			// afterwards using the ID index.
			startJuggling := time.Now()
			// ANd now we do need the signed binary data, so add only this back in;
			for _, selected := range response.Selected {
				// Juggle some data to remove the signed body, it is disregarded at egress anyways.
				// Saves a bit of bw, hopefully not at the expense of too much comp

				irecASEntries := make([]*cppb.ASEntry, 0, len(request.Beacons[selected.Id].PathSeg.AsEntries))
				for _, asEntry := range request.Beacons[selected.Id].PathSeg.AsEntries {
					irecASEntries = append(irecASEntries, &cppb.ASEntry{
						Signed:   asEntry.Signed,
						Unsigned: asEntry.Unsigned,
					})
				}
				selectedBeacons = append(selectedBeacons, &cppb.StoredBeaconAndIntf{Beacon: &cppb.StoredBeacon{
					PathSeg: &cppb.PathSegment{
						SegmentInfo: request.Beacons[selected.Id].PathSeg.SegmentInfo,
						AsEntries:   irecASEntries,
					},
					InIfId: request.Beacons[selected.Id].InIfId,
				}, Egress: selected.EgressIntfs})
			}
			timeJuggling := time.Since(startJuggling)
			startGrpc3 := time.Now()
			// Write the selected beacons to the egress gateway
			err := s.Writer.writeBeacons(ctx, selectedBeacons)
			if err != nil {
				log.Info("err", "msg", err)
				return
			}
			timeGrpc3 := time.Since(startGrpc3)
			fmt.Printf("%d, grpc3=%d, %d\n", our_num, timeGrpc3.Nanoseconds(), timeJuggling.Nanoseconds())
		})

	return &rac.ExecutionResponse{}, nil
}
