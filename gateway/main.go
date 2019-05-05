package gateway

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-yaml/yaml"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"io/ioutil"
	"log"
	"net/http"
)

type Route struct {
	Route   string `yaml:"route"`
	BackEnd string `yaml:"backend"`
}
type BaseConf struct {
	Base struct {
		Login string `yaml:"login"`
		Port  string `yaml:"port"`
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

func LoadRoute(path string) error {
	return nil
}
func LoadRoutes() error {
	dir, err := ioutil.ReadDir("conf.d/")
	if err != nil {
		return nil
	}
	for _, i := range dir {
		if !i.IsDir() {
			err := LoadRoute(i.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func LoadConfigure() error {
	yamlFile, err := ioutil.ReadFile("src/conf.yaml")
	if err != nil {
		return err
	}
	conf = new(BaseConf)
	err = yaml.Unmarshal(yamlFile, conf)
	if err != nil {
		return err
	}
	return nil

}

func Verify(session string) {

}

func ManualGateWay(c echo.Context) error {

	if _, ok := c.Request().Header["x-headers-session"]; ok {
		// 驗證session
	} else {
		// 返回403
		return c.String(http.StatusForbidden, "Please.Login!")
	}
	return nil //明明理論上執行不到這裏，卻還是給我報提示mmm
}

func main() {
	// Echo instance
	err := LoadConfigure()
	if err != nil {
		log.Println(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
	}
	defer watcher.Close()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
				}
				// 在這裏加載新路由
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	err = watcher.Add("conf.d/") //也可以监听文件夹
	if err != nil {
		log.Println("listen folder error:", err)
	}
	<-done

	e = echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	err = LoadRoutes()
	if err != nil {
		log.Println("Load routes error", err)
	}

	// Start server
	e.Logger.Fatal(e.Start(":" + conf.Base.Port))
}
