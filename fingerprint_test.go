package simian

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"math"
	"testing"
)

func TestFingerprint(t *testing.T) {

	testImage := func() image.Image {
		img := image.NewNRGBA(image.Rectangle{Max: image.Point{X: 256, Y: 256}})

		for i := img.Bounds().Min.Y; i < img.Bounds().Max.Y; i++ {
			for j := img.Bounds().Min.X; j < img.Bounds().Max.X; j++ {
				img.Set(j, i, color.RGBA{uint8(i), uint8(j), uint8(i), 255})
			}
		}

		return img
	}

	t.Run("Bytes() serialises to packed bytes", func(t *testing.T) {
		f := Fingerprint{samples: []byte{0x00, 0x00, 0xF0, 0xF0}}

		actualString := fmt.Sprintf("%x", f.Bytes())

		if actualString != "00ff" {
			t.Errorf("Fingerprint '%s' doesn't match expected", actualString)
		}
	})

	t.Run("Difference() returns zero for same fingerprint", func(t *testing.T) {
		f1 := Fingerprint{samples: []byte{0, 1, 2, 3, 130, 255}}
		f2 := Fingerprint{samples: []byte{0, 1, 2, 3, 130, 255}}

		diff := f1.Difference(f2)

		if diff != 0.0 {
			t.Errorf("Difference %f doesn't match expected", diff)
		}

		diff = f2.Difference(f1)

		if diff != 0.0 {
			t.Errorf("Difference %f doesn't match expected", diff)
		}
	})

	t.Run("Difference() returns one for completely different fingerprint", func(t *testing.T) {
		f1 := Fingerprint{samples: []byte{0, 0, 0, 255, 255, 255}}
		f2 := Fingerprint{samples: []byte{255, 255, 255, 0, 0, 0}}

		diff := f1.Difference(f2)

		if diff != 1.0 {
			t.Errorf("Difference %f doesn't match expected", diff)
		}

		diff = f2.Difference(f1)

		if diff != 1.0 {
			t.Errorf("Difference %f doesn't match expected", diff)
		}
	})

	t.Run("Difference() returns one for differently sized fingerprint", func(t *testing.T) {
		f1 := Fingerprint{samples: []byte{255, 255, 255}}
		f2 := Fingerprint{samples: []byte{255, 255, 255, 255}}

		diff := f1.Difference(f2)

		if diff != 1.0 {
			t.Errorf("Difference %f doesn't match expected", diff)
		}

		diff = f2.Difference(f1)

		if diff != 1.0 {
			t.Errorf("Difference %f doesn't match expected", diff)
		}
	})

	t.Run("Distance() returns componentwise absolute difference", func(t *testing.T) {
		f1 := Fingerprint{samples: []byte{0, 1, 2, 3, 130, 255}}
		f2 := Fingerprint{samples: []byte{1, 3, 6, 11, 146, 0}}

		dist := f1.Distance(f2)

		if dist != 286 {
			t.Errorf("Distance %d doesn't match expected", dist)
		}

		dist = f2.Distance(f1)

		if dist != 286 {
			t.Errorf("Distance %d doesn't match expected", dist)
		}
	})

	t.Run("Distance() returns max value for mismatched length", func(t *testing.T) {
		f1 := Fingerprint{samples: []byte{0, 0, 0}}
		f2 := Fingerprint{samples: []byte{0, 0, 0, 0}}

		dist := f1.Distance(f2)

		if dist != math.MaxUint64 {
			t.Errorf("Distance %d wasn't max uint64", dist)
		}
	})

	t.Run("MarshalText() serialises to packed hex string bytes", func(t *testing.T) {
		f := Fingerprint{samples: []byte{0x00, 0x00, 0xFF, 0xFF}}

		actual, err := f.MarshalText()

		if err != nil {
			t.Errorf("Error while marshalling: %s", err)
		}
		if string(actual) != "00ff" {
			t.Errorf("Fingerprint '%s' doesn't match expected", actual)
		}
	})

	t.Run("Size() returns correct side length", func(t *testing.T) {
		img := testImage()

		f := NewFingerprint(img, 3)
		size := f.Size()

		if size != 3 {
			t.Errorf("Size %d doesn't match expected", size)
		}

		f = NewFingerprint(img, 7)
		size = f.Size()

		if size != 7 {
			t.Errorf("Size %d doesn't match expected", size)
		}

		f = Fingerprint{samples: make([]byte, 5*5)}
		size = f.Size()

		if size != 5 {
			t.Errorf("Size %d doesn't match expected", size)
		}
	})

	t.Run("String() serialises to packed hex string", func(t *testing.T) {
		f := Fingerprint{samples: []byte{
			0xF0, 0xF0, 0xF0, 0xF0, 0xF0,
			0xF0, 0xF0, 0xF0, 0xF0, 0xF0,
			0xF0, 0xF0, 0xF0, 0xF0, 0xF0,
			0xF0, 0xF0, 0xF0, 0xF0, 0xF0,
			0xF0, 0xF0, 0xF0, 0xF0, 0xF0,
		}}

		actualString := fmt.Sprintf("%s", f)

		if actualString != "fffffffffffffffffffffffff0" {
			t.Errorf("Fingerprint '%s' doesn't match expected", actualString)
		}
	})

	t.Run("UnmarshalBytes() deserialises from packed bytes", func(t *testing.T) {
		b := []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xF0}

		f := Fingerprint{}
		f.UnmarshalBytes(b)

		if len(f.samples) != 25 {
			t.Fatalf("Fingerprint length %d doesn't match expected", len(f.samples))
		}
		for i := 0; i < 25; i++ {
			if f.samples[i] != 0xF0 {
				t.Errorf("Fingerprint byte '%d' doesn't match expected", f.samples[i])
			}
		}
	})

	t.Run("UnmarshalText() deserialises from packed hex string bytes", func(t *testing.T) {
		text := []byte("fffffffffffffffffffffffff0")

		f := Fingerprint{}
		f.UnmarshalText(text)

		if len(f.samples) != 25 {
			t.Fatalf("Fingerprint length %d doesn't match expected", len(f.samples))
		}
		for i := 0; i < 25; i++ {
			if f.samples[i] != 0xF0 {
				t.Errorf("Fingerprint byte '%d' doesn't match expected", f.samples[i])
			}
		}
	})

	t.Run("NewFingerprint() generates binary representation", func(t *testing.T) {
		f := NewFingerprint(testImage(), 3)

		expected, _ := hex.DecodeString("3060805080a070a0c0")

		expectedString := hex.EncodeToString(expected)
		actualString := hex.EncodeToString(f.samples)

		if expectedString != actualString {
			t.Fatalf("Fingerprint '%s' doesn't match expected '%s'", actualString, expectedString)
		}
	})
}
