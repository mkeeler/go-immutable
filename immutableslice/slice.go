package immutableslice

import (
	"cmp"
	"slices"
)

// Append will create a new slice with the elements of s followed by the elements of e.
// This is a shorthand for a call to Concat(s, e) that allows passing the elements to
// append as variadic arguments.
func Append[S ~[]E, E any](s S, e ...E) S {
	return Concat(s, e)
}

// Compact is an immutable variant of the standard libraries slices.Compact function.
// It will replace consecutive runs of equal elements with the first instance.
func Compact[S ~[]E, E comparable](s S) S {
	if len(s) == 0 {
		return nil
	}

	newS := make(S, len(s))
	j := -1
	for _, v := range s {
		if j >= 0 && v == newS[j] {
			continue
		}

		newS[j+1] = v
		j += 1
	}

	clear(newS[j+1:])
	return newS[:j+1]
}

// CompactFunc is an immutable variant of the standard libraries slices.CompactFunc function.
// It will replace consecutive runs of elements for which the eq function returns true with
// with the first instance.
func CompactFunc[S ~[]E, E any](s S, eq func(E, E) bool) S {
	if len(s) == 0 {
		return nil
	}
	newS := make(S, len(s))
	j := -1
	for _, v := range s {
		if j >= 0 && eq(v, newS[j]) {
			continue
		}

		newS[j+1] = v
		j += 1
	}

	clear(newS[j+1:])
	return newS[:j+1]
}

// Concat is an immutable variant of the stdlibs slices.Concat function. It will create
// a new slice (with a newly allocated backing array) that stores all elements of the
// specified slices in the order that they were specified.
func Concat[S ~[]E, E any](slices ...S) S {
	size := 0
	for _, s := range slices {
		size += len(s)
		if size < 0 {
			panic("len out of range")
		}
	}

	if size == 0 {
		return nil
	}

	output := make(S, size)
	index := 0
	for _, s := range slices {
		index += copy(output[index:], s)
	}
	return output
}

// Delete is an immutable variant of the standard libraries slices.Delete function.
// It will delete all elements at indexes from i up to but excluding j from the slice.
// The underlying slice and its backing array will not be modified in any way. When
// this function returns a non-nil value, the returned slice will be backed by a fresh
// array and so modifications to the output will not affect the input slice. Note that
// this will not copy the slice elements so care must still be taken to ensure no
// modifications to those elements if that is the desired constraint.
func Delete[S ~[]E, E any](s S, i, j int) S {
	if len(s) == 0 {
		// the input has nothing to delete so the output is always nil
		return nil
	}

	if i == j {
		// There is nothing to delete but the property we want to keep is that the output
		// of this function can always be freely modified without affecting the input slice
		// so we return a clone.
		return slices.Clone(s)
	}

	// Bounds check the values of i and j to ensure that
	// * i < j ( i == j is handled above )
	// * i & j are within the bounds of the slice
	//
	// This will panic if given invalid indices.
	_ = s[i:j]

	// Handle the case where all slice elements are being deleted
	if i == 0 && j == len(s) {
		return nil
	}

	newS := make(S, len(s)-(j-i))
	if i == 0 {
		// Handle the case where elements are being deleted from the start of the slice
		copy(newS, s[j:])
	} else if j == len(s) {
		// Handle the case where elements are being deleted from the end of the slice
		copy(newS, s[:i])
	} else {
		// Handle the case where elements are being deleted from the middle of the slice
		copy(newS, s[:i])
		copy(newS[i:], s[j:])
	}

	return newS
}

// DeleteFunc is an immutable variant of the standard libraries slices.DeleteFunc function.
// This function will return a new slice with all elements of s for which the del function
// returns false. The underlying slice and its backing array will not be modified.
func DeleteFunc[S ~[]E, E any](s S, del func(E) bool) S {
	if len(s) == 0 {
		return nil
	}

	newS := make(S, 0, len(s))
	for _, v := range s {
		if !del(v) {
			newS = append(newS, v)
		}
	}

	if len(newS) == 0 {
		return nil
	}

	return newS
}

// Insert is an immutable variant of the standard libraries slices.Insert function.
// It will insert the elements of v at index i in the slice s. The original slices
// will not be modified in any way.
func Insert[S ~[]E, E any](s S, i int, v ...E) S {
	return Concat(s[:i], v, s[i:])
}

// Prepend will create a new slice with the elements of e followed by the elements of s.
// This is a shorthand for a call to Concat(e, s) that allows passing the elements to
// prepend as variadic arguments.
func Prepend[S ~[]E, E any](s S, e ...E) S {
	return Concat(e, s)
}

// Replace is an immutable variant of the standard libraries slices.Replace function.
// It will replace the elements represented by s[i:j] with the elements of v.
func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
	return Concat(s[:i], v, s[j:])
}

// Reverse is an immutable variant of the standard libraries slices.Reverse function.
// It will create a new slice with a fresh backing array and populate it with elements
// from s in reverse order.
func Reverse[S ~[]E, E any](s S) S {
	if len(s) < 1 {
		return nil
	}

	newS := make(S, len(s))

	for i, j := 0, len(s)-1; i <= j; i, j = i+1, j-1 {
		newS[i], newS[j] = s[j], s[i]
	}

	return newS
}

// Sort is an immutable variant of the standard libraries slices.Sort function. It is
// a small wrapper which will clone the slice and then sort it to ensure that the original
// is not modified.
func Sort[S ~[]E, E cmp.Ordered](s S) S {
	if len(s) < 1 {
		return nil
	}
	newS := slices.Clone(s)
	slices.Sort(newS)
	return newS
}

// SortFunc is an immutable variant of the standard libraries slices.SortFunc function. It is
// a small wrapper which will clone the slice and then sort it to ensure that the original
// is not modified.
func SortFunc[S ~[]E, E any](s S, cmp func(a, b E) int) S {
	if len(s) < 1 {
		return nil
	}
	newS := slices.Clone(s)
	slices.SortFunc(newS, cmp)
	return newS
}

// SortStableFunc is an immutable variant of the standard libraries slices.SortStableFunc function. It is
// a small wrapper which will clone the slice and then sort it to ensure that the original
// is not modified.
func SortStableFunc[S ~[]E, E any](s S, cmp func(a, b E) int) S {
	if len(s) < 1 {
		return nil
	}
	newS := slices.Clone(s)
	slices.SortStableFunc(newS, cmp)
	return newS
}
