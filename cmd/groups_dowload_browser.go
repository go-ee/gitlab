package cmd

import (
	"fmt"
	"github.com/go-ee/gitlab/core"
	"github.com/go-ee/utils/cliu"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
	"strings"
)

type GroupsDownloaderByBrowser struct {
	*cli.Command
	groups, ignores, url, groupsFolder, urlApiPart *cliu.StringFlag
	waitForAuth                                    *cliu.IntFlag
	installBrowsers                                *cliu.BoolFlag
}

func NewGroupsDownloaderByBrowser() (o *GroupsDownloaderByBrowser) {
	o = &GroupsDownloaderByBrowser{
		groups:          NewGroupsFlag(),
		ignores:         NewIgnoresFlag(),
		url:             NewUrlFlag(),
		groupsFolder:    NewGroupsFolderFlag(),
		urlApiPart:      NewApiUrlPart(),
		waitForAuth:     NewWaitForAuthInteractive(),
		installBrowsers: NewInstallBrowsers(),
	}

	o.Command = &cli.Command{
		Name:  "groups-download-browser",
		Usage: "Download Gitlab groups JSON files by browser automation to group JSON files",
		Flags: []cli.Flag{
			o.url, o.groups, o.ignores, o.waitForAuth, o.groupsFolder, o.urlApiPart, o.installBrowsers,
		},
		Action: func(c *cli.Context) (err error) {
			logrus.Debugf("execute %v for %v", c.Command.Name, o.groups.CurrentValue)

			var gitlabLite *core.GitlabLiteByBrowser
			if gitlabLite, err = o.gitlabLiteByBrowser(); err != nil {
				return
			}
			modelWriter := &core.ModelWriter{GroupsFolder: o.groupsFolder.CurrentValue}
			if err = gitlabLite.AuthInteractive(o.waitForAuth.CurrentValue); err == nil {
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
			}
			return
		},
	}
	return
}

func (o *GroupsDownloaderByBrowser) gitlabLiteByBrowser() (ret *core.GitlabLiteByBrowser, err error) {
	if err = os.MkdirAll(o.groupsFolder.CurrentValue, 0755); err == nil {
		ret, err = core.NewGitlabLiteByBrowser(o.buildBrowserAccess(), o.installBrowsers.CurrentValue)
	}
	return
}

func (o *GroupsDownloaderByBrowser) buildBrowserAccess() *core.BrowserAccess {
	return &core.BrowserAccess{
		UrlAuth: o.url.CurrentValue,
		UrlApi:  fmt.Sprintf("%v/%v", o.url.CurrentValue, o.urlApiPart.CurrentValue)}
}
