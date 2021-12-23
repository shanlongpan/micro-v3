/**
* @Author:Tristan
* @Date: 2021/6/24 9:27 下午
 */

package subscriber

import (
	"context"
	"github.com/asim/go-micro/v3/util/log"
	helloworld "github.com/shanlongpan/micro-v3-pub/idl/micro-grpc"
)

type Apollo struct{}

func (e *Apollo) Handle(ctx context.Context, msg *helloworld.Message) error {
	log.Info("Handler Received message: ", msg.Say)
	return nil
}

func Handler(ctx context.Context, msg *helloworld.Message) error {
	log.Info("Function Received message: ", msg.Say)
	return nil
}