package cmd

import (
	"fmt"
	"github.com/go-ee/gitlab/core"
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
	"github.com/urfave/cli/v2"
	"os"
)

type GroupModelBrowser struct {
	*GroupModelBase
	url, groupsFolder, urlApiPart *cliu.StringFlag
	waitForAuth                   *cliu.IntFlag
	installBrowsers               *cliu.BoolFlag
}

func NewGroupModelByBrowser() (o *GroupModelBrowser) {
	o = &GroupModelBrowser{
		GroupModelBase:  NewGroupModelBase(),
		url:             NewUrlFlag(),
		groupsFolder:    NewGroupsFolderFlag(),
		urlApiPart:      NewApiUrlPart(),
		waitForAuth:     NewWaitForAuthInteractive(),
		installBrowsers: NewInstallBrowsers(),
	}

	o.Command = &cli.Command{
		Name:  "group-model-browser",
		Usage: "Build group model by browser automation to a JSON file",
		Flags: []cli.Flag{
			o.url, o.group, o.ignores, o.waitForAuth, o.jsonFile, o.groupsFolder, o.urlApiPart, o.installBrowsers,
		},
		Action: func(c *cli.Context) (err error) {
			if err = o.prepareJsonFile(c); err != nil {
				return
			}
			lg.LOG.Debugf("execute %v to %v", c.Command.Name, o.jsonFile)

			var gitlabLite *core.GitlabLiteByBrowser
			if gitlabLite, err = o.gitlabLiteByBrowser(); err != nil {
				return
			}

			if err = gitlabLite.AuthInteractive(o.waitForAuth.CurrentValue); err == nil {
				var groupNode *core.GroupNode
				if groupNode, err = o.extract(gitlabLite); err == nil {
					err = o.writeJsonFile(groupNode)
				} else {
					lg.LOG.Errorf("error %v by %v to %v", err, c.Command.Name, o.jsonFile)
				}
			}
			return
		},
	}
	return
}

func (o *GroupModelBrowser) gitlabLiteByBrowser() (ret *core.GitlabLiteByBrowser, err error) {
	if err = os.MkdirAll(o.groupsFolder.CurrentValue, 0755); err == nil {
		ret, err = core.NewGitlabLiteByBrowser(o.buildBrowserAccess(), o.installBrowsers.CurrentValue)
	}
	return
}

func (o *GroupModelBrowser) buildBrowserAccess() *core.BrowserAccess {
	return &core.BrowserAccess{
		UrlAuth: o.url.CurrentValue,
		UrlApi:  fmt.Sprintf("%v/%v", o.url.CurrentValue, o.urlApiPart.CurrentValue)}
}
