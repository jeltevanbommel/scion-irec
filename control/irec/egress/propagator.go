package egress

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"github.com/scionproto/scion/control/beacon"
	"github.com/scionproto/scion/control/ifstate"
	"github.com/scionproto/scion/pkg/log"
	cppb "github.com/scionproto/scion/pkg/proto/control_plane"
	seg "github.com/scionproto/scion/pkg/segment"
	"sync"
	"time"
)

// The propagator takes beacons from the beaconDb,
// Propagator checks the egress db, if not present for intf, extend and send immediately.
// It also keeps track when last sent beacon, such that an originate script can send beacon if necessary.

type Propagator struct {
	Store DB
	//Tick                  Tick
	AllInterfaces         *ifstate.Interfaces
	PropagationInterfaces func() []*ifstate.Interface
	Interfaces            map[uint32]*ifstate.Interface
	Extender              Extender
	PropagationFilter     func(*ifstate.Interface) bool
	Peers                 []uint16 // TODO(jvb): static or reindex everytime?
	SenderFactory         SenderFactory
	Writers               []Writer
}

func HashBeacon(ctx context.Context, segment *seg.PathSegment) []byte {
	h := sha256.New()
	binary.Write(h, binary.BigEndian, segment.Info.SegmentID)
	binary.Write(h, binary.BigEndian, segment.Info.Timestamp)
	binary.Write(h, binary.BigEndian, segment.Info.Raw)
	for _, ase := range segment.ASEntries {
		binary.Write(h, binary.BigEndian, ase.Local)
		if ase.Extensions.Irec != nil {
			binary.Write(h, binary.BigEndian, ase.Extensions.Irec.AlgorithmHash)
			binary.Write(h, binary.BigEndian, ase.Extensions.Irec.AlgorithmId)
			binary.Write(h, binary.BigEndian, ase.Extensions.Irec.InterfaceGroup)
		}
		//binary.Write(h, binary.BigEndian, ase.Extensions.StaticInfo.Bandwidth.Inter)
		binary.Write(h, binary.BigEndian, ase.HopEntry.HopField.ConsIngress)
		binary.Write(h, binary.BigEndian, ase.HopEntry.HopField.ConsEgress)
		for _, peer := range ase.PeerEntries {
			binary.Write(h, binary.BigEndian, peer.Peer)
			binary.Write(h, binary.BigEndian, peer.HopField.ConsIngress)
			binary.Write(h, binary.BigEndian, peer.HopField.ConsEgress)
		}
	}
	return h.Sum(nil)
}

const defaultNewSenderTimeout = 5 * time.Second

func (p *Propagator) RequestPropagation(ctx context.Context, request *cppb.PropagationRequest) (*cppb.PropagationRequestResponse, error) {
	var wg sync.WaitGroup
	// Write the beacons to path servers in a seperate goroutine

	go func() {
		defer log.HandlePanic()
		// A non-core AS can have multiple writers, core only one, write to all:
		for _, writer := range p.Writers {
			//TODO(jvb): not nice, streamline this...
			for _, bcn := range request.Beacon {
				segment, err := seg.BeaconFromPB(bcn.Beacon.PathSeg)
				if err != nil {
					log.Error("error occurred", "err", err)
					continue
				}
				// writer has side-effects for beacon, therefore recreate beacon arr for each writer
				err = writer.Write(context.Background(), []beacon.Beacon{{Segment: segment, InIfId: uint16(bcn.Beacon.InIfId)}}, p.Peers)
				if err != nil {
					log.Error("error occurred", "err", err)
				}
			}
		}
	}()
	for _, beacon := range request.Beacon {
		// Every beacon is to be propagated on a set of interfaces
		for _, intfId := range beacon.Egress {
			wg.Add(1)
			// Copy to have the vars in the goroutine
			intfId := intfId
			beacon := beacon
			go func() {
				defer log.HandlePanic()
				defer wg.Done()
				intf := p.Interfaces[intfId]
				if intf == nil {
					log.Error("Attempt to send beacon on non-existent interface", "egress_interface", intfId)
					return
				}
				if !p.PropagationFilter(intf) {
					log.Error("Attempt to send beacon on filtered egress interface", "egress_interface", intfId)
					return
				}

				segment, err := seg.BeaconFromPB(beacon.Beacon.PathSeg)
				if err != nil {
					log.Error("Beacon DB propagation segment failed", "err", err)
					return
				}

				beaconHash := HashBeacon(ctx, segment)

				log.Info("Irec Propagation requested for", hex.EncodeToString(beaconHash), segment)
				// Check if beacon is already propagated before using egress db
				propagated, _, err := p.Store.IsBeaconAlreadyPropagated(ctx, beaconHash, intf)
				if err != nil {
					log.Error("Beacon DB Propagation check failed", "err", err)
					return
				}
				// If so, don't propagate
				if propagated {
					log.Info("Beacon is known in egress database, skipping propagation.")
					return
				}
				// If not, mark it as propagated
				err = p.Store.MarkBeaconAsPropagated(ctx, beaconHash, intf, time.Now().Add(time.Hour))
				if err != nil {
					log.Error("Beacon DB Propagation add failed", "err", err)
					return
				}
				// If the Origin-AS used Irec, we copy the algorithmID and hash from the first as entry
				if segment.ASEntries[0].Extensions.Irec != nil {
					err = p.Extender.Extend(ctx, segment, uint16(beacon.Beacon.InIfId), 0,
						intf.TopoInfo().ID,
						segment.ASEntries[0].Extensions.Irec.AlgorithmId,
						segment.ASEntries[0].Extensions.Irec.AlgorithmHash, true,
						p.Peers)
				} else {
					// Otherwise, default values.
					err = p.Extender.Extend(ctx, segment, uint16(beacon.Beacon.InIfId), 0,
						intf.TopoInfo().ID, 0, nil, false, p.Peers)
				}

				if err != nil {
					log.Error("Extending failed", "err", err)
					return
				}

				//	Propagate to ingress gateway
				senderCtx, cancel := context.WithTimeout(ctx, defaultNewSenderTimeout)
				defer cancel()
				//TODO(jvb): Reuse these senders
				sender, err := p.SenderFactory.NewSender(
					senderCtx,
					intf.TopoInfo().IA,
					intf.TopoInfo().ID,
					intf.TopoInfo().InternalAddr.UDPAddr(),
				)
				if err != nil {
					log.Error("Creating sender failed", "err", err)
					return
				}
				if err := sender.Send(ctx, segment); err != nil {
					log.Error("Sending beacon failed", "err", err)
					return
				}
				// Here we keep track of the last time a beacon has been sent on an interface per algorithm hash.
				// Such that the egress gateway can plan origination scripts for those algorithms that need origination.
				if segment.ASEntries[0].Extensions.Irec != nil {
					intf.Propagate(time.Now(), HashToString(segment.ASEntries[0].Extensions.Irec.AlgorithmHash))
				} else {
					intf.Propagate(time.Now(), "")
				}
			}()
		}
	}
	wg.Wait()
	return &cppb.PropagationRequestResponse{}, nil
	//TODO(jvb): PULL-BASED ROUTING
}

func HashToString(hash []byte) string {
	return hex.EncodeToString(hash)
}
