package cmd

import (
	"fmt"
	"github.com/go-ee/gitlab/lite"
	"github.com/go-ee/gitlab/script"
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
	"github.com/urfave/cli/v2"
	"strings"
)

type ScriptsForGroups struct {
	*GroupModelGitLabAPI
	groups, scriptsFolder, reposFolder *cliu.StringFlag
}

func NewGroupsScriptsByAPI() (o *ScriptsForGroups) {
	o = &ScriptsForGroups{
		GroupModelGitLabAPI: NewGroupModelByGitLabAPI(),
		groups:              NewGroupsFlag(),
	}

	o.Command = &cli.Command{
		Name:  "groups-scripts",
		Usage: "Generate groups scripts",
		Flags: []cli.Flag{
			o.token, o.url, o.groups, o.ignores,
		},
		Action: func(c *cli.Context) (err error) {
			lg.LOG.Debugf("execute %v for %v", c.Command.Name, o.groups.CurrentValue)

			var gitlabLite *lite.GitlabLiteByAPI
			if gitlabLite, err = o.gitlabLiteByAPI(); err != nil {
				return
			}

			groups := strings.Split(o.groups.CurrentValue, ",")
			for _, group := range groups {
				if groupNode, groupErr := lite.FetchGroupModel(&lite.GroupModelParams{
					Group:            group,
					IgnoreGroupNames: buildIgnoresMap(o.ignores.CurrentValue),
				}, gitlabLite); groupErr == nil {
					if err = script.Generate(groupNode,
						o.buildScriptDirForGroup(groupNode), o.buildReposDirForGroup(groupNode)); err != nil {
						lg.LOG.Warnf("error scripts generation for group '%v'", group)
					}
				} else {
					lg.LOG.Warnf("error group retrieving '%v'", group)
				}
			}
			return
		},
	}
	return
}

func (o *ScriptsForGroups) buildScriptDirForGroup(groupNode *lite.GroupNode) string {
	return fmt.Sprintf("%v/%v", groupNode.Group.Name, o.scriptsFolder)
}

func (o *ScriptsForGroups) buildReposDirForGroup(groupNode *lite.GroupNode) string {
	return fmt.Sprintf("%v/%v", groupNode.Group.Name, o.reposFolder)
}
