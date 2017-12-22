package simian

import (
	"fmt"
	"image"
	"image/color"
	"math"

	"golang.org/x/image/draw"
)

const fingerprintDCTSideLength = 8
const fingerprintACShift = 7
const fingerprintDCShift = 5

const SamplesPerFingerprint = fingerprintDCTSideLength * fingerprintDCTSideLength

type Fingerprint [SamplesPerFingerprint]int16

func (f *Fingerprint) Difference(other *Fingerprint) float64 {
	result := 0.0
	for i := 0; i < SamplesPerFingerprint; i++ {
		result += math.Abs(float64(f[i] - other[i]))
	}

	return result / float64(SamplesPerFingerprint*12)
}

func NewFingerprintFromImage(src image.Image) *Fingerprint {
	scaled := image.NewNRGBA(image.Rectangle{Max: image.Point{X: fingerprintDCTSideLength, Y: fingerprintDCTSideLength}})
	draw.BiLinear.Scale(scaled, scaled.Bounds(), src, src.Bounds(), draw.Src, nil)

	samples := make([]int8, SamplesPerFingerprint)
	offset := 0

	for i := scaled.Bounds().Min.Y; i < scaled.Bounds().Max.Y; i++ {
		for j := scaled.Bounds().Min.X; j < scaled.Bounds().Max.X; j++ {
			r, g, b, _ := scaled.At(j, i).RGBA()
			y, _, _ := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))

			samples[offset] = int8(y - 128)
			offset++
		}
	}

	dct := DCT(fingerprintDCTSideLength, fingerprintDCTSideLength, samples)

	fmt.Printf("DCT:\n")
	for i := 0; i < len(dct); i++ {
		if i == 0 {
			dct[i] >>= fingerprintACShift
		} else {
			dct[i] = dct[i] >> fingerprintDCShift
		}

		if i > 0 && i%fingerprintDCTSideLength == 0 {
			fmt.Println()
		}
		fmt.Printf(" %5d", dct[i])
	}
	fmt.Println()
	fmt.Println()

	return dctToFingerprint(dct)
}

func dctToFingerprint(squareMatrix []int16) (f *Fingerprint) {
	f = &Fingerprint{}

	level := 0
	offset := 0

	for i := 0; i != SamplesPerFingerprint; {
		if offset == level {

			// Sample the last corner of the current square
			f[i] = squareMatrix[level*fingerprintDCTSideLength+level]
			i++

			// Start the next larger square
			offset = 0
			level++

		} else {

			// Sample one from the right and one from the bottom
			f[i] = squareMatrix[offset*fingerprintDCTSideLength+level]
			i++
			f[i] = squareMatrix[level*fingerprintDCTSideLength+offset]
			i++

			offset++
		}
	}

	return
}
