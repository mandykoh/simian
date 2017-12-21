package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/mandykoh/simian"
	"golang.org/x/image/draw"
)

const fingerprintDCTSideLength = 8
const samplesPerFingerprint = fingerprintDCTSideLength * fingerprintDCTSideLength

func makeFingerprint(src image.Image) []int16 {
	scaled := image.NewNRGBA(image.Rectangle{Max: image.Point{X: fingerprintDCTSideLength, Y: fingerprintDCTSideLength}})
	draw.BiLinear.Scale(scaled, scaled.Bounds(), src, src.Bounds(), draw.Src, nil)

	fingerprintSamples := make([]int8, samplesPerFingerprint)
	offset := 0

	for i := scaled.Bounds().Min.Y; i < scaled.Bounds().Max.Y; i++ {
		for j := scaled.Bounds().Min.X; j < scaled.Bounds().Max.X; j++ {
			r, g, b, _ := scaled.At(j, i).RGBA()
			y, _, _ := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))

			fingerprintSamples[offset] = int8(y - 128)
			offset++
		}
	}

	dct := simian.DCT(fingerprintDCTSideLength, fingerprintDCTSideLength, fingerprintSamples)

	fmt.Printf("DCT:\n")
	for i := 0; i < len(dct); i++ {
		if i == 0 {
			dct[i] >>= 7
		} else {
			dct[i] = dct[i] >> 5
		}

		if i > 0 && i%fingerprintDCTSideLength == 0 {
			fmt.Println()
		}
		fmt.Printf(" %5d", dct[i])
	}
	fmt.Println()
	fmt.Println()

	fingerprint := simian.FlattenRecursiveSquares(dct)

	return fingerprint
}

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Usage: simian-fingerprint <image.jpg>\n")
		return
	}

	imageFile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fingerprint := makeFingerprint(img)

	for i := 0; i < len(fingerprint); i++ {
		fmt.Printf("%02x", fingerprint[i]+128)
	}
	fmt.Println()
}
