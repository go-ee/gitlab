package cmd

import (
	"github.com/go-ee/gitlab/lite"
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
	"github.com/urfave/cli/v2"
)

type GroupModelGitLabAPI struct {
	*GroupModelBase
	token, url *cliu.StringFlag
}

func NewGroupModelByGitLabAPI() (o *GroupModelGitLabAPI) {
	o = &GroupModelGitLabAPI{
		GroupModelBase: NewGroupModelBase(),
		token:          NewTokenFlag(),
		url:            NewUrlFlag(),
	}

	o.Command = &cli.Command{
		Name:  "group-model",
		Usage: "Build group model from over Gitlab API to a JSON file",
		Flags: []cli.Flag{
			o.token, o.url,
			o.group, o.ignores, o.groupModelFile,
		},
		Action: func(c *cli.Context) (err error) {
			if err = o.prepareJsonFile(c); err != nil {
				return
			}
			lg.LOG.Debugf("execute %v to %v", c.Command.Name, o.groupModelFile)

			var gitlabLite *lite.GitlabLiteByAPI
			if gitlabLite, err = o.gitlabLiteByAPI(); err == nil {
				var groupNode *lite.GroupNode
				if groupNode, err = o.extract(gitlabLite); err == nil {
					err = o.writeJsonFile(groupNode)
				} else {
					lg.LOG.Errorf("error %v by %v to %v", err, c.Command.Name, o.groupModelFile)
				}
			}
			return
		},
	}
	return
}

func (o *GroupModelGitLabAPI) gitlabLiteByAPI() (ret *lite.GitlabLiteByAPI, err error) {
	ret, err = lite.NewGitlabLiteByAPI(&lite.ServerAccess{Url: o.url.CurrentValue, Token: o.token.CurrentValue})
	return
}
