package main

import (
	"github.com/go-ee/gitlab"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var token, url, group, target, ignores, devBranch string
	app := cli.NewApp()
	app.Usage = "Gitlab helper"
	app.Version = "1.0"

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug log level",
		},
	}

	commonFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Required:    true,
			Usage:       "Gitlab token",
			Value:       ".",
			Destination: &token,
		}, &cli.StringFlag{
			Name:        "url",
			Required:    true,
			Usage:       "Base Gitlab server url",
			Destination: &url,
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "generateScripts",
			Usage: "GenerateScripts for clone, pull all projects of a group",
			Flags: append(commonFlags,
				&cli.StringFlag{
					Name:        "group",
					Usage:       "Gitlab group",
					Required:    true,
					Destination: &group,
				}, &cli.StringFlag{
					Name:        "target",
					Usage:       "Target dir",
					Destination: &target,
				}, &cli.StringFlag{
					Name:        "ignores",
					Usage:       "Ignores the comma separated groups",
					Destination: &ignores,
				}, &cli.StringFlag{
					Name:        "devBranch",
					Usage:       "The default development branch, that will be used in 'devBranch' script",
					Value:       "development",
					Destination: &devBranch,
				},
			),
			Action: func(c *cli.Context) (err error) {
				if target, err = filepath.Abs(target); err != nil {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, target)
					return
				}

				logrus.Infof("execute %v to %v", c.Command.Name, target)

				ignoresMap := make(map[string]bool)
				if ignores != "" {
					for _, name := range strings.Split(c.String(ignores), ",") {
						ignoresMap[name] = true
					}
				}

				if err = gitlab.Generate(&gitlab.Params{
					Url:       url,
					GroupName: group,
					Target:    target,
					Token:     token,
					DevBranch: devBranch,
					Ignores:   ignoresMap}); err != nil {

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
