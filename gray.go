package pixl

import (
	"image"
	"image/color"
)

type grayAlgoName string

type grayAlgoList struct {
	Lightness  grayAlgoName
	Average    grayAlgoName
	Luminosity grayAlgoName
}

// GrayAlgorithms consists of a list of algorithms that can be used as
// algorithm type in pixl.Gray struct. e.g.
// pixl.Gray{Algorithm: pixl.GrayAlgorithms.Lithtness}
var GrayAlgorithms = &grayAlgoList{
	Lightness:  "lightness",
	Average:    "average",
	Luminosity: "luminosity",
}

//Gray is a config struct
type Gray struct {
	Algorithm grayAlgoName
}

//Convert takes an image as an input and returns grayscale of the image
func (config Gray) Convert(input image.Image) *image.Gray {
	output := image.NewGray(input.Bounds())

	if config.Algorithm == "lightness" {
		traverseImage(input, output, grayLightness{})
	} else if config.Algorithm == "average" {
		traverseImage(input, output, grayAverage{})
	} else { //if conf.Algorithm == "luminosity"
		traverseImage(input, output, grayLuminosity{})
	}
	return output
}

type grayLightness struct{}

func (config grayLightness) transform(input color.Color) color.Color {
	r, g, b, _ := input.RGBA()

	max := func(a, b uint32) uint32 {
		if a > b {
			return a
		}
		return b
	}

	min := func(a, b uint32) uint32 {
		if a < b {
			return a
		}
		return b
	}

	_max := max(r, max(g, b))
	_min := min(r, min(g, b))
	result := (_max + _min) / 2
	return color.Gray{
		Y: uint8(result >> 8),
	}
}

type grayAverage struct{}

func (config grayAverage) transform(input color.Color) color.Color {
	r, g, b, _ := input.RGBA()
	result := (r + g + b) / 3

	return color.Gray{
		Y: uint8(result >> 8),
	}
}

type grayLuminosity struct{}

func (config grayLuminosity) transform(input color.Color) color.Color {
	r, g, b, _ := input.RGBA()
	result := uint32(0.21*float32(r) + 0.72*float32(g) + 0.07*float32(b))

	return color.Gray{
		Y: uint8(result >> 8),
	}
}
