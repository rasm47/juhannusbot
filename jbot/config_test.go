package jbot

import (
	"testing"
)

// configsAreSimilar returns true if configs a and b
// share the same APIKey and DatabaseURL.
func configsAreSimilar(a config, b config) bool {
	return a.APIKey == b.APIKey && a.DatabaseURL == b.DatabaseURL
}

func TestConfigureFromWorkingFile(t *testing.T) {
	actualResult, err := configureFromFile("tests/working_config.json")
	if err != nil {
		t.Error(err)
	}

	expectedResult := config{
		APIKey:      "TestKey123",
		DatabaseURL: "Poirot",
		Features:    []byte("some raw bytes"),
	}

	if !configsAreSimilar(actualResult, expectedResult) {
		t.Errorf("expected %v, got %v", expectedResult, actualResult)
	}
	if len(actualResult.Features) == 0 {
		t.Error("Failed to read features portion of tests/working_config.json")
	}
	return
}

func TestConfigureFromNonExistentFile(t *testing.T) {
	_, err := configureFromFile("tests/not_found_config.json")
	if err == nil {
		t.Error("Opening a non-existent file did not produce any errors")
	}
	return
}

func TestConfigureFromBrokenFile(t *testing.T) {
	_, err := configureFromFile("tests/broken_config.json")
	if err == nil {
		t.Error("Opening a broken file did not produce any errors")
	}
	return
}

func TestConfigureFromModifiedFile(t *testing.T) {
	actualResult, err := configureFromFile("tests/config_extra_options.json")
	if err != nil {
		t.Error(err)
		t.Fail()
	}

	expectedResult := config{
		APIKey:      "TestKey123",
		DatabaseURL: "Poirot",
		Features:    []byte("some raw bytes"),
	}

	if !configsAreSimilar(actualResult, expectedResult) {
		t.Errorf("expected %v, got %v", expectedResult, actualResult)
	}
	return
}
