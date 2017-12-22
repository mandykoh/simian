package simian

import (
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
}
