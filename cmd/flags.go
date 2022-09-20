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
	flagName = "ignore-groups"
	flagSet.StringSliceVar(p, flagName, nil, "Ignore Gitlab groups (ID or name), semicolon separated or multiple flags")
	return
}

func FlagOutputDir(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "output-dir"
	flagSet.StringVarP(p, flagName, "o", ".", "output directory")
	return
}

func FlagOfflineGroupsDir(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "offline-groups-dir"
	flagSet.StringVarP(p, flagName, "", offlineGroupsDir, "Directory for single Gitlab groups JSON files")
	return
}

func FlagOfflineSupport(flagSet *pflag.FlagSet, p *bool) (flagName string) {
	flagName = "offline-support"
	flagSet.BoolVarP(p, flagName, "", offlineModeSupport,
		"Some operations of the tool could be executed without connection to Gitlab server. In order to support it, single Gitlab groups representation have to be downloaded if connection to Gitlab server is available.")
	return
}

func FlagStoreGroupModel(flagSet *pflag.FlagSet, p *bool) (flagName string) {
	flagName = "store-model"
	flagSet.BoolVarP(p, flagName, "", storeGroupModel,
		"The locally stored group models can be reused (e.g. to improve Git script generation) without having to read them in from scratch.")
	return
}

func FlagGroupModelFileName(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "group-model-fileName"
	flagSet.StringVarP(p, flagName, "f", groupsModelFileName, "JSON ")
	return
}

func FlagGitlabUrl(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "url"
	flagSet.StringVarP(p, flagName, "u", "", "Gitlab URL")
	return
}

func FlagGitlabUrlApiPart(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "url-api-part"
	flagSet.StringVarP(p, flagName, "", gitlabUrlApiPart, "Gitlab URL API part as suffix for Gitlab URL")
	return
}

func FlagGitlabAccessToken(flagSet *pflag.FlagSet, p *string) (flagName string) {
	flagName = "token"
	flagSet.StringVarP(p, flagName, "t", "", "Gitlab token")
	return
}

func FlagWaitForAuthInteractive(flagSet *pflag.FlagSet, p *int) (flagName string) {
	flagName = "wait-for-auth"
	flagSet.IntVarP(p, flagName, "", waitForAuthInteractive, "Wait duration for interactive authentication delay (in seconds)")
	return
}

func FlagInstallEmbeddedBrowsers(flagSet *pflag.FlagSet, p *bool) (flagName string) {
	flagName = "drivers-and-embedded-browsers"
	flagSet.BoolVarP(p, flagName, "", installDriversAndEmbeddedBrowsers, "Install browser drivers and compatible embedded browsers or reuse browser of the system.")
	return
}
