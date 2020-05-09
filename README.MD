# Pixl 
Pixl is a tool for image processing.

Currently, it allows us to process images using the following algorithms:
- dithering
  - Floyd–Steinberg 
- halftone
- grayscale
  - average
  - luminosity
  - lightness
- normalize
- threshold
  - static
  - Otsu's Method


## Install

With a [correctly configured](https://golang.org/doc/install#testing) Go toolchain:

```sh
go get -u github.com/bednarc/pixl
```
## Examples

Here is an example of main.go file which gets filename.jpg as input, uses threshold function to process an image and save image as output.png

```go
package main

import (
	"image"
        _ "image/jpeg"
	"image/png"
	"os"
	
	"github.com/bednarc/pixl"
)

func main() {
   file, err := os.Open("filename.jpg")
   if err != nil {
	panic(err.Error())
   }
   defer file.Close()
   
   image, _, err := image.Decode(file)
   if err != nil {
   	panic(err.Error())
   }
	
   output := pixl.Threshold{Algorithm: pixl.ThresholdAlgorithms.Otsu}.Convert(image)
		
   newfile, err := os.Create("output.png")
   if err != nil {
	panic(err.Error())
   }
   defer newfile.Close()
   png.Encode(newfile, output)
}
```

### dithering

oryginal             |  dithering
:-------------------------:|:-------------------------:
<img src="https://pixxler.s3.eu-central-1.amazonaws.com/faceWoman.jpg" width="450">  |  <img src="https://pixxler.s3.eu-central-1.amazonaws.com/dithering.jpg" width="450">

```go
output := pixl.Dithering{}.Convert(input)
```

### halfltone

oryginal             |  halfltone
:-------------------------:|:-------------------------:
<img src="https://pixxler.s3.eu-central-1.amazonaws.com/eye.jpg" width="450">  |  <img src="https://pixxler.s3.eu-central-1.amazonaws.com/output_eye.jpg" width="450">

```go
output := pixl.Halftone{
	ColorBackground:       "#fffff0",
	ColorFront:            "#000000",
	ElementsHorizontaly:   100,
	OffsetSize:            20,
	Shift:                 50,
	MaxBoxSize:            10,
	TransparentBackground: false,
	Normalize:             true,
}.Convert(input)
```

### normalize

oryginal             |  normalize
:-------------------------:|:-------------------------:
<img src="https://pixxler.s3.eu-central-1.amazonaws.com/dark.jpg" width="450">  |  <img src="https://pixxler.s3.eu-central-1.amazonaws.com/dark_normalized.jpg" width="450">
```go
output := pixl.Normalize{}.Convert(input)
```

### threshold
oryginal             |  threshold
:-------------------------:|:-------------------------:
<img src="https://pixxler.s3.eu-central-1.amazonaws.com/face.jpg" width="450">  |  <img src="https://pixxler.s3.eu-central-1.amazonaws.com/face_otsu.jpg" width="450">
```go
output := pixl.Threshold{Algorithm: pixl.ThresholdAlgorithms.Otsu}.Convert(input)
```
### gray
oryginal             |  gray
:-------------------------:|:-------------------------:
<img src="https://pixxler.s3.eu-central-1.amazonaws.com/animal.jpg" width="450">  |  <img src="https://pixxler.s3.eu-central-1.amazonaws.com/gray_luminosity.jpg" width="450">
```go
output := pixl.Gray{Algorithm: pixl.GrayAlgorithms.Luminosity}.Convert(input)
```

## Contribute

If you want to contribute to a project and make it better, your help is very welcome. Contributing is also a great way to learn more about social coding on Github, new technologies, and their ecosystems.

There are some ideas for new features
- resize image 
- rotate image
- square halftone pattern
https://engineering.purdue.edu/~bouman/ece637/notes/pdf/Halftoning.pdf p.5
- new dithering algorithms
  - Jarvis, Judice, Ninke dithering
  - Stucki dithering
  - Burkes dithering
- median cut algorithm
- normalize colorful image

### benchmarks
```bash
go test -cpu 1 -bench=.
```

### tests 
```bash
go test 
```

## License

[![License](http://img.shields.io/:license-mit-blue.svg?style=flat-square)](http://badges.mit-license.org)

- **[MIT license](http://opensource.org/licenses/mit-license.php)**
