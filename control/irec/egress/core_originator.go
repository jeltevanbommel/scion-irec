package egress

import (
	"context"
	"crypto/rand"
	"github.com/scionproto/scion/control/ifstate"
	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/pkg/log"
	"github.com/scionproto/scion/pkg/private/serrors"
	seg "github.com/scionproto/scion/pkg/segment"
	"math/big"
	"sort"
	"sync"
	"time"
)

// The Core Originator originates for each algorithm and intfgroup.
type OriginationAlgorithm struct {
	ID   uint32
	Hash []byte
}

// Originator originates beacons. It should only be used by core ASes.
type Originator struct {
	Extender              Extender
	SenderFactory         SenderFactory
	IA                    addr.IA
	Intfs                 []*ifstate.Interface
	OriginatePerIntfGroup bool
	OriginationAlgorithms []OriginationAlgorithm
}

// Name returns the tasks name.
func (o *Originator) Name() string {
	return "control_beaconing_originator"
}

// Run originates core and downstream beacons.
func (o *Originator) Run(ctx context.Context) {
	o.originateBeacons(ctx)
}
func (o *Originator) originateBeacons(ctx context.Context) {
	intfs := o.Intfs
	sort.Slice(intfs, func(i, j int) bool {
		return intfs[i].TopoInfo().ID < intfs[j].TopoInfo().ID
	})
	if len(intfs) == 0 || len(o.OriginationAlgorithms) == 0 {
		return
	}

	var wg sync.WaitGroup
	wg.Add(len(intfs) * len(o.OriginationAlgorithms))
	for _, alg := range o.OriginationAlgorithms {
		for _, intf := range intfs {
			b := intfOriginator{
				Originator: o,
				intf:       intf,
				timestamp:  time.Now(),
				algId:      alg.ID,
				algHash:    alg.Hash,
			}
			go func() {
				defer log.HandlePanic()
				defer wg.Done()

				if err := b.originateMessage(ctx); err != nil {
					log.Info("Unable to originate on interface",
						"egress_interface", b.intf.TopoInfo().ID, "err", err)
				}
			}()
		}
	}
	wg.Wait()
}

// intfOriginator originates one beacon on the given interface.
type intfOriginator struct {
	*Originator
	intf      *ifstate.Interface
	timestamp time.Time
	algId     uint32
	algHash   []byte
}

// originateBeacon originates a beacon on the given ifid.
func (o *intfOriginator) originateMessage(ctx context.Context) error {
	topoInfo := o.intf.TopoInfo()

	senderStart := time.Now()
	duration := time.Duration(5 * len(topoInfo.Groups))
	senderCtx, cancelF := context.WithTimeout(ctx, duration*time.Second)
	defer cancelF()

	sender, err := o.SenderFactory.NewSender(
		senderCtx,
		topoInfo.IA,
		topoInfo.ID,
		topoInfo.InternalAddr.UDPAddr(),
	)
	if err != nil {
		return serrors.WrapStr("getting beacon sender", err,
			"waited_for", time.Since(senderStart).String())
	}
	defer sender.Close()
	for _, intfGroup := range topoInfo.Groups {
		log.Debug("Originating using group", "grp", intfGroup)
		sendStart := time.Now()
		beacon, err := o.createBeacon(ctx, intfGroup)
		if err != nil {
			return serrors.WrapStr("creating beacon", err)
		}
		if err := sender.Send(ctx, beacon); err != nil {
			return serrors.WrapStr("sending beacon", err,
				"waited_for", time.Since(sendStart).String(),
			)
		}
	}
	return nil
}
func (o *intfOriginator) createBeacon(ctx context.Context, intfGroup uint16) (*seg.PathSegment, error) {
	segID, err := rand.Int(rand.Reader, big.NewInt(1<<16))
	if err != nil {
		return nil, err
	}
	beacon, err := seg.CreateSegment(o.timestamp, uint16(segID.Uint64()))
	if err != nil {
		return nil, serrors.WrapStr("creating segment", err)
	}

	if err := o.Extender.Extend(ctx, beacon, 0, intfGroup, o.intf.TopoInfo().ID, o.algId, o.algHash, true, nil); err != nil {
		return nil, serrors.WrapStr("extending segment", err)
	}
	return beacon, nil
}
