package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
)

const (
	chunkSize          = 10
	numberOfCharacters = 8
	numberOfWorkers    = 5
)

var characters = [numberOfCharacters]byte{'.', ':', 'c', 'o', 'P', 'O', '@', '$'}

type processingMessage struct {
	inputPath  string
	outputPath string
	err        error
}

func startWorkers(taskChan <-chan string, outChan chan<- processingMessage, options *options, fontContext *freetype.Context) {
	for range numberOfWorkers {
		go func() {
			for imagePath := range taskChan {
				outputPath, err := processImage(imagePath, options, fontContext)

				outChan <- processingMessage{
					inputPath:  imagePath,
					outputPath: outputPath,
					err:        err,
				}
			}
		}()
	}
}

func processDirectory(options *options, fontContext *freetype.Context) error {
	entries, err := os.ReadDir(options.inputDirPath)
	if err != nil {
		return fmt.Errorf("failed to read input directory '%s': %v", options.inputDirPath, err)
	}

	var imagePaths []string
	for _, entry := range entries {
		if !entry.IsDir() {
			imagePaths = append(imagePaths, path.Join(options.inputDirPath, entry.Name()))
		}
	}

	taskChan := make(chan string)
	outChan := make(chan processingMessage)

	startWorkers(taskChan, outChan, options, fontContext)

	go func() {
		for _, imagePath := range imagePaths {
			taskChan <- imagePath
		}
		close(taskChan)
	}()

	for range imagePaths {
		result := <-outChan

		if result.err != nil {
			fmt.Printf("failed to process image '%s': %v\n", result.inputPath, result.err)
		} else {
			fmt.Printf("processed '%s' -> wrote new image to '%s'\n", result.inputPath, result.outputPath)
		}
	}

	return nil
}

func processImage(imagePath string, options *options, fontContext *freetype.Context) (string, error) {
	img, err := readImage(imagePath)
	if err != nil {
		return "", fmt.Errorf("failed to process image '%s': %v", options.inputFilePath, err)
	}

	outImg := convertToAscii(img, options, fontContext)

	var outputPath string
	if !options.isProcessDir() {
		outputPath = options.outputFilePath
	} else {
		fileNameElements := strings.Split(path.Base(imagePath), ".")
		outputPath = path.Join(options.outputDirPath, fileNameElements[0]+"_ascii."+fileNameElements[1])
	}

	err = writeImage(outImg, outputPath)
	if err != nil {
		return "", fmt.Errorf("failed to process image '%s': %v", options.inputFilePath, err)
	}

	return outputPath, nil
}

func readImage(inputPath string) (img image.Image, err error) {
	file, err := os.Open(inputPath)
	if err != nil {
		return img, fmt.Errorf("failed to read image: %v", err)
	}
	defer file.Close()

	img, _, err = image.Decode(file)
	if err != nil {
		return img, fmt.Errorf("failed to decode image: %v", err)
	}

	return img, nil
}

func writeImage(img image.Image, outputPath string) error {
	outFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer outFile.Close()

	newEncodingError := func(err error) error {
		return fmt.Errorf("failed to encode image: %v", err)
	}

	switch filepath.Ext(outputPath) {
	case ".jpg", ".jpeg":
		if err = jpeg.Encode(outFile, img, nil); err != nil {
			return newEncodingError(err)
		}
	case ".png":
		if err = png.Encode(outFile, img); err != nil {
			return newEncodingError(err)
		}
	default:
		panic("Unknown output format")
	}

	return nil
}

func convertToAscii(img image.Image, options *options, c *freetype.Context) image.Image {
	numberOfRowChunks, numberOfColChunks := getNumberOfChunksFromImage(img)

	outImg := image.NewRGBA(img.Bounds())
	fillImageBgColor(outImg, options.bg)

	c.SetClip(outImg.Bounds())
	c.SetDst(outImg)
	c.SetSrc(image.NewUniform(options.fg))

	for rowChunk := 0; rowChunk < numberOfRowChunks; rowChunk++ {
		for colChunk := 0; colChunk < numberOfColChunks; colChunk++ {
			val := getGrayscaleValueFromChunk(img, rowChunk, colChunk)
			char := getCharFromGrayscaleValue(val)

			pt := freetype.Pt(colChunk*chunkSize, rowChunk*chunkSize)
			c.DrawString(string(char), pt)
		}
	}

	return outImg
}

func fillImageBgColor(img *image.RGBA, bg color.Color) {
	draw.Draw(img, img.Bounds(), &image.Uniform{bg}, image.Point{}, draw.Src)
}

func getNumberOfChunksFromImage(img image.Image) (int, int) {
	numberOfRowChunks := int(math.Ceil(float64(img.Bounds().Max.Y) / float64(chunkSize)))
	numberOfColChunks := int(math.Ceil(float64(img.Bounds().Max.X) / float64(chunkSize)))
	return numberOfRowChunks, numberOfColChunks
}

func getGrayscaleValueFromChunk(img image.Image, rowChunk, colChunk int) uint8 {
	var total float64
	var count int

	for y := 0; y < 10; y++ {
		for x := 0; x < 10; x++ {
			row := y + rowChunk*chunkSize
			col := x + colChunk*chunkSize

			if row >= img.Bounds().Max.Y || col >= img.Bounds().Max.X {
				continue
			}

			r, g, b, _ := img.At(col, row).RGBA()
			total += float64(r)*0.299 + float64(g)*0.587 + float64(b)*0.114
			count++
		}
	}

	return uint8(total / 256 / float64(count))
}

func getCharFromGrayscaleValue(val uint8) byte {
	bucket := int(math.Floor((float64(val) / 256) * numberOfCharacters))
	return characters[bucket]
}
