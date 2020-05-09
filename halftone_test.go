package pixl

import (
	"image"
	"image/color"
	"testing"
)

func Test_averageSquare(t *testing.T) {

	width := 20
	height := 20
	input := image.NewNRGBA(image.Rect(0, 0, width, height))

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			input.Set(i, j, color.Gray{Y: uint8(i)})
		}
	}

	average := averageColor(0, 0, 7, input)

	if average != 0xFF-3 {
		t.Errorf("Average colors of square is incorrect, got: %d, want: %d.", average, 0xFF-3)
	}
}

func TestHalftoneConvert(t *testing.T) {

	width := 10
	height := 10
	input := image.NewNRGBA(image.Rect(0, 0, width, height))

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			input.Set(i, j, color.Gray{Y: 0x0F})
		}
	}

	halftoneNormalize := func(normalize bool) image.Image {
		out := Halftone{
			ColorBackground:       "#fffff0",
			ColorFront:            "#000000",
			ElementsHorizontaly:   1,
			OffsetSize:            0,
			Shift:                 50,
			MaxBoxSize:            10,
			TransparentBackground: false,
			Normalize:             normalize,
		}.Convert(input)

		return out
	}
	out := halftoneNormalize(true)

	blackPixelsCounter := func() int {
		blackPixels := 0
		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				if r, _, _, _ := out.At(i, j).RGBA(); (r >> 8) == 0xFF {
					blackPixels++
				}
			}
		}
		return blackPixels
	}

	blackPixels := blackPixelsCounter()

	if expected := 30; blackPixels != expected {
		t.Errorf("Invalid pixels amount in circle, got: %d, want: %d.", blackPixels, expected)
	}

	out = halftoneNormalize(false)
	blackPixels = blackPixelsCounter()

	if expected := 60; blackPixels != expected {
		t.Errorf("Invalid pixels amount in circle, got: %d, want: %d.", blackPixels, expected)
	}

}
