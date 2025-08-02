package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func main() {

	var rootCmd = &cobra.Command{
		Use:   "upscayl",
		Short: "upscayl-cli is a tool to Upscayl images using command line or server",
		Long:  "upscayl-cli tool is for users who wants to Upscayl (https://github.com/upscayl/upscayl) images using command line or a server.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Welcome to upscayl cli! Start upscayling your images using command line or start a server. Use --help for available commands.")
		},
	}
	rootCmd.AddCommand(getServeCommand())
	rootCmd.AddCommand(getRunCommand())
	rootCmd.AddCommand(getResetCommand())
	rootCmd.AddCommand(getDownloadCommand())
	rootCmd.AddCommand(getListCommand())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
