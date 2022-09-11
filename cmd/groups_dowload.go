package cmd

import (
	"github.com/go-ee/gitlab/lite"
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
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
			lg.LOG.Debugf("execute %v for %v", c.Command.Name, o.groups.CurrentValue)

			var gitlabLite *lite.GitlabLiteByAPI
			if gitlabLite, err = o.gitlabLiteByAPI(); err != nil {
				return
			}

			modelWriter := &lite.ModelWriter{GroupsFolder: o.groupsFolder.CurrentValue}
			groups := strings.Split(o.groups.CurrentValue, ",")
			for _, group := range groups {
				if groupNode, groupErr := lite.FetchGroupModel(&lite.GroupModelParams{
					Group:            group,
					IgnoreGroupNames: buildIgnoresMap(o.ignores.CurrentValue),
				}, gitlabLite); groupErr != nil {
					lg.LOG.Warnf("error at downloading of JSON for group %v", group)
				} else {
					lg.LOG.Warnf("error at writing of JSON for group %v", group)
					if groupWriter := modelWriter.OnGroupNode(groupNode); groupWriter != nil {
					}
				}
			}
			return
		},
	}
	return
}
