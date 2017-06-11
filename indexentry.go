package simian

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"image"
	"image/png"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

const keyBitLength = 256

type IndexEntry struct {
	key            string
	Thumbnail      image.Image            `json:"-"`
	MaxFingerprint Fingerprint            `json:"maxFingerprint"`
	Attributes     map[string]interface{} `json:"attributes"`
}

func (entry *IndexEntry) FingerprintForSize(size int) Fingerprint {
	return NewFingerprint(entry.Thumbnail, size)
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

func (entry *IndexEntry) save(path string, thumbnailPath string) error {
	jsonFile := filepath.Join(path, entry.key+".entry")
	jsonOut, err := os.Create(jsonFile)
	if err != nil {
		return err
	}
	defer jsonOut.Close()

	jsonEncoder := json.NewEncoder(jsonOut)
	err = jsonEncoder.Encode(entry)
	if err != nil {
		return err
	}

	return entry.saveThumbnail(thumbnailPath)
}

func NewIndexEntry(image image.Image, maxFingerprintSize int) (*IndexEntry, error) {
	key, err := makeKey()
	if err != nil {
		return nil, err
	}

	entry := &IndexEntry{
		key:        key,
		Thumbnail:  makeThumbnail(image, maxFingerprintSize*2),
		Attributes: make(map[string]interface{}),
	}

	entry.MaxFingerprint = entry.FingerprintForSize(maxFingerprintSize)

	return entry, nil
}

func NewIndexEntryFromFile(file string) (*IndexEntry, error) {
	jsonFile, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer jsonFile.Close()

	key := filepath.Base(file)
	key = key[:(len(key) - len(filepath.Ext(key)))]

	entry := &IndexEntry{
		key: key,
	}

	jsonDecoder := json.NewDecoder(jsonFile)
	err = jsonDecoder.Decode(entry)
	if err != nil {
		return nil, err
	}

	return entry, nil
}

func makeKey() (string, error) {
	b := make([]byte, keyBitLength/8)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
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
