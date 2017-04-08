package simian

import "math"

func DCT(width int, height int, values []int8) (result []int16) {

	doubleWidth := 2.0 * float64(width)
	doubleHeight := 2.0 * float64(height)

	result = make([]int16, len(values))

	for u := 0; u < height; u++ {
		for v := 0; v < width; v++ {
			sum := 0.0

			for i := 0; i < height; i++ {
				for j := 0; j < width; j++ {

					sum += float64(values[i*width+j]) *
						math.Cos(((math.Pi*float64(u))/doubleHeight)*(2*float64(i)+1)) *
						math.Cos(((math.Pi*float64(v))/doubleWidth)*(2*float64(j)+1))
				}
			}

			result[u*width+v] = int16(sum)
		}
	}

	return
}
