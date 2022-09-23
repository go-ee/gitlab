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
	"fmt"
	"github.com/go-ee/gitlab/lite"
	"github.com/spf13/cobra"
)

var waitForAuthInteractive = 40000
var installDriversAndEmbeddedBrowsers = false

// groupModelBrowserCmd represents the groupModelBrowser command
var groupModelBrowserCmd = &cobra.Command{
	Use:   "by-browser",
	Short: "Use browser automation for Gitlab group reading instead of GitLab API",
	Long: `If Gitlab API is not accessible (or there are some problems) use browser automation for reading of Gitlab model. 
            This operation requires GUI desktop for manual interaction for authentication to Gitlab.`,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		return modelByBrowser()
	},
}

func modelByBrowser() (err error) {
	var gitlabLiteNyBrowser *lite.GitlabLiteByBrowser
	if gitlabLiteNyBrowser, err = lite.NewGitlabLiteByBrowser(buildBrowserAccess(), installDriversAndEmbeddedBrowsers); err != nil {
		return
	}
	gitlabLite = gitlabLiteNyBrowser

	if modelHandler, err = newJsonModelWriter(); err != nil {
		return
	}

	if err = gitlabLiteNyBrowser.AuthInteractive(waitForAuthInteractive); err == nil {
		err = readGroupsModels()
	}
	return
}

func buildBrowserAccess() *lite.BrowserAccess {
	return &lite.BrowserAccess{
		UrlAuth: gitlabUrl,
		UrlApi:  fmt.Sprintf("%v/%v", gitlabUrl, gitlabUrlApiPart)}
}

func init() {
	groupModelCmd.AddCommand(groupModelBrowserCmd)

	_ = groupModelBrowserCmd.MarkPersistentFlagRequired(
		FlagGitlabUrl(groupModelBrowserCmd.Flags(), &gitlabUrl))

	FlagWaitForAuthInteractive(groupModelBrowserCmd.PersistentFlags(), &waitForAuthInteractive)
	FlagInstallEmbeddedBrowsers(groupModelBrowserCmd.PersistentFlags(), &installDriversAndEmbeddedBrowsers)
	FlagGitlabUrlApiPart(groupModelBrowserCmd.PersistentFlags(), &gitlabUrlApiPart)
}
