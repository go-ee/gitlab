package cmd

import (
	"github.com/go-ee/gitlab/core"
	"github.com/go-ee/utils/cliu"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type ModelGitLabAPI struct {
	*ModelBase
	token, url *cliu.StringFlag
}

func NewModelGitLabAPI() (o *ModelGitLabAPI) {
	o = &ModelGitLabAPI{
		ModelBase: NewModelBase(),
		token:     NewTokenFlag(),
		url:       NewUrlFlag(),
	}

	o.Command = &cli.Command{
		Name:  "model",
		Usage: "Build group model from over Gitlab API to a JSON file",
		Flags: []cli.Flag{
			o.token, o.url,
			o.group, o.ignores, o.jsonFile,
		},
		Action: func(c *cli.Context) (err error) {
			if err = o.prepareJsonFile(c); err != nil {
				return
			}
			logrus.Debugf("execute %v to %v", c.Command.Name, o.jsonFile)

			var gitlabLite *core.GitlabLiteByAPI
			if gitlabLite, err = o.gitlabLiteByAPI(); err == nil {
				var groupNode *core.GroupNode
				if groupNode, err = o.extract(gitlabLite); err == nil {
					err = o.writeJsonFile(groupNode)
				} else {
					logrus.Errorf("error %v by %v to %v", err, c.Command.Name, o.jsonFile)
				}
			}
			return
		},
	}
	return
}

func (o *ModelGitLabAPI) gitlabLiteByAPI() (ret *core.GitlabLiteByAPI, err error) {
	ret, err = core.NewGitlabLiteByAPI(&core.ServerAccess{Url: o.url.CurrentValue, Token: o.token.CurrentValue})
	return
}