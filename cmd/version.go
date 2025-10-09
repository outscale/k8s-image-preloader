package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   = "DEV"
	buildDate = ""
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v"},
	Short:   "Display k8s-image-preloader version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("k8s-image-preloader version %s built on %s\n", version, buildDate)
	},
}
