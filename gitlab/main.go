package main

import (
	"encoding/json"
	"github.com/go-ee/gitlab"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	var token, url, group, source, target, ignores, devBranch string
	app := cli.NewApp()
	app.Usage = "Gitlab helper"
	app.Version = "1.0"

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "Enable debug log level",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:  "extract",
			Usage: "Extract group recursively to a JSON file",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "token",
					Required:    true,
					Usage:       "Gitlab token",
					Destination: &token,
				}, &cli.StringFlag{
					Name:        "url",
					Required:    true,
					Usage:       "Base Gitlab server url",
					Destination: &url,
				},
				&cli.StringFlag{
					Name:        "group",
					Usage:       "Gitlab group",
					Required:    true,
					Destination: &group,
				}, &cli.StringFlag{
					Name:        "target",
					Usage:       "Target JSON file name",
					Value:       "gitlab.json",
					Destination: &target,
				}, &cli.StringFlag{
					Name:        "ignores",
					Usage:       "Ignore group names the comma separated groups",
					Destination: &ignores,
				},
			},
			Action: func(c *cli.Context) (err error) {
				if target, err = filepath.Abs(target); err != nil {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, target)
					return
				}

				logrus.Infof("execute %v to %v", c.Command.Name, target)

				ignoreGroupNames := make(map[string]bool)
				if ignores != "" {
					for _, name := range strings.Split(c.String(ignores), ",") {
						ignoreGroupNames[name] = true
					}
				}

				var groupNode *gitlab.GroupNode
				if groupNode, err = gitlab.Extract(&gitlab.ExtractParams{
					Url:              url,
					Token:            token,
					GroupName:        group,
					IgnoreGroupNames: ignoreGroupNames,
				}); err == nil {
					data, _ := json.MarshalIndent(groupNode, "", " ")
					_ = ioutil.WriteFile(target, data, 0644)
				} else {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, target)
					return
				}
				return
			},
		},
		{
			Name:  "scripts",
			Usage: "Generate scripts for clone, pull.. and structure creation for all projects of a group recursively",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:        "source",
					Usage:       "Source JSON file name",
					Value:       "gitlab.json",
					Destination: &source,
				}, &cli.StringFlag{
					Name:        "target",
					Usage:       "Target dir",
					Value:       ".",
					Destination: &target,
				},
			},
			Action: func(c *cli.Context) (err error) {
				if target, err = filepath.Abs(target); err != nil {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, target)
					return
				}

				logrus.Infof("execute %v to %v", c.Command.Name, target)

				data, _ := ioutil.ReadFile(source)

				groupNode := gitlab.GroupNode{}

				_ = json.Unmarshal(data, &groupNode)

				err = gitlab.Generate(&groupNode, target, devBranch)

				return
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Warn("exit because of error.")
	}
}
