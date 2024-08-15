package main

func main() {
	options := getArgs()
	fontContext := loadFont()

	img := readImage(options.inputPath)
	outImg := convertToAscii(img, options, fontContext)
	writeImage(outImg, options.outputPath)
}
