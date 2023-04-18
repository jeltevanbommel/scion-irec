package ingress

import (
	"bytes"
	"context"
	"github.com/opentracing/opentracing-go"
	"github.com/scionproto/scion/control/beacon"
	"github.com/scionproto/scion/control/ifstate"
	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/pkg/log"
	"github.com/scionproto/scion/pkg/private/serrors"
	seg "github.com/scionproto/scion/pkg/segment"
	"github.com/scionproto/scion/pkg/snet"
	"github.com/scionproto/scion/private/segment/segverifier"
	infra "github.com/scionproto/scion/private/segment/verifier"
	"github.com/scionproto/scion/private/svc"
	"github.com/scionproto/scion/private/topology"
)

type Handler struct {
	LocalIA     addr.IA
	IngressDB   IngressStore
	Verifier    infra.Verifier
	Interfaces  *ifstate.Interfaces
	AlgorithmDB AlgorithmStorage
}

func (h Handler) HandleBeacon(ctx context.Context, b beacon.Beacon, peer *snet.UDPAddr) error {
	span := opentracing.SpanFromContext(ctx)
	intf := h.Interfaces.Get(b.InIfId)
	if intf == nil {
		err := serrors.New("received beacon on non-existent interface",
			"ingress_interface", b.InIfId)
		return err
	}

	upstream := intf.TopoInfo().IA
	if span != nil {
		span.SetTag("ingress_interface", b.InIfId)
		span.SetTag("upstream", upstream)
	}
	logger := log.FromCtx(ctx).New("beacon", b, "upstream", upstream)
	ctx = log.CtxWith(ctx, logger)

	logger.Debug("Received beacon", "bcn", b)
	//if err := h.IngressDB.PreFilter(b); err != nil {
	//	logger.Debug("Beacon pre-filtered", "err", err)
	//	return err
	//}
	if err := h.validateASEntry(b, intf); err != nil {
		logger.Info("Beacon validation failed", "err", err)
		return err
	}
	if err := h.verifySegment(ctx, b.Segment, peer); err != nil {
		logger.Info("Beacon verification failed", "err", err)
		return serrors.WrapStr("verifying beacon", err)
	}
	if len(b.Segment.ASEntries) == 0 { // Should not happen
		logger.Info("Not enough AS entries to process")
		return serrors.New("Not enough AS entries to process")
	}
	// Check if all algorithm ids in the as entry extensions are equal
	// Not needed anymore, as it is possible for hops to not have Irec.
	//if err := h.validateAlgorithmHash(b.Segment); err != nil {
	//	logger.Info("Beacon verification failed", "err", err)
	//	return serrors.WrapStr("verifying beacon", err)
	//}

	// Verification checks passed, now check if the algorithm is known
	if err := h.checkAndFetchAlgorithm(b.Segment, peer); err != nil {
		logger.Info("Beacon verification failed", "err", err)
		return serrors.WrapStr("verifying beacon", err)
	}
	// Insert with algorithm id and origin intfgroup
	if _, err := h.IngressDB.InsertBeacon(ctx, b); err != nil {
		logger.Debug("Failed to insert beacon", "err", err)
		return serrors.WrapStr("inserting beacon", err)
	}
	logger.Debug("Inserted beacon")
	return nil
}

func (h Handler) checkAndFetchAlgorithm(segment *seg.PathSegment, peer *snet.UDPAddr) error {
	//exists, err := h.AlgorithmDB.Exists(segment.ASEntries[0].Extensions.Irec.AlgorithmHash)
	//if err != nil {
	//	return serrors.WrapStr("checking algorithm exists failed", err)
	//}
	//if exists {
	//	return nil
	//}
	//Schedule async retrieval of algorithm. Also ensure we don't fetch multiple times (e.g. through a timeout)

	// We can't just hand-over algorithm over hop-by-hop, as a hop may not support Irec..
	// reverse the path of the beacon: see scion/private/svc/svc.go
	new_path, _ := svc.ReversePath(peer.Path)
	log.Info("Reversed path", "algo", new_path)
	//log.Info(peer.Path)
	return nil
}

func (h Handler) validateAlgorithmHash(segment *seg.PathSegment) error {
	originASEntry := segment.ASEntries[0]
	for _, asEntry := range segment.ASEntries {
		if asEntry.Extensions.Irec == nil {
			return serrors.New(" algorithm extension not found")
		}
		if !bytes.Equal(asEntry.Extensions.Irec.AlgorithmHash, originASEntry.Extensions.Irec.AlgorithmHash) {
			return serrors.New("algorithm hash is different between AS entries")
		}
		if asEntry.Extensions.Irec.AlgorithmId != originASEntry.Extensions.Irec.AlgorithmId {
			return serrors.New("algorithm id is different between AS entries")
		}
	}
	return nil
}

func (h Handler) validateASEntry(b beacon.Beacon, intf *ifstate.Interface) error {
	topoInfo := intf.TopoInfo()
	if topoInfo.LinkType != topology.Parent && topoInfo.LinkType != topology.Core {
		return serrors.New("beacon received on invalid link",
			"ingress_interface", b.InIfId, "link_type", topoInfo.LinkType)
	}
	asEntry := b.Segment.ASEntries[b.Segment.MaxIdx()]
	if !asEntry.Local.Equal(topoInfo.IA) {
		return serrors.New("invalid upstream ISD-AS",
			"expected", topoInfo.IA, "actual", asEntry.Local)
	}
	if !asEntry.Next.Equal(h.LocalIA) {
		return serrors.New("next ISD-AS of upstream AS entry does not match local ISD-AS",
			"expected", h.LocalIA, "actual", asEntry.Next)
	}
	return nil
}

func (h Handler) verifySegment(ctx context.Context, segment *seg.PathSegment,
	peer *snet.UDPAddr) error {

	peerPath, err := peer.GetPath()
	if err != nil {
		return err
	}
	svcToQuery := &snet.SVCAddr{
		IA:      peer.IA,
		Path:    peerPath.Dataplane(),
		NextHop: peerPath.UnderlayNextHop(),
		SVC:     addr.SvcCS,
	}
	return segverifier.VerifySegment(ctx, h.Verifier, svcToQuery, segment)
}
