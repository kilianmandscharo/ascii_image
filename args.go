package main

import (
	"flag"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	defaultBgColor = "#000000"
	defaultFgColor = "#FFCF75"
)

type options struct {
	inputPath    string
	outputPath   string
	fg           color.Color
	bg           color.Color
	inputIsDir   bool
}

func getArgs() *options {
	var inputPathArg = flag.String("f", "", "the file input path")
	var outputPathArg = flag.String("o", "", "the output path")
	var fgArg = flag.String("fg", defaultFgColor, "foreground color in HEX / RGB format")
	var bgArg = flag.String("bg", defaultBgColor, "background color in HEX / RGB format")

	flag.Parse()

	if len(*inputPathArg) == 0 || len(*outputPathArg) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	inputIsDir := isDir(*inputPathArg)
	outputIsDir := inputIsDir
	if inputIsDir {
		outputIsDir = isDir(*outputPathArg)
	}

	if inputIsDir != outputIsDir {
		log.Fatalf("input and output have to be both either a file or a directory")
	}

	if !inputIsDir {
		ensureAllowedOutputFormat(filepath.Ext(*outputPathArg))
	}

	fg, err := parseColorString(*fgArg)
	if err != nil {
		log.Fatalf("error parsing fg color string: %v", err)
	}

	bg, err := parseColorString(*bgArg)
	if err != nil {
		log.Fatalf("error parsing bg color string: %v", err)
	}

	return &options{
		inputPath:    *inputPathArg,
		outputPath:   *outputPathArg,
		fg:           fg,
		bg:           bg,
		inputIsDir:   inputIsDir,
	}
}

func isDir(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("failed to get info for '%s': %v", path, err)
	}

	return info.IsDir()
}

func ensureAllowedOutputFormat(format string) {
	if !isAllowedOuputFormat(format) {
		log.Fatalf(
			"format '%s' is not an allowed output format, allowed formats: %s",
			format,
			strings.Join(allowedOuputFormats, " | "),
		)
	}
}
