package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang/freetype"
)

const (
	chunkSize          = 10
	numberOfCharacters = 8
	fontHeight         = 13
)

var (
	characters          = [numberOfCharacters]byte{'.', ':', 'c', 'o', 'P', 'O', '@', '$'}
	allowedOuputFormats = []string{".jpg", ".jpeg", ".png"}
)

func readImage(inputPath string) image.Image {
	file, err := os.Open(inputPath)
	if err != nil {
		log.Fatalf("failed to read image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("failed to decode image: %v", err)
	}

	return img
}

func writeImage(img image.Image, outputPath string) {
	fileExtension := filepath.Ext(outputPath)
	if !isAllowedOuputFormat(fileExtension) {
		log.Fatalf("Format %s is not an allowed output format", fileExtension)
		log.Fatalf("Allowed formats: %s", strings.Join(allowedOuputFormats, ", "))
	}

	outFile, err := os.Create(outputPath)
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	switch fileExtension {
	case ".jpg", ".jpeg":
		if err = jpeg.Encode(outFile, img, nil); err != nil {
			log.Fatalf("Failed to encode image: %v", err)
		}
	case ".png":
		if err = png.Encode(outFile, img); err != nil {
			log.Fatalf("Failed to encode image: %v", err)
		}
	default:
		panic("Unknown output format")
	}

	log.Println("Image saved to:", outputPath)
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

			pt := freetype.Pt(colChunk*chunkSize, rowChunk*chunkSize+fontHeight)
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

func isAllowedOuputFormat(format string) bool {
	for _, allowedFormat := range allowedOuputFormats {
		if format == allowedFormat {
			return true
		}
	}
	return false
}
