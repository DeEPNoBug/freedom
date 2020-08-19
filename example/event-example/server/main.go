// Code generated by 'freedom new-project event-example'
package main

import (
	"fmt"
	"sort"
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
	installLogger(app)
	installMiddleware(app)
	addrRunner := app.CreateH2CRunner(conf.Get().App.Other["listen_addr"].(string))

	// Obtain and install the kafka infrastructure for domain events
	app.InstallDomainEventInfra(kafka.GetDomainEventInfra())
	app.Run(addrRunner, *conf.Get().App)
}

func installMiddleware(app freedom.Application) {
	app.InstallMiddleware(middleware.NewRecover())
	app.InstallMiddleware(middleware.NewTrace("x-request-id"))
	app.InstallMiddleware(middleware.NewRequestLogger("x-request-id"))

	app.InstallBusMiddleware(middleware.NewBusFilter())
	requests.InstallPrometheus(conf.Get().App.Other["service_name"].(string), freedom.Prometheus())
}

func installLogger(app freedom.Application) {
	//logger中间件，每一行日志都会触发回调，返回true停止。
	app.Logger().Handle(func(value *freedom.LogRow) bool {
		fieldKeys := []string{}
		for k := range value.Fields {
			fieldKeys = append(fieldKeys, k)
		}
		sort.Strings(fieldKeys)
		for i := 0; i < len(fieldKeys); i++ {
			fieldMsg := value.Fields[fieldKeys[i]]
			if value.Message != "" {
				value.Message += " "
			}
			value.Message += fmt.Sprintf("%s:%v", fieldKeys[i], fieldMsg)
		}
		return false

		/*
			logrus.WithFields(value.Fields).Info(value.Message)
			return true
		*/
		/*
			zapLogger, _ := zap.NewProduction()
			zapLogger.Info(value.Message)
			return true
		*/
	})
}
