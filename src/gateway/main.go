package gateway

import (
	"github.com/urfave/cli"
	"log"
	"os"
)

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
				err := WriteRoute(route, backend, name)
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
				err := LoadConfigure()
				if err != nil {
					return err
				}
				go StartFolderHandle()
				StartEchoHandle()
				return nil
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
