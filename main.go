package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dece2183/webpage-compile/compiler"

	"gopkg.in/yaml.v3"
)

const (
	DEFAULT_CFG_PATH    = "webpage.yaml"
	DEFAULT_OUTPUT_NAME = "webpage.bin"
)

type Config struct {
	Input       string            `yaml:"input"`
	Output      string            `yaml:"output"`
	Compress    bool              `yaml:"compress"`
	OmitHeaders bool              `yaml:"omit-headers"`
	Extensions  map[string]string `yaml:"extensions"`
	Ignore      []string          `yaml:"ignore"`
}

const helpMessage = `Web pages compiler.

Usage:

	webpage-compile [flags] [path to config]

Flags:

	-h, --help: display this message
	-c, --gen-config: generate default config

The default config looks like:
	input:        .             # input directory
	output:       webpage.bin   # output file
	compress:     false         # compress web page content with gzip
	omit-headers: false         # omit http headers
	extensions:   {}            # add files with specified extentions
	ignore:       []            # ignore files and folders

All paths in the config file are relative to the directory
in which the config file is located.

The 'extensions' field is a map where the key is the file extension
and the value is the MIME type for that file type.

`

func main() {
	var genCfgMode bool

	cfgPath := DEFAULT_CFG_PATH
	for i := 1; i < len(os.Args); i++ {
		arg := os.Args[i]
		if arg[0] == '-' {
			switch arg {
			case "-h", "--help":
				fmt.Print(helpMessage)
				os.Exit(0)
			case "-c", "--gen-config":
				genCfgMode = true
			}
		} else if i == len(os.Args)-1 {
			cfgPath = arg
			cfgPath = strings.ReplaceAll(cfgPath, "\"", "")
			cfgPath = strings.ReplaceAll(cfgPath, "\\", "/")
			cfgPath = filepath.Clean(cfgPath)
		}
	}

	if genCfgMode {
		config := Config{
			Input:  ".",
			Output: DEFAULT_OUTPUT_NAME,
		}

		configBytes, err := yaml.Marshal(config)
		if err != nil {
			fmt.Println("Config serialization error:", err)
			os.Exit(4)
		}

		err = os.WriteFile(cfgPath, configBytes, 0755)
		if err != nil {
			fmt.Println("Config write error:", err)
			os.Exit(5)
		}

		os.Exit(0)
	}

	file, err := os.Open(cfgPath)
	if err != nil {
		fmt.Println("Config error:", err)
		os.Exit(1)
	}

	var config Config
	dec := yaml.NewDecoder(file)
	err = dec.Decode(&config)
	if err != nil {
		fmt.Println("Config decoding error:", err)
		os.Exit(3)
	}

	configDir := filepath.Dir(cfgPath)
	inputPath := config.Input
	if len(inputPath) == 0 {
		inputPath = configDir
	}
	if filepath.IsLocal(inputPath) {
		inputPath = filepath.Join(configDir, inputPath)
	}

	outputPath := config.Output
	if len(outputPath) == 0 {
		outputPath = DEFAULT_OUTPUT_NAME
	}
	if filepath.IsLocal(outputPath) {
		outputPath = filepath.Join(configDir, outputPath)
	}

	err = compiler.Compile(inputPath, outputPath, config.Compress, !config.OmitHeaders, config.Extensions, config.Ignore)
	if err != nil {
		fmt.Println("Compilataion fail:", err)
		os.Exit(5)
	}
}
