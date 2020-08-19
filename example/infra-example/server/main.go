// Code generated by 'freedom new-project infra-example'
package main

import (
	"fmt"
	"sort"
	"time"

	"github.com/8treenet/freedom"
	_ "github.com/8treenet/freedom/example/infra-example/adapter/controllers"
	"github.com/8treenet/freedom/example/infra-example/server/conf"
	"github.com/8treenet/freedom/infra/requests"
	"github.com/8treenet/freedom/middleware"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

func main() {
	app := freedom.NewApplication()
	installDatabase(app)
	installLogger(app)

	installMiddleware(app)
	addrRunner := app.CreateRunner(conf.Get().App.Other["listen_addr"].(string))
	app.Run(addrRunner, *conf.Get().App)
}

func installMiddleware(app freedom.Application) {
	app.InstallMiddleware(middleware.NewRecover())
	app.InstallMiddleware(middleware.NewTrace("x-request-id"))
	app.InstallMiddleware(middleware.NewRequestLogger("x-request-id"))

	app.InstallBusMiddleware(middleware.NewBusFilter())
	requests.InstallPrometheus(conf.Get().App.Other["service_name"].(string), freedom.Prometheus())
}

func installDatabase(app freedom.Application) {
	app.InstallDB(func() interface{} {
		conf := conf.Get().DB
		db, e := gorm.Open("mysql", conf.Addr)
		if e != nil {
			freedom.Logger().Fatal(e.Error())
		}

		db.DB().SetMaxIdleConns(conf.MaxIdleConns)
		db.DB().SetMaxOpenConns(conf.MaxOpenConns)
		db.DB().SetConnMaxLifetime(time.Duration(conf.ConnMaxLifeTime) * time.Second)
		return db
	})
}

func installLogger(app freedom.Application) {
	//logger中间件，每一行日志都会触发回调。如果返回true，将停止中间件遍历回调。
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
