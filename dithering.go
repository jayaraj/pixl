package pixl

import (
	"image"
	"image/color"
)

//Dithering is a config struct
type Dithering struct {
}

//Convert takes an image as an input and returns dithered image
func (dithering Dithering) Convert(input image.Image) image.Image {
	grayInput := Gray{
		Algorithm: GrayAlgorithms.Luminosity,
	}.Convert(input)
	bounds := grayInput.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	matrix := make([][]float32, w)
	for x := 0; x < w; x++ {
		matrix[x] = make([]float32, h)
		for y := 0; y < h; y++ {
			r, _, _, _ := grayInput.At(x, y).RGBA()
			matrix[x][y] = float32(r >> 8)
		}
	}
	threshold := calculateThreshold(grayInput)

	for y := 0; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			oldpixel := matrix[x][y]
			newpixel := toBlackOrWhite(oldpixel, threshold)
			matrix[x][y] = newpixel
			quantError := oldpixel - newpixel
			matrix[x+1][y] = matrix[x+1][y] + quantError*7/16
			matrix[x-1][y+1] = matrix[x-1][y+1] + quantError*3/16
			matrix[x][y+1] = matrix[x][y+1] + quantError*5/16
			matrix[x+1][y+1] = matrix[x+1][y+1] + quantError*1/16
		}
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			col := color.Gray{Y: uint8(matrix[x][y])}
			grayInput.Set(x, y, col)
		}
	}

	return grayInput
}

func toBlackOrWhite(in float32, threshold uint8) float32 {
	if in < float32(threshold) {
		return 0
	}
	return 255
}
