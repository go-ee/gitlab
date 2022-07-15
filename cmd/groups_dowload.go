package cmd

import (
	"github.com/go-ee/gitlab/core"
	"github.com/go-ee/utils/cliu"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"strings"
)

type GroupsDownloaderByAPI struct {
	*GroupModelGitLabAPI
	groups, groupsFolder *cliu.StringFlag
}

func NewGroupsDownloaderByAPI() (o *GroupsDownloaderByAPI) {
	o = &GroupsDownloaderByAPI{
		GroupModelGitLabAPI: NewGroupModelByGitLabAPI(),
		groups:              NewGroupsFlag(),
		groupsFolder:        NewGroupsFolderFlag(),
	}

	o.Command = &cli.Command{
		Name:  "groups-download",
		Usage: "Download Gitlab groups JSON files to group JSON files",
		Flags: []cli.Flag{
			o.token, o.url, o.groups, o.ignores, o.groupsFolder,
		},
		Action: func(c *cli.Context) (err error) {
			logrus.Debugf("execute %v for %v", c.Command.Name, o.groups.CurrentValue)

			var gitlabLite *core.GitlabLiteByAPI
			if gitlabLite, err = o.gitlabLiteByAPI(); err != nil {
				return
			}

			modelWriter := &core.ModelWriter{GroupsFolder: o.groupsFolder.CurrentValue}
			groups := strings.Split(o.groups.CurrentValue, ",")
			for _, group := range groups {
				if groupNode, groupErr := core.Extract(&core.ExtractParams{
					Group:            group,
					IgnoreGroupNames: buildIgnoresMap(o.ignores.CurrentValue),
				}, gitlabLite); groupErr != nil {
					logrus.Warnf("error at downloading of JSON for group %v", group)
				} else {
					if groupWriter := modelWriter.OnGroupNode(groupNode); groupWriter != nil {
						logrus.Warnf("error at writing of JSON for group %v", group)
					}
				}
			}
			return
		},
	}
	return
}
