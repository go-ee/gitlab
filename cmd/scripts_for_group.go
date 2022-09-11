package cmd

import (
	"encoding/json"
	"github.com/go-ee/gitlab/lite"
	"github.com/go-ee/gitlab/script"
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
	"github.com/urfave/cli/v2"
	"os"
	"path/filepath"
)

type ScriptsForGroup struct {
	*cli.Command
	groupFile, scriptsFolder, reposFolder *cliu.StringFlag
}

func NewScriptsForGroup() (o *ScriptsForGroup) {
	o = &ScriptsForGroup{
		groupFile:     NewGroupModelFileFlag(),
		scriptsFolder: NewsScriptsFolderFlag(),
		reposFolder:   NewsReposFolderFlag(),
	}

	o.Command = &cli.Command{
		Name:  "scripts",
		Usage: "Generate scripts for clone, pull.. and structure creation for all projects of a group recursively",
		Flags: []cli.Flag{
			o.groupFile, o.scriptsFolder, o.reposFolder,
		},
		Action: func(c *cli.Context) (err error) {
			if o.scriptsFolder.CurrentValue, err = filepath.Abs(o.scriptsFolder.CurrentValue); err != nil {
				lg.LOG.Errorf("error %v by %v to %v", err, c.Command.Name, o.scriptsFolder)
				return
			}

			lg.LOG.Debugf("execute %v to %v", c.Command.Name, o.scriptsFolder)

			var groupNode *lite.GroupNode
			if groupNode, err = o.loadJsonFile(); err != nil {
				return
			}

			err = script.Generate(groupNode, o.scriptsFolder.CurrentValue, o.reposFolder.CurrentValue)

			return
		},
	}
	return
}

func (o *ScriptsForGroup) loadJsonFile() (ret *lite.GroupNode, err error) {
	data, _ := os.ReadFile(o.groupFile.CurrentValue)

	groupNode := lite.GroupNode{}
	err = json.Unmarshal(data, &groupNode)
	ret = &groupNode

	return
}
