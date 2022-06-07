package main

import (
	"github.com/go-ee/gitlab/cmd"
	"github.com/go-ee/utils/cliu"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
	"os"
)

func main() {
	app := NewCli()

	if err := app.Run(os.Args); err != nil {
		logrus.WithFields(logrus.Fields{"err": err}).Warn("exit because of error.")
	}
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
		if o.debug.CurrentValue {
			logrus.SetLevel(logrus.DebugLevel)
		}
		logrus.Debugf("execute %v, %v", c.Command.Name, c.Args())
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
