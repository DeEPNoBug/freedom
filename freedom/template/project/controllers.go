package project

import "strings"

func init() {
	content["/adapter/controller/default.go"] = controllerTemplate()
}

func controllerTemplate() string {
	result := `
	// Code generated by 'freedom new-project {{.PackagePath}}'
	package controller

	import (
		"github.com/8treenet/freedom"
		"{{.PackagePath}}/domain"
		"{{.PackagePath}}/infra"
	)
	
	func init() {
		freedom.Prepare(func(initiator freedom.Initiator) {
			initiator.BindController("/", &Default{})
		})
	}
	
	type Default struct {
		Sev     *domain.Default
		Worker freedom.Worker
		Request *infra.Request
	}
	
	// Get handles the GET: / route.
	func (c *Default) Get() freedom.Result {
		c.Worker.Logger().Infof("我是控制器")
		remote := c.Sev.RemoteInfo()
		return &infra.JSONResponse{Object: remote}
	}
	
	// GetHello handles the GET: /hello route.
	func (c *Default) GetHello() string {
		field := freedom.LogFields{
			"framework": "freedom",
			"like":      "DDD",
		}
		c.Worker.Logger().Info("hello", field)
		c.Worker.Logger().Infof("hello %s", "format", field)
		c.Worker.Logger().Debug("hello", field)
		c.Worker.Logger().Debugf("hello %s", "format", field)
		c.Worker.Logger().Error("hello", field)
		c.Worker.Logger().Errorf("hello %s", "format", field)
		c.Worker.Logger().Warn("hello", field)
		c.Worker.Logger().Warnf("hello %s", "format", field)
		c.Worker.Logger().Print("hello")
		c.Worker.Logger().Println("hello")
		return "hello"
	}

	// PutHello handles the PUT: /hello route.
	func (c *Default) PutHello() freedom.Result {
		return &infra.JSONResponse{Object: "putHello"}
	}

	// PostHello handles the POST: /hello route.
	func (c *Default) PostHello() freedom.Result {
		var postJsonData struct {
			UserName     string $$wavejson:"userName" validate:"required"$$wave
			UserPassword string $$wavejson:"userPassword" validate:"required"$$wave
		}
		if err := c.Request.ReadJSON(&postJsonData); err != nil {
			return &infra.JSONResponse{Error: err}
		}
		
		return &infra.JSONResponse{Object: "postHello"}
	}

	func (m *Default) BeforeActivation(b freedom.BeforeActivation) {
		b.Handle("ANY", "/custom", "CustomHello")
		//b.Handle("GET", "/custom", "CustomHello")
		//b.Handle("PUT", "/custom", "CustomHello")
		//b.Handle("POST", "/custom", "CustomHello")
	}

	// PostHello handles the POST: /hello route.
	func (c *Default) CustomHello() freedom.Result {
		method := c.Worker.IrisContext().Request().Method
		c.Worker.Logger().Info("CustomHello", freedom.LogFields{"method": method})
		return &infra.JSONResponse{Object: method + "CustomHello"}
	}
	
	// GetUserBy handles the GET: /user/{username:string} route.
	func (c *Default) GetUserBy(username string) freedom.Result {
		var query struct {
			Token string $$waveurl:"token" validate:"required"$$wave
			Id    int64  $$waveurl:"id" validate:"required"$$wave
		}
		if err := c.Request.ReadQuery(&query); err != nil {
			return &infra.JSONResponse{Error: err}
		}
		var data struct {
			Name  string
			Token string
			Id    int64
		}
		data.Id = query.Id
		data.Token = query.Token
		data.Name = username
		return &infra.JSONResponse{Object: data}
	}
	
	// GetAgeByUserBy handles the GET: /age/{age:int}/user/{user:string} route.
	func (c *Default) GetAgeByUserBy(age int, user string) freedom.Result {
		var result struct {
			User string
			Age  int
		}
		result.Age = age
		result.User = user

		return &infra.JSONResponse{Object: result}
	}
	`

	result = strings.ReplaceAll(result, "$$wave", "`")
	return result
}
