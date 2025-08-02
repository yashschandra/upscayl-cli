package main

import (
	"github.com/spf13/cobra"
	"github.com/yashschandra/upscayl-cli/upscayl"
	"log"
)

func getRunCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Upscayl single image using command line options",
		Run: func(cmd *cobra.Command, args []string) {
			image, _ := cmd.Flags().GetString("input")
			url, _ := cmd.Flags().GetString("url")
			model, _ := cmd.Flags().GetString("model-name")
			input := upscayl.Input{
				ImagePath: image,
				ImageURL:  url,
				Model:     model,
			}
			outputPath, err := upscayl.Upscayl(input)
			if err != nil {
				log.Fatal("error while upscayling", err.Error())
			}
			log.Println("output image at", outputPath)
		},
	}
	cmd.Flags().StringP("input", "i", "", "Input image path (jpg/png/webp) or directory")
	cmd.Flags().StringP("url", "u", "", "Input image url (jpg/png/webp)")
	cmd.Flags().StringP("output", "o", "", "Output image path (jpg/png/webp) or directory")
	cmd.Flags().IntP("model-scale", "z", 4, "Scale according to the model (can be 2, 3, 4)")
	cmd.Flags().IntP("output-scale", "s", 4, "Custom output scale (can be 2, 3, 4)")
	cmd.Flags().StringP("resize", "r", "default", "Resize output to dimension (default=WxH:default), use '-r help' for more details")
	cmd.Flags().IntP("width", "w", 0, "Resize output to a width (default=W:default), use '-r help' for more details")
	cmd.Flags().IntP("compress", "c", 0, "Compression of the output image, default 0 and varies to 100")
	cmd.Flags().StringP("tile-size", "t", "0", "Tile size (>=32/0=auto, default=0) can be 0,0,0 for multi-gpu")
	cmd.Flags().StringP("model-path", "m", "models", "Folder path to the pre-trained models")
	cmd.Flags().StringP("model-name", "n", "realesrgan-x4plus", "Model name")
	cmd.Flags().StringP("gpu-id", "g", "auto", "GPU device to use (default=auto) can be 0,1,2 for multi-gpu")
	cmd.Flags().StringP("threads", "j", "1:2:2", "Thread count for load/proc/save (default=1:2:2)")
	cmd.Flags().BoolP("tta", "x", false, "Enable TTA mode")
	cmd.Flags().StringP("format", "f", "ext/png", "Output image format (jpg/png/webp)")
	cmd.Flags().BoolP("verbose", "v", false, "Verbose output")
	return cmd
}
