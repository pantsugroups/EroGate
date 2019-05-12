package main

import "github.com/labstack/echo"

var TemplateConfig string = `base:
  secret: this is a secret
  port: 80
route:
  login: /login
  backend: http://127.0.0.1/`

type Verify struct {
	Secret string   `json:"secret"`
	U      UserInfo `json:"u"`
}

type UserInfo struct {
	ID       int    `json:"ID"`
	Username string `json:"username"`
}

type Route struct {
	Route   string `yaml:"route"`
	BackEnd string `yaml:"backend"`
	Auth    bool   `yaml:"auth"`
}
type BaseConf struct {
	Base struct {
		Secret string `yaml:"secret"`
		Port   string `yaml:"port"`
	}
	Route struct {
		Login   string `yaml:"login"`
		Backend string `yaml:"backend"`
	}
}

//type Request struct {
//	Code     int      `json:"code"`
//	Secret   string   `json:"secret"`
//	UserInfo UserInfo `json:"userinfo"`
//}
type ForwardRequest struct {
	UserInfo      UserInfo `json:"userinfo"`
	Method        string   `json:"method"`
	RequestHeader string   `json:"requestheader"`
	RequestBody   string   `json:"requestbody"`
	RequestForm   string   `json:"requestform"`
}

var pathMap = make(map[string]string)
var authMap = make(map[string]bool)
var conf *BaseConf

var e *echo.Echo
