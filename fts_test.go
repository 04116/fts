package main

import (
	"testing"
)

func equal[T comparable](a, b []T) bool {
	if len(a) != len(b) {
		return false
	}

	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func TestTextToTokens(t *testing.T) {
	actual, _ := textToTokens("The wildcat a ).")
	expect := []string{"wildcat"}

	if !equal(actual, expect) {
		t.Errorf("need %v but got %v", expect, actual)
	}
}

func TestIntersectionSelect(t *testing.T) {
	actual := intersectionSelect([][]int{
		{1, 2, 3},
		{2, 3, 5},
	})
	expect := []int{2, 3}

	if !equal(actual, expect) {
		t.Errorf("need %v but got %v", expect, actual)
	}
}

func TestQuery(t *testing.T) {
	idx := make(index)

	if idx.query("Small wild cat") != nil {
		t.Error("not yet indexed must be empty")
	}

	idx.add([]Documment{{Text: "teacher here"}})
	if idx.query("Small wild cat") != nil {
		t.Error("not found if not yet indexed")
	}
	idx.add([]Documment{{Text: "The wildcat is a species complex comprising two small wild cat species, the European wildcat (Felis silvestris) and the African wildcat (F. lybica)."}})
	if idx.query("Small wild cat") == nil {
		t.Error("must found after index")
	}
}
