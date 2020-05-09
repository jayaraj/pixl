package pixl

import (
	"image"
	_ "image/jpeg"
	"testing"
)

func TestGrayLightness(t *testing.T) {

	image := image.NewNRGBA(image.Rectangle{Max: image.Point{X: 1, Y: 1}})
	inputColor, _ := parseHexColor("#FF000F")
	image.Set(0, 0, inputColor)
	out := Gray{Algorithm: GrayAlgorithms.Lightness}.Convert(image)

	r, _, _, _ := out.At(0, 0).RGBA()
	r = r >> 8

	if expected := uint32(127); r != expected {
		t.Errorf("Invalid output color, got: %d, want: %d.", r, expected)
	}
}

func TestGrayAverage(t *testing.T) {
	image := image.NewNRGBA(image.Rectangle{Max: image.Point{X: 1, Y: 1}})
	inputColor, _ := parseHexColor("#FF000F")
	image.Set(0, 0, inputColor)
	out := Gray{Algorithm: GrayAlgorithms.Average}.Convert(image)

	r, _, _, _ := out.At(0, 0).RGBA()
	r = r >> 8

	if expected := uint32(90); r != expected {
		t.Errorf("Invalid output color, got: %d, want: %d.", r, expected)
	}
}

func TestGrayLuminosity(t *testing.T) {
	image := image.NewNRGBA(image.Rectangle{Max: image.Point{X: 1, Y: 1}})
	inputColor, _ := parseHexColor("#7F7F7F")
	image.Set(0, 0, inputColor)
	out := Gray{Algorithm: GrayAlgorithms.Luminosity}.Convert(image)

	r, _, _, _ := out.At(0, 0).RGBA()
	r = r >> 8

	if expected := uint32(0.21*127 + 0.72*127 + 0.07*127); r != expected {
		t.Errorf("Invalid output color, got: %d, want: %d.", r, expected)
	}
}

func generateImage() image.Image {
	size := 10
	image := image.NewNRGBA(image.Rectangle{Max: image.Point{X: size, Y: size}})

	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			color, _ := parseHexColor("#0F011F")
			image.Set(i, j, color)
		}
	}

	return image
}

func BenchmarkGrayAverage(b *testing.B) {
	b.StopTimer()
	input := generateImage()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Gray{
			Algorithm: GrayAlgorithms.Average,
		}.Convert(input)
	}
}

func BenchmarkGrayLuminosity(b *testing.B) {
	b.StopTimer()
	input := generateImage()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Gray{
			Algorithm: GrayAlgorithms.Luminosity,
		}.Convert(input)
	}
}

func BenchmarkGrayLightness(b *testing.B) {
	b.StopTimer()
	input := generateImage()
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		Gray{
			Algorithm: GrayAlgorithms.Lightness,
		}.Convert(input)
	}
}
