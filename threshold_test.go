package pixl

import (
	"image"
	"image/color"
	_ "image/jpeg"
	"testing"
)

func TestThresholdStatic(t *testing.T) {
	size := 10
	image := image.NewNRGBA(image.Rectangle{Max: image.Point{X: size, Y: size}})
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if j < 5 {
				image.Set(i, j, color.Gray{Y: 200})
			} else if j < 9 {
				image.Set(i, j, color.Gray{Y: 50})
			} else {
				image.Set(i, j, color.Gray{Y: 160})
			}
		}
	}

	out := Threshold{Algorithm: ThresholdAlgorithms.Static, StaticLevel: 155}.Convert(image)
	blackPixels := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if r, _, _, _ := out.At(i, j).RGBA(); (r >> 8) == 0xFF {
				blackPixels++
			}
		}
	}

	if expected := 60; blackPixels != expected {
		t.Errorf("Invalid number of black pixels in image, got: %d, want: %d.", blackPixels, expected)
	}
}

func TestThresholdStaticInversion(t *testing.T) {
	size := 16
	image := image.NewNRGBA(image.Rectangle{Max: image.Point{X: size, Y: size}})
	for i := 0; i < size; i++ {
		counter := uint8(0)
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				image.Set(i, j, color.Gray{Y: counter})
				counter++
			}
		}
	}

	out := Threshold{Algorithm: ThresholdAlgorithms.Static, StaticLevel: 99, InvertColors: true}.Convert(image)
	blackPixels := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if r, _, _, _ := out.At(i, j).RGBA(); (r >> 8) == 0xFF {
				blackPixels++
			}
		}
	}

	if expected := 100; blackPixels != expected {
		t.Errorf("Invalid number of black pixels in image, got: %d, want: %d.", blackPixels, expected)
	}
}

func TestThresholdStaticDefaultLevel(t *testing.T) {
	size := 16
	image := image.NewNRGBA(image.Rectangle{Max: image.Point{X: size, Y: size}})
	counter := uint8(0)
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			image.Set(i, j, color.Gray{Y: counter})
			counter++
		}
	}

	out := Threshold{Algorithm: ThresholdAlgorithms.Static, InvertColors: false}.Convert(image)
	blackPixels := 0
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if r, _, _, _ := out.At(i, j).RGBA(); (r >> 8) == 0xFF {
				blackPixels++
			}
		}
	}

	if expected := 128; blackPixels != expected {
		t.Errorf("Invalid number of black pixels in image, got: %d, want: %d.", blackPixels, expected)
	}
}

func BenchmarkThresholdStatic(b *testing.B) {
	b.StopTimer()
	input := generateImage()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Threshold{
			Algorithm: ThresholdAlgorithms.Static,
		}.Convert(input)
	}
}

func BenchmarkThresholdOtsu(b *testing.B) {
	b.StopTimer()
	input := generateImage()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Threshold{
			Algorithm: ThresholdAlgorithms.Otsu,
		}.Convert(input)
	}
}

func BenchmarkThresholdStaticWithInvert(b *testing.B) {
	b.StopTimer()
	input := generateImage()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Threshold{
			Algorithm: ThresholdAlgorithms.Static, InvertColors: true,
		}.Convert(input)
	}
}
