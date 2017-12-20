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

	t.Run("flattenZigZag()", func(t *testing.T) {

		t.Run("produces a zig-zag rearrangement of a 2D matrix", func(t *testing.T) {

			m1 := []int16{
				0, 1, 5,
				2, 4, 6,
				3, 7, 10,
				8, 9, 11,
			}

			result := flattenZigZag(3, 4, m1)

			if expected, actual := len(m1), len(result); expected != actual {
				t.Fatalf("Expected result to be of length %d but got %d", expected, actual)
			}
			for i := 0; i < len(m1); i++ {
				if result[i] != int16(i) {
					t.Errorf("Expected element %d but got %d", i, result[i])
				}
			}

			m2 := []int16{
				0, 1, 5, 6, 13,
				2, 4, 7, 12, 14,
				3, 8, 11, 15, 18,
				9, 10, 16, 17, 19,
			}

			result = flattenZigZag(5, 4, m2)

			if expected, actual := len(m2), len(result); expected != actual {
				t.Fatalf("Expected result to be of length %d but got %d", expected, actual)
			}
			for i := 0; i < len(m2); i++ {
				if result[i] != int16(i) {
					t.Errorf("Expected element %d but got %d", i, result[i])
				}
			}

			m3 := []int16{
				0, 1,
				2, 3,
			}

			result = flattenZigZag(2, 2, m3)

			if expected, actual := len(m3), len(result); expected != actual {
				t.Fatalf("Expected result to be of length %d but got %d", expected, actual)
			}
			for i := 0; i < len(m3); i++ {
				if result[i] != int16(i) {
					t.Errorf("Expected element %d but got %d", i, result[i])
				}
			}

			m4 := []int16{
				1,
			}

			result = flattenZigZag(1, 1, m4)

			if expected, actual := len(m4), len(result); expected != actual {
				t.Fatalf("Expected result to be of length %d but got %d", expected, actual)
			}
			if result[0] != int16(1) {
				t.Errorf("Expected element %d but got %d", 1, result[0])
			}
		})
	})
}
