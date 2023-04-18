package ingress

import (
	"context"
	"fmt"
	libgrpc "github.com/scionproto/scion/pkg/grpc"
	"github.com/scionproto/scion/pkg/log"
	racservice "github.com/scionproto/scion/pkg/proto/rac"
	"google.golang.org/grpc"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type racstate int32

const (
	// Registered etc.
	Free int32 = iota
	Busy
)

type RACInfo struct {
	State int32
	Addr  net.Addr
	// TODO: algorithms known etc.
}

// Deprecated construct, used to keep the state of which RACs are connected to the ingress. Kept for future usage.
type RACManager struct {
	RACList       *sync.Map
	RACListUnsafe map[string]*RACInfo
	RACCount      *atomic.Uint64
	libgrpc.Dialer
}
type JobQueue struct {
	sync.Mutex
	Items []*racservice.ExecutionRequest
}

func (q *JobQueue) Push(job *racservice.ExecutionRequest) {
	q.Lock()
	defer q.Unlock()
	q.Items = append(q.Items, job)
}

func (q *JobQueue) Pop() *racservice.ExecutionRequest {
	q.Lock()
	defer q.Unlock()
	if len(q.Items) == 0 {
		return nil
	}
	item := q.Items[0]
	q.Items = q.Items[1:]
	return item
}

func (r *RACManager) AddRAC(id string, addr net.Addr) {
	if _, ok := r.RACList.Load(id); !ok {
		ri := RACInfo{State: Free, Addr: addr}

		r.RACList.Store("id", &ri)
		r.RACCount.Add(1)
		r.RACListUnsafe["id"] = &ri
	}
}

func (r *RACManager) Queue(ctx context.Context, request *racservice.ExecutionRequest) {
	//r
}

func (r *RACManager) RunImm(ctx context.Context, request *racservice.ExecutionRequest, wg *sync.WaitGroup) {
	wg.Add(1) //todo maybe this waitgroup handover is the issue?
	go func() {
		defer log.HandlePanic()
		defer wg.Done()
		completed := false
		for !completed {
			r.RACList.Range(func(k, v any) bool {
				rac := v.(*RACInfo)
				if rac.State == Busy {
					return true // Try next rac
				}
				rac.State = Busy

				fmt.Println("\t[RAC]=d", time.Now().UnixNano())
				conn, err := grpc.DialContext(ctx, rac.Addr.String(),
					grpc.WithInsecure(),
					grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*12)),
				)

				if err != nil {
					log.Error("error dialing RAC registration", "err", err)
					return true // try next rac
				}
				defer conn.Close()
				client := racservice.NewRACServiceClient(conn)
				opts := make([]grpc.CallOption, 0)
				opts = append(opts, grpc.MaxCallRecvMsgSize(1024*1024*12))
				opts = append(opts, libgrpc.RetryProfile...)
				fmt.Println("\t[RAC]=dc", time.Now().UnixNano())
				_, err = client.Execute(ctx, request, opts...)
				if err != nil {
					log.Debug("Was not queued, due to ", "err", err)
				}
				rac.State = Free
				completed = true
				return false // do not try next rac, we are done.
			})
			//for _, v := range r.RACListUnsafe {
			//	if !atomic.CompareAndSwapInt32(&v.State, Free, Busy) {
			//		continue
			//	}
			//	conn, err := grpc.DialContext(ctx, v.Addr.String(),
			//		grpc.WithInsecure(),
			//		grpc.WithDefaultCallOptions(grpc.MaxCallRecvMsgSize(1024*1024*12)),
			//	)
			//
			//	if err != nil {
			//		log.Error("error dialing RAC registration", "err", err)
			//		return
			//	}
			//	defer conn.Close()
			//	client := racservice.NewRACServiceClient(conn)
			//	opts := make([]grpc.CallOption, 0)
			//	opts = append(opts, grpc.MaxCallRecvMsgSize(1024*1024*12))
			//	opts = append(opts, libgrpc.RetryProfile...)
			//	_, err = client.Execute(ctx, request, opts...)
			//	if err != nil {
			//		log.Debug("Was not queued, due to ", "err", err)
			//	}
			//	v.State = Free
			//	completed = true
			//	break
			//
			//}
		}
	}()
}
