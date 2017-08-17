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

func flattenZigZag(width, height int, values []int16) []int16 {
	result := make([]int16, width*height)

	x := 0
	y := 0

	pAxis := &x
	sAxis := &y
	pBound := width
	sBound := height

	for i := 0; i < len(result); i++ {
		result[i] = values[y*width+x]

		if *pAxis+1 < pBound && *sAxis-1 >= 0 {

			// Unobstructed diagonal traversal
			*pAxis++
			*sAxis--
			continue

		} else if *pAxis+1 < pBound {

			// Obstructed at the top/left; move right/down
			*pAxis++

		} else {

			// Obstructed at the bottom/right; move right/down
			*sAxis++
		}

		// Swap direction (obstructed)
		tmpAxis := pAxis
		pAxis = sAxis
		sAxis = tmpAxis
		tmpBound := pBound
		pBound = sBound
		sBound = tmpBound
	}

	return result
}
