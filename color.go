package main

import (
	"fmt"
	"image/color"
	"regexp"
	"strconv"
	"strings"
)

type colorParserError struct {
	input    string
	position int
	length   int
	message  string
}

func (e colorParserError) Error() string {
	if e.position != -1 {
		pointerLen := len(e.input) + 15
		pointerBytes := make([]byte, pointerLen, pointerLen)
		for i := range len(e.input) + 15 {
			pointerBytes[i] = ' '
		}
		for i := range e.length {
			pointerBytes[e.position+i+15] = '^'
		}
		return fmt.Sprintf("invalid input '%s': %s\n%s", e.input, e.message, pointerBytes)
	}

	if len(e.input) > 0 {
		return fmt.Sprintf("invalid input '%s': %s", e.input, e.message)
	}

	return fmt.Sprintf("%s", e.message)
}

func parseColorString(colorString string) (c color.RGBA, err error) {
	panic("not implemented")
	return c, err
}

func parseHexColorString(colorString string) (c color.RGBA, err error) {
	if len(colorString) == 0 {
		return c, colorParserError{
			input:    colorString,
			position: -1,
			length:   0,
			message:  "empty color string",
		}
	}

	if colorString[0] != '#' {
		return c, colorParserError{
			input:    colorString,
			position: 0,
			length:   1,
			message:  "hex color string has to start with '#'",
		}
	}

	hexByteToDecimalByte := func(b byte) (byte, error) {
		switch {
		case b >= '0' && b <= '9':
			return b - '0', nil
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10, nil
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10, nil
		default:
			return 0, fmt.Errorf("invalid byte '%s' in hex string", string(b))
		}
	}

	newByteParsingError := func(position int, err error) error {
		return colorParserError{
			input:    colorString,
			position: position,
			length:   1,
			message:  err.Error(),
		}
	}

	switch len(colorString) {
	case 7:
		for i, color := range []*uint8{&c.R, &c.G, &c.B} {
			firstByte, err := hexByteToDecimalByte(colorString[2*i+1])
			if err != nil {
				return c, newByteParsingError(2*i+1, err)
			}
			secondByte, err := hexByteToDecimalByte(colorString[2*i+2])
			if err != nil {
				return c, newByteParsingError(2*i+2, err)
			}
			*color = firstByte<<4 + secondByte
		}
		c.A = 255
	case 4:
		for i, color := range []*uint8{&c.R, &c.G, &c.B} {
			firstByte, err := hexByteToDecimalByte(colorString[i+1])
			if err != nil {
				return c, newByteParsingError(i+1, err)
			}
			*color = firstByte * 17
		}
		c.A = 255
	default:
		return c, colorParserError{
			input:    colorString,
			position: 0,
			length:   len(colorString),
			message:  fmt.Sprintf("invalid length %d of hex string, has to be 4 or 7", len(colorString)),
		}
	}

	c.A = 255
	return c, nil
}

var rgbRegex = regexp.MustCompile("^rgb\\((.+)\\)$")

func parseRgbColorString(colorString string) (c color.RGBA, err error) {
	matches := rgbRegex.FindStringSubmatch(colorString)
	if len(matches) <= 1 {
		return c, colorParserError{
			input:    colorString,
			position: 0,
			length:   len(colorString),
			message:  "invalid rgb color string",
		}
	}

	values := strings.Split(matches[1], ",")
	if len(values) != 3 {
		return c, colorParserError{
			input:    colorString,
			position: 4,
			length:   len(matches[1]),
			message:  fmt.Sprintf("invalid number of values: found %d, want 3", len(values)),
		}
	}

	for i, color := range []*uint8{&c.R, &c.G, &c.B} {
		val, err := strconv.Atoi(strings.TrimSpace(values[i]))
		if err != nil || val > 255 {
			position := strings.Index(colorString, values[i])
			return c, colorParserError{
				input:    colorString,
				position: position,
				length:   len(values[i]),
				message:  fmt.Sprintf("invalid color value '%s', want 0-255", values[i]),
			}
		}
		*color = uint8(val)
	}

	c.A = 255

	return c, err
}
