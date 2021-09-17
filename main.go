package main

import (
	"encoding/json"
	"github.com/go-ee/gitlab/core"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	app := NewCli()

	if err := app.Run(os.Args); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Warn("exit because of error.")
	}
}

type Cli struct {
	*cli.App
	debug                                                       bool
	token, url, group, jsonFile, scriptsDir, ignores, devBranch string
	ignoreGroupNames                                            map[string]bool
}

func NewCli() (ret *Cli) {
	ret = &Cli{}
	ret.init()
	return
}

func (o *Cli) init() {
	o.App = cli.NewApp()
	o.Usage = "Gitlab helper"
	o.Version = "1.0"

	o.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:        "debug",
			Usage:       "Enable debug log level",
			Destination: &o.debug,
		},
	}

	o.Before = func(c *cli.Context) (err error) {
		if o.debug {
			logrus.SetLevel(logrus.DebugLevel)
		}
		logrus.Debugf("execute %v, %v", c.Command.Name, c.Args())
		return
	}

	extractFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "token",
			Required:    true,
			Usage:       "Gitlab token",
			Destination: &o.token,
		}, &cli.StringFlag{
			Name:        "url",
			Required:    true,
			Usage:       "Base Gitlab server url",
			Destination: &o.url,
		},
		&cli.StringFlag{
			Name:        "group",
			Usage:       "Gitlab group",
			Required:    true,
			Destination: &o.group,
		}, &cli.StringFlag{
			Name:        "ignores",
			Usage:       "Ignore group names the comma separated groups",
			Destination: &o.ignores,
		}, &cli.StringFlag{
			Name:        "jsonFile",
			Usage:       "Target JSON file name",
			Value:       "gitlab.json",
			Destination: &o.jsonFile,
		},
	}

	o.Commands = []cli.Command{
		{
			Name:  "extract",
			Usage: "Extract group recursively to a JSON file",
			Flags: extractFlags,
			Action: func(c *cli.Context) (err error) {
				if err = o.prepareJsonFile(c); err != nil {
					return
				}
				logrus.Debugf("execute %v to %v", c.Command.Name, o.jsonFile)

				var groupNode *core.GroupNode
				if groupNode, err = o.extract(); err == nil {
					err = o.writeJsonFile(groupNode)
				} else {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, o.jsonFile)
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
					Name:        "jsonFile",
					Usage:       "Source JSON file name",
					Value:       "gitlab.json",
					Destination: &o.jsonFile,
				}, &cli.StringFlag{
					Name:        "scriptsFolder",
					Usage:       "Folder where scripts are generated",
					Value:       ".",
					Destination: &o.scriptsDir,
				},
			},
			Action: func(c *cli.Context) (err error) {
				if o.scriptsDir, err = filepath.Abs(o.scriptsDir); err != nil {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, o.scriptsDir)
					return
				}

				logrus.Debugf("execute %v to %v", c.Command.Name, o.scriptsDir)

				var groupNode *core.GroupNode
				if groupNode, err = o.loadJsonFile(); err != nil {
					return
				}

				err = core.Generate(groupNode, o.scriptsDir, o.devBranch)

				return
			},
		},
		{
			Name:  "extract-scripts",
			Usage: "Extract and generate scripts",
			Flags: append(extractFlags,
				&cli.StringFlag{
					Name:        "scriptsFolder",
					Usage:       "Folder where scripts are generated",
					Value:       ".",
					Destination: &o.scriptsDir,
				},
			),
			Action: func(c *cli.Context) (err error) {
				if err = o.prepareJsonFile(c); err != nil {
					return
				}
				if err = o.prepareScriptsDir(c); err != nil {
					return
				}
				logrus.Debugf("execute %v to %v and %v", c.Command.Name, o.jsonFile, o.scriptsDir)

				var groupNode *core.GroupNode
				if groupNode, err = o.extract(); err == nil {
					if err = o.writeJsonFile(groupNode); err != nil {
						return
					}
				} else {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, o.jsonFile)
					return
				}
				err = core.Generate(groupNode, o.scriptsDir, o.devBranch)

				return
			},
		},
	}
}

func (o *Cli) writeJsonFile(groupNode *core.GroupNode) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(groupNode, "", " "); err != nil {
		return
	}
	err = ioutil.WriteFile(o.jsonFile, data, 0644)
	return err
}

func (o *Cli) loadJsonFile() (ret *core.GroupNode, err error) {
	data, _ := ioutil.ReadFile(o.jsonFile)

	groupNode := core.GroupNode{}
	err = json.Unmarshal(data, &groupNode)
	ret = &groupNode

	return
}

func (o *Cli) prepareJsonFile(c *cli.Context) (err error) {
	if o.jsonFile, err = filepath.Abs(o.jsonFile); err != nil {
		logrus.Errorf("error %v by %v to %v", err, c.Command.Name, o.jsonFile)
	}
	return
}

func (o *Cli) prepareScriptsDir(c *cli.Context) (err error) {
	if o.scriptsDir, err = filepath.Abs(o.scriptsDir); err != nil {
		logrus.Errorf("error %v by %v to %v", err, c.Command.Name, o.scriptsDir)
	}
	return
}

func (o *Cli) extract() (ret *core.GroupNode, err error) {

	o.buildIgnoresMap()

	ret, err = core.Extract(&core.ExtractParams{
		Url:              o.url,
		Token:            o.token,
		GroupName:        o.group,
		IgnoreGroupNames: o.ignoreGroupNames,
	})
	return
}

func (o *Cli) buildIgnoresMap() {
	o.ignoreGroupNames = make(map[string]bool)
	if o.ignores != "" {
		for _, name := range strings.Split(o.ignores, ",") {
			o.ignoreGroupNames[name] = true
		}
	}
	return
}
