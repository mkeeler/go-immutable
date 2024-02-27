package immutableslice

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	zeroLengthSlice    = make([]int, 0, 10)
	nonZeroLengthSlice = []int{1, 2, 3, 4, 5}
)

func TestAppend(t *testing.T) {
	type testCase struct {
		slice    []int
		toAppend []int
		expected []int
	}

	cases := map[string]testCase{
		"empty slice zero append elements": {
			slice:    nil,
			toAppend: nil,
			expected: nil,
		},
		"empty slice non-zero append elements": {
			slice:    nil,
			toAppend: []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		"non-empty slice with zero append elements": {
			slice:    nonZeroLengthSlice,
			toAppend: nil,
			expected: nonZeroLengthSlice,
		},
		"non-empty slice with non-zero append elements": {
			slice:    []int{1, 2, 3},
			toAppend: []int{4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for name, tcase := range cases {
		tcase := tcase

		t.Run(name, func(t *testing.T) {
			// clone the value to isolate any immutability issues to a single test case
			original := slices.Clone(tcase.slice)

			actual := Append(original, tcase.toAppend...)
			// check the correctness of the operation
			require.Equal(t, tcase.expected, actual)

			// check the immutability of the input slice.
			if len(tcase.expected) == 0 {
				// when the expected slice should have no elements we want
				// to ensure a proper nil is returned as it will have no
				// backing array that could possibly be shared with the input.
				require.Nil(t, actual)
			} else {
				// Modify a value in the new slice and then check the original
				// to ensure it is unmodified.
				actual[0] = 42
				require.Equal(t, original, tcase.slice)
			}
		})
	}
}

func TestCompact(t *testing.T) {
	type testCase struct {
		slice    []int
		expected []int
	}

	cases := map[string]testCase{
		"empty slice": {
			slice:    nil,
			expected: nil,
		},
		"already compact": {
			slice:    []int{1, 2, 3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		"needs compaction": {
			slice:    []int{1, 1, 1, 2, 3, 3, 4, 5, 5, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
	}

	for name, tcase := range cases {
		tcase := tcase

		t.Run(name, func(t *testing.T) {
			validateCompact := func(t *testing.T, original, actual []int) {
				t.Helper()

				// check the correctness of the operation
				require.Equal(t, tcase.expected, actual)

				// check the immutability of the input slice.
				if len(tcase.expected) == 0 {
					// when the expected slice should have no elements we want
					// to ensure a proper nil is returned as it will have no
					// backing array that could possibly be shared with the input.
					require.Nil(t, actual)
				} else {
					// Modify a value in the new slice and then check the original
					// to ensure it is unmodified.
					actual[0] = 42
					require.Equal(t, original, tcase.slice)
				}
			}

			t.Run("Compact", func(t *testing.T) {
				// clone the value to isolate any immutability issues to a single test case
				original := slices.Clone(tcase.slice)
				actual := Compact(original)
				validateCompact(t, original, actual)
			})

			t.Run("CompactFunc", func(t *testing.T) {
				// clone the value to isolate any immutability issues to a single test case
				original := slices.Clone(tcase.slice)
				actual := CompactFunc(original, func(i, j int) bool { return i == j })
				validateCompact(t, original, actual)
			})
		})
	}
}

func TestConcat(t *testing.T) {
	type testCase struct {
		slices   [][]int
		expected []int
	}

	cases := map[string]testCase{
		"empty or nil": {
			slices:   [][]int{zeroLengthSlice, nil, zeroLengthSlice, nil},
			expected: nil,
		},
		"some non empty": {
			slices:   [][]int{{1, 2, 3}, zeroLengthSlice, {8, 9, 10}},
			expected: []int{1, 2, 3, 8, 9, 10},
		},
	}

	for name, tcase := range cases {
		tcase := tcase

		t.Run(name, func(t *testing.T) {
			// clone the value to isolate any immutability issues to a single test case
			var original [][]int
			for _, s := range tcase.slices {
				original = append(original, slices.Clone(s))
			}

			actual := Concat(original...)
			// check the correctness of the operation
			require.Equal(t, tcase.expected, actual)

			// check the immutability of the input slices.
			if len(tcase.expected) == 0 {
				// when the expected slice should have no elements we want
				// to ensure a proper nil is returned as it will have no
				// backing array that could possibly be shared with the input.
				require.Nil(t, actual)
			} else {
				// Modify a value in the new slice and then check the original
				// to ensure it is unmodified.
				actual[0] = 42
				for i, s := range original {
					require.Equal(t, tcase.slices[i], s)
				}
			}
		})
	}
}

func TestDelete(t *testing.T) {

	type testCase struct {
		slice        []int
		i, j         int
		expected     []int
		panics       bool
		shouldDelete func(int) bool
	}

	shouldDeleteWithValues := func(values ...int) func(int) bool {
		return func(i int) bool {
			return slices.Contains(values, i)
		}
	}

	cases := map[string]testCase{
		"empty slice": {
			slice:    zeroLengthSlice,
			expected: nil,
			i:        2,
			j:        3,
			shouldDelete: func(int) bool {
				return true
			},
		},
		"delete nothing": {
			slice:    nonZeroLengthSlice,
			expected: slices.Clone(nonZeroLengthSlice),
			i:        2,
			j:        2,
		},
		"delete all": {
			slice:        nonZeroLengthSlice,
			expected:     nil,
			i:            0,
			j:            len(nonZeroLengthSlice),
			shouldDelete: func(int) bool { return true },
		},
		"delete from start": {
			slice:        []int{1, 2, 3, 4, 5},
			expected:     []int{4, 5},
			i:            0,
			j:            3,
			shouldDelete: shouldDeleteWithValues(1, 2, 3),
		},
		"delete from end": {
			slice:        []int{1, 2, 3, 4, 5},
			expected:     []int{1, 2, 3},
			i:            3,
			j:            5,
			shouldDelete: shouldDeleteWithValues(4, 5),
		},
		"delete from middle": {
			slice:        []int{1, 2, 3, 4, 5},
			expected:     []int{1, 2, 5},
			i:            2,
			j:            4,
			shouldDelete: shouldDeleteWithValues(3, 4),
		},
		"out of bounds i less than j": {
			slice:  nonZeroLengthSlice,
			panics: true,
			i:      3,
			j:      2,
		},
		"out of bounds i greater than len": {
			slice:  nonZeroLengthSlice,
			panics: true,
			i:      10,
			j:      11,
		},
		"out of bounds i negative": {
			slice:  nonZeroLengthSlice,
			panics: true,
			i:      -3,
			j:      1,
		},
	}

	for name, tcase := range cases {
		tcase := tcase

		t.Run(name, func(t *testing.T) {
			original := slices.Clone(tcase.slice)

			validateDelete := func(t *testing.T, original, actual []int) {
				t.Helper()
				// check the correctness of the operation
				require.Equal(t, tcase.expected, actual)

				// check the immutability of the input slice.
				_ = append(actual, 1, 2, 3)
				require.Equal(t, original, tcase.slice)
			}

			t.Run("Delete", func(t *testing.T) {
				if tcase.panics {
					require.Panics(t, func() {
						Delete(original, tcase.i, tcase.j)
					})
				} else {
					actual := Delete(original, tcase.i, tcase.j)
					validateDelete(t, original, actual)
				}
			})

			if tcase.shouldDelete != nil {
				t.Run("DeleteFunc", func(t *testing.T) {
					actual := DeleteFunc(original, tcase.shouldDelete)
					validateDelete(t, original, actual)
				})
			}
		})
	}
}

func TestInsert(t *testing.T) {
	type testCase struct {
		slice    []int
		insertAt int
		values   []int
		expected []int
		panics   bool
	}

	cases := map[string]testCase{
		"empty slice": {
			slice:    nil,
			insertAt: 0,
			values:   []int{1, 2, 3},
			expected: []int{1, 2, 3},
		},
		"insert at start": {
			slice:    []int{4, 5},
			insertAt: 0,
			values:   []int{1, 2, 3},
			expected: []int{1, 2, 3, 4, 5},
		},
		"insert at end": {
			slice:    []int{1, 2},
			insertAt: 2,
			values:   []int{3, 4, 5},
			expected: []int{1, 2, 3, 4, 5},
		},
		"insert in middle": {
			slice:    []int{1, 4, 5},
			insertAt: 1,
			values:   []int{2, 3},
			expected: []int{1, 2, 3, 4, 5},
		},
		"out of bounds": {
			slice:    []int{1, 2, 3},
			insertAt: 4,
			expected: nil,
			panics:   true,
		},
	}

	for name, tcase := range cases {
		tcase := tcase

		t.Run(name, func(t *testing.T) {
			// clone the value to isolate any immutability issues to a single test case
			original := slices.Clone(tcase.slice)

			if tcase.panics {
				require.Panics(t, func() {
					Insert(original, tcase.insertAt, tcase.values...)
				})
			} else {
				actual := Insert(original, tcase.insertAt, tcase.values...)

				// check the correctness of the operation
				require.Equal(t, tcase.expected, actual)

				// check the immutability of the input slices.
				if len(tcase.expected) == 0 {
					// when the expected slice should have no elements we want
					// to ensure a proper nil is returned as it will have no
					// backing array that could possibly be shared with the input.
					require.Nil(t, actual)
				} else {
					// Modify a value in the new slice and then check the original
					// to ensure it is unmodified.
					actual[0] = 42
					require.Equal(t, tcase.slice, original)
				}
			}
		})
	}
}

func TestPrepend(t *testing.T) {
	type testCase struct {
		toPrepend []int
		slice     []int
		expected  []int
	}

	cases := map[string]testCase{
		"empty slice zero prepend elements": {
			toPrepend: nil,
			slice:     nil,
			expected:  nil,
		},
		"non-empty slice zero prepend elements": {
			toPrepend: nil,
			slice:     []int{1, 2, 3},
			expected:  []int{1, 2, 3},
		},
		"empty slice with non-zero prepend elements": {
			toPrepend: nonZeroLengthSlice,
			slice:     nil,
			expected:  nonZeroLengthSlice,
		},
		"non-empty slice with non-zero prepend elements": {
			toPrepend: []int{1, 2, 3},
			slice:     []int{4, 5},
			expected:  []int{1, 2, 3, 4, 5},
		},
	}

	for name, tcase := range cases {
		tcase := tcase

		t.Run(name, func(t *testing.T) {
			// clone the value to isolate any immutability issues to a single test case
			original := slices.Clone(tcase.slice)

			actual := Prepend(original, tcase.toPrepend...)
			// check the correctness of the operation
			require.Equal(t, tcase.expected, actual)

			// check the immutability of the input slice.
			if len(tcase.expected) == 0 {
				// when the expected slice should have no elements we want
				// to ensure a proper nil is returned as it will have no
				// backing array that could possibly be shared with the input.
				require.Nil(t, actual)
			} else {
				// Modify a value in the new slice and then check the original
				// to ensure it is unmodified.
				actual[0] = 42
				require.Equal(t, original, tcase.slice)
			}
		})
	}
}

func TestReplace(t *testing.T) {
	type testCase struct {
		slice        []int
		replaceStart int
		replaceEnd   int
		values       []int
		expected     []int
		panics       bool
	}

	cases := map[string]testCase{
		"empty slice": {
			slice:        nil,
			replaceStart: 0,
			replaceEnd:   0,
			values:       []int{1, 2, 3},
			expected:     []int{1, 2, 3},
		},
		"prepend": {
			slice:        []int{4, 5},
			replaceStart: 0,
			replaceEnd:   0,
			values:       []int{1, 2, 3},
			expected:     []int{1, 2, 3, 4, 5},
		},
		"append": {
			slice:        []int{1, 2},
			replaceStart: 2,
			replaceEnd:   2,
			values:       []int{3, 4, 5},
			expected:     []int{1, 2, 3, 4, 5},
		},
		"insert": {
			slice:        []int{1, 4, 5},
			replaceStart: 1,
			replaceEnd:   1,
			values:       []int{2, 3},
			expected:     []int{1, 2, 3, 4, 5},
		},
		"remove and replace at head": {
			slice:        []int{9, 8, 4, 5},
			replaceStart: 0,
			replaceEnd:   2,
			values:       []int{1, 2, 3},
			expected:     []int{1, 2, 3, 4, 5},
		},
		"remove and replace at tail": {
			slice:        []int{1, 2, 3, 9, 8},
			replaceStart: 3,
			replaceEnd:   5,
			values:       []int{4, 5},
			expected:     []int{1, 2, 3, 4, 5},
		},
		"remove and replace in middle": {
			slice:        []int{1, 2, 9, 8, 5},
			replaceStart: 2,
			replaceEnd:   4,
			values:       []int{3, 4},
			expected:     []int{1, 2, 3, 4, 5},
		},
		"out of bounds": {
			slice:        []int{1, 2, 3},
			replaceStart: 4,
			replaceEnd:   5,
			expected:     nil,
			panics:       true,
		},
	}

	for name, tcase := range cases {
		tcase := tcase

		t.Run(name, func(t *testing.T) {
			// clone the value to isolate any immutability issues to a single test case
			original := slices.Clone(tcase.slice)

			if tcase.panics {
				require.Panics(t, func() {
					Replace(original, tcase.replaceStart, tcase.replaceEnd, tcase.values...)
				})
			} else {
				actual := Replace(original, tcase.replaceStart, tcase.replaceEnd, tcase.values...)

				// check the correctness of the operation
				require.Equal(t, tcase.expected, actual)

				// check the immutability of the input slices.
				if len(tcase.expected) == 0 {
					// when the expected slice should have no elements we want
					// to ensure a proper nil is returned as it will have no
					// backing array that could possibly be shared with the input.
					require.Nil(t, actual)
				} else {
					// Modify a value in the new slice and then check the original
					// to ensure it is unmodified.
					actual[0] = 42
					require.Equal(t, tcase.slice, original)
				}
			}
		})
	}
}

func TestReverse(t *testing.T) {
	type testCase struct {
		slice    []int
		expected []int
	}

	cases := map[string]testCase{
		"empty slice": {
			slice:    nil,
			expected: nil,
		},
		"single element": {
			slice:    []int{1},
			expected: []int{1},
		},
		"multiple elements of elements": {
			slice:    []int{1, 2, 3, 4},
			expected: []int{4, 3, 2, 1},
		},
	}

	for name, tcase := range cases {
		tcase := tcase

		t.Run(name, func(t *testing.T) {
			// clone the value to isolate any immutability issues to a single test case
			original := slices.Clone(tcase.slice)

			actual := Reverse(original)
			// check the correctness of the operation
			require.Equal(t, tcase.expected, actual)

			// check the immutability of the input slice.
			if len(tcase.expected) == 0 {
				// when the expected slice should have no elements we want
				// to ensure a proper nil is returned as it will have no
				// backing array that could possibly be shared with the input.
				require.Nil(t, actual)
			} else {
				// Modify a value in the new slice and then check the original
				// to ensure it is unmodified.
				actual[0] = 42
				require.Equal(t, original, tcase.slice)
			}
		})
	}
}

func TestSort(t *testing.T) {
	type testCase struct {
		slice []int
	}

	cases := map[string]testCase{
		"empty": {
			slice: nil,
		},
		"non-empty": {
			slice: []int{5, 3, 1, 4, 2},
		},
	}

	for name, tcase := range cases {
		tcase := tcase

		t.Run(name, func(t *testing.T) {
			// clone the value to isolate any immutability issues to a single test case
			original := slices.Clone(tcase.slice)

			validateSort := func(t *testing.T, original, actual []int) {
				t.Helper()

				// check the correctness of the operation
				require.True(t, slices.IsSorted(actual))

				// check the immutability of the input slice.
				if len(tcase.slice) == 0 {
					// when the expected slice should have no elements we want
					// to ensure a proper nil is returned as it will have no
					// backing array that could possibly be shared with the input.
					require.Nil(t, actual)
				} else {
					// Modify a value in the new slice and then check the original
					// to ensure it is unmodified.
					actual[0] = 42
					require.Equal(t, original, tcase.slice)
				}
			}

			t.Run("Sort", func(t *testing.T) {
				actual := Sort(original)

				validateSort(t, original, actual)
			})

			t.Run("SortFunc", func(t *testing.T) {
				actual := SortFunc(original, func(a, b int) int {
					if a < b {
						return -1
					} else if a > b {
						return 1
					}
					return 0
				})

				validateSort(t, original, actual)
			})

			t.Run("SortStableFunc", func(t *testing.T) {
				actual := SortStableFunc(original, func(a, b int) int {
					if a < b {
						return -1
					} else if a > b {
						return 1
					}
					return 0
				})

				validateSort(t, original, actual)
			})
		})
	}
}
