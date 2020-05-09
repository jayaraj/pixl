package pixl

import (
	"image"
	"image/color"
	"image/draw"
	"math"
	"sync"
)

//Halftone is a config struct
//Configuration contains:
//  TransparentBackground - if true then background of png image will be transparent
//  ColorBackground - background color in hex format (e.g. #b690d9)
//  ColorBackground - frontend color in hex format (e.g. #b690d9)
// 	Shift - shift between adjacent rows
//  ElementsHorizontaly - amount of elements in the x axis
//  OffsetSize - increases or decreases output pattern size
//  MaxBoxSize - maximum size of output pattern
//  Normalize - if true then image will be normalized before conversion to halftone
type Halftone struct {
	TransparentBackground bool
	ColorBackground       string
	ColorFront            string
	Shift                 int8 /* -100(%) to 100(%) */
	ElementsHorizontaly   uint16
	OffsetSize            int8 /* -50(%) to 50(%) */
	MaxBoxSize            uint8
	Normalize             bool
}

//Convert takes an image as an input and returns halftone image
func (config Halftone) Convert(input image.Image) image.Image {
	trim(&config)
	boxAmountHorizont := int(config.ElementsHorizontaly)
	boxAmountVertical := input.Bounds().Max.Y * boxAmountHorizont /
		input.Bounds().Max.X
	outputBoxSize := int(config.MaxBoxSize)

	newWidth := boxAmountHorizont * outputBoxSize
	newHeight := boxAmountVertical * outputBoxSize

	output := image.NewNRGBA(image.Rect(0, 0, newWidth, newHeight))

	if !config.TransparentBackground {
		color, _ := parseHexColor(config.ColorBackground)
		traverseImage(output, output, paintAll{color: color})
	}

	var grayInput *image.Gray
	if config.Normalize {
		grayInput = Normalize{}.Convert(input).(*image.Gray)
	} else {
		grayInput = Gray{}.Convert(input)
	}

	colorFront, _ := parseHexColor(config.ColorFront)
	shift := int(config.Shift) * outputBoxSize / 100
	scale := float32(input.Bounds().Max.X) / float32(newWidth)

	var waitgroup sync.WaitGroup

	for j := 0; j < boxAmountVertical; j++ {
		for i := 0; i < boxAmountHorizont; i++ {
			waitgroup.Add(1)
			go func(ii, jj int) {
				defer waitgroup.Done()

				offset := (jj * shift) % outputBoxSize

				x := ii*outputBoxSize + offset
				y := jj * outputBoxSize
				_x := int(scale * float32(x))
				_y := int(scale * float32(y))
				orygBoxSize := scale * float32(outputBoxSize)
				blackIntensity := averageColor(_x, _y, int(orygBoxSize), grayInput)

				x0, y0, rMax, size := getCircleProperties(x, y, outputBoxSize, blackIntensity)

				drawCircle(output,
					x0,
					y0,
					rMax,
					size,
					config.OffsetSize,
					colorFront)
			}(i, j)
		}
	}
	waitgroup.Wait()

	return output
}
func getCircleProperties(x, y, outputBoxSize, blackIntensity int) (x0, y0, r int, size float32) {
	r = outputBoxSize / 2
	x0 = x + r
	y0 = y + r
	size = float32(blackIntensity) / 255.0
	return
}

func averageColor(x, y, side int, img image.Image) int {
	colorSum := 0
	for i := x; i < x+side; i++ {
		for j := y; j < y+side; j++ {
			r, _, _, _ := img.At(i, j).RGBA()
			colorSum += int(r >> 8)
		}
	}
	return int(0xFF - uint8(colorSum/(side*side)))
}

func trim(config *Halftone) {
	if config.OffsetSize < -50 {
		config.OffsetSize = -50
	}
	if config.OffsetSize > 50 {
		config.OffsetSize = 50
	}
	if config.Shift < -100 {
		config.Shift = -100
	}
	if config.Shift > 100 {
		config.Shift = 100
	}
}

func drawCircle(img draw.Image, x0 int, y0 int, rMax int, size float32, offset int8, color color.Color) {
	maxArea := 3.14 * float32(rMax*rMax)

	desireArea := maxArea * size
	rPower2 := (desireArea + (float32(offset)/100)*maxArea) / 3.14

	r := int(math.Sqrt(float64(rPower2)))

	for x := -r; x < r; x++ {
		height := math.Sqrt(float64(r*r - x*x))
		for y := int(-height); y < int(height); y++ {
			img.Set(x0+x, y0+y, color)
		}
	}

}
