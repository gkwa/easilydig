package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/gkwa/easilydig/version"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of easilydig",
	Long:  `All software has versions. This is easilydig's`,
	Run: func(cmd *cobra.Command, args []string) {
		buildInfo := version.GetBuildInfo()
		fmt.Println(buildInfo)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
