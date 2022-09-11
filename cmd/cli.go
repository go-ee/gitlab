package cmd

import (
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
	"github.com/urfave/cli/v2"
)

type Cli struct {
	*cli.App
	debug *cliu.BoolFlag
}

func NewCli() (ret *Cli) {
	ret = &Cli{}
	ret.init()
	return
}

func (o *Cli) init() {
	o.App = cli.NewApp()
	o.Usage = "Gitlab automation"
	o.Version = "1.0"

	o.debug = NewDebugFlag()

	o.Flags = []cli.Flag{
		o.debug,
	}

	o.Before = func(c *cli.Context) (err error) {
		lg.InitLOG(o.debug.CurrentValue)
		lg.LOG.Debugf("execute %v, %v", c.Command.Name, c.Args())
		return
	}

	o.Commands = []*cli.Command{
		NewGroupsScriptsByAPI().Command,
		NewGroupsDownloaderByAPI().Command,
		NewGroupsDownloaderByBrowser().Command,
		NewGroupModelFromJsonFiles().Command,
		NewGroupModelByGitLabAPI().Command,
		NewGroupModelByBrowser().Command,
		NewScriptsForGroup().Command,
	}
}
