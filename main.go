package main

import (
	"github.com/go-ee/gitlab/cmd"
	"github.com/go-ee/utils/lg"
	"os"
)

func main() {
	app := cmd.NewCli()

	if err := app.Run(os.Args); err != nil {
		lg.LOG.Warnf("exit because of error, %v", err)
		os.Exit(1)
	}
	_ = lg.LOG.Sync()
}
