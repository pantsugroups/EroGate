package gateway

import (
	"github.com/fsnotify/fsnotify"
	"github.com/go-yaml/yaml"
	"io/ioutil"
	"log"
)

func LoadRoute(path string) error {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	route := new(Route)
	err = yaml.Unmarshal(yamlFile, route)
	if err != nil {
		return err
	}
	e.Any(route.Route, ManualGateWay)
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
func WriteRoute(route string, backend string, name string) error {
	s, err := yaml.Marshal(&Route{Route: route, BackEnd: backend})
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
}
