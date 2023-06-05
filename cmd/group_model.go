/*
Copyright © 2022 Eugen Eisler <eoeisler@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package cmd

import (
	"github.com/go-ee/gitlab/lite"
	"github.com/go-ee/utils/lg"
	"github.com/spf13/cobra"
	"path/filepath"
)

var gitlabUrl string
var gitlabUrlApiPart = "api/v4"
var groups []string
var ignoreGroups []string
var outputDir = "."
var offlineModeSupport = false
var offlineGroupsDir = ".gitlab"
var offlineGroupsFilesPattern = ".+?\\.json$"
var groupsModelFileName = ".gitlab.json"
var storeGroupModel = true

var gitlabLite lite.GitlabLite

// groupModelCmd represents the groupsModel command
var groupModelCmd = &cobra.Command{
	Use:              "group-model",
	Short:            "Reads groups models",
	TraverseChildren: true,
}

func readGroupsModels(modelHandler lite.ModelHandler) (err error) {
	modelReader := &lite.ModelReader{
		Client:           gitlabLite,
		IgnoreGroupNames: SliceToMap(ignoreGroups),
	}

	for _, groupNameOrId := range groups {
		if groupNode, groupErr := modelReader.ReadGroupModelByGroup(groupNameOrId); groupErr == nil {
			if err = modelHandler.OnGroupNode(groupNode); err != nil {
				lg.LOG.Warnf("error at writing of JSON for groupNameOrId %v: %v", groupNameOrId, err)
				return
			}
		} else {
			lg.LOG.Warnf("error at downloading of JSON for groupNameOrId %v", groupNameOrId)
		}
	}
	return
}

func newJsonModelWriter() (ret *lite.JsonWriterModelHandler, err error) {
	var absOutputDir string
	if absOutputDir, err = filepath.Abs(outputDir); err != nil {
		return
	}
	ret = &lite.JsonWriterModelHandler{
		OutputDir:           absOutputDir,
		OfflineGroupsDir:    offlineGroupsDir,
		GroupsModelFileName: groupsModelFileName,
		WriteGroup:          offlineModeSupport,
		WriteGroupNode:      storeGroupModel,
		WriteSubGroup:       storeGroupModel,
	}
	return
}

func init() {
	rootCmd.AddCommand(groupModelCmd)

	_ = groupModelCmd.MarkPersistentFlagRequired(
		FlagGroups(groupModelCmd.PersistentFlags(), &groups))

	FlagGroupModelFileName(groupModelCmd.PersistentFlags(), &groupsModelFileName)
	FlagOutputDir(groupModelCmd.PersistentFlags(), &outputDir)
	FlagIgnoreGroups(groupModelCmd.PersistentFlags(), &ignoreGroups)

	FlagStoreGroupModel(groupModelCmd.PersistentFlags(), &storeGroupModel)

	FlagOfflineSupport(groupModelCmd.PersistentFlags(), &offlineModeSupport)
	FlagOfflineGroupsDir(groupModelCmd.PersistentFlags(), &offlineGroupsDir)
}
