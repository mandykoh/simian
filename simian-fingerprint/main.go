package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/mandykoh/simian"
)

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

	fingerprint := simian.FingerprintFromImage(img)

	for i := 0; i < len(fingerprint); i++ {
		fmt.Printf("%02x", fingerprint[i]+128)
	}
	fmt.Println()
}
