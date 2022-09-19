package lite

import (
	"encoding/json"
	"fmt"
	"github.com/go-ee/utils/lg"
	"github.com/xanzy/go-gitlab"
	"os"
	"path/filepath"
	"regexp"
)

type FileLoader interface {
	LoadFileOrFolder(fileOfFolder string) (err error)
}

type FileLoaderJson struct {
	onGroup     func(group *gitlab.Group)
	fileMatcher *regexp.Regexp
}

func NewFileLoaderJson(filePattern string, onGroup func(group *gitlab.Group)) (ret *FileLoaderJson, err error) {
	var fileMatcher *regexp.Regexp
	if fileMatcher, err = regexp.Compile(filePattern); err == nil {
		ret = &FileLoaderJson{
			onGroup:     onGroup,
			fileMatcher: fileMatcher,
		}
	}
	return
}

func (o *FileLoaderJson) loadFile(jsonFile string, onGroup func(group *gitlab.Group)) (err error) {
	file, _ := os.ReadFile(jsonFile)
	group := gitlab.Group{}

	if err = json.Unmarshal(file, &group); err == nil {
		onGroup(&group)
	}
	return
}

func (o *FileLoaderJson) LoadFileOrFolder(fileOfFolder string) (err error) {
	var fileInfo os.FileInfo
	if fileInfo, err = os.Stat(fileOfFolder); err != nil {
		return
	}

	if fileInfo.IsDir() {
		err = filepath.Walk(fileOfFolder, func(path string, child os.FileInfo, err error) (retErr error) {
			if err != nil {
				retErr = err
				return
			}

			// ignore nested folder
			if child.IsDir() {
				if fileInfo.Name() != child.Name() {
					retErr = filepath.SkipDir
				}
				return
			}

			if o.matchFileName(filepath.Base(path)) {
				retErr = o.loadFile(path, o.onGroup)
			}
			return
		})
	} else if o.matchFileName(fileOfFolder) {
		err = o.loadFile(fileOfFolder, o.onGroup)
	}
	return
}

func (o *FileLoaderJson) matchFileName(file string) (ret bool) {
	if ret = o.fileMatcher.MatchString(file); !ret {
		lg.LOG.Infof("ignore file '%v'", file)
	}
	return ret
}

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
