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
	baseFolderName = ".upscayl-cli"
	baseModelPath  = "/resources/models"
	basebinaryPath = "/resources/bin/upscayl-bin"

	defaultModel = "upscayl-standard-4x"

	binaryUrlFmt     = "https://raw.githubusercontent.com/upscayl/upscayl/main/resources/%s/bin/upscayl-bin"
	modelBinUrlFmt   = "https://raw.githubusercontent.com/upscayl/upscayl/main/resources/models/%s.bin"
	modelParamUrlFmt = "https://raw.githubusercontent.com/upscayl/upscayl/main/resources/models/%s.param"
)

var (
	rootDir           string
	modelsPath        string
	binaryPath        string
	defaultOutputPath string
	osNameMap         = map[string]string{
		"darwin": "mac",
		"linux":  "linux",
	}
)

func isFileExist(fpath string) bool {
	_, err := os.OpenFile(fpath, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0755)
	return os.IsExist(err)
}

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

type Input struct {
	ImagePath      string `json:"imagePath"`
	ImageURL       string `json:"imageUrl"`
	OutputPath     string `json:"outputPath"`
	ModelPath      string `json:"modelPath"`
	Model          string `json:"model"`
	SaveImageAs    string `json:"saveImageAs"`
	GPUId          *int   `json:"gpuId"`
	Scale          string `json:"scale"`
	Overwrite      bool   `json:"overwrite"`
	Compression    string `json:"compression"`
	CustomWidth    int    `json:"customWidth"`
	UseCustomWidth bool   `json:"useCustomWidth"`
	TileSize       *int   `json:"tileSize"`
	TTAMode        bool   `json:"ttaMode"`
}

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

func Reset() error {
	return os.RemoveAll(rootDir)
}

func Download(model string) error {
	return downloadModel(model)
}
