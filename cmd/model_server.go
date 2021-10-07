package cmd

import (
	"github.com/go-ee/gitlab/core"
	"github.com/go-ee/utils/cliu"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

type ModelServer struct {
	*ModelBase
	token, url *cliu.StringFlag
}

func NewModelServer() (o *ModelServer) {
	o = &ModelServer{
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

			var gitlabLite *core.GitlabLiteByServer
			if gitlabLite, err = o.gitlabLiteByServer(); err == nil {
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

func (o *ModelServer) gitlabLiteByServer() (ret *core.GitlabLiteByServer, err error) {
	ret, err = core.NewGitlabLiteServer(&core.ServerAccess{Url: o.url.CurrentValue, Token: o.token.CurrentValue})
	return
}
