package simian

import (
	"encoding/json"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

const keyBitLength = 256

type IndexEntry struct {
	Thumbnail      image.Image
	MaxFingerprint Fingerprint
	Attributes     map[string]interface{}
}

func (entry *IndexEntry) FingerprintForSize(size int) Fingerprint {
	return NewFingerprint(entry.Thumbnail, size)
}

func (entry *IndexEntry) MarshalJSON() ([]byte, error) {
	return json.Marshal(&indexEntryJSON{
		MaxFingerprint: entry.MaxFingerprint.Bytes(),
		Attributes:     entry.Attributes,
	})
}

func (entry *IndexEntry) UnmarshalJSON(b []byte) error {
	var value indexEntryJSON
	err := json.Unmarshal(b, &value)
	if err != nil {
		return err
	}

	var fingerprint Fingerprint
	err = fingerprint.UnmarshalBytes(value.MaxFingerprint)
	if err != nil {
		return err
	}

	entry.MaxFingerprint = fingerprint
	entry.Attributes = value.Attributes

	return nil
}

func (entry *IndexEntry) loadThumbnail(path string) error {
	thumbnailFile, err := os.Open(path)
	if err != nil {
		return err
	}
	defer thumbnailFile.Close()

	entry.Thumbnail, err = png.Decode(thumbnailFile)
	return err
}

func (entry *IndexEntry) saveThumbnail(path string) error {
	thumbnailDir := filepath.Dir(path)
	os.MkdirAll(thumbnailDir, os.FileMode(0700))

	thumbnailOut, err := os.Create(path)
	if err != nil {
		return err
	}
	defer thumbnailOut.Close()

	pngEncoder := png.Encoder{}
	return pngEncoder.Encode(thumbnailOut, entry.Thumbnail)
}

func NewIndexEntry(image image.Image, maxFingerprintSize int, attributes map[string]interface{}) (*IndexEntry, error) {
	entry := &IndexEntry{
		Thumbnail:  makeThumbnail(image, maxFingerprintSize*2),
		Attributes: attributes,
	}

	entry.MaxFingerprint = entry.FingerprintForSize(maxFingerprintSize)

	return entry, nil
}

func makeThumbnail(src image.Image, size int) image.Image {
	width := float64(src.Bounds().Max.X - src.Bounds().Min.X)
	height := float64(src.Bounds().Max.Y - src.Bounds().Min.Y)
	target := float64(size)

	if width > height {
		width /= height / target
		height = target
	} else {
		height /= width / target
		width = target
	}

	thumbnail := image.NewNRGBA(image.Rect(0, 0, int(width), int(height)))
	draw.BiLinear.Scale(thumbnail, thumbnail.Bounds(), src, src.Bounds(), draw.Src, nil)

	return thumbnail
}

type indexEntryJSON struct {
	MaxFingerprint []byte                 `json:"maxFingerprint"`
	Attributes     map[string]interface{} `json:"attributes"`
}
