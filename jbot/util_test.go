package jbot

import (
	"testing"
)

func TestStringSlicesAreEqualNormal(t *testing.T) {

	stringSliceA := []string{"a", "b", "c"}
	stringSliceB := []string{"a", "b", "c"}

	if !stringSlicesAreEqual(stringSliceA, stringSliceB) {
		t.Fatalf("Equal string slices did not return true")
	}
}

func TestStringSlicesAreEqualSizeWrong(t *testing.T) {

	stringSliceA := []string{"a", "b", "c"}
	stringSliceB := []string{"a", "b"}

	if stringSlicesAreEqual(stringSliceA, stringSliceB) {
		t.Fatalf("string slices with different sizes returned true")
	}
}

func TestStringSlicesAreEqualContentWrong(t *testing.T) {

	stringSliceA := []string{"a", "b", "c"}
	stringSliceB := []string{"a", "b", "C"}

	if !stringSlicesAreEqual(stringSliceA, stringSliceB) {
		t.Fatalf("A difference in string content did not cause the function to return false")
	}
}

func TestStringSlicesAreEqualEmptySlices(t *testing.T) {

	stringSliceA := []string{}
	stringSliceB := []string{}

	if !stringSlicesAreEqual(stringSliceA, stringSliceB) {
		t.Fatalf("Two empty slices were not seen as equal")
	}
}
