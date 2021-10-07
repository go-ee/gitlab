package cmd

import (
	"github.com/go-ee/gitlab/core"
	"github.com/go-ee/utils/cliu"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type ModelJsonFiles struct {
	*ModelBase
	groupsFolder, filePattern *cliu.StringFlag
}

func NewModelJsonFiles() (o *ModelJsonFiles) {
	o = &ModelJsonFiles{
		ModelBase:    NewModelBase(),
		groupsFolder: NewGroupsFolderFlag(),
		filePattern:  NewFilePatternFlag(),
	}

	o.Command = &cli.Command{
		Name:  "model-files",
		Usage: "Build group model from Gitlab group JSON files",
		Flags: []cli.Flag{
			o.groupsFolder, o.filePattern,
			o.group, o.ignores, o.jsonFile,
		},
		Action: func(c *cli.Context) (err error) {
			if err = o.prepareJsonFile(c); err != nil {
				return
			}
			logrus.Debugf("execute %v to %v", c.Command.Name, o.jsonFile)

			var gitlabLite *core.GitlabLiteMem
			if gitlabLite, err = core.NewGitlabLiteMemJson(
				o.groupsFolder.CurrentValue, o.filePattern.CurrentValue); err != nil {
				return
			}

			var groupNode *core.GroupNode
			if groupNode, err = o.extract(gitlabLite); err == nil {
				err = o.writeJsonFile(groupNode)
			} else {
				logrus.Errorf("error %v by %v to %v", err, c.Command.Name, o.jsonFile)
				return
			}
			return
		},
	}
	return
}
