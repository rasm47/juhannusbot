package jbot

import "testing"

func TestFilterWordsTypical(t *testing.T) {
	testWords := []string{
		"word1",
		"word2",
		"word3",
		"word4",
		"word5",
	}

	wordsToFilter := []string{
		"word2",
		"word99",
	}

	filtered := filterWords(testWords, wordsToFilter)
	if len(filtered) != 4 {
		t.Fatalf("%v did not reduce its length to 4 when filtered with %v",
			testWords,
			wordsToFilter)
	}
	return
}

func TestFilterWordsEmpty(t *testing.T) {
	testWords := []string{
		"word1",
		"word2",
		"word3",
		"word4",
		"word5",
	}

	wordsToFilter := []string{}

	filtered := filterWords(testWords, wordsToFilter)
	if len(filtered) != 5 {
		t.Fatalf("length of %v changed with empty filtering", testWords)
	}
	return
}

func TestDuplicateWordsTypical(t *testing.T) {
	testWords := []string{
		"word1",
		"word2",
		"word3",
		"word4",
		"word5",
	}

	wordsToDuplicate := []string{
		"word2",
		"word99",
	}

	duplicated := duplicateWords(testWords, wordsToDuplicate)
	if len(duplicated) != 6 {
		t.Fatalf("%v did not increase its length to 6 when duplicated with %v",
			testWords,
			wordsToDuplicate)
	}
	return
}

func TestDuplicateWordsEmpty(t *testing.T) {
	testWords := []string{
		"word1",
		"word2",
		"word3",
		"word4",
		"word5",
	}

	wordsToDuplicate := []string{}

	duplicated := duplicateWords(testWords, wordsToDuplicate)
	if len(duplicated) != 5 {
		t.Fatalf("length of %v changed with no duplicating", testWords)
	}
	return
}
