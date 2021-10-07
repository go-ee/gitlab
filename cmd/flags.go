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

func NewUrlFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:     "url",
		Required: true,
		Usage:    "Base Gitlab server url",
	})
}

func NewGroupFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:     "group",
		Usage:    "Gitlab group",
		Required: true,
	})
}

func NewGroupsFolderFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "groups-folder",
		Usage: "Folder of group JSON files",
		Value: "gitlab",
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
		Usage: "Ignore group names the comma separated groups",
		Value: "develop",
	})
}

func NewsScriptsFolderFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "scripts-folder",
		Usage: "Folder where scripts are generated",
		Value: ".",
	})
}

func NewJsonFileFlag() *cliu.StringFlag {
	return cliu.NewStringFlag(&cli.StringFlag{
		Name:  "json-file",
		Usage: "JSON file name",
		Value: "gitlab.json",
	})

}
