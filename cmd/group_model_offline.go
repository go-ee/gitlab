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

// groupModelOfflineCmd represents the groupModelOffline command
var groupModelOfflineCmd = &cobra.Command{
	Use:              "offline",
	Short:            "Use groups files downloaded before before.",
	Long:             `This operation can be used without connection to Gitlab when Gitlab groups files were downloaded before`,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return modelByOffline()
	},
}

func modelByOffline() (err error) {
	if gitlabLite, err = lite.NewGitlabLiteMemJson(offlineGroupsDir, offlineGroupsFilesPattern); err != nil {
		return
	}
	if modelHandler, err = newJsonModelWriter(); err != nil {
		return
	}

	err = readGroupsModels()
	return
}

func init() {
	groupModelCmd.AddCommand(groupModelOfflineCmd)
}
