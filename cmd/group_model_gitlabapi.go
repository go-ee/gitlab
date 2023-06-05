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
	"github.com/spf13/cobra"
)

var gitlabAccessToken string

// groupModelGitlabApiCmd represents the groupModelGitlabApi command
var groupModelGitlabApiCmd = &cobra.Command{
	Use:              "by-api",
	Short:            "Use Gitlab API for Gitlab group reading",
	Long:             `This is the recommended way when Gitlab API is accessible.`,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		_, err = modelByApi()
		return
	},
}

func modelByApi() (ret []string, err error) {
	if gitlabLite, err = newGitlabLiteByApi(); err != nil {
		return
	}

	var jsonModelWriter *lite.JsonWriterModelHandler
	if jsonModelWriter, err = newJsonModelWriter(); err != nil {
		return
	}

	err = readGroupsModels(jsonModelWriter)
	ret = jsonModelWriter.Files
	return
}

func newGitlabLiteByApi() (*lite.GitlabLiteByAPI, error) {
	return lite.NewGitlabLiteByAPI(&lite.ServerAccess{Url: gitlabUrl, Token: gitlabAccessToken})
}

func init() {
	groupModelCmd.AddCommand(groupModelGitlabApiCmd)

	_ = groupModelGitlabApiCmd.MarkPersistentFlagRequired(
		FlagGitlabUrl(groupModelGitlabApiCmd.PersistentFlags(), &gitlabUrl))

	_ = groupModelGitlabApiCmd.MarkPersistentFlagRequired(
		FlagGitlabAccessToken(groupModelGitlabApiCmd.PersistentFlags(), &gitlabAccessToken))
}
