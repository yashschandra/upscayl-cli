// Package upscayl provides a lightweight Go-based CLI wrapper to download,
// configure, and run the Upscayl image upscaling tool. It supports downloading
// models and binaries, managing user config directories, and running the
// Upscayl binary with custom parameters for local or remote images.
package upscayl

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"os/user"
	"path"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	// baseFolderName is the main directory under user's home for storing Upscayl resources.
	baseFolderName = ".upscayl-cli"

	// baseModelPath defines the relative path for storing models.
	baseModelPath = "/resources/models"

	// basebinaryPath is the relative path for storing the Upscayl binary.
	basebinaryPath = "/resources/bin/upscayl-bin"

	// defaultModel is the default model used for image upscaling.
	defaultModel = "upscayl-standard-4x"

	// binaryUrlFmt is the format URL for downloading the binary based on OS.
	binaryUrlFmt = "https://raw.githubusercontent.com/upscayl/upscayl/main/resources/%s/bin/upscayl-bin"

	// modelBinUrlFmt is the format URL for downloading model .bin files.
	modelBinUrlFmt = "https://raw.githubusercontent.com/upscayl/upscayl/main/resources/models/%s.bin"

	// modelParamUrlFmt is the format URL for downloading model .param files.
	modelParamUrlFmt = "https://raw.githubusercontent.com/upscayl/upscayl/main/resources/models/%s.param"
)

var (
	rootDir           string
	modelsPath        string
	binaryPath        string
	defaultOutputPath string

	// osNameMap maps runtime.GOOS to corresponding folder names for binary download.
	osNameMap = map[string]string{
		"darwin": "mac",
		"linux":  "linux",
	}
)

// isFileExist checks whether a file already exists by trying to exclusively create it.
// Returns true if the file already exists.
func isFileExist(fpath string) bool {
	_, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
	return os.IsExist(err)
}

// download retrieves a file from the given URL and writes it to the specified file path.
// It returns an error if the HTTP request fails or writing to file fails.
func download(url, fpath string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: %s", resp.Status)
	}
	outFile, err := os.Create(fpath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	_, err = io.Copy(outFile, resp.Body)
	return err
}

// downloadModel downloads both .bin and .param files for a given model name
// if they are not already present locally.
func downloadModel(model string) error {
	binPath := filepath.Join(modelsPath, model+".bin")
	if !isFileExist(binPath) {
		log.Printf("downloading %s bin file\n", model)
		err := download(fmt.Sprintf(modelBinUrlFmt, model), binPath)
		if err != nil {
			return fmt.Errorf("error downloading bin file %s", err.Error())
		}
	}
	paramPath := filepath.Join(modelsPath, model+".param")
	if !isFileExist(paramPath) {
		log.Printf("downloading %s param file\n", model)
		err := download(fmt.Sprintf(modelParamUrlFmt, model), paramPath)
		if err != nil {
			return fmt.Errorf("error downloading param file %s", err.Error())
		}
	}
	return nil
}

// listModels prints the names of all .bin model files in the root directory.
func listModels() error {
	return filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".bin") {
			fmt.Println(info.Name())
		}
		return nil
	})
}

// init sets up the default directories and downloads the binary and default model
// if they are not already present.
func init() {
	usr, err := user.Current()
	if err != nil {
		log.Fatal("Error getting current user:", err)
		return
	}
	rootDir = filepath.Join(usr.HomeDir, baseFolderName)
	modelsPath = filepath.Join(rootDir, baseModelPath)
	binaryPath = filepath.Join(rootDir, basebinaryPath)
	defaultOutputPath = usr.HomeDir

	if !isFileExist(filepath.Join(modelsPath, defaultModel+".bin")) {
		err := os.MkdirAll(modelsPath, 0755)
		if err != nil {
			log.Fatal("error creating models dir", err.Error())
		}
		err = downloadModel(defaultModel)
		if err != nil {
			log.Fatal("error downloading model", err.Error())
		}
	}

	if !isFileExist(binaryPath) {
		err := os.MkdirAll(filepath.Dir(binaryPath), 0755)
		if err != nil {
			log.Fatal("error creating binary dir", err.Error())
		}
		log.Println("downloading binary file")
		err = download(fmt.Sprintf(binaryUrlFmt, osNameMap[runtime.GOOS]), binaryPath)
		if err != nil {
			log.Fatal("error downloading bin file", err.Error())
		}
		err = os.Chmod(binaryPath, 0755)
		if err != nil {
			log.Fatal("error in giving permission", err.Error())
		}
	}
}

// Input holds all parameters required to perform image upscaling with Upscayl.
type Input struct {
	ImagePath      string `json:"imagePath"`      // Path to local image file
	ImageURL       string `json:"imageUrl"`       // URL of the image to upscale
	OutputPath     string `json:"outputPath"`     // Path to save the upscaled image
	ModelPath      string `json:"modelPath"`      // Path to model directory
	Model          string `json:"model"`          // Model name to use (without extension)
	SaveImageAs    string `json:"saveImageAs"`    // Output image format (png, jpg, etc.)
	GPUId          *int   `json:"gpuId"`          // Optional GPU ID to use
	Scale          string `json:"scale"`          // Scale factor (default is "4")
	Overwrite      bool   `json:"overwrite"`      // Whether to overwrite existing file
	Compression    string `json:"compression"`    // Output compression level
	CustomWidth    int    `json:"customWidth"`    // Desired width if UseCustomWidth is true
	UseCustomWidth bool   `json:"useCustomWidth"` // Flag to enable custom width
	TileSize       *int   `json:"tileSize"`       // Optional tile size for upscaling
	TTAMode        bool   `json:"ttaMode"`        // Enable test-time augmentation mode
}

// Upscayl validates and processes the input configuration, downloads remote image
// if needed, and then runs the upscaling process. Returns the path to the output image.
func Upscayl(input Input) (string, error) {
	if input.ImagePath == "" && input.ImageURL == "" {
		return "", errors.New("input path or url not set")
	}
	if input.ImagePath != "" && input.ImageURL != "" {
		return "", errors.New("only 1 of the input allowed")
	}
	if input.ImageURL != "" {
		parsedURL, err := url.Parse(input.ImageURL)
		if err != nil {
			return "", err
		}
		ext := path.Ext(parsedURL.Path)
		if ext != ".png" && ext != ".jpg" && ext != ".jpeg" && ext != ".webp" {
			return "", errors.New("invalid extension in url expected .png, .jpg, .jpeg, .webp")
		}
		tmpFilePath := filepath.Join(rootDir, "input"+ext)
		err = download(input.ImageURL, tmpFilePath)
		if err != nil {
			return "", err
		}
		defer os.Remove(tmpFilePath)
		input.ImagePath = tmpFilePath
	}
	return upscaylImage(input)
}

// upscaylImage builds the CLI command and executes the Upscayl binary with the provided input.
// Returns the final output path or an error if the command fails.
func upscaylImage(input Input) (string, error) {
	args := make([]string, 0)
	args = append(args, fmt.Sprintf("-i %s", input.ImagePath))
	if input.OutputPath == "" {
		input.OutputPath = filepath.Join(defaultOutputPath, filepath.Base(input.ImagePath))
	}
	args = append(args, fmt.Sprintf("-o %s", input.OutputPath))
	if input.Model == "" {
		input.Model = defaultModel
	}
	args = append(args, fmt.Sprintf("-n %s", input.Model))
	if input.ModelPath == "" {
		input.ModelPath = modelsPath
	}
	args = append(args, fmt.Sprintf("-m %s", input.ModelPath))
	if input.SaveImageAs == "" {
		input.SaveImageAs = filepath.Ext(input.ImagePath)[1:]
	}
	args = append(args, fmt.Sprintf("-f %s", input.SaveImageAs))
	if input.GPUId != nil {
		args = append(args, fmt.Sprintf("-g %d", input.GPUId))
	}
	if input.Scale == "" {
		input.Scale = "4"
	}
	args = append(args, fmt.Sprintf("-s %s", input.Scale))
	if input.Compression == "" {
		input.Compression = "0"
	}
	args = append(args, fmt.Sprintf("-c %s", input.Compression))
	if input.UseCustomWidth {
		args = append(args, fmt.Sprintf("-w %d", input.CustomWidth))
	}
	if input.TileSize != nil {
		args = append(args, fmt.Sprintf("-t %d", input.TileSize))
	}
	if input.TTAMode {
		args = append(args, "-x")
	}
	bashCommand := fmt.Sprintf("%s %s", binaryPath, strings.Join(args, " "))
	fmt.Println("bash command", bashCommand)
	cmd := exec.Command("bash", "-c", bashCommand)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return input.OutputPath, cmd.Run()
}

// Reset deletes the entire root directory along with downloaded models and binary.
func Reset() error {
	return os.RemoveAll(rootDir)
}

// Download manually downloads the specified model's .bin and .param files.
func Download(model string) error {
	return downloadModel(model)
}

// List prints all available model .bin files found in the local directory.
func List() error {
	return listModels()
}
