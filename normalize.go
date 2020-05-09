package pixl

import (
	"image"
	"image/color"
)

//Normalize is a config struct
type Normalize struct {
}

//Convert takes an image as an input and returns a normalized image
func (config Normalize) Convert(input image.Image) (output image.Image) {

	output = Gray{
		Algorithm: GrayAlgorithms.Luminosity,
	}.Convert(input)

	oldMax, oldMin := uint8(0), uint8(255)

	max := func(a, b uint8) uint8 {
		if a > b {
			return a
		}
		return b
	}

	min := func(a, b uint8) uint8 {
		if a < b {
			return a
		}
		return b
	}

	for x := 0; x < output.Bounds().Max.X; x++ {
		for y := 0; y < output.Bounds().Max.Y; y++ {
			r, _, _, _ := output.At(x, y).RGBA()
			oldMax = max(uint8(r>>8), oldMax)
			oldMin = min(uint8(r>>8), oldMin)
		}
	}

	traverseImage(output, output,
		normalizeParameters{
			newMax: 255,
			newMin: 0,
			oldMax: oldMax,
			oldMin: oldMin,
		})
	return
}

type normalizeParameters struct {
	newMax, newMin, oldMax, oldMin uint8
}

func (config normalizeParameters) transform(input color.Color) color.Color {
	r, _, _, _ := input.RGBA()
	rr := uint8(r)
	result := float64(rr-config.oldMin)*
		(float64(config.newMax-config.newMin)/float64(config.oldMax-config.oldMin)) +
		float64(config.newMin)

	return color.Gray{
		Y: uint8(result),
	}
}
