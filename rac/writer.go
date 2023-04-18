package rac

import (
	"context"
	"github.com/scionproto/scion/pkg/grpc"
	cppb "github.com/scionproto/scion/pkg/proto/control_plane"
	grpc2 "google.golang.org/grpc"
)

type Writer struct {
	Conn *grpc2.ClientConn
}

func (w *Writer) writeBeacons(ctx context.Context, beacons []*cppb.StoredBeaconAndIntf) error {
	client := cppb.NewEgressServiceClient(w.Conn)
	_, err := client.RequestPropagation(ctx, &cppb.PropagationRequest{
		Beacon: beacons,
	}, grpc.RetryProfile...)
	return err
}
