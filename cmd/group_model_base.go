package cmd

import (
	"encoding/json"
	"github.com/go-ee/gitlab/lite"
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type GroupModelBase struct {
	*cli.Command
	group, ignores, groupModelFile *cliu.StringFlag
}

func NewGroupModelBase() (o *GroupModelBase) {
	o = &GroupModelBase{
		group:          NewGroupFlag(),
		ignores:        NewIgnoresFlag(),
		groupModelFile: NewGroupModelFileFlag(),
	}
	return
}

func (o *GroupModelBase) prepareJsonFile(c *cli.Context) (err error) {
	if o.groupModelFile.CurrentValue, err = filepath.Abs(o.groupModelFile.CurrentValue); err != nil {
		lg.LOG.Errorf("error %v by %v to %v", err, c.Command.Name, o.groupModelFile)
	}
	return
}

func (o *GroupModelBase) extract(client lite.GitlabLite) (ret *lite.GroupNode, err error) {
	ret, err = lite.FetchGroupModel(&lite.GroupModelParams{
		Group:            o.group.CurrentValue,
		IgnoreGroupNames: buildIgnoresMap(o.ignores.CurrentValue),
	}, client)
	return
}

func (o *GroupModelBase) writeJsonFile(groupNode *lite.GroupNode) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(groupNode, "", " "); err != nil {
		return
	}
	targetFile := o.groupModelFile.CurrentValue
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
