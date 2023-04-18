package ingress

import (
	"context"
	"github.com/scionproto/scion/control/beacon"
	"github.com/scionproto/scion/pkg/addr"
	"github.com/scionproto/scion/pkg/log"
	"github.com/scionproto/scion/pkg/private/serrors"
	cppb "github.com/scionproto/scion/pkg/proto/control_plane"
)

type IngressStore interface {
	GetBeacons(ctx context.Context, req *cppb.RACBeaconRequest) (*cppb.RACBeaconResponse, error)
	GetAndMarkBeacons(ctx context.Context, req *cppb.RACBeaconRequest) (*cppb.RACBeaconResponse, error)
	InsertBeacon(ctx context.Context, b beacon.Beacon) (beacon.InsertStats, error)
	BeaconSources(ctx context.Context, ignoreIntfGroup bool) ([]*cppb.RACBeaconSource, error)
}

type AlgorithmStorage interface {
	Exists(algHash []byte) (bool, error)
	Store(algHash []byte, algorithm []byte) error
	Retrieve(algHash []byte) ([]byte, error)
}

// Beaconstore
type Store struct {
	baseStore
	policies beacon.Policies
}

// NewBeaconStore creates a new beacon store for the ingress gateway
func NewIngressDB(policies beacon.Policies, db DB) (*Store, error) {
	policies.InitDefaults()
	if err := policies.Validate(); err != nil {
		return nil, err
	}
	s := &Store{
		baseStore: baseStore{
			db: db,
		},
		policies: policies,
	}
	s.baseStore.usager = &s.policies
	return s, nil
}

type usager interface {
	Filter(beacon beacon.Beacon) error
	Usage(beacon beacon.Beacon) beacon.Usage
}

type BeaconSource struct {
	IA            addr.IA
	IntfGroup     uint16
	AlgorithmHash []byte
	AlgorithmId   uint32
}

// DB defines the interface that all beacon DB backends have to implement.
type DB interface {
	// CandidateBeacons returns up to `setSize`. The beacons in the slice are ordered by segment length from
	// shortest to longest.
	GetBeacons(ctx context.Context, maximum uint32, algHash []byte, algID uint32, originAS addr.IA, originIntfGroup uint32, ignoreIntfGroup bool) ([]*cppb.IRECBeacon, error)
	GetAndMarkBeacons(ctx context.Context, maximum uint32, algHash []byte, algID uint32, originAS addr.IA, originIntfGroup uint32, ignoreIntfGroup bool, marker uint32) ([]*cppb.IRECBeacon, error)
	BeaconSources(ctx context.Context, ignoreIntfGroup bool) ([]*cppb.RACBeaconSource, error)
	// Insert inserts a beacon with its allowed usage into the database.
	InsertBeacon(ctx context.Context, b beacon.Beacon, usage beacon.Usage) (beacon.InsertStats, error)
}

type baseStore struct {
	db     DB
	usager usager
}

func (s *baseStore) InsertBeacon(ctx context.Context, b beacon.Beacon) (beacon.InsertStats, error) {
	usage := s.usager.Usage(b)
	if usage.None() {
		return beacon.InsertStats{Filtered: 1}, nil
	}
	return s.db.InsertBeacon(ctx, b, usage)
}

func (s *baseStore) GetAndMarkBeacons(ctx context.Context, req *cppb.RACBeaconRequest) (*cppb.RACBeaconResponse, error) {
	log.Info("getting beacons + marking for", "as", addr.IA(req.OriginAS))
	// TODO(jvb) different markers?
	beacons, err := s.db.GetAndMarkBeacons(ctx, req.Maximum, req.AlgorithmHash, req.AlgorithmID, addr.IA(req.OriginAS), req.OriginIntfGroup, req.IgnoreIntfGroup, 1)
	if err != nil {
		return nil, serrors.WrapStr("retrieving beacons failed", err)
	}

	return &cppb.RACBeaconResponse{Beacons: beacons}, nil
}

func (s *baseStore) GetBeacons(ctx context.Context, req *cppb.RACBeaconRequest) (*cppb.RACBeaconResponse, error) {
	log.Info("getting beacons for", "as", addr.IA(req.OriginAS))
	beacons, err := s.db.GetBeacons(ctx, req.Maximum, req.AlgorithmHash, req.AlgorithmID, addr.IA(req.OriginAS), req.OriginIntfGroup, req.IgnoreIntfGroup)
	if err != nil {
		return nil, serrors.WrapStr("retrieving beacons failed", err)
	}

	return &cppb.RACBeaconResponse{Beacons: beacons}, nil
}

func (s *baseStore) BeaconSources(ctx context.Context, ignoreIntfGroup bool) ([]*cppb.RACBeaconSource, error) {
	return s.db.BeaconSources(ctx, ignoreIntfGroup)
}
func (s *Store) MaxExpTime(policyType beacon.PolicyType) uint8 {
	switch policyType {
	case beacon.UpRegPolicy:
		return *s.policies.UpReg.MaxExpTime
	case beacon.DownRegPolicy:
		return *s.policies.DownReg.MaxExpTime
	case beacon.PropPolicy:
		return *s.policies.Prop.MaxExpTime
	}
	return beacon.DefaultMaxExpTime
}
