package irec

import (
	cppb "github.com/scionproto/scion/pkg/proto/control_plane"
)

type Irec struct {
	AlgorithmHash  []byte
	AlgorithmId    uint32
	InterfaceGroup uint16
}

func ExtensionFromPB(d *cppb.IRECExtension) *Irec {
	if d == nil {
		return nil
	}

	//if d.AlgorithmHash
	return &Irec{
		AlgorithmHash:  d.AlgorithmHash,
		AlgorithmId:    d.AlgorithmId,
		InterfaceGroup: uint16(d.InterfaceGroup),
	}
}

func ExtensionToPB(d *Irec) *cppb.IRECExtension {
	if d == nil {
		return nil
	}

	return &cppb.IRECExtension{
		AlgorithmId:    d.AlgorithmId,
		AlgorithmHash:  d.AlgorithmHash,
		InterfaceGroup: uint32(d.InterfaceGroup),
	}
}
