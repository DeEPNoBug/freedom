// Code generated by 'freedom new-project event-example'
package main

import (
	"fmt"
	"time"

	"github.com/8treenet/freedom"
	_ "github.com/8treenet/freedom/example/event-example/adapter/controllers"
	"github.com/8treenet/freedom/example/event-example/server/conf"
	"github.com/8treenet/freedom/infra/kafka"
	"github.com/8treenet/freedom/infra/requests"
	"github.com/8treenet/freedom/middleware"
	"github.com/Shopify/sarama"
)

// mac: Start kafka: zookeeper-server-start /usr/local/etc/kafka/zookeeper.properties & kafka-server-start /usr/local/etc/kafka/server.properties
func main() {
	// If you use the default Kafka configuration, no need to set
	kafka.SettingConfig(func(conf *sarama.Config, other map[string]interface{}) {
		conf.Producer.Retry.Max = 3
		conf.Producer.Retry.Backoff = 5 * time.Second
		conf.Consumer.Offsets.Initial = sarama.OffsetOldest
		fmt.Println(other)
	})
	app := freedom.NewApplication()
	installMiddleware(app)
	addrRunner := app.CreateH2CRunner(conf.Get().App.Other["listen_addr"].(string))

	// Obtain and install the kafka infrastructure for domain events
	app.InstallDomainEventInfra(kafka.GetDomainEventInfra())
	app.Run(addrRunner, *conf.Get().App)
}

func installMiddleware(app freedom.Application) {
	//Recover中间件
	app.InstallMiddleware(middleware.NewRecover())
	//Trace链路中间件
	app.InstallMiddleware(middleware.NewTrace("x-request-id"))
	//日志中间件，每个请求一个logger
	app.InstallMiddleware(middleware.NewRequestLogger("x-request-id"))
	//logRow中间件，每一行日志都会触发回调。如果返回true，将停止中间件遍历回调。
	app.Logger().Handle(middleware.DefaultLogRowHandle)
	//HttpClient 普罗米修斯中间件，监控下游的API请求。
	requests.InstallPrometheus(conf.Get().App.Other["service_name"].(string), freedom.Prometheus())
	//总线中间件，处理上下游透传的Header
	app.InstallBusMiddleware(middleware.NewBusFilter())
}
