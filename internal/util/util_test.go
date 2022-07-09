package util

import "testing"

func TestCheckID(t *testing.T) {
	if !CheckID("20") {
		t.Error("failed to validate `20` as real id")
	}

	if CheckID("10l") {
		t.Error("failed to validate `10l` is real id")
	}
}

func TestContains(t *testing.T) {
	food := []string{"Pizza", "Sushi", "Tapas"}
	if !Contains(food, "Tapas") {
		t.Error("Tapas does exist in array")
	}

	if Contains(food, "tapas") {
		t.Error("tapas does not exist but still validate")
	}
}
