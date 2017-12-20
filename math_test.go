package simian

import "testing"

func TestMath(t *testing.T) {

	t.Run("flattenRecursiveSquares()", func(t *testing.T) {

		t.Run("produces a recursive square traversal of a square 2D matrix", func(t *testing.T) {
			m := []int16{
				0, 1, 4,
				2, 3, 6,
				5, 7, 8,
			}

			result := flattenRecursiveSquares(m)

			if expected, actual := len(m), len(result); expected != actual {
				t.Fatalf("Expected result to be of length %d but got %d", expected, actual)
			}

			for i := 0; i < len(result); i++ {
				if result[i] != int16(i) {
					t.Errorf("Expected element %d but got %d", i, result[i])
				}
			}
		})

		t.Run("produces an empty result for empty input", func(t *testing.T) {
			result := flattenRecursiveSquares([]int16{})

			if actual := len(result); actual != 0 {
				t.Fatalf("Expected zero length result but got %d", actual)
			}
		})
	})
}
