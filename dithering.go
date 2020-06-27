package pixl

import (
	"image"
	"image/color"
)

type ditheringAlgoType string

const (
	FloydSteinberg    ditheringAlgoType = "floydsteinberg"
	JarvisJudiceNinke ditheringAlgoType = "jarvisjudiceninke"
	Stucki            ditheringAlgoType = "stucki"
	Atkinson          ditheringAlgoType = "atkinson"
	Burkes            ditheringAlgoType = "burkes"
	Sierra            ditheringAlgoType = "sierra"
	TwoRowSierra      ditheringAlgoType = "tworowsierra"
	SierraLite        ditheringAlgoType = "sierralite"
)

//Dithering is a config struct
type Dithering struct {
	Algorithm ditheringAlgoType
}

type cell struct {
	x int
	y int
	m int16
}

type mask struct {
	divisor int16
	cells   []cell
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
	mask := dithering.getMask()

	for y := 0; y < h-1; y++ {
		for x := 1; x < w-1; x++ {
			oldpixel := matrix[x][y]
			newpixel := toBlackOrWhite(oldpixel, threshold)
			matrix[x][y] = newpixel
			quantError := oldpixel - newpixel
			for _, c := range mask.cells {
				if (x+c.x) < 0 || (y+c.y) < 0 || (x+c.x >= w) || (y+c.y >= h) {
					continue
				}
				matrix[x+c.x][y+c.y] = matrix[x+c.x][y+c.y] + (quantError*float32(c.m))/float32(mask.divisor)
			}
			// matrix[x+1][y] = matrix[x+1][y] + quantError*7/16
			// matrix[x-1][y+1] = matrix[x-1][y+1] + quantError*3/16
			// matrix[x][y+1] = matrix[x][y+1] + quantError*5/16
			// matrix[x+1][y+1] = matrix[x+1][y+1] + quantError*1/16
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

func (dithering Dithering) getMask() mask {
	switch dithering.Algorithm {

	case FloydSteinberg:
		{
			return mask{
				divisor: 16,
				cells: []cell{
					{1, 0, 7},
					{-1, 1, 3},
					{0, 1, 5},
					{1, 1, 1},
				},
			}
		}
	case JarvisJudiceNinke:
		{
			return mask{
				divisor: 48,
				cells: []cell{
					{1, 0, 7},
					{2, 0, 5},
					{-2, 1, 3},
					{-1, 1, 5},
					{0, 1, 7},
					{1, 1, 5},
					{2, 1, 3},
					{-2, 2, 1},
					{-1, 2, 3},
					{0, 2, 5},
					{1, 2, 3},
					{2, 2, 1},
				},
			}
		}
	case Stucki:
		{
			return mask{
				divisor: 42,
				cells: []cell{
					{1, 0, 8},
					{2, 0, 4},
					{-2, 1, 2},
					{-1, 1, 4},
					{0, 1, 8},
					{1, 1, 4},
					{2, 1, 2},
				},
			}
		}
	case Atkinson:
		{
			return mask{
				divisor: 8,
				cells: []cell{
					{1, 0, 1},
					{2, 0, 1},
					{-1, 1, 1},
					{0, 1, 1},
					{1, 1, 1},
					{0, 2, 1},
				},
			}
		}
	case Burkes:
		{
			return mask{
				divisor: 32,
				cells: []cell{
					{1, 0, 8},
					{2, 0, 4},
					{-2, 1, 2},
					{-1, 1, 4},
					{0, 1, 8},
					{1, 1, 4},
					{2, 1, 2},
				},
			}
		}
	case Sierra:
		{
			return mask{
				divisor: 32,
				cells: []cell{
					{1, 0, 5},
					{2, 0, 3},
					{-2, 1, 2},
					{-1, 1, 4},
					{0, 1, 5},
					{1, 1, 4},
					{2, 1, 2},
					{-1, 2, 2},
					{0, 2, 3},
					{1, 2, 2},
				},
			}
		}
	case TwoRowSierra:
		{
			return mask{
				divisor: 16,
				cells: []cell{
					{1, 0, 4},
					{2, 0, 3},
					{-2, 1, 2},
					{-1, 1, 2},
					{0, 1, 3},
					{1, 1, 2},
					{2, 1, 1},
				},
			}
		}
	case SierraLite:
		{
			return mask{
				divisor: 4,
				cells: []cell{
					{1, 0, 1},
					{-1, 1, 1},
					{0, 1, 1},
				},
			}
		}
	default:
		{
			//FloydSteinberg
			return mask{
				divisor: 16,
				cells: []cell{
					{1, 0, 7},
					{-1, 1, 3},
					{0, 1, 5},
					{1, 1, 1},
				},
			}
		}
	}
}

func toBlackOrWhite(in float32, threshold uint8) float32 {
	if in < float32(threshold) {
		return 0
	}
	return 255
}
