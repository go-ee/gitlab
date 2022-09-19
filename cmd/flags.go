package cmd

import (
	"github.com/spf13/pflag"
)

func FlagGroups(flagSet *pflag.FlagSet, p *[]string) (flagName string) {
	flagName = "groups"
	flagSet.StringSliceVar(p, flagName, nil, "Gitlab groups (ID or name), semicolon separated or multiple flags")
	return
}

func FlagIgnoreGroups(flagSet *pflag.FlagSet, p *[]string) (flagName string) {
	flagName = "ignoreGroups"
	flagSet.StringSliceVar(p, flagName, nil, "Ignore Gitlab groups (ID or name), semicolon separated or multiple flags")
	return
}

func FlagOutputDir(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "outputDir"
	flagSet.StringVarP(p, flagName, "o", ".", "output directory")
	return
}

func FlagOfflineGroupsDir(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "OfflineGroupsDir"
	flagSet.StringVarP(p, flagName, "", offlineGroupsDir, "Directory for single Gitlab groups JSON files")
	return
}

func FlagOfflineModeSupport(flagSet *pflag.FlagSet, p *bool) (flagName string) {
	flagName = "OfflineModeSupport"
	flagSet.BoolVarP(p, flagName, "", offlineModeSupport,
		"Some operations of the tool could be executed without connection to Gitlab server. In order to support it, single Gitlab groups representation have to be downloaded if connection to Gitlab server is available.")
	return
}

func FlagGroupModelFileName(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "groupModelFileName"
	flagSet.StringVarP(p, flagName, "f", groupsModelFileName, "JSON ")
	return
}

func FlagGitlabUrl(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "url"
	flagSet.StringVarP(p, flagName, "u", "", "Gitlab URL")
	return
}

func FlagGitlabUrlApiPart(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "urlApiPart"
	flagSet.StringVarP(p, flagName, "", gitlabUrlApiPart, "Gitlab URL API part as suffix for Gitlab URL")
	return
}

func FlagGitlabAccessToken(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "token"
	flagSet.StringVarP(p, flagName, "t", "", "Gitlab token")
	return
}

func FlagWaitForAuthInteractive(flagSet *pflag.FlagSet, p *int) (flagName string) {
	flagName = "waitForAuth"
	flagSet.IntVarP(p, flagName, "", waitForAuthInteractive, "Wait duration for interactive authentication delay (in seconds)")
	return
}

func FlagInstallEmbeddedBrowsers(flagSet *pflag.FlagSet, p *bool) (flagName string) {
	flagName = "installEmbeddedBrowsers"
	flagSet.BoolVarP(p, flagName, "", installBrowsers, "Install compatible embedded browsers or reuse browser of the system.")
	return
}
