package main

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/fsnotify/fsnotify"
	"github.com/go-yaml/yaml"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

var TemplateConfig string = `base:
  secret: this is a sercret
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
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}
func PathExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}
func CheckConfFolder(){
	if !PathExist("conf.d/") && !IsDir("conf.d/"){
		err:=os.Mkdir("conf.d/",os.ModePerm)
		if err!=nil{
			log.Println(err)
		}
	}
}
func secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.Base.Secret), nil
	}
}

func CreateToken(user *UserInfo) (tokens string, err error) {

	claim := jwt.MapClaims{
		"id":       user.ID,
		"username": user.Username,
		"nbf":      time.Now().Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	tokens, err = token.SignedString([]byte(conf.Base.Secret))
	if err != nil {
		return "", nil
	}
	return tokens, nil
}

func ParseToken(tokens string) (user *UserInfo, err error) {
	user = &UserInfo{}
	token, err := jwt.Parse(tokens, secret())
	if err != nil {
		return user, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		err = errors.New("cannot convert claim to mapclaim")
		return user, err
	}
	if !token.Valid {
		err = errors.New("token is invalid")
		return user, err
	}

	user.ID = claim["id"].(int)
	user.Username = claim["username"].(string)
	return user, nil
}

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
	e.Any(route.Route,ManualGateWay)
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
func StartEchoHandle() {
	// Echo instance

	e = echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	err := LoadRoutes()
	if err != nil {
		log.Println("Load routes error", err)
	}

	// Start server
	e.Logger.Fatal(e.Start(":" + conf.Base.Port))
}

func main() {
	var route string
	var backend string
	var name string
	app := cli.NewApp()
	app.Name = "gateway"
	app.Usage = "Ero Gateway. https://ero.ink"
	app.Commands = []cli.Command{
		{
			Name:    "setup",
			Aliases: []string{"s"},
			Usage:   "setup the configure to ./conf.yaml",
			Action: func(c *cli.Context) error {
				//log.Println("added task: ", c.Args().First())
				log.Println("installing...")
				if PathExist("yaml.conf") {
					log.Println("the configure is exits.")
				} else {
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

				}
				log.Println("installed.")
				return nil
			},
		},
		{
			Name:    "add",
			Aliases: []string{"a"},
			Usage:   "--name name1 --route /website --backend http://localhost:8080/ \n\t add the route to routes.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Value:       "name1",
					Usage:       "route's name",
					Destination: &name,
				},
				cli.StringFlag{
					Name:        "route",
					Value:       "/website",
					Usage:       "gateway listening address",
					Destination: &route,
				},
				cli.StringFlag{
					Name:        "backend",
					Value:       "http://localhost:8080/",
					Usage:       "backend address",
					Destination: &backend,
				},
			},
			Action: func(c *cli.Context) error {
				CheckConfFolder()
				s, err := yaml.Marshal(&Route{Route: route, BackEnd: backend})
				if err != nil {
					return err
				}
				err = ioutil.WriteFile("conf.d/"+name+".yaml", []byte(s), 0644)
				if err != nil {

					return err
				}
				return nil
			},
		},
		{
			Name:    "del",
			Aliases: []string{"d"},
			Usage:   "--name name1    del the route to routes.",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "name",
					Value:       "name1",
					Usage:       "route's name",
					Destination: &name,
				},
			},
			Action: func(c *cli.Context) error {
				err:= os.Remove("conf.d/"+name+".yaml")
				if err != nil{
					return err
				}
				return nil
			},
		},
		{
			Name:    "run",
			Aliases: []string{"r"},
			Usage:   "running the server.",
			Action: func(c *cli.Context) error {
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
