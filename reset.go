package main

import (
	"github.com/spf13/cobra"
	"github.com/yashschandra/upscayl-cli/upscayl"
	"log"
)

func getResetCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "reset",
		Short: "reset all upscayl cli data",
		Run: func(cmd *cobra.Command, args []string) {
			err := upscayl.Reset()
			if err != nil {
				log.Fatal("error while resetting", err.Error())
			}
		},
	}
	return cmd
}
