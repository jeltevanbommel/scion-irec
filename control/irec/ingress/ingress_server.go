package ingress

import (
	"context"
	"fmt"
	"github.com/scionproto/scion/control/beacon"
	"github.com/scionproto/scion/pkg/log"
	"github.com/scionproto/scion/pkg/private/common"
	"github.com/scionproto/scion/pkg/private/serrors"
	cppb "github.com/scionproto/scion/pkg/proto/control_plane"
	seg "github.com/scionproto/scion/pkg/segment"
	"github.com/scionproto/scion/pkg/slayers/path/scion"
	"github.com/scionproto/scion/pkg/snet"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"net"
)

type IncomingBeaconHandler interface {
	HandleBeacon(ctx context.Context, b beacon.Beacon, peer *snet.UDPAddr) error
}

type IngressServer struct {
	IncomingHandler       IncomingBeaconHandler
	IngressDB             IngressStore
	RACManager            RACManager
	PropagationInterfaces []uint32
}

func (i *IngressServer) GetJob(ctx context.Context, request *cppb.RACBeaconRequest) (*cppb.RACJob, error) {
	bcns, err := i.IngressDB.GetAndMarkBeacons(ctx, &cppb.RACBeaconRequest{
		Maximum:         0,
		AlgorithmHash:   request.AlgorithmHash,
		AlgorithmID:     request.AlgorithmID,
		OriginAS:        request.OriginAS,
		OriginIntfGroup: request.OriginIntfGroup,
		IgnoreIntfGroup: false,
	})
	if err != nil {
		log.Error("An error occurred when retrieving beacons from db", "err", err)
		return &cppb.RACJob{}, err
	}
	log.Debug("Queueing to RAC", "beacons", len(bcns.Beacons))
	return &cppb.RACJob{
		Beacons:       bcns.Beacons,
		AlgorithmHash: request.AlgorithmHash,
		//Intfs:         nil,
		PropIntfs: i.PropagationInterfaces,
	}, nil
}

func (i *IngressServer) Handle(ctx context.Context, req *cppb.IncomingBeacon) (*cppb.IncomingBeaconResponse, error) {
	gPeer, ok := peer.FromContext(ctx)
	if !ok {
		return nil, serrors.New("peer must exist")
	}
	logger := log.FromCtx(ctx)
	logger.Debug("Received Beacon from " + fmt.Sprintf("%T", gPeer.Addr.(*snet.UDPAddr)) + ", msg:" + req.String())

	peer, ok := gPeer.Addr.(*snet.UDPAddr)
	if !ok {
		logger.Debug("peer must be *snet.UDPAddr", "actual", fmt.Sprintf("%T", gPeer))
		return nil, serrors.New("peer must be *snet.UDPAddr", "actual", fmt.Sprintf("%T", gPeer))
	}
	ingress, err := extractIngressIfID(peer.Path)
	if err != nil {
		logger.Debug("Failed to extract ingress interface", "peer", peer, "err", err)
		return nil, status.Error(codes.InvalidArgument, "failed to extract ingress interface")
	}
	ps, err := seg.FullBeaconFromPB(req.Segment)
	if err != nil {
		logger.Debug("Failed to parse beacon", "peer", peer, "err", err)
		return nil, status.Error(codes.InvalidArgument, "failed to parse beacon")
	}
	b := beacon.Beacon{
		Segment: ps,
		InIfId:  ingress,
	}
	if err := i.IncomingHandler.HandleBeacon(ctx, b, peer); err != nil {
		logger.Debug("Failed to handle beacon", "peer", peer, "err", err)
		return nil, serrors.WrapStr("handling beacon", err)
	}
	return &cppb.IncomingBeaconResponse{}, nil
}

func (i *IngressServer) RegisterRAC(ctx context.Context, req *cppb.RegisterRACRequest) (*cppb.RegisterRACResponse, error) {
	logger := log.FromCtx(ctx)
	logger.Debug("Registration request incoming for RAC", "addr", req.Addr)
	// Resolve the register request address
	ip, err := net.ResolveTCPAddr("tcp", req.Addr)
	if err != nil {
		log.Error("error resolving", "err", err)
		return &cppb.RegisterRACResponse{}, err
	}

	// Check if the rac register request comes from the registering rac
	gPeer, ok := peer.FromContext(ctx)
	if !ok {
		return &cppb.RegisterRACResponse{}, serrors.New("peer must exist")
	}

	_, host, err := net.SplitHostPort(gPeer.Addr.String())
	if err != nil {
		log.Error("error resolving", "err", err)
		return &cppb.RegisterRACResponse{}, err
	}
	_, receivedHost, err := net.SplitHostPort(req.Addr)
	if err != nil {
		log.Error("error resolving", "err", err)
		return &cppb.RegisterRACResponse{}, err
	}
	// TODO(jvb) Fix this check
	if host != receivedHost && false {
		log.Error("peer address is not the same as received address for registration", req.Addr, gPeer)
		return &cppb.RegisterRACResponse{}, serrors.New("peer address is not the same as received registration address", "actual", fmt.Sprintf("%T %T", gPeer, req.Addr))
	}

	i.RACManager.AddRAC(req.Addr, ip)
	return &cppb.RegisterRACResponse{}, nil
}

func (i *IngressServer) BeaconSources(ctx context.Context, request *cppb.RACBeaconSourcesRequest) (*cppb.RACBeaconSources, error) {
	log.Info("request incoming for beacons", "req", request.IgnoreIntfGroup)

	beaconSources, err := i.IngressDB.BeaconSources(ctx, request.IgnoreIntfGroup)
	if err != nil {
		log.Error("error happened oh no", "err", err)
	}
	log.Info("beacon sources", "src", beaconSources)
	return &cppb.RACBeaconSources{Sources: beaconSources}, err
}

func (i *IngressServer) RequestBeacons(ctx context.Context, req *cppb.RACBeaconRequest) (*cppb.RACBeaconResponse, error) {
	return i.IngressDB.GetBeacons(ctx, req)
}

// extractIngressIfID extracts the ingress interface ID from a path.
func extractIngressIfID(path snet.DataplanePath) (uint16, error) {
	invertedPath, ok := path.(snet.RawReplyPath)
	if !ok {
		return 0, serrors.New("unexpected path", "type", common.TypeOf(path))
	}
	rawScionPath, ok := invertedPath.Path.(*scion.Raw)
	if !ok {
		return 0, serrors.New("unexpected path", "type", common.TypeOf(path))
	}
	hf, err := rawScionPath.GetCurrentHopField()
	if err != nil {
		return 0, serrors.WrapStr("getting current hop field", err)
	}
	return hf.ConsIngress, nil
}
