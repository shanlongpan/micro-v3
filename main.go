/**
* @Author:Tristan
* @Date: 2021/6/15 9:55 下午
 */
package main

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/plugins/wrapper/ratelimiter/uber/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/util/log"
	"github.com/shanlongpan/micro-v3-pub/idl/grpc/microv3"
	"github.com/shanlongpan/micro-v3/handler"
	"github.com/shanlongpan/micro-v3/subscriber"
	"time"
)

func main() {
	// New Service
	reg := etcd.NewRegistry(registry.Addrs("http://127.0.0.1:2377", "http://127.0.0.1:2378", "http://127.0.0.1:2379"))

	service := micro.NewService(
		micro.Name("micro-v3-learn"),
		micro.Version("latest"),
		micro.RegisterTTL(10*time.Second),
		micro.RegisterInterval(5*time.Second),
		micro.WrapHandler(ratelimit.NewHandlerWrapper(1000)), // 限流1000每s,
		micro.Registry(reg),
		//micro.Server(grpc.NewServer()),
	)

	service.Init()
	// Register Handler
	err := microv3.RegisterMicroV3ServiceHandler(service.Server(), new(handler.Microv3))
	if err != nil {
		log.Errorf("init fail %s", err.Error())
		return
	}
	// Register Struct as Subscriber  可以不注册
	err = micro.RegisterSubscriber("palfish.com.service.apollo", service.Server(), new(subscriber.Subscribe))
	if err != nil {
		log.Errorf("sub fail %s", err.Error())
		return
	}
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
