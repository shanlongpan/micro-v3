/**
* @Author:Tristan
* @Date: 2021/6/15 9:55 下午
 */
package main

import (
	"github.com/asim/go-micro/plugins/registry/etcd/v3"
	"github.com/asim/go-micro/v3"
	"github.com/asim/go-micro/v3/registry"
	"github.com/asim/go-micro/v3/util/log"
	helloworld "github.com/shanlongpan/micro-v3-pub/idl/micro-grpc"
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
		micro.Registry(reg),
	)

	service.Init()
	// Register Handler
	err := helloworld.RegisterHelloworldHandler(service.Server(), new(handler.HelloWorld))
	if err != nil {
		log.Errorf("init fail %s", err.Error())
		return
	}
	// Register Struct as Subscriber  可以不注册
	micro.RegisterSubscriber("palfish.com.service.apollo", service.Server(), new(subscriber.Apollo))
	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
