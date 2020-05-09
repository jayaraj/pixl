package pixl

import (
	"image"
	"image/color"
	"testing"
)

func Test_paintAll(t *testing.T) {
	size := 10
	out := image.NewGray(image.Rect(0, 0, size, size))

	traverseImage(out, out, paintAll{color: color.Gray{Y: 0x0F}})

	result := true
	for x := 0; x < size; x++ {
		for y := 0; y < size; y++ {
			if (out.At(x, y) != color.Gray{Y: 0x0F}) {
				result = false
			}
		}
	}

	if result == false {
		t.Errorf("Not a whole picture has been painted.")
	}

}

func Test_histogramGray(t *testing.T) {
	grayInput := image.NewGray(image.Rect(0, 0, 10, 10))
	for i := 0; i < 10; i++ {
		for j := 0; j < 10; j++ {
			if i >= 3 && i <= 6 && j >= 2 && j <= 3 {
				grayInput.Set(i, j, color.Black)
			} else {
				grayInput.Set(i, j, color.White)
			}
		}
	}

	h := histogramGray(grayInput)

	if h[0] != 8 {
		t.Errorf("Amount of black pixels is incorect, got: %d, want: %d.", h[0], 8)
	}

	if h[0xFF] != 92 {
		t.Errorf("Amount of white pixels is incorect, got: %d, want: %d.", h[0xFF], 92)
	}
}

func Test_parseHexColor(t *testing.T) {
	_, err := parseHexColor("err")
	if err == nil {
		t.Errorf("Err is nil. Should be: parseHexColor: invalid format")
	}

	err = nil
	_, err = parseHexColor("#ff00Fx")
	if err == nil {
		t.Errorf("Err is nil. Should be: parseHexColor: invalid format")
	}

	err = nil
	_, err = parseHexColor("#ff00F")
	if err == nil {
		t.Errorf("Err is nil. Should be: parseHexColor: invalid format")
	}

	c, _ := parseHexColor("#Ff00f0")

	if c.R != 0xFF {
		t.Errorf("R value of pixel is incorect, got: %d, want: %d.", c.R, 0xFF)
	}
	if c.G != 0x00 {
		t.Errorf("G value of pixel is incorect, got: %d, want: %d.", c.G, 0x00)
	}
	if c.B != 0xF0 {
		t.Errorf("B value of pixel is incorect, got: %d, want: %d.", c.B, 0xF0)
	}

}
