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
	"github.com/go-ee/filegen/gen"
	"github.com/go-ee/gitlab/templates"
	"github.com/spf13/cobra"
	"path/filepath"
)

var gitScriptsCmdUse = "git-scripts"
var gitScriptsCmdShort = "Generate Git scripts (clone, pull, push, etc.) based on group model"

// gitScriptsCmd represents the generateGitScripts command
var gitScriptsCmd = &cobra.Command{
	Use:              gitScriptsCmdUse,
	Short:            gitScriptsCmdShort,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		err = generateGitScripts()
		return
	},
}

var gitScriptsByApiCmd = &cobra.Command{
	Use:              gitScriptsCmdUse,
	Short:            gitScriptsCmdShort,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = modelByApi(); err != nil {
			return
		}
		err = generateGitScripts()
		return
	},
}

var gitScriptsByBrowserCmd = &cobra.Command{
	Use:              gitScriptsCmdUse,
	Short:            gitScriptsCmdShort,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = modelByBrowser(); err != nil {
			return
		}
		err = generateGitScripts()
		return
	},
}

var gitScriptsByOfflineCmd = &cobra.Command{
	Use:              gitScriptsCmdUse,
	Short:            gitScriptsCmdShort,
	TraverseChildren: true,
	RunE: func(cmd *cobra.Command, args []string) (err error) {
		if err = modelByOffline(); err != nil {
			return
		}
		err = generateGitScripts()
		return
	},
}

func generateGitScripts() (err error) {
	var templateProvider *gen.NextTemplateProvider
	if templateProvider, err = gen.NewNextTemplateProviderFromText(
		templates.Templates(), templates.MacrosTemplates()); err != nil {
		return
	}

	var absOutDir string
	if absOutDir, err = filepath.Abs(outputDir); err != nil {
		return
	}

	var dataFiles []string
	if dataFiles, err = gen.CollectFilesRecursive(filepath.Join(absOutDir, groupsModelFileName)); err != nil {
		return
	}
	templateDataProvider := &gen.ArrayNextProvider[gen.DataLoader]{
		Items: gen.FilesToTemplateDataLoaders(dataFiles),
	}
	generator := &gen.Generator{
		FileNameBuilder: &gen.DefaultsFileNameBuilder{
			OutputPath: absOutDir, RelativeToTemplate: false, RelativeToData: true},
		NextTemplateLoader:     templateProvider,
		NextTemplateDataLoader: templateDataProvider,
	}
	err = generator.Generate()
	return
}

func init() {
	groupModelGitlabApiCmd.AddCommand(gitScriptsByApiCmd)
	groupModelBrowserCmd.AddCommand(gitScriptsByBrowserCmd)
	groupModelOfflineCmd.AddCommand(gitScriptsByOfflineCmd)

	rootCmd.AddCommand(gitScriptsCmd)

	FlagGroupModelFileName(gitScriptsCmd.PersistentFlags(), &groupsModelFileName)
	FlagOutputDir(gitScriptsCmd.PersistentFlags(), &outputDir)
}
