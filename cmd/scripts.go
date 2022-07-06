package cmd

import (
	"encoding/json"
	"github.com/go-ee/gitlab/core"
	"github.com/go-ee/utils/cliu"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"path/filepath"
)

type Scripts struct {
	*cli.Command
	jsonFile, scriptsFolder, reposFolder, devBranch *cliu.StringFlag
}

func NewScripts() (o *Scripts) {
	o = &Scripts{
		jsonFile:      NewJsonFileFlag(),
		scriptsFolder: NewsScriptsFolderFlag(),
		reposFolder:   NewsReposFolderFlag(),
		devBranch:     NewDevBranchFlag(),
	}

	o.Command = &cli.Command{
		Name:  "scripts",
		Usage: "Generate scripts for clone, pull.. and structure creation for all projects of a group recursively",
		Flags: []cli.Flag{
			o.jsonFile, o.scriptsFolder, o.reposFolder,
		},
		Action: func(c *cli.Context) (err error) {
			if o.scriptsFolder.CurrentValue, err = filepath.Abs(o.scriptsFolder.CurrentValue); err != nil {
				logrus.Errorf("error %v by %v to %v", err, c.Command.Name, o.scriptsFolder)
				return
			}

			logrus.Debugf("execute %v to %v", c.Command.Name, o.scriptsFolder)

			var groupNode *core.GroupNode
			if groupNode, err = o.loadJsonFile(); err != nil {
				return
			}

			err = core.Generate(
				groupNode, o.scriptsFolder.CurrentValue, o.reposFolder.CurrentValue, o.devBranch.CurrentValue)

			return
		},
	}
	return
}

func (o *Scripts) loadJsonFile() (ret *core.GroupNode, err error) {
	data, _ := ioutil.ReadFile(o.jsonFile.CurrentValue)

	groupNode := core.GroupNode{}
	err = json.Unmarshal(data, &groupNode)
	ret = &groupNode

	return
}
