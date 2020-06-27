package pixl

import (
	"image"
	"math"
	"runtime"
	"sync"
)

type Resize struct {
	Width  uint
	Height uint
}

// A new image with the given dimensions will be returned.
// If one of the parameters width or height is set to 0, its size will be calculated so that
// the aspect ratio is that of the originating image.
// The resizing algorithm uses channels for parallel computation.
// If the input image has width or height of 0, it is returned unchanged.
func (config Resize) Convert(input image.Image) image.Image {

	scaleX, scaleY := aspectRatio(config.Width, config.Height, float64(input.Bounds().Dx()), float64(input.Bounds().Dy()))
	if config.Width == 0 {
		config.Width = uint(0.7 + float64(input.Bounds().Dx())/scaleX)
	}
	if config.Height == 0 {
		config.Height = uint(0.7 + float64(input.Bounds().Dy())/scaleY)
	}

	if int(config.Width) == input.Bounds().Dx() && int(config.Height) == input.Bounds().Dy() {
		return input
	}

	if input.Bounds().Dx() <= 0 || input.Bounds().Dy() <= 0 {
		return input
	}

	temp := image.NewRGBA64(image.Rect(0, 0, input.Bounds().Dy(), int(config.Width)))
	result := image.NewRGBA64(image.Rect(0, 0, int(config.Width), int(config.Height)))
	cpus := runtime.GOMAXPROCS(0)
	wg := sync.WaitGroup{}

	coeffs, offset, filterLength := createWeightsNearest(temp.Bounds().Dy(), scaleX)
	wg.Add(cpus)
	for i := 0; i < cpus; i++ {
		slice := makeSlice(temp, i, cpus).(*image.RGBA64)
		go func() {
			defer wg.Done()
			nearestGeneric(input, slice, scaleX, coeffs, offset, filterLength)
		}()
	}
	wg.Wait()

	coeffs, offset, filterLength = createWeightsNearest(result.Bounds().Dy(), scaleY)
	wg.Add(cpus)
	for i := 0; i < cpus; i++ {
		slice := makeSlice(result, i, cpus).(*image.RGBA64)
		go func() {
			defer wg.Done()
			nearestRGBA64(temp, slice, scaleY, coeffs, offset, filterLength)
		}()
	}
	wg.Wait()
	return result
}

func nearestRGBA64(in *image.RGBA64, out *image.RGBA64, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		row := in.Pix[x*in.Stride:]
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case uint(xi) < uint(maxX):
						xi *= 8
					case xi >= maxX:
						xi = 8 * maxX
					default:
						xi = 0
					}
					rgba[0] += float32(uint16(row[xi+0])<<8 | uint16(row[xi+1]))
					rgba[1] += float32(uint16(row[xi+2])<<8 | uint16(row[xi+3]))
					rgba[2] += float32(uint16(row[xi+4])<<8 | uint16(row[xi+5]))
					rgba[3] += float32(uint16(row[xi+6])<<8 | uint16(row[xi+7]))
					sum++
				}
			}

			xo := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := floatToUint16(rgba[0] / sum)
			out.Pix[xo+0] = uint8(value >> 8)
			out.Pix[xo+1] = uint8(value)
			value = floatToUint16(rgba[1] / sum)
			out.Pix[xo+2] = uint8(value >> 8)
			out.Pix[xo+3] = uint8(value)
			value = floatToUint16(rgba[2] / sum)
			out.Pix[xo+4] = uint8(value >> 8)
			out.Pix[xo+5] = uint8(value)
			value = floatToUint16(rgba[3] / sum)
			out.Pix[xo+6] = uint8(value >> 8)
			out.Pix[xo+7] = uint8(value)
		}
	}
}

func nearestGeneric(in image.Image, out *image.RGBA64, scale float64, coeffs []bool, offset []int, filterLength int) {
	newBounds := out.Bounds()
	maxX := in.Bounds().Dx() - 1

	for x := newBounds.Min.X; x < newBounds.Max.X; x++ {
		for y := newBounds.Min.Y; y < newBounds.Max.Y; y++ {
			var rgba [4]float32
			var sum float32
			start := offset[y]
			ci := y * filterLength
			for i := 0; i < filterLength; i++ {
				if coeffs[ci+i] {
					xi := start + i
					switch {
					case xi < 0:
						xi = 0
					case xi >= maxX:
						xi = maxX
					}
					r, g, b, a := in.At(xi+in.Bounds().Min.X, x+in.Bounds().Min.Y).RGBA()
					rgba[0] += float32(r)
					rgba[1] += float32(g)
					rgba[2] += float32(b)
					rgba[3] += float32(a)
					sum++
				}
			}

			offset := (y-newBounds.Min.Y)*out.Stride + (x-newBounds.Min.X)*8
			value := floatToUint16(rgba[0] / sum)
			out.Pix[offset+0] = uint8(value >> 8)
			out.Pix[offset+1] = uint8(value)
			value = floatToUint16(rgba[1] / sum)
			out.Pix[offset+2] = uint8(value >> 8)
			out.Pix[offset+3] = uint8(value)
			value = floatToUint16(rgba[2] / sum)
			out.Pix[offset+4] = uint8(value >> 8)
			out.Pix[offset+5] = uint8(value)
			value = floatToUint16(rgba[3] / sum)
			out.Pix[offset+6] = uint8(value >> 8)
			out.Pix[offset+7] = uint8(value)
		}
	}
}

func floatToUint16(x float32) uint16 {
	if x > 0xfffe {
		return 0xffff
	}
	return uint16(x)
}

func aspectRatio(width, height uint, oldWidth, oldHeight float64) (scaleX, scaleY float64) {
	if width == 0 {
		if height == 0 {
			scaleX = 1.0
			scaleY = 1.0
		} else {
			scaleY = oldHeight / float64(height)
			scaleX = scaleY
		}
	} else {
		scaleX = oldWidth / float64(width)
		if height == 0 {
			scaleY = scaleX
		} else {
			scaleY = oldHeight / float64(height)
		}
	}
	return
}

type imageWithSubImage interface {
	image.Image
	SubImage(image.Rectangle) image.Image
}

func makeSlice(img imageWithSubImage, i, n int) image.Image {
	return img.SubImage(image.Rect(img.Bounds().Min.X, img.Bounds().Min.Y+i*img.Bounds().Dy()/n, img.Bounds().Max.X, img.Bounds().Min.Y+(i+1)*img.Bounds().Dy()/n))
}

func createWeightsNearest(dy int, scale float64) ([]bool, []int, int) {
	filterLength := 2 * int(math.Max(math.Ceil(scale), 1))
	filterFactor := math.Min(1./(scale), 1)

	coeffs := make([]bool, dy*filterLength)
	start := make([]int, dy)
	for y := 0; y < dy; y++ {
		interpX := scale*(float64(y)+0.5) - 0.5
		start[y] = int(interpX) - filterLength/2 + 1
		interpX -= float64(start[y])
		for i := 0; i < filterLength; i++ {
			in := (interpX - float64(i)) * filterFactor
			if in >= -0.5 && in < 0.5 {
				coeffs[y*filterLength+i] = true
			} else {
				coeffs[y*filterLength+i] = false
			}
		}
	}

	return coeffs, start, filterLength
}
