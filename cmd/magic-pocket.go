package main

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/fjboy/magic-pocket/cmd/commands"
	"github.com/fjboy/magic-pocket/pkg/global/gitutils"
	"github.com/fjboy/magic-pocket/pkg/global/logging"
)

var Version string
var (
	debug bool
)

func getVersion() string {
	if Version == "" {
		return gitutils.GetVersion()
	}
	return fmt.Sprint(Version)
}

func main() {
	var rootCmd = &cobra.Command{
		Use:     "magic-pocket",
		Short:   "Golang 工具集",
		Long:    "Golang 实现的工具合集",
		Version: getVersion(),
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			level := logging.INFO
			if debug {
				level = logging.DEBUG
			}
			logging.BasicConfig(logging.LogConfig{Level: level})
		},
	}

	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "显示Debug信息")

	rootCmd.AddCommand(
		commands.BingImgDownloadCmd,
		commands.HttpDownloadCmd,
		commands.SimpleHttpFS,
		commands.IniCrud,
	)

	rootCmd.Execute()
}
