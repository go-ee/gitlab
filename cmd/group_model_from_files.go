package cmd

import (
	"github.com/go-ee/gitlab/lite"
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
	"github.com/urfave/cli/v2"
)

type GroupModelFromJsonFiles struct {
	*GroupModelBase
	groupsFolder, filePattern *cliu.StringFlag
}

func NewGroupModelFromJsonFiles() (o *GroupModelFromJsonFiles) {
	o = &GroupModelFromJsonFiles{
		GroupModelBase: NewGroupModelBase(),
		groupsFolder:   NewGroupsFolderFlag(),
		filePattern:    NewFilePatternFlag(),
	}

	o.Command = &cli.Command{
		Name:  "group-model-from-files",
		Usage: "Build group model from Gitlab group JSON files",
		Flags: []cli.Flag{
			o.groupsFolder, o.filePattern,
			o.group, o.ignores, o.groupModelFile,
		},
		Action: func(c *cli.Context) (err error) {
			if err = o.prepareJsonFile(c); err != nil {
				return
			}
			lg.LOG.Debugf("execute %v to %v", c.Command.Name, o.groupModelFile)

			var gitlabLite *lite.GitlabLiteMem
			if gitlabLite, err = lite.NewGitlabLiteMemJson(
				o.groupsFolder.CurrentValue, o.filePattern.CurrentValue); err != nil {
				return
			}

			var groupNode *lite.GroupNode
			if groupNode, err = o.extract(gitlabLite); err == nil {
				err = o.writeJsonFile(groupNode)
			} else {
				lg.LOG.Errorf("error %v by %v to %v", err, c.Command.Name, o.groupModelFile)
				return
			}
			return
		},
	}
	return
}
