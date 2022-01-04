/**
* @Author:Tristan
* @Date: 2021/6/15 9:56 下午
 */

package handler

import (
	"context"
	"fmt"
	"github.com/asim/go-micro/v3/metadata"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/shanlongpan/micro-v3-pub/idl/grpc/microv3"
)

type Microv3 struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *Microv3) Call(ctx context.Context, req *microv3.CallRequest, rsp *microv3.CallResponse) error {
	log.Info("Received Microv3.Call request")
	traceId, _ := metadata.Get(ctx, "trace_id")
	rsp.Msg = fmt.Sprintf("name %s trace_id %s", req.Name, traceId)
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *Microv3) Stream(ctx context.Context, req *microv3.StreamingRequest, stream microv3.MicroV3Service_StreamStream) error {
	log.Infof("Received Microv3.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&microv3.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *Microv3) PingPong(ctx context.Context, stream microv3.MicroV3Service_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&microv3.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
