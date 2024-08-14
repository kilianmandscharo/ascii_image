package main

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseHexColorString(t *testing.T) {
	tests := []struct {
		colorString string
		wantError   bool
		color       color.RGBA
	}{
		{colorString: "", wantError: true, color: color.RGBA{R: 0, G: 0, B: 0, A: 0}},
		{colorString: "FFFFFF", wantError: true, color: color.RGBA{R: 0, G: 0, B: 0, A: 0}},
		{colorString: "#FF", wantError: true, color: color.RGBA{R: 0, G: 0, B: 0, A: 0}},
		{colorString: "#F", wantError: true, color: color.RGBA{R: 0, G: 0, B: 0, A: 0}},
		{colorString: "#xFFFFF", wantError: true, color: color.RGBA{R: 0, G: 0, B: 0, A: 0}},
		{colorString: "#FxFFFF", wantError: true, color: color.RGBA{R: 0, G: 0, B: 0, A: 0}},
		{colorString: "#FFxFFF", wantError: true, color: color.RGBA{R: 255, G: 0, B: 0, A: 0}},
		{colorString: "#FFFxFF", wantError: true, color: color.RGBA{R: 255, G: 0, B: 0, A: 0}},
		{colorString: "#FFFFxF", wantError: true, color: color.RGBA{R: 255, G: 255, B: 0, A: 0}},
		{colorString: "#FFFFFx", wantError: true, color: color.RGBA{R: 255, G: 255, B: 0, A: 0}},
		{colorString: "#FFFFFF", wantError: false, color: color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{colorString: "#ffffff", wantError: false, color: color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{colorString: "#FFFfff", wantError: false, color: color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{colorString: "#FFF", wantError: false, color: color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{colorString: "#fff", wantError: false, color: color.RGBA{R: 255, G: 255, B: 255, A: 255}},
		{colorString: "#F00", wantError: false, color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		{colorString: "#f00", wantError: false, color: color.RGBA{R: 255, G: 0, B: 0, A: 255}},
		{colorString: "#AA4A44", wantError: false, color: color.RGBA{R: 170, G: 74, B: 68, A: 255}},
		{colorString: "#aa4a44", wantError: false, color: color.RGBA{R: 170, G: 74, B: 68, A: 255}},
	}

	for _, tt := range tests {
		c, err := parseHexColorString(tt.colorString)
		assert.Equal(t, tt.wantError, err != nil)
		assert.Equal(t, tt.color, c)
	}
}

func TestParseRgbColorString(t *testing.T) {
	tests := []struct {
		colorString string
		wantError   bool
		color       color.RGBA
	}{
		{colorString: "(0, 0, 0)", wantError: true},
		{colorString: "rgb(0, 0, 0, 0)", wantError: true},
		{colorString: "rgb(0, xx, 0)", wantError: true},
		{colorString: "rgb(0, 12, 400)", wantError: true, color: color.RGBA{R: 0, G: 12, B: 0, A: 0}},
		{colorString: "rgb(0, 0, 0)", wantError: false, color: color.RGBA{R: 0, G: 0, B: 0, A: 255}},
		{colorString: "rgb(0, 0, 0)", wantError: false, color: color.RGBA{R: 0, G: 0, B: 0, A: 255}},
		{colorString: "rgb(15, 246, 233)", wantError: false, color: color.RGBA{R: 15, G: 246, B: 233, A: 255}},
	}

	for _, tt := range tests {
		c, err := parseRgbColorString(tt.colorString)
		assert.Equal(t, tt.wantError, err != nil)
		assert.Equal(t, tt.color, c)
	}
}
