/**
* @Author:Tristan
* @Date: 2021/6/15 9:56 下午
 */

package handler

import (
	"context"
	"github.com/asim/go-micro/v3/util/log"
	helloworld "github.com/shanlongpan/micro-v3-pub/idl/micro-grpc"
)

type HelloWorld struct{}

// Call is a single request handler called via client.Call or the generated client code
func (e *HelloWorld) Call(ctx context.Context, req *helloworld.Request, rsp *helloworld.Response) error {
	log.Info("Received helloworld.Call request")
	rsp.Msg = "Hello " + req.Name
	return nil
}

// Stream is a server side stream handler called via client.Stream or the generated client code
func (e *HelloWorld) Stream(ctx context.Context, req *helloworld.StreamingRequest, stream helloworld.Helloworld_StreamStream) error {
	log.Infof("Received helloworld.Stream request with count: %d", req.Count)

	for i := 0; i < int(req.Count); i++ {
		log.Infof("Responding: %d", i)
		if err := stream.Send(&helloworld.StreamingResponse{
			Count: int64(i),
		}); err != nil {
			return err
		}
	}

	return nil
}

// PingPong is a bidirectional stream handler called via client.Stream or the generated client code
func (e *HelloWorld) PingPong(ctx context.Context, stream helloworld.Helloworld_PingPongStream) error {
	for {
		req, err := stream.Recv()
		if err != nil {
			return err
		}
		log.Infof("Got ping %v", req.Stroke)
		if err := stream.Send(&helloworld.Pong{Stroke: req.Stroke}); err != nil {
			return err
		}
	}
}
