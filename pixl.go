package pixl

import (
	"errors"
	"image"
	"image/color"
	"image/draw"
	"sync"
)

type paintAll struct {
	color color.Color
}

func (config paintAll) transform(input color.Color) color.Color {
	return config.color
}

func parseHexColor(s string) (c color.RGBA, err error) {
	c.A = 0xff

	errMsg := errors.New("parseHexColor: invalid format")
	if len(s) == 0 || s[0] != '#' {
		return c, errMsg
	}

	hexToByte := func(b byte) byte {
		switch {
		case b >= '0' && b <= '9':
			return b - '0'
		case b >= 'a' && b <= 'f':
			return b - 'a' + 10
		case b >= 'A' && b <= 'F':
			return b - 'A' + 10
		}
		err = errMsg
		return 0
	}

	switch len(s) {
	case 7:
		c.R = hexToByte(s[1])<<4 + hexToByte(s[2])
		c.G = hexToByte(s[3])<<4 + hexToByte(s[4])
		c.B = hexToByte(s[5])<<4 + hexToByte(s[6])
	default:
		err = errMsg
	}
	return
}

func histogramGray(image *image.Gray) map[int]int {
	histogram := make(map[int]int)
	bounds := image.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			r, _, _, _ := image.At(x, y).RGBA()
			histogram[int(r>>8)]++
		}
	}
	return histogram
}

type transformer interface {
	transform(color.Color) color.Color
}

func traverseImage(in image.Image, out image.Image, t transformer) {
	var waitgroup sync.WaitGroup
	if output, ok := out.(draw.Image); ok {
		w, h := output.Bounds().Max.X, output.Bounds().Max.Y
		for x := 0; x < w; x++ {
			waitgroup.Add(1)
			go func(_x int) {
				for y := 0; y < h; y++ {
					color := t.transform(in.At(_x, y))
					output.Set(_x, y, color)
				}
				waitgroup.Done()
			}(x)
		}
	}
	waitgroup.Wait()
}
