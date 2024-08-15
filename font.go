package main

import (
	"log"
	"os"

	"github.com/golang/freetype"
)

const fontPath = "fonts/OpenSans-VariableFont_wdth,wght.ttf"

func loadFont() *freetype.Context {
	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		log.Fatalf("failed to read font %s: %v", fontPath, err)
	}

	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Fatalf("failed to parse font: %v", err)
	}

	c := freetype.NewContext()
	c.SetDPI(72)
	c.SetFontSize(12)
	c.SetFont(font)

	return c
}
