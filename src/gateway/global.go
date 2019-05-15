package main

import "github.com/labstack/echo"

var TemplateConfig string = `base:
  secret: this is a secret
  port: 80
API:
  login: /login
  register: /register`

type Verify struct {
	Secret string   `json:"secret"`
	U      UserInfo `json:"u"`
}

type UserInfo struct {
	ID       int    `json:"ID"`
	Username string `json:"username"`
}

type Route struct {
	Route   string `yaml:"route" json:"route"`
	BackEnd string `yaml:"backend" json:"backend"`
	Auth    bool   `yaml:"auth" json:"auth"`
}
type BaseConf struct {
	Base struct {
		Secret string `yaml:"secret"`
		Port   string `yaml:"port"`
	}
	API struct {
		Login    string `yaml:"login"`
		Register string `yaml:"register"`
	}
}

type HttpResponse struct {
	Successful bool   `json:"successful"`
	Message    string `json:"message"`
}
type Register struct {
	Secret string `json:"secret"`
	Route  Route  `json:"route"`
}
type ForwardRequest struct {
	UserInfo      UserInfo `json:"UserInfo"`
	Method        string   `json:"method"`
	RequestHeader string   `json:"RequestHeader"`
	RequestBody   string   `json:"RequestBody"`
	RequestForm   string   `json:"RequestForm"`
}

var pathMap = make(map[string]string)
var authMap = make(map[string]bool)
var conf *BaseConf

var e *echo.Echo
