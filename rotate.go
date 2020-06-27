package pixl

import (
	"image"
	"image/color"
	"math"
	"runtime"
	"sync"
)

type Rotate struct {
	Angle   float64
	BGColor color.Color
}

func (config Rotate) Convert(input image.Image) image.Image {
	angle := config.Angle - math.Floor(config.Angle/360)*360

	inputWidth := input.Bounds().Max.X
	inputHeight := input.Bounds().Max.Y
	width, height := rotatedSize(inputWidth, inputHeight, angle)
	result := image.NewNRGBA(image.Rect(0, 0, width, height))

	if width <= 0 || height <= 0 {
		return result
	}
	src := toNRGBA(input)
	srcXOff := float64(inputWidth)/2 - 0.5
	srcYOff := float64(inputHeight)/2 - 0.5
	dstXOff := float64(width)/2 - 0.5
	dstYOff := float64(height)/2 - 0.5

	bgColorNRGBA := color.NRGBAModel.Convert(config.BGColor).(color.NRGBA)
	sin, cos := math.Sincos(math.Pi * angle / 180)

	parallel(0, height, func(ys <-chan int) {
		for dstY := range ys {
			for dstX := 0; dstX < width; dstX++ {
				xf, yf := rotatePoint(float64(dstX)-dstXOff, float64(dstY)-dstYOff, sin, cos)
				xf, yf = xf+srcXOff, yf+srcYOff
				interpolatePoint(result, dstX, dstY, src, xf, yf, bgColorNRGBA)
			}
		}
	})
	return result
}

func rotatePoint(x, y, sin, cos float64) (float64, float64) {
	return x*cos - y*sin, x*sin + y*cos
}

func rotatedSize(w, h int, angle float64) (int, int) {
	if w <= 0 || h <= 0 {
		return 0, 0
	}

	sin, cos := math.Sincos(math.Pi * angle / 180)
	x1, y1 := rotatePoint(float64(w-1), 0, sin, cos)
	x2, y2 := rotatePoint(float64(w-1), float64(h-1), sin, cos)
	x3, y3 := rotatePoint(0, float64(h-1), sin, cos)

	minx := math.Min(x1, math.Min(x2, math.Min(x3, 0)))
	maxx := math.Max(x1, math.Max(x2, math.Max(x3, 0)))
	miny := math.Min(y1, math.Min(y2, math.Min(y3, 0)))
	maxy := math.Max(y1, math.Max(y2, math.Max(y3, 0)))

	neww := maxx - minx + 1
	if neww-math.Floor(neww) > 0.1 {
		neww++
	}
	newh := maxy - miny + 1
	if newh-math.Floor(newh) > 0.1 {
		newh++
	}

	return int(neww), int(newh)
}

func interpolatePoint(dst *image.NRGBA, dstX, dstY int, src *image.NRGBA, xf, yf float64, bgColor color.NRGBA) {
	j := dstY*dst.Stride + dstX*4
	d := dst.Pix[j : j+4 : j+4]

	x0 := int(math.Floor(xf))
	y0 := int(math.Floor(yf))
	bounds := src.Bounds()
	if !image.Pt(x0, y0).In(image.Rect(bounds.Min.X-1, bounds.Min.Y-1, bounds.Max.X, bounds.Max.Y)) {
		d[0] = bgColor.R
		d[1] = bgColor.G
		d[2] = bgColor.B
		d[3] = bgColor.A
		return
	}

	xq := xf - float64(x0)
	yq := yf - float64(y0)
	points := [4]image.Point{
		{x0, y0},
		{x0 + 1, y0},
		{x0, y0 + 1},
		{x0 + 1, y0 + 1},
	}
	weights := [4]float64{
		(1 - xq) * (1 - yq),
		xq * (1 - yq),
		(1 - xq) * yq,
		xq * yq,
	}

	var r, g, b, a float64
	for i := 0; i < 4; i++ {
		p := points[i]
		w := weights[i]
		if p.In(bounds) {
			i := p.Y*src.Stride + p.X*4
			s := src.Pix[i : i+4 : i+4]
			wa := float64(s[3]) * w
			r += float64(s[0]) * wa
			g += float64(s[1]) * wa
			b += float64(s[2]) * wa
			a += wa
		} else {
			wa := float64(bgColor.A) * w
			r += float64(bgColor.R) * wa
			g += float64(bgColor.G) * wa
			b += float64(bgColor.B) * wa
			a += wa
		}
	}
	if a != 0 {
		aInv := 1 / a
		d[0] = clamp(r * aInv)
		d[1] = clamp(g * aInv)
		d[2] = clamp(b * aInv)
		d[3] = clamp(a)
	}
}

func clamp(x float64) uint8 {
	v := int64(x + 0.5)
	if v > 255 {
		return 255
	}
	if v > 0 {
		return uint8(v)
	}
	return 0
}

func toNRGBA(img image.Image) *image.NRGBA {
	if img, ok := img.(*image.NRGBA); ok {
		return &image.NRGBA{
			Pix:    img.Pix,
			Stride: img.Stride,
			Rect:   img.Rect.Sub(img.Rect.Min),
		}
	}
	return Clone(img)
}

func Clone(img image.Image) *image.NRGBA {
	src := newScanner(img)
	dst := image.NewNRGBA(image.Rect(0, 0, src.w, src.h))
	size := src.w * 4
	parallel(0, src.h, func(ys <-chan int) {
		for y := range ys {
			i := y * dst.Stride
			src.scan(0, y, src.w, y+1, dst.Pix[i:i+size])
		}
	})
	return dst
}

func parallel(start, stop int, fn func(<-chan int)) {
	count := stop - start
	if count < 1 {
		return
	}
	procs := runtime.GOMAXPROCS(0)

	c := make(chan int, count)
	for i := start; i < stop; i++ {
		c <- i
	}
	close(c)

	var wg sync.WaitGroup
	for i := 0; i < procs; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			fn(c)
		}()
	}
	wg.Wait()
}
