package simian

import (
	"bytes"
	"encoding/hex"
	"image"
	"image/color"
	"math"

	"golang.org/x/image/draw"
)

const bitsPerSample = 4
const sampleBitsMask = (2 << bitsPerSample) - 1
const samplesPerByte = 8 / bitsPerSample

type Fingerprint struct {
	samples []uint8
}

func (f *Fingerprint) Bytes() []byte {
	packed := bytes.Buffer{}
	current := byte(0)
	bits := uint(8)
	i := 0

	for ; i < len(f.samples); i++ {
		y := f.samples[i]

		bits -= bitsPerSample
		current = (current << bitsPerSample) | (y >> (8 - bitsPerSample))

		if bits == 0 {
			packed.WriteByte(current)
			current = 0
			bits = 8
		}
	}

	if bits < 8 {
		current <<= bits
		packed.WriteByte(current)
	}

	return packed.Bytes()
}

func (f *Fingerprint) Difference(to Fingerprint) (diff float64) {
	return math.Min(float64(f.Distance(to))/float64(len(to.samples)*255), 1.0)
}

func (f *Fingerprint) Distance(to Fingerprint) (dist uint64) {
	if len(f.samples) != len(to.samples) {
		return math.MaxUint64
	}

	for i := 0; i < len(f.samples); i++ {
		if f.samples[i] > to.samples[i] {
			dist += uint64(f.samples[i] - to.samples[i])
		} else {
			dist += uint64(to.samples[i] - f.samples[i])
		}
	}

	return dist
}

func (f *Fingerprint) MarshalText() (text []byte, err error) {
	bytes := f.Bytes()
	result := make([]byte, hex.EncodedLen(len(bytes)))

	hex.Encode(result, bytes)
	return result, nil
}

func (f *Fingerprint) Size() int {
	return int(math.Sqrt(float64(len(f.samples))))
}

func (f Fingerprint) String() string {
	return hex.EncodeToString(f.Bytes())
}

func (f *Fingerprint) UnmarshalBytes(fingerprintBytes []byte) error {
	sampleCount := int(math.Sqrt(float64(len(fingerprintBytes) * samplesPerByte)))
	sampleCount *= sampleCount
	f.samples = make([]uint8, sampleCount)

	for i := 0; i < sampleCount; i++ {
		b := fingerprintBytes[i/samplesPerByte]
		shift := uint(8 - bitsPerSample - (i%samplesPerByte)*bitsPerSample)
		bits := b >> shift & sampleBitsMask
		f.samples[i] = bits << (8 - bitsPerSample)
	}

	return nil
}

func (f *Fingerprint) UnmarshalText(text []byte) error {
	hexBytes := make([]byte, hex.DecodedLen(len(text)))
	_, err := hex.Decode(hexBytes, text)
	if err != nil {
		return err
	}

	return f.UnmarshalBytes(hexBytes)
}

func NewFingerprint(src image.Image, size int) Fingerprint {
	scaled := image.NewNRGBA(image.Rectangle{Max: image.Point{X: size, Y: size}})
	draw.BiLinear.Scale(scaled, scaled.Bounds(), src, src.Bounds(), draw.Src, nil)

	fingerprintSamples := make([]uint8, size*size)
	offset := 0

	for i := scaled.Bounds().Min.Y; i < scaled.Bounds().Max.Y; i++ {
		for j := scaled.Bounds().Min.X; j < scaled.Bounds().Max.X; j++ {
			r, g, b, _ := scaled.At(j, i).RGBA()
			y, _, _ := color.RGBToYCbCr(uint8(r>>8), uint8(g>>8), uint8(b>>8))

			fingerprintSamples[offset] = y & 0xC0
			offset++
		}
	}

	return Fingerprint{samples: fingerprintSamples}
}
