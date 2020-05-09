package pixl

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"testing"
)

func TestNormalize(t *testing.T) {
	size := 2
	image := image.NewNRGBA(image.Rectangle{Max: image.Point{X: size, Y: size}})
	image.Set(0, 0, color.Gray{Y: 0})
	image.Set(0, 1, color.Gray{Y: 25})
	image.Set(1, 0, color.Gray{Y: 50})
	image.Set(1, 1, color.Gray{Y: 127})

	out := Normalize{}.Convert(image)

	c := out.At(0, 0)
	r, _, _, _ := c.RGBA()
	if expected := uint32(0); r>>8 != expected {
		t.Errorf("Invalid value of pixel, got: %d, want: %d.", r>>8, expected)
	}

	c = out.At(0, 1)
	r, _, _, _ = c.RGBA()
	if expected := uint32(50); r>>8 != expected {
		t.Errorf("Invalid value of pixel, got: %d, want: %d.", r>>8, expected)
	}

	c = out.At(1, 0)
	r, _, _, _ = c.RGBA()
	if expected := uint32(100); r>>8 != expected {
		t.Errorf("Invalid value of pixel, got: %d, want: %d.", r>>8, expected)
	}

	c = out.At(1, 1)
	r, _, _, _ = c.RGBA()
	if expected := uint32(255); r>>8 != expected {
		t.Errorf("Invalid value of pixel, got: %d, want: %d.", r>>8, expected)
	}
}

func BenchmarkNormalize(b *testing.B) {
	b.StopTimer()
	input := generateImage()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Normalize{}.Convert(input)
	}
}
