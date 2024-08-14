package main

func main() {
	options := getArgs()
	img := readImage(options.inputPath)
	outImg := convertToAscii(img, options.fg)
	writeImage(outImg, options.outputPath)
}
