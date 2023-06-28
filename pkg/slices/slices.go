package slices

import "errors"

type Function[E any, S any] func(item E) S
type Comparator[E any, S any] func(item1 E, item2 S) bool

// Map function converts items in the list to the output list.
func Map[E any, S any](source []E, function Function[E, S]) (output []S) {
	for _, item := range source {
		output = append(output, function(item))
	}

	return output
}

func HasCommon[E comparable](slice1 []E, slice2 []E) bool {
	elementMap := make(map[E]bool)
	for _, element := range slice1 {
		elementMap[element] = true
	}

	for _, element := range slice2 {
		if _, ok := elementMap[element]; ok {
			return true
		}
	}

	return false
}

// MatchOrder takes two lists, 'source' and 'second' and returns new list with elements from 'second' with preserved
// order of elements from 'source'.
func MatchOrder[E any, S any](source []E, second []S, comparator Comparator[E, S]) (result []S, err error) {

	if len(source) != len(second) {
		return nil, errors.New("provided lists have different length")
	}

	for _, elem1 := range source {
		matched := false
		for _, elem2 := range second {
			if comparator(elem1, elem2) {
				result = append(result, elem2)
				matched = true
				continue
			}
		}

		if !matched {
			var zeroVal S
			result = append(result, zeroVal)
		}
	}

	return result, nil
}
