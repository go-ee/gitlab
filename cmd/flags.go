package cmd

import (
	"github.com/go-ee/utils/cliu"
	"github.com/urfave/cli/v2"
)

func NewDebugFlag() *cliu.BoolFlag {
	return cliu.NewBoolFlag(&cli.BoolFlag{
		Name:  "debug",
		Usage: "Enable debug log level",
	})
}

func NewTokenFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:     "token",
		Required: true,
		Usage:    "Gitlab token",
	})
}

func NewApiUrlPart() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "urlApiPart",
		Usage: "URL API part",
		Value: "api/v4",
	})
}

func NewWaitForAuthInteractive() *cliu.IntFlag {
	return cliu.NewIntFlag(&cli.IntFlag{
		Name:  "waitForAuth",
		Usage: "Wait duration for interactive authentication delay (in seconds)",
		Value: 20000,
	})
}

func NewInstallBrowsers() *cliu.BoolFlag {
	return cliu.NewBoolFlag(&cli.BoolFlag{
		Name:  "installEmbeddedBrowsers",
		Usage: "Install compatible embedded browsers",
		Value: false,
	})
}

func NewUrlFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:     "url",
		Required: true,
		Usage:    "Gitlab server url",
	})
}

func NewGroupFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:     "group",
		Usage:    "Gitlab group",
		Required: true,
	})
}

func NewGroupsFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:     "groups",
		Usage:    "Gitlab groups (comma separated)",
		Required: true,
	})
}

func NewGroupsFolderFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "groups-folder",
		Usage: "Folder of group JSON files",
		Value: "__gitlab",
	})
}

func NewFilePatternFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "file-pattern",
		Usage: "Log file Name regular expression pattern",
		Value: ".+?\\.json$",
	})
}

func NewIgnoresFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "ignores",
		Usage: "Ignore group names the comma separated groups",
	})
}

func NewDevBranchFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "dev-branch",
		Usage: "Development branch",
		Value: "dev",
	})
}

func NewsScriptsFolderFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "scripts-folder",
		Usage: "Folder where scripts are generated",
		Value: ".",
	})
}

func NewsReposFolderFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "repos-folder",
		Usage: "Folder where repositories will be cloned and pulled",
		Value: "src",
	})
}

func NewJsonFileFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "json-file",
		Usage: "JSON file name",
		Value: "gitlab.json",
	})

}
