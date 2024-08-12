package main

import (
	"errors"
	"fmt"
	"image/color"
	"strconv"
	"strings"
)

type colorParserError struct {
	input    string
	position int
	message  string
}

func (e colorParserError) Error() string {
	if e.position != -1 {
		pointerLen := len(e.input) + 15
		pointerBytes := make([]byte, pointerLen, pointerLen)
		for i := range len(e.input) + 15 {
			if i != e.position+15 {
				pointerBytes[i] = ' '
			} else {
				pointerBytes[i] = '^'
			}
		}
		return fmt.Sprintf("invalid input '%s': %s\n%s", e.input, e.message, pointerBytes)
	}

	if len(e.input) > 0 {
		return fmt.Sprintf("invalid input '%s': %s", e.input, e.message)
	}

	return fmt.Sprintf("%s", e.message)
}

func parseHexColorString(colorString string) (c color.RGBA, err error) {
	if len(colorString) == 0 {
		return c, colorParserError{
			input:    colorString,
			position: -1,
			message:  "empty color string",
		}
	}

	if colorString[0] != '#' {
		return c, colorParserError{
			input:    colorString,
			position: 0,
			message:  "hex color string has to start with '#'",
		}
	}

	hexToDecimal := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		default:
			return 0
		}
	}

	switch len(colorString) {
	case 7:
		c.R = hexToDecimal(colorString[1])<<4 + hexToDecimal(colorString[2])
		c.G = hexToDecimal(colorString[3])<<4 + hexToDecimal(colorString[4])
		c.B = hexToDecimal(colorString[5])<<4 + hexToDecimal(colorString[6])
		c.A = 255
	case 4:
		c.R = hexToDecimal(colorString[1]) * 17
		c.G = hexToDecimal(colorString[2]) * 17
		c.B = hexToDecimal(colorString[3]) * 17
		c.A = 255
	default:
		return c, colorParserError{
			input:    colorString,
			position: -1,
			message:  "hex color string length has to be 3 or 6",
		}
	}

	return c, nil
}

// func parseRgbColorString()

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
