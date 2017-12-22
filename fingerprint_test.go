package simian

import (
	"image"
	"image/color"
	"testing"
)

func TestFingerprint(t *testing.T) {

	t.Run("dctToFingerprint()", func(t *testing.T) {

		t.Run("produces a recursive square traversal of a square 2D matrix", func(t *testing.T) {
			m := []int16{
				0, 1, 4, 9, 16, 25, 36, 49,
				2, 3, 6, 11, 18, 27, 38, 51,
				5, 7, 8, 13, 20, 29, 40, 53,
				10, 12, 14, 15, 22, 31, 42, 55,
				17, 19, 21, 23, 24, 33, 44, 57,
				26, 28, 30, 32, 34, 35, 46, 59,
				37, 39, 41, 43, 45, 47, 48, 61,
				50, 52, 54, 56, 58, 60, 62, 63,
			}

			result := dctToFingerprint(m)

			if expected, actual := len(m), len(result); expected != actual {
				t.Fatalf("Expected result to be of length %d but got %d", expected, actual)
			}

			for i := 0; i < len(result); i++ {
				if result[i] != int16(i) {
					t.Errorf("Expected element %d but got %d", i, result[i])
				}
			}
		})
	})

	t.Run("FingerprintFromImage()", func(t *testing.T) {

		t.Run("should product correct fingerprint from DCT of white image", func(t *testing.T) {
			img := image.NewNRGBA(image.Rectangle{Max: image.Point{X: 256, Y: 256}})
			for i := img.Bounds().Min.Y; i < img.Bounds().Max.Y; i++ {
				for j := img.Bounds().Min.X; j < img.Bounds().Max.X; j++ {
					img.Set(j, i, color.RGBA{uint8(255), uint8(255), uint8(255), 255})
				}
			}

			f := FingerprintFromImage(img)

			if expected, actual := int16(8064>>fingerprintACShift), f[0]; actual != expected {
				t.Errorf("Expected value %d but found %d at position 0", expected, actual)
			}

			for i := 1; i < len(f); i++ {
				if expected, actual := int16(0), f[i]; actual != expected {
					t.Errorf("Expected value %d but found %d at position %d", expected, actual, i)
				}
			}
		})

		t.Run("should product correct fingerprint from DCT of checkered image", func(t *testing.T) {
			img := image.NewNRGBA(image.Rectangle{Max: image.Point{X: fingerprintDCTSideLength, Y: fingerprintDCTSideLength}})
			offset := 0
			for i := img.Bounds().Min.Y; i < img.Bounds().Max.Y; i++ {
				for j := img.Bounds().Min.X; j < img.Bounds().Max.X; j++ {
					if offset%2 == 0 {
						img.Set(j, i, color.RGBA{uint8(255), uint8(255), uint8(255), 255})
					} else {
						img.Set(j, i, color.RGBA{uint8(0), uint8(0), uint8(0), 255})
					}
					offset++
				}
				offset++
			}

			f := FingerprintFromImage(img)

			expected := Fingerprint{
				-1, 0, 0, 4, 0, 0, 0, 0,
				0, 0, 0, 4, 4, 0, 0, 5,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 7, 7, 0, 0, 8,
				8, 0, 0, 12, 0, 0, 0, 0,
				0, 0, 0, 0, 0, 0, 0, 0,
				0, 0, 0, 20, 20, 0, 0, 24,
				24, 0, 0, 36, 36, 0, 0, 104,
			}

			for i := 0; i < len(expected); i++ {
				if expected[i] != f[i] {
					t.Errorf("Expected value %d but found %d at position %d", expected[i], f[i], i)
				}
			}
		})
	})
}
