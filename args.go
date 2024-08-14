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
	defaultBgColor = "rgb(0, 0, 0)"
	defaultFgColor = "rgb(255, 207, 117)"
)

type options struct {
	inputPath  string
	outputPath string
	fg         color.Color
	bg         color.Color
}

func getArgs() *options {
	var inputPathArg = flag.String("f", "", "the image input path")
	var outputPathArg = flag.String("o", "", "the image output path")
	var fgArg = flag.String("fg", defaultFgColor, "the ascii font color in HEX / RGB format")
	var bgArg = flag.String("bg", defaultBgColor, "the ascii font color in HEX / RGB format")

	flag.Parse()

	if len(*inputPathArg) == 0 || len(*outputPathArg) == 0 {
		flag.PrintDefaults()
		os.Exit(1)
	}

	outputFileExtension := filepath.Ext(*outputPathArg)
	if !isAllowedOuputFormat(outputFileExtension) {
		log.Fatalf(
			"format '%s' is not an allowed output format, allowed formats: %s",
			outputFileExtension,
			strings.Join(allowedOuputFormats, " | "),
		)
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
		inputPath:  *inputPathArg,
		outputPath: *outputPathArg,
		fg:         fg,
		bg:         bg,
	}
}
