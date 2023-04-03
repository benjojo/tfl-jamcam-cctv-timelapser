package main

import (
	"bytes"
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"testing"
)

func TestFFMPEGTimeLapse(t *testing.T) {
	// Make some jpegs
	jpegs := make([][]byte, 0)

	for i := 0; i < 200; i++ {
		myImage := image.NewRGBA(image.Rect(0, 0, 512, 512))

		for x := 0; x < 400; x++ {
			for y := 0; y < 400; y++ {
				myImage.Set(x, y, color.RGBA{
					R: uint8(x + i),
					G: uint8(x - y - 10),
					B: uint8(y),
					A: 255,
				})
			}
		}

		buf := bytes.NewBuffer(nil)
		jpeg.Encode(buf, myImage, nil)
		jpegs = append(jpegs, buf.Bytes())
	}

	video := createTimeLapseVideo(jpegs)
	if len(video) == 0 {
		t.FailNow()
	}

	os.WriteFile("test.mp4", video, 0644)
}
