package main

import (
	"fmt"
	"os"
)

func main() {
	options := getArgs()
	fontContext := loadFont()

	if options.isProcessDir() {
		processDirectory(options, fontContext)
	} else {
		outputPath, err := processImage(options.inputFilePath, options, fontContext)
		if err != nil {
			fmt.Printf("[ERROR] failed to process image '%s': %v", options.inputFilePath, err)
			os.Exit(1)
		} else {
			fmt.Printf("[INFO] image saved to '%s'", outputPath)
		}
	}
}
