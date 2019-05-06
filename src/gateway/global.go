package gateway

import "github.com/labstack/echo"

var TemplateConfig string = `base:
  secret: this is a secret
  login: /login
  port: 80
mysql:
  user : root
  password : sa
  host : 127.0.0.1
  port : 3306
  db : database`

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
	Mysql struct {
		User     string `yaml:"user"`
		Host     string `yaml:"host"`
		Password string `yaml:"password"`
		Port     string `yaml:"port"`
		Database string `yaml:"db"`
	}
}

var conf *BaseConf

var e *echo.Echo
