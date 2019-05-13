package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
	"os"
)

func CheckConfFolder() {
	if !PathExist("conf.d/") && !IsDir("conf.d/") {
		err := os.Mkdir("conf.d/", os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}
}
func LoadRoute(path string) error {

	yamlFile, err := ioutil.ReadFile("conf.d/" + path)
	if err != nil {
		return err
	}
	route := new(Route)
	err = yaml.Unmarshal(yamlFile, route)
	if err != nil {
		return err
	}
	log.Println("Route add:", route.Route)
	log.Println("Route Backend:", route.BackEnd)
	log.Println("Route Auth:", route.Auth)
	e.Any(route.Route, ManualGateWay)
	pathMap[route.Route] = route.BackEnd
	authMap[route.Route] = route.Auth
	return nil
}
func LoadRoutes() error {
	e.POST(conf.Route.Login, ManualLogin)
	//pathMap[conf.Route.Login] = conf.Route.Backend
	e.Any(Secret2Route(conf.Base.Secret), ShowAllRoutes)
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
	CheckConfFolder()
	yamlFile, err := ioutil.ReadFile("conf.yaml")
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
func WriteRoute(route string, backend string, name string, auth string) error {
	var b bool
	if auth == "false" {
		b = false
	} else {
		b = true
	}
	s, err := yaml.Marshal(&Route{Route: route, BackEnd: backend, Auth: b})
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("conf.d/"+name+".yaml", []byte(s), 0644)
	if err != nil {

		return err
	}
	return nil
}
func WriteBaseConfigure() error {
	//TemplateConfig,_ := yaml.Marshal(&BaseConf{Base: struct {
	//	Secret string
	//	Login  string
	//	Port   string
	//}{Secret:"lalala" , Login:"/login" , Port:"8080" } ,
	//Mysql: struct {
	//	User     string
	//	Host     string
	//	Password string
	//	Port     string
	//	Database string
	//}{User:"root" , Host:"localhost" , Password: "123456", Port:"3306" , Database:"db" },
	//})
	//爲了避免這種蛋疼的寫法我選擇硬編碼
	err := ioutil.WriteFile("conf.yaml", []byte(TemplateConfig), 0644)
	if err != nil {

		return err
	}
	return nil
}
func StartFolderHandle() {
	err := LoadConfigure()
	if err != nil {
		log.Println(err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Println(err)
	}
	defer func() {
		if err = watcher.Close(); err != nil {
			log.Println(err)
		}
	}()
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					log.Println("modified file:", event.Name)
					err := LoadRoute(event.Name)
					if err != nil {
						log.Println("error:", err)
					}
				}

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
}
