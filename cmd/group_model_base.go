package cmd

import (
	"encoding/json"
	"github.com/go-ee/gitlab/core"
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type GroupModelBase struct {
	*cli.Command
	group, ignores, jsonFile *cliu.StringFlag
}

func NewGroupModelBase() (o *GroupModelBase) {
	o = &GroupModelBase{
		group:    NewGroupFlag(),
		ignores:  NewIgnoresFlag(),
		jsonFile: NewJsonFileFlag(),
	}
	return
}

func (o *GroupModelBase) prepareJsonFile(c *cli.Context) (err error) {
	if o.jsonFile.CurrentValue, err = filepath.Abs(o.jsonFile.CurrentValue); err != nil {
		lg.LOG.Errorf("error %v by %v to %v", err, c.Command.Name, o.jsonFile)
	}
	return
}

func (o *GroupModelBase) extract(client core.GitlabLite) (ret *core.GroupNode, err error) {
	ret, err = core.Extract(&core.ExtractParams{
		Group:            o.group.CurrentValue,
		IgnoreGroupNames: buildIgnoresMap(o.ignores.CurrentValue),
	}, client)
	return
}

func (o *GroupModelBase) writeJsonFile(groupNode *core.GroupNode) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(groupNode, "", " "); err != nil {
		return
	}
	targetFile := o.jsonFile.CurrentValue
	lg.LOG.Infof("write gitlab model '%v(%v)' to '%v'", groupNode.Group.Name, groupNode.Group.ID, targetFile)
	err = ioutil.WriteFile(targetFile, data, 0644)
	return err
}

func buildIgnoresMap(ignores string) (ret map[string]bool) {
	ret = make(map[string]bool)
	if ignores != "" {
		for _, name := range strings.Split(ignores, ",") {
			ret[name] = true
		}
	}
	return
}
