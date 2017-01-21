package simian

import (
	"encoding/hex"
	"fmt"
	"image"
	"image/color"
	"math"
	"testing"
)

func TestBytesSerialisesToPackedBytes(t *testing.T) {
	f := Fingerprint{samples: []byte{0x00, 0x00, 0xF0, 0xF0}}

	actualString := fmt.Sprintf("%x", f.Bytes())

	if actualString != "00ff" {
		t.Errorf("Fingerprint '%s' doesn't match expected", actualString)
	}
}

func TestDifferenceReturnsZeroForSameFingerprint(t *testing.T) {
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
}

func TestDifferenceReturnsOneForCompletelyDifferentFingerprint(t *testing.T) {
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
}

func TestDifferenceReturnsOneForDifferentlySizedFingerprint(t *testing.T) {
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
}

func TestDistanceReturnsComponentwiseAbsoluteDifference(t *testing.T) {
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
}

func TestDistanceReturnsMaxValueForMismatchedLength(t *testing.T) {
	f1 := Fingerprint{samples: []byte{0, 0, 0}}
	f2 := Fingerprint{samples: []byte{0, 0, 0, 0}}

	dist := f1.Distance(f2)

	if dist != math.MaxUint64 {
		t.Errorf("Distance %d wasn't max uint64", dist)
	}
}

func TestMarshalTextSerialisesToPackedHexStringBytes(t *testing.T) {
	f := Fingerprint{samples: []byte{0x00, 0x00, 0xFF, 0xFF}}

	actual, err := f.MarshalText()

	if err != nil {
		t.Errorf("Error while marshalling: %s", err)
	}
	if string(actual) != "00ff" {
		t.Errorf("Fingerprint '%s' doesn't match expected", actual)
	}
}

func TestSizeReturnsCorrectSideLength(t *testing.T) {
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
}

func TestStringSerialisesToPackedHexString(t *testing.T) {
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
}

func TestUnmarshalBytesDeserialisesFromPackedBytes(t *testing.T) {
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
}

func TestUnmarshalTextDeserialisesFromPackedHexStringBytes(t *testing.T) {
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
}

func TestNewFingerprintGeneratesBinaryRepresentation(t *testing.T) {
	f := NewFingerprint(testImage(), 3)

	expected, _ := hex.DecodeString("0040804080804080c0")

	expectedString := hex.EncodeToString(expected)
	actualString := hex.EncodeToString(f.samples)

	if expectedString != actualString {
		t.Fatalf("Fingerprint '%s' doesn't match expected '%s'", actualString, expectedString)
	}
}

func testImage() image.Image {
	img := image.NewNRGBA(image.Rectangle{Max: image.Point{X: 256, Y: 256}})

	for i := img.Bounds().Min.Y; i < img.Bounds().Max.Y; i++ {
		for j := img.Bounds().Min.X; j < img.Bounds().Max.X; j++ {
			img.Set(j, i, color.RGBA{uint8(i), uint8(j), uint8(i), 255})
		}
	}

	return img
}
