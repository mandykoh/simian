package main

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"math"
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

	for i := 0; i < len(dct); i++ {
		if i == 0 {
			dct[i] >>= 7
		} else {
			dct[i] = dct[i] >> 5
		}
	}

	fingerprint := simian.FlattenRecursiveSquares(dct)

	return fingerprint
}

func makeFingerprintFromImageFile(imageFileName string) ([]int16, error) {
	imageFile, err := os.Open(imageFileName)
	if err != nil {
		return nil, err
	}
	defer imageFile.Close()

	img, _, err := image.Decode(imageFile)
	if err != nil {
		return nil, err
	}

	return makeFingerprint(img), nil
}

func main() {
	if len(os.Args) < 3 {
		fmt.Printf("Usage: simian-compare <image1.jpg> <image2.jpg>\n")
		return
	}

	fingerprint1, err := makeFingerprintFromImageFile(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	fingerprint2, err := makeFingerprintFromImageFile(os.Args[2])
	if err != nil {
		fmt.Printf("Error: %s\n", err)
		return
	}

	difference := 0.0
	for i := 0; i < len(fingerprint1); i++ {
		difference += math.Abs(float64(fingerprint1[i] - fingerprint2[i]))
	}
	difference /= samplesPerFingerprint * 12

	var judgment string
	switch {
	case difference < 0.05:
		judgment = "duplicate"
	case difference < 0.1:
		judgment = "variation"
	case difference < 0.2:
		judgment = "similar"
	case difference < 0.3:
		judgment = "tonally/texturally similar"
	default:
		judgment = "different"
	}

	fmt.Printf("%.4f (%s)\n", difference, judgment)
}
