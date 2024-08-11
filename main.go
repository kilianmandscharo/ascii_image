package main

import (
	"cmp"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
)

const chunkSize = 10
const numberOfCharacters = 9
const fontHeight = 13

var defaultFontColor = "255, 207, 117, 255"
var characters = [numberOfCharacters]rune{'.', ':', 'c', 'o', 'P', 'O', '@', '■'}
var allowedOuputFormats = []string{".jpg", ".jpeg", ".png"}

func main() {
	options := getArgs()
	img := readImage(options.inputPath)
	outImg := convertToAscii(img, options.color)
	writeImage(outImg, options.outputPath)
}

type options struct {
	inputPath  string
	outputPath string
	color      color.Color
}

func getArgs() *options {
	var inputPathArg = flag.String("f", "", "the image input path")
	var outputPathArg = flag.String("o", "", "the image output path")
	var colorArg = flag.String("c", defaultFontColor, "the ascii font color in RGB / RGBA format")

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

	color, err := parseRgbaColorString(cmp.Or(*colorArg, defaultFontColor))
	if err != nil {
		log.Fatalf("error parsing color string: %v", err)
	}

	return &options{
		inputPath:  *inputPathArg,
		outputPath: *outputPathArg,
		color:      color,
	}
}

func parseRgbaColorString(colorString string) (color.Color, error) {
	color := color.RGBA{}
	values := strings.Split(strings.Trim(colorString, "()"), ",")

	if len(values) < 3 {
		return color, errors.New("not enough color values")
	}

	r, err := strconv.Atoi(strings.TrimSpace(values[0]))
	if err != nil {
		return color, fmt.Errorf("invalid value '%s' for red: %v", values[0], err)
	}

	g, err := strconv.Atoi(strings.TrimSpace(values[1]))
	if err != nil {
		return color, fmt.Errorf("invalid value '%s' for green: %v", values[1], err)
	}

	b, err := strconv.Atoi(strings.TrimSpace(values[2]))
	if err != nil {
		return color, fmt.Errorf("invalid value '%s' for blue: %v", values[2], err)
	}

	var a int
	if len(values) > 3 {
		a, err = strconv.Atoi(strings.TrimSpace(values[3]))
		if err != nil {
			return color, fmt.Errorf("invalid value '%s' for alpha: %v", values[3], err)
		}
	} else {
		a = 255
	}

	color.R = uint8(r)
	color.G = uint8(g)
	color.B = uint8(b)
	color.A = uint8(a)

	return color, nil
}

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
	case ".jpg":
	case ".jpeg":
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

func convertToAscii(img image.Image, color color.Color) image.Image {
	numberOfRowChunks := int(math.Ceil(float64(img.Bounds().Max.Y) / float64(chunkSize)))
	numberOfColChunks := int(math.Ceil(float64(img.Bounds().Max.X) / float64(chunkSize)))

	outImg := image.NewRGBA(img.Bounds())

	d := &font.Drawer{
		Dst:  outImg,
		Src:  image.NewUniform(color),
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

func isAllowedOuputFormat(format string) bool {
	for _, allowedFormat := range allowedOuputFormats {
		if format == allowedFormat {
			return true
		}
	}
	return false
}
