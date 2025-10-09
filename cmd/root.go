package cmd

import (
	"github.com/spf13/cobra"
)

var (
	ctrPath  string
	ctrFlags string
	debug    bool
)

func init() {
	rootCmd.PersistentFlags().StringVar(&ctrPath, "ctr-path", "/usr/local/bin/ctr", "ctr binary path")
	rootCmd.PersistentFlags().StringVar(&ctrFlags, "ctr-flags", "-a /var/run/containerd/containerd.sock", "ctr flags")
	rootCmd.PersistentFlags().BoolVar(&debug, "debug", false, "log ctr command output")
}

var rootCmd = &cobra.Command{
	Use:   "preloader",
	Short: "",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	exitOnError(err)
}
