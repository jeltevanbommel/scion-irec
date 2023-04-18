package wasm

import (
	"context"
	"fmt"
	"github.com/bytecodealliance/wasmtime-go"
	"github.com/scionproto/scion/pkg/log"
	"github.com/scionproto/scion/pkg/proto/control_plane"
	"github.com/scionproto/scion/pkg/proto/rac"
	"google.golang.org/protobuf/proto"
	"time"
	"unsafe"
)

type WasmEnv struct {
	config *wasmtime.Config
	Engine *wasmtime.Engine
}

func (w *WasmEnv) Initialize() {
	w.config = wasmtime.NewConfig()
	w.Engine = wasmtime.NewEngineWithConfig(w.config)
}

func guestWriteBytes(instance *wasmtime.Instance, store *wasmtime.Store, bytes []byte) (int32, error) {
	mem := instance.GetExport(store, "memory").Memory()
	// Use alloc in the WASM module to allocate memory for the data.
	ptr, err := instance.GetExport(store, "__alloc").Func().Call(store, len(bytes), 0)
	if err != nil {
		return 0, err
	}
	// This yields a pointer, which is where the data will be placed.
	//fmt.Println("Pointer: ", ptr)
	dst := unsafe.Pointer(uintptr(mem.Data(store)) + uintptr(ptr.(int32)))
	// Write the data to the linear memory at this pointer.
	copy(unsafe.Slice((*byte)(dst), len(bytes)), bytes)
	return ptr.(int32), nil
}

// AdditionalInfo is given when called
type AdditionalInfo struct {
	PropagationInterfaces []uint32
}

func (w *WasmEnv) Execute(ctx context.Context, count int32, beacons []*control_plane.IRECBeacon, module *wasmtime.Module, info AdditionalInfo, callback func(response *rac.RACResponse)) {
	startModule := time.Now()
	store := wasmtime.NewStore(w.Engine) //TODO can't reuse these if we use fuel mechanism.
	linker := wasmtime.NewLinker(w.Engine)

	wasiconfig := wasmtime.NewWasiConfig()
	wasiconfig.SetEnv([]string{"WASMTIME"}, []string{"GO"})
	store.SetWasi(wasiconfig)
	if err := linker.DefineWasi(); err != nil {
		log.Error("initializing wasi env", "err", err)
		return
	}
	//hookKvFns(linker, store)
	//err := linker.DefineFunc(store, "env", "__internal_print", func(caller *wasmtime.Caller, offset int32, length int32) {
	//	if offset < 0 {
	//		log.Info("Invalid offset")
	//	}
	//	if length < 0 {
	//		log.Info("Invalid length")
	//	}
	//	mem := caller.GetExport("memory").Memory()
	//	log.Info(fmt.Sprintf("[RAC]: %s", string(mem.UnsafeData(store)[offset:offset+length])))
	//})
	//if err != nil {
	//	log.Error("hooking functions", "err", err)
	//	return
	//}
	startInstance := time.Now()
	instance1, err := linker.Instantiate(store, module)

	if err != nil {
		log.Error("instantiating module", "err", err)
		return
	}

	// We trust that the signed fields we receive are accurate and do not verify these ourselves
	// for performance reasons.
	// We therefore disregard the signed entry, and only include the content of the signed body
	startJuggling := time.Now()
	irecBeacons := rac.RACRequest{
		Beacons:     make([]*rac.IRECBeaconSB, 0, len(beacons)),
		EgressIntfs: info.PropagationInterfaces,
	}
	// 'Juggle' the beacons, such that we do not send the signed bytes to the RACs.
	for _, bcn := range beacons {
		irecASEntries := make([]*rac.IRECASEntry, 0, len(bcn.PathSeg.AsEntries))
		for _, asEntry := range bcn.PathSeg.AsEntries {
			irecASEntries = append(irecASEntries, &rac.IRECASEntry{
				Signed:   asEntry.SignedBody,
				Unsigned: asEntry.Unsigned,
			})
		}
		beacon := &rac.IRECBeaconSB{
			Segment: &rac.IRECPathSegment{
				SegmentInfo: bcn.PathSeg.SegmentInfo,
				AsEntries:   irecASEntries,
			},
			Id: bcn.Id,
		}
		irecBeacons.Beacons = append(irecBeacons.Beacons, beacon)
	}
	startMarshal := time.Now()
	// Now turn it into a protobuf binary message
	bytes, err := proto.Marshal(&irecBeacons)
	if err != nil {
		log.Error("marshalling beacons", "err", err)
		return
	}

	startMemory := time.Now()
	// And write to the guest memory, which will then be able to read it
	ptr, err := guestWriteBytes(instance1, store, bytes)
	if err != nil {
		log.Error("writing beacons to env", "err", err)
		return
	}

	startExec := time.Now()
	// Through the provided pointer to the bytes of the protobuf message.
	run := instance1.GetExport(store, "run").Func()
	_, err = run.Call(store, ptr, len(bytes))

	if err != nil {
		log.Error("running env", "err", err)
		return
	}
	startReadingRac := time.Now()
	// Len and ptr functions are a bypass to prevent a memory leak in WASMtime.
	length, err := instance1.GetExport(store, "len").Func().Call(store)
	if err != nil {
		log.Error("running env", "err", err)
		return
	}
	offset, err := instance1.GetExport(store, "ptr").Func().Call(store)
	if err != nil {
		log.Error("running env", "err", err)
		return
	}
	// Get the response from the memory;
	mem := instance1.GetExport(store, "memory").Memory()
	racResult := new(rac.RACResponse)
	err = proto.Unmarshal(mem.UnsafeData(store)[offset.(int32):offset.(int32)+length.(int32)], racResult)
	if err != nil {
		log.Info("Invalid rac response")
	}
	startEgressCall := time.Now()
	// And write them to the egress.
	callback(racResult) //TODO(jvb): callback is now unneccessary.
	fmt.Printf("%d, wasm=%d, %d, %d, %d, %d, %d, %d\n", count, startInstance.Sub(startModule).Nanoseconds(), startJuggling.Sub(startInstance).Nanoseconds(), startMarshal.Sub(startJuggling).Nanoseconds(), startMemory.Sub(startMarshal).Nanoseconds(), startExec.Sub(startMemory).Nanoseconds(), startReadingRac.Sub(startExec).Nanoseconds(), startEgressCall.Sub(startReadingRac).Nanoseconds())
}
