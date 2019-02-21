package jbot

import (
    "testing"
)

func TestConfigureFromWorkingFile(t *testing.T) {
    actualResult, err := configureFromFile("tests/working_config.json")
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    
    expectedResult := config{
        apiKey:      "TestKey123",
        debug:       false,
        databaseURL: "Poirot.txt",
    }
    
    if actualResult.apiKey != expectedResult.apiKey {
        t.Fatalf("Unmarshaling json produced unexpected values")
    }
}

func TestConfigureFromNonExistentFile(t *testing.T) {
    _, err := configureFromFile("tests/not_found_config.json")
    if err == nil {
        t.Fatalf("Opening a non-existent file did not produce any errors")
    }
}

func TestConfigureFromBrokenFile(t *testing.T) {
    _, err := configureFromFile("tests/broken_config.json")
    if err == nil {
        t.Fatalf("Opening a broken file did not produce any errors")
    }
}

func TestConfigureFromModifiedFile(t *testing.T) {
    actualResult, err := configureFromFile("tests/config_extra_options.json")
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    
    expectedResult := config{
        apiKey:      "TestKey123",
        debug:       false,
        databaseURL: "Poirot.txt",
    }
    
    if actualResult.apiKey != expectedResult.apiKey {
        t.Fatalf("Unmarshaling json produced unexpected values")
    }
}
