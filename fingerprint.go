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

const SamplesPerFingerprint = fingerprintDCTSideLength * fingerprintDCTSideLength

type Fingerprint [SamplesPerFingerprint]int16

func (f *Fingerprint) Difference(other *Fingerprint) float64 {
	result := 0.0
	for i := 0; i < SamplesPerFingerprint; i++ {
		result += math.Abs(float64(f[i] - other[i]))
	}

	return result / float64(SamplesPerFingerprint*22)
}

func (f *Fingerprint) Prefix(level int) []int16 {
	return f[:level*level]
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

			val := int8(y - 128)
			samples[offset] = val
			offset++
		}
	}

	dct := DCT(fingerprintDCTSideLength, fingerprintDCTSideLength, samples)

	min := int16(math.MaxInt16)
	max := int16(math.MinInt16)

	for i := 1; i < len(dct); i++ {
		if dct[i] < min {
			min = dct[i]
		}
		if dct[i] > max {
			max = dct[i]
		}
	}

	scale := 127.0 / float64(max-min) / 2.0

	fmt.Printf("DCT:\n")

	dct[0] >>= fingerprintACShift
	for i := 0; i < len(dct); i++ {
		if i != 0 {
			dct[i] = int16(float64(dct[i]) * scale)
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
