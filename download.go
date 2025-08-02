package main

import (
	"github.com/spf13/cobra"
	"github.com/yashschandra/upscayl-cli/upscayl"
	"log"
)

func getDownloadCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "download",
		Short: "download any upscayl model",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			model := args[0]
			err := upscayl.Download(model)
			if err != nil {
				log.Fatal("error while downloading", err.Error())
			}
		},
	}
	return cmd
}
