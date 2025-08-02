package main

import (
	"github.com/spf13/cobra"
	"github.com/yashschandra/upscayl-cli/upscayl"
	"log"
)

func getListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list all models available for use locally",
		Run: func(cmd *cobra.Command, args []string) {
			err := upscayl.List()
			if err != nil {
				log.Fatal("error while listing", err.Error())
			}
		},
	}
	return cmd
}
