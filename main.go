package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"log"
	"math"
	"os"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const chunkSize = 10
const numberOfCharacters = 9
const fontHeight = 13

var fontColor = color.RGBA{255, 207, 117, 255}
var characters = [numberOfCharacters]rune{'.', ':', 'c', 'o', 'P', 'O', '@', '■'}

func main() {
	img := readImage()
	outImg := convertToAscii(img)
	writeImage(outImg)
}

func readImage() image.Image {
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <image path>", os.Args[0])
	}

	path := os.Args[1]

	file, err := os.Open(path)
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

func writeImage(img image.Image) {
	outFile, err := os.Create("output.jpg")
	if err != nil {
		log.Fatalf("Failed to create output file: %v", err)
	}
	defer outFile.Close()

	err = jpeg.Encode(outFile, img, nil)
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}

	log.Println("Image saved to output.jpg")
}

func convertToAscii(img image.Image) image.Image {
	numberOfRowChunks := int(math.Ceil(float64(img.Bounds().Max.Y) / float64(chunkSize)))
	numberOfColChunks := int(math.Ceil(float64(img.Bounds().Max.X) / float64(chunkSize)))

	outImg := image.NewRGBA(img.Bounds())

	d := &font.Drawer{
		Dst:  outImg,
		Src:  image.NewUniform(fontColor),
		Face: basicfont.Face7x13,
	}

	for rowChunk := 0; rowChunk < numberOfRowChunks; rowChunk++ {
		for colChunk := 0; colChunk < numberOfColChunks; colChunk++ {
			val := getGrayscaleValueFromChunk(img, rowChunk, colChunk)
			char := getCharFromGrayscaleValue(val)

			d.Dot = fixed.P(colChunk*chunkSize, rowChunk*chunkSize+fontHeight)
			d.DrawString(string(char))
		}
	}

	return outImg
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

func getCharFromGrayscaleValue(val uint8) rune {
	bucket := int(math.Floor((float64(val) / 256) * numberOfCharacters))
	return characters[bucket]
}
