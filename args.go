package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"path"
	"path/filepath"
	"strings"
)

type imageFormat string

const (
	imageFormatJpg  = ".jpg"
	imageFormatJpeg = ".jpeg"
	imageFormatPng  = ".png"
)

var (
	allowedInputFormats  = []string{imageFormatJpg, imageFormatJpeg, imageFormatPng}
	allowedOutputFormats = []string{imageFormatJpg, imageFormatJpeg, imageFormatPng}
)

const (
	defaultBgColor = "#000000"
	defaultFgColor = "#FFCF75"
)

type options struct {
	inputFilePath  string
	outputFilePath string
	inputDirPath   string
	outputDirPath  string
	fg             color.Color
	bg             color.Color
}

func (opts *options) isProcessDir() bool {
	return len(opts.inputDirPath) > 0
}

func getArgs() *options {
	var (
		inputFilePathArg  = flag.String("f", "", "the file input path")
		outputFilePathArg = flag.String("of", "", "the file output path")

		inputDirPathArg  = flag.String("d", "", "the directory input path")
		outputDirPathArg = flag.String("od", "", "the directory output path")

		fgArg = flag.String("fg", defaultFgColor, "foreground color in HEX / RGB format")
		bgArg = flag.String("bg", defaultBgColor, "background color in HEX / RGB format")

		helpArg = flag.Bool("h", false, "usage information")
	)

	flag.Parse()

	if *helpArg {
		flag.PrintDefaults()
		os.Exit(0)
	}

	ensureSomeInputPathProvided(inputFilePathArg, inputDirPathArg)
	ensureNotBothInputPathsProvided(inputFilePathArg, inputDirPathArg)

	if len(*inputFilePathArg) > 0 {
		ensureInputFormatAllowed(inputFilePathArg)
		if len(*outputFilePathArg) > 0 {
			ensureOutputFormatAllowed(outputFilePathArg)
		} else {
			outputFilePath := path.Join(
				path.Dir(*inputFilePathArg),
				"out"+filepath.Ext(*inputFilePathArg),
			)
			fmt.Printf("[INFO] no output path provided, using %s as output file path\n", outputFilePath)
			*outputFilePathArg = outputFilePath
		}
	}

	if len(*inputDirPathArg) > 0 {
		if len(*outputDirPathArg) > 0 {
			createOutputDirIfNotExists(outputDirPathArg)
		} else {
			fmt.Printf("[INFO] no ouput path provided, using %s as output directory\n", *inputDirPathArg)
			*outputDirPathArg = *inputDirPathArg
		}
	}

	fg, err := parseColorString(*fgArg)
	if err != nil {
		fmt.Printf("[ERROR] error parsing fg color string: %v\n", err)
		os.Exit(1)
	}

	bg, err := parseColorString(*bgArg)
	if err != nil {
		fmt.Printf("[ERROR] error parsing fg color string: %v\n", err)
		os.Exit(1)
	}

	return &options{
		inputFilePath:  *inputFilePathArg,
		outputFilePath: *outputFilePathArg,
		inputDirPath:   *inputDirPathArg,
		outputDirPath:  *outputDirPathArg,
		fg:             fg,
		bg:             bg,
	}
}

func ensureSomeInputPathProvided(inputFilePathArg, inputDirPathArg *string) {
	if len(*inputFilePathArg) == 0 && len(*inputDirPathArg) == 0 {
		fmt.Printf("[ERROR] no input file or directory path provided\n")
		flag.PrintDefaults()
		os.Exit(1)
	}
}

func ensureNotBothInputPathsProvided(inputFilePathArg, inputDirPathArg *string) {
	if len(*inputFilePathArg) > 0 && len(*inputDirPathArg) > 0 {
		fmt.Printf("[ERROR] both input file and directory path provided\n")
		os.Exit(1)
	}
}

func ensureInputFormatAllowed(inputFilePathArg *string) {
	inputFormat := filepath.Ext(*inputFilePathArg)
	if !isAllowedInputFormat(inputFormat) {
		fmt.Printf(
			"[ERROR] format '%s' is not an allowed input format, allowed formats: %s\n",
			inputFormat,
			strings.Join(allowedInputFormats, " | "),
		)
		os.Exit(1)
	}
}

func ensureOutputFormatAllowed(outputFilePathArg *string) {
	outputFormat := filepath.Ext(*outputFilePathArg)
	if !isAllowedOutputFormat(outputFormat) {
		fmt.Printf(
			"[ERROR] format '%s' is not an allowed output format, allowed formats: %s\n",
			outputFormat,
			strings.Join(allowedOutputFormats, " | "),
		)
		os.Exit(1)
	}
}

func createOutputDirIfNotExists(outputDirPathArg *string) {
	dirExists, err := exists(*outputDirPathArg)
	if err != nil {
		fmt.Printf("[ERROR] failed to get info on '%s': %v\n", *outputDirPathArg, err)
		os.Exit(1)
	}

	if !dirExists {
		err := os.MkdirAll(*outputDirPathArg, 0777)
		if err != nil {
			fmt.Printf("[ERROR] failed to create directory '%s': %v\n", *outputDirPathArg, err)
			os.Exit(1)
		}
		fmt.Printf("[INFO] created output directory '%s'\n", *outputDirPathArg)
	}
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func contains[T comparable](a []T, item T) bool {
	for _, el := range a {
		if el == item {
			return true
		}
	}
	return false
}

func isAllowedInputFormat(format string) bool {
	return contains(allowedInputFormats, format)
}

func isAllowedOutputFormat(format string) bool {
	return contains(allowedOutputFormats, format)
}
