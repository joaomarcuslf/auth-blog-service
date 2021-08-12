package helpers

import "testing"

func TestContains(t *testing.T) {
	sarr := []string{"a", "b", "c"}

	if Contains(sarr, "d") {
		t.Error("Contains 01 failed")
	} else {
		t.Log("Contains 01 passed")
	}

	if Contains(sarr, "a") {
		t.Log("Contains 02 passed")
	} else {
		t.Error("Contains 02 failed")
	}

	if Contains(sarr, "b") {
		t.Log("Contains 03 passed")
	} else {
		t.Error("Contains 03 failed")
	}
}

func TestContainsSubSLice(t *testing.T) {
	s1 := []string{"a", "b", "c"}
	s2 := []string{"a", "b"}

	if ContainsSubSLice(s1, s2) {
		t.Error("ContainsSubSLice 01 failed")
	} else {
		t.Log("ContainsSubSLice 01 passed")
	}

	if ContainsSubSLice(s2, s1) {
		t.Log("ContainsSubSLice 02 passed")
	} else {
		t.Error("ContainsSubSLice 02 failed")
	}

	s1 = []string{"a"}
	s2 = []string{"a"}

	if ContainsSubSLice(s1, s2) {
		t.Log("ContainsSubSLice 03 passed")
	} else {
		t.Error("ContainsSubSLice 03 failed")
	}
}
