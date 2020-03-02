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
		t.Fatalf("%v did not reduce its length to 4 when filtered with %v", testWords, wordsToFilter)
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
		t.Fatalf("%v changed with empty filtering", testWords)
	}
	return
}
