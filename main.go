package main

func main() {
	options := getArgs()
	fontContext := loadFont()

	if options.inputIsDir {
		processDirectory(options, fontContext)
	} else {
		processImagesConcurrently(options, fontContext, options.inputPath)
	}
}
