package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/yashschandra/upscayl-cli/upscayl"
	"log"
	"net/http"
)

func getServeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "start server",
		Run: func(cmd *cobra.Command, args []string) {
			port, _ := cmd.Flags().GetInt16("port")
			serve(port)
		},
	}
	cmd.Flags().Int16P("port", "p", 0, "server port")
	_ = cmd.MarkFlagRequired("port")
	return cmd
}

func serve(port int16) {
	log.Printf("starting server at %d", port)
	http.HandleFunc("/upscayl", func(w http.ResponseWriter, r *http.Request) {
		var req upscayl.Input
		err := json.NewDecoder(r.Body).Decode(&req)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(err.Error()))
			return
		}
		outputPath, err := upscayl.Upscayl(req)
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(err.Error()))
			return
		}
		w.WriteHeader(200)
		w.Write([]byte(fmt.Sprintf("image saved at %s", outputPath)))
	})
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Fatal("error in running http server ", err.Error())
	}
}
