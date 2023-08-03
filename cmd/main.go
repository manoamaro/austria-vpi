package main

import (
	"github.com/manoamaro/vpi2015/internal"
	"github.com/spf13/cobra"
)

const (
	// Version is the current version of the application
	Version = "0.0.1"
)

var rootCmd = &cobra.Command{
	Use:     "vpi2015",
	Short:   "A command line tool to query the Austrian consumer price index VPI 2015",
	Long:    "A command line tool to query the Austrian consumer price index VPI 2015",
	Version: Version,
}

func main() {
	rootCmd.AddCommand(internal.DiffCmd())

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}
