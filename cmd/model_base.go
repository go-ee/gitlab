package cmd

import (
	"encoding/json"
	"github.com/go-ee/gitlab/core"
	"github.com/go-ee/utils/cliu"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"io/ioutil"
	"path/filepath"
	"strings"
)

type ModelBase struct {
	*cli.Command
	group, ignores, jsonFile *cliu.StringFlag
}

func NewModelBase() (o *ModelBase) {
	o = &ModelBase{
		group:    NewGroupFlag(),
		ignores:  NewIgnoresFlag(),
		jsonFile: NewJsonFileFlag(),
	}
	return
}

func (o *ModelBase) prepareJsonFile(c *cli.Context) (err error) {
	if o.jsonFile.CurrentValue, err = filepath.Abs(o.jsonFile.CurrentValue); err != nil {
		logrus.Errorf("error %v by %v to %v", err, c.Command.Name, o.jsonFile)
	}
	return
}

func (o *ModelBase) extract(client core.GitlabLite) (ret *core.GroupNode, err error) {
	ret, err = core.Extract(&core.ExtractParams{
		Group:            o.group.CurrentValue,
		IgnoreGroupNames: buildIgnoresMap(o.ignores.CurrentValue),
	}, client)
	return
}

func (o *ModelBase) writeJsonFile(groupNode *core.GroupNode) (err error) {
	var data []byte
	if data, err = json.MarshalIndent(groupNode, "", " "); err != nil {
		return
	}
	err = ioutil.WriteFile(o.jsonFile.CurrentValue, data, 0644)
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
