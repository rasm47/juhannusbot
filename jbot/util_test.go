package jbot

import (
	"testing"
)

func TestStringHasAnyPrefix(t *testing.T) {

	testStrings := []string{
		"cats and dogs",
		"Cats And Dogs",
		"CATS AND DOGS",
		"dogs and CATS",
		"and",
		"",
	}

	testPrefixes := [][]string{
		[]string{"cats"},
		[]string{"and"},
		[]string{"dogs"},
		[]string{"cats", "and", "dogs"},
		[]string{"dogs", "and", "cats"},
		[]string{"mice", "swans", "elephants"},
		[]string{""},
	}

	expectedResults := [][]bool{
		[]bool{true, false, false, true, true, false, true},
		[]bool{true, false, false, true, true, false, true},
		[]bool{true, false, false, true, true, false, true},
		[]bool{false, false, true, true, true, false, true},
		[]bool{false, true, false, true, true, false, true},
		[]bool{false, false, false, false, false, false, true},
	}

	for i, testString := range testStrings {
		for j, testPrefixesSlice := range testPrefixes {
			testResult := stringHasAnyPrefix(testString, testPrefixesSlice) 
			if testResult != expectedResults[i][j] {
				t.Errorf("stringHasAnyPrefix with inputs (%v; %v)"+
					"did not return the expected value of %v",
					testString,
					testPrefixesSlice,
					expectedResults[i][j])
			}
		}
	}
}
