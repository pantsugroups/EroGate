package main

import (
	"github.com/urfave/cli"
	"log"
	"os"
	"strconv"
)

func main() {
	var route string
	var backend string
	var name string
	var tokens string
	var uid string
	var auth string
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
					err := WriteBaseConfigure()
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
			Usage:   "--name name1 --route /website --backend http://localhost:8080/ --auth false \n\t add the route to routes.",
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
				cli.StringFlag{
					Name:        "auth",
					Value:       "false",
					Usage:       "need login?",
					Destination: &auth,
				},
			},
			Action: func(c *cli.Context) error {
				CheckConfFolder()
				err := WriteRoute(route, backend, name, auth)
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
				err := os.Remove("conf.d/" + name + ".yaml")
				if err != nil {
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
				e := LoadConfigure()
				if e != nil {
					log.Fatal("Loading BaseConfigure Failed.Please setup.")
				}
				go StartFolderHandle()
				StartEchoHandle()
				return nil
			},
		},
		{
			Name:    "test",
			Aliases: []string{"t"},
			Usage:   "inline test interface",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:        "token",
					Value:       "",
					Usage:       "token",
					Destination: &tokens,
				},
				cli.StringFlag{
					Name:        "username",
					Value:       "test",
					Usage:       "ai xie sha jiu xie sha.",
					Destination: &name,
				},
				cli.StringFlag{
					Name:        "uid",
					Value:       "1",
					Usage:       "ai xie sha jiu xie sha.",
					Destination: &uid,
				},
			},
			Action: func(c *cli.Context) error {
				e := LoadConfigure()
				if e != nil {
					log.Fatal("Loading BaseConfigure Failed.Please setup.")
				}
				if tokens == "" {
					uuid, _ := strconv.Atoi(uid)
					token, err := CreateToken(&UserInfo{ID: uuid, Username: name})
					if err != nil {
						return err
					}
					log.Println(token)
				} else {
					u, err := ParseToken(tokens)
					if err != nil {
						return err
					}
					log.Println("UId:", u.ID, " name:", u.Username)
				}
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
