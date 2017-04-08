package simian

import "math"

func dct(width int, height int, values []uint8) (result []float32) {

	result = make([]float32, len(values))

	for u := 0; u < height; u++ {
		for v := 0; v < width; v++ {
			sum := 0.0

			for i := 0; i < height; i++ {
				for j := 0; j < width; j++ {

					sum += float64(values[i*width+j]) *
						math.Cos(((math.Pi*float64(u))/(2*float64(height)))*(2*float64(i)+1)) *
						math.Cos(((math.Pi*float64(v))/(2*float64(width)))*(2*float64(j)+1))
				}
			}

			result[u*width+v] = float32(sum)
		}
	}

	return
}
