/**
* @Author:Tristan
* @Date: 2021/6/24 9:27 下午
 */

package subscriber

import (
	"context"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/shanlongpan/micro-v3-pub/idl/grpc/microv3"
)

type Subscribe struct{}

func (e *Subscribe) Handle(ctx context.Context, msg *microv3.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *microv3.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}