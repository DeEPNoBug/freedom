package project

func init() {
	content["/server/conf/app.toml"] = appTomlConf()
	content["/server/conf/db.toml"] = dbTomlConf()
	content["/server/conf/redis.toml"] = redisConf()
	content["/server/main.go"] = mainTemplate()
}

func appTomlConf() string {
	return `[other]
listen_addr = ":8000"
service_name = "{{.PackageName}}"
repository_request_timeout = 10
prometheus_listen_addr = ":9090"
# "fatal" "error" "warn" "info"  "debug"
logger_level = "debug"
# shutdown_second : Elegant lying off for the longest time
shutdown_second = 3`
}

func dbTomlConf() string {
	return `addr = "root:123123@tcp(127.0.0.1:3306)/xxxx?charset=utf8&parseTime=True&loc=Local"
max_open_conns = 16
max_idle_conns = 8
conn_max_life_time = 300
`
}

func redisConf() string {
	return `#地址
addr = "127.0.0.1:6379"
#密码
password = ""
#redis 库
db = 0
#重试次数, 默认不重试
max_retries = 0
#连接池大小
pool_size = 32
#读取超时时间 3秒
read_timeout = 3
#写入超时时间 3秒
write_timeout = 3
#连接空闲时间 300秒
idle_timeout = 300
#检测死连接,并清理 默认60秒
idle_check_frequency = 60
#连接最长时间，300秒
max_conn_age = 300
#如果连接池已满 等待可用连接的时间默认 8秒
pool_timeout = 8`
}

func mainTemplate() string {
	return `
	// Code generated by 'freedom new-project {{.PackagePath}}'
	package main

	import (
		"fmt"
		"sort"
		"time"
		_ "github.com/jinzhu/gorm/dialects/mysql"
		"github.com/8treenet/freedom"
		_ "{{.PackagePath}}/adapter/repository" //引入输出适配器 repository资源库
		_ "{{.PackagePath}}/adapter/controller" //引入输入适配器 http路由
		"{{.PackagePath}}/server/conf"
		"github.com/go-redis/redis"
		"github.com/jinzhu/gorm"
		"github.com/8treenet/freedom/middleware"
		"github.com/8treenet/freedom/infra/requests"
	)
	
	func main() {
		app := freedom.NewApplication()
		installLogger(app)
		/*
			installDatabase(app) //安装数据库
			installRedis(app) //安装redis

			http2 h2c 服务
			h2caddrRunner := app.CreateH2CRunner(conf.Get().App.Other["listen_addr"].(string))
		*/
		installMiddleware(app)
		addrRunner := app.CreateRunner(conf.Get().App.Other["listen_addr"].(string))
		app.Run(addrRunner, *conf.Get().App)
	}

	func installMiddleware(app freedom.Application) {
		app.InstallMiddleware(middleware.NewRecover())
		app.InstallMiddleware(middleware.NewTrace("x-request-id"))
		app.InstallMiddleware(middleware.NewRequestLogger("x-request-id"))
		
		requests.InstallPrometheus(conf.Get().App.Other["service_name"].(string), freedom.Prometheus())
		app.InstallBusMiddleware(middleware.NewBusFilter())
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
	
	func installRedis(app freedom.Application) {
		app.InstallRedis(func() (client redis.Cmdable) {
			cfg := conf.Get().Redis
			opt := &redis.Options{
				Addr:               cfg.Addr,
				Password:           cfg.Password,
				DB:                 cfg.DB,
				MaxRetries:         cfg.MaxRetries,
				PoolSize:           cfg.PoolSize,
				ReadTimeout:        time.Duration(cfg.ReadTimeout) * time.Second,
				WriteTimeout:       time.Duration(cfg.WriteTimeout) * time.Second,
				IdleTimeout:        time.Duration(cfg.IdleTimeout) * time.Second,
				IdleCheckFrequency: time.Duration(cfg.IdleCheckFrequency) * time.Second,
				MaxConnAge:         time.Duration(cfg.MaxConnAge) * time.Second,
				PoolTimeout:        time.Duration(cfg.PoolTimeout) * time.Second,
			}
			redisClient := redis.NewClient(opt)
			if e := redisClient.Ping().Err(); e != nil {
				freedom.Logger().Fatal(e.Error())
			}
			client = redisClient
			return
		})
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
	`
}
