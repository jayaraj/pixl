package pixl

import (
	"image"
	"image/color"
)

type thresholdAlgoName string

type thresholdAlgoList struct {
	Static thresholdAlgoName
	Otsu   thresholdAlgoName
}

// ThresholdAlgorithms consists of a list of algorithms that can be used as
// algorithm type in pixl.Threshold struct. for ex:
// pixl.Threshold{Algorithm: pixl.ThresholdAlgorithms.Static}
var ThresholdAlgorithms = &thresholdAlgoList{
	Static: "static",
	Otsu:   "otsu",
}

//Threshold is a config struct
//Configuration contains:
//  Algorithm - grayscale Algorithm used to convert image
//  StaticLevel - threshold level is used only with Static Algorithm type
//  InvertColors - if true then change all white pixel with black pixels
type Threshold struct {
	Algorithm    thresholdAlgoName
	StaticLevel  uint8
	InvertColors bool
}

//Convert takes an image as an input and returns thresholded image
func (config Threshold) Convert(img image.Image) *image.Gray {
	gray := Gray{
		Algorithm: GrayAlgorithms.Luminosity,
	}
	out := gray.Convert(img)

	if config.Algorithm == "static" {
		level := uint8(127)
		if config.StaticLevel != 0 {
			level = config.StaticLevel
		}
		traverseImage(out, out, threshold{level: level, invertColors: config.InvertColors})
		return out
	}
	//else if config.Algorithm == "otsu" {
	traverseImage(out, out, threshold{level: calculateThreshold(out), invertColors: config.InvertColors})
	return out
}

type threshold struct {
	level        uint8
	invertColors bool
}

func (config threshold) transform(input color.Color) color.Color {
	r, _, _, _ := input.RGBA()

	var result uint8

	if uint8(r>>8) <= config.level {
		result = 0x00
	} else {
		result = 0xFF
	}

	if config.invertColors {
		result = 255 - result
	}

	return color.Gray{
		Y: result,
	}
}

func calculateThreshold(img *image.Gray) uint8 {
	hist := histogramGray(img)

	pixelAmount := img.Bounds().Max.X * img.Bounds().Max.Y
	sum := 0

	for t := 0; t < 256; t++ {
		sum += t * hist[t]
	}

	sumB := 0.0
	weightB := 0
	weightF := 0

	varMax := 0.0
	threshold := 0
	for t := 0; t < 256; t++ {
		weightB += hist[t]
		if weightB == 0 {
			continue
		}
		weightF = pixelAmount - weightB
		if weightF == 0 {
			break
		}
		sumB += (float64)(t * hist[t])

		meanB := sumB / float64(weightB)
		meanF := (float64(sum) - sumB) / float64(weightF)

		varBetween := float64(weightB) * float64(weightF) * (meanB - meanF) * (meanB - meanF)

		if varBetween > varMax {
			varMax = varBetween
			threshold = t
		}
	}
	return uint8(threshold)
}
