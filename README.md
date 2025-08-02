# 📦 Upscayl-cli

A command line tool to run [Upscayl](https://github.com/upscayl/upscayl) without GUI

---

## ✨ Features

- ✅ Upscayl your images with command line
- 🌐 Optionally run a server to integrate with other tools
- 🔧 Most of the original settings supported


## 💾 Installation

Currently, the only way to run this tool is to either build locally or download from releases (packages coming soon!)

### Build locally

#### Prerequisites

`go 1.17`

#### Commands

```
cd go/src/github.com
git clone https://github.com/yashschandra/upscayl-cli.git
cd upscayl-cli
go mod download
make local # executable will appear in release/build/local directory
```

### Download a release

You can download the releases from [here](https://github.com/yashschandra/upscayl-cli/releases). Only Mac (intel/silicon) and Linux supported.
Current latest release is version `v0.0.5`

## 📚 Usage

By default `upscayl-standard-4x` model is used.

### Basic usage

To Upscayl an image, either pass the path of the image -

```
./path/to/upscayl run -i /path/to/input-image -o /path/to/output-image
```

OR

pass the url of the image - 

```
./path/to/upscayl run -u https://your/image/url -o /path/to/output-image
```

### Download models

To download a particular Upscayl model use the download command -

```
./path/to/upscayl download [MODEL NAME]
```

### Run a server

```
./path/to/upscayl serve -p [PORT]
```

### Use server api

```
curl -X POST http://localhost:[PORT]/upscayl \
     -H "Content-Type: application/json" \
     -d '{
           "imagePath": "/path/to/input-image",
           "outputPath": "/path/to/output-image"
         }'
```

### Advance usage

We support a variety of settings that are supported originally. Just use help option to check how to set them and if not then what are the default values.

```
./path/to/upscayl run --help
```

## 🤝 Contributing

We welcome contributions from developers all around the world to help evolve this project! 🌍✨

Whether you're fixing bugs, suggesting new features, improving documentation, or just sharing feedback — every bit counts and is truly appreciated.

### 🧭 How to Contribute

1. Fork the repository

2. Create a branch for your feature/fix

3. Make your changes with clear commits

4. Submit a Pull Request with a short description of your work

5. Wait for review and feedback

>💡 No contribution is too small — even fixing a typo helps!

## 📝 In progress

1. Provide the tool as package for Mac and Linux
2. Provide commands to list available models
3. Improve documentation