/**
* @Author:Tristan
* @Date: 2021/6/15 9:55 下午
 */
package main

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/server/grpc/v3"
	ratelimit "github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/server"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/shanlongpan/micro-v3-pub/idl/grpc/microv3"
	"github.com/shanlongpan/micro-v3/consts"
	"github.com/shanlongpan/micro-v3/handler"
	"github.com/shanlongpan/micro-v3/subscriber"
)

func main() {
	// New Service
	reg := etcd.NewRegistry(registry.Addrs("http://127.0.0.1:2377", "http://127.0.0.1:2378", "http://127.0.0.1:2379"))

	service := micro.NewService(
		micro.Server(
			grpc.NewServer(
				server.Name("micro-v3-learn"),
				server.Registry(reg),
				server.Version("latest")),
		),
	)
	service.Init(
		micro.WrapHandler(ratelimit.NewHandlerWrapper(consts.RateLimit)), // 针对单个服务限流1000每s,如果多台，能相应的请求数*N
	)

	// Register Handler
	err := microv3.RegisterMicroV3ServiceHandler(service.Server(), new(handler.Microv3))
	if err != nil {
		log.Errorf("init fail %s", err.Error())
		return
	}
	// Register Struct as Subscriber  可以不注册
	err = micro.RegisterSubscriber("micro-v3-learn-subscribe", service.Server(), new(subscriber.Subscribe))
	if err != nil {
		log.Errorf("sub fail %s", err.Error())
		return
	}
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
