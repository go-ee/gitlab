package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-ee/gitlab"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"strings"
)

const flagDebug = "debug"
const flagToken = "token"
const flagURL = "url"
const flagGroup = "group"
const flagIgnores = "ignores"
const flagTarget = "target"

func main() {
	app := cli.NewApp()
	app.Usage = "Gitlab helper"
	app.Version = "1.0"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  flagToken,
			Usage: "Gitlab token",
		}, cli.BoolFlag{
			Name:  flagDebug,
			Usage: "Enable debug log level",
		}, cli.StringFlag{
			Name:  flagURL,
			Usage: "Base url",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "generateScripts",
			Usage: "GenerateScripts for clone, pull all projects of a group",
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  flagGroup,
					Usage: "Gitlab group",
				}, cli.StringFlag{
					Name:  flagTarget,
					Usage: "Target dir",
				}, cli.StringFlag{
					Name:  flagIgnores,
					Usage: "Ignores the comma separated groups",
				},
			},
			Action: func(c *cli.Context) (err error) {
				var target string
				if target, err = filepath.Abs(c.String(flagTarget)); err != nil {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, target)
					return
				}

				logrus.Infof("execute %v to %v", c.Command.Name, target)

				ignores := make(map[string]bool)
				if c.GlobalIsSet(flagIgnores) {
					for _, name := range strings.Split(c.GlobalString(flagIgnores), ",") {
						ignores[name] = true
					}
				}

				if err = gitlab.Generate(&gitlab.Params{
					Url:       c.GlobalString(flagURL),
					GroupName: c.String(flagGroup),
					Target:    target,
					Token:     c.GlobalString(flagToken),
					Ignores:   ignores}); err != nil {

					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, target)
					return
				}

				return
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Warn("exit because of error.")
	}
}
