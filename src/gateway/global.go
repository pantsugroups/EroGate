package main

import "github.com/labstack/echo"

var TemplateConfig string = `base:
  secret: this is a secret
  login: /login
  port: 80`

type UserInfo struct {
	ID       int
	Username string
}

type Route struct {
	Route   string `yaml:"route"`
	BackEnd string `yaml:"backend"`
}
type BaseConf struct {
	Base struct {
		Secret string `yaml:"secret"`
		Login  string `yaml:"login"`
		Port   string `yaml:"port"`
	}
}

var conf *BaseConf

var e *echo.Echo
