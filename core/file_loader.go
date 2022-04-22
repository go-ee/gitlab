package core

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
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
	file, _ := ioutil.ReadFile(jsonFile)
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
		logrus.Infof("ignore file '%v'", file)
	}
	return ret
}