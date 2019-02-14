package main

import (
	"github.com/Sirupsen/logrus"
	"github.com/go-ee/gitlab"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
)

const flagDebug = "debug"
const flagToken = "token"
const flagURL = "url"
const flagGroup = "group"
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
				},
			},
			Action: func(c *cli.Context) (err error) {
				var target string
				if target, err = filepath.Abs(c.String(flagTarget)); err != nil {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, target)
					return
				}

				logrus.Infof("execute %v to %v", c.Command.Name, target)

				if err = gitlab.Generate(&gitlab.Params{
					Url:       c.GlobalString(flagURL),
					GroupName: c.String(flagGroup),
					Target:    target,
					Token:     c.GlobalString(flagToken)}); err != nil {

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
