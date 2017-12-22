package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"

	"github.com/mandykoh/simian"
)

func makeFingerprintFromImageFile(imageFileName string) (f simian.Fingerprint, err error) {
	var imageFile *os.File
	imageFile, err = os.Open(imageFileName)
	if err != nil {
		return
	}
	defer imageFile.Close()

	var img image.Image
	img, _, err = image.Decode(imageFile)
	if err != nil {
		return
	}

	return simian.FingerprintFromImage(img), nil
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
	difference /= float64(len(fingerprint1) * 12)

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
