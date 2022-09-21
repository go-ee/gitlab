package lite

import (
	"encoding/json"
	"fmt"
	"github.com/go-ee/utils/lg"
	"github.com/xanzy/go-gitlab"
	"os"
	"path/filepath"
)

type ModelHandler interface {
	OnGroup(group *gitlab.Group) error
	OnGroupNode(groupNode *GroupNode) error
}

type JsonWriterModelHandler struct {
	OutputDir           string
	GroupsModelFileName string
	WriteGroup          bool
	WriteGroupNode      bool
	WriteSubGroup       bool
}

func (o *JsonWriterModelHandler) OnGroup(group *gitlab.Group) (err error) {
	if o.WriteGroup {
		targetFile := o.buildGroupFilePath(group.ID)
		lg.LOG.Infof("write Gitlab group '%v' to '%v'", group.Name, targetFile)
		err = o.writeJsonFile(group, targetFile)
	}
	return
}

func (o *JsonWriterModelHandler) OnGroupNode(groupNode *GroupNode) (err error) {
	if o.WriteGroupNode {
		targetFile := filepath.Join(o.OutputDir, groupNode.RelativeRootPath, o.GroupsModelFileName)
		lg.LOG.Infof("write Gitlab model '%v(%v)' to '%v'", groupNode.Group.Name, groupNode.Group.ID, targetFile)
		if err = o.writeJsonFile(groupNode, targetFile); err != nil {
			return
		}
	}

	if o.WriteGroup {
		if err = o.OnGroup(groupNode.Group); err != nil {
			return
		}
	}

	if o.WriteSubGroup {
		for _, subGroupNode := range groupNode.Children {
			if err = o.OnGroupNode(subGroupNode); err != nil {
				break
			}
		}
	}
	return
}

func (o *JsonWriterModelHandler) writeJsonFile(content interface{}, targetFile string) (err error) {
	if err = CreateDirIfNotExists(targetFile); err != nil {
		return
	}

	var data []byte
	if data, err = json.MarshalIndent(content, "", " "); err != nil {
		return
	}
	err = os.WriteFile(targetFile, data, 0700)
	return
}

func (o *JsonWriterModelHandler) buildGroupFilePath(groupNameOrId interface{}) string {
	return filepath.Join(o.OutputDir, fmt.Sprintf("%v", groupNameOrId)+".json")
}

func (o *JsonWriterModelHandler) buildGroupModelFilePath(groupName string) string {
	return filepath.Join(o.OutputDir, groupName+".json")
}

func CreateDirIfNotExists(targetFile string) (err error) {
	dir := filepath.Dir(targetFile)
	err = os.MkdirAll(dir, 0700)
	return
}
