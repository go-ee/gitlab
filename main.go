package main

import (
	"github.com/go-ee/gitlab/cmd"
	"github.com/go-ee/utils/cliu"
	"github.com/go-ee/utils/lg"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := NewCli()

	if err := app.Run(os.Args); err != nil {
		lg.LOG.Warnf("exit because of error, %v", err)
		os.Exit(1)
	}
	_ = lg.LOG.Sync()
}

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

	o.debug = cmd.NewDebugFlag()

	o.Flags = []cli.Flag{
		o.debug,
	}

	o.Before = func(c *cli.Context) (err error) {
		lg.InitLOG(o.debug.CurrentValue)
		lg.LOG.Debugf("execute %v, %v", c.Command.Name, c.Args())
		return
	}

	o.Commands = []*cli.Command{
		cmd.NewGroupsDownloaderByAPI().Command,
		cmd.NewGroupsDownloaderByBrowser().Command,
		cmd.NewGroupModelByJsonFiles().Command,
		cmd.NewGroupModelByGitLabAPI().Command,
		cmd.NewGroupModelByBrowser().Command,
		cmd.NewScripts().Command,
	}
}
