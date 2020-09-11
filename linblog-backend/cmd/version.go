package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "blog-api version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("blog-api v0.2")
	},
}
