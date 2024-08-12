package main

import (
	"fmt"
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
		{colorString: "#FF", wantError: true, color: color.RGBA{R: 0, G: 0, B: 0, A: 0}},
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
		fmt.Println(err)
	}
}
