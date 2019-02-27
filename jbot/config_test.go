package jbot

import (
    "testing"
)

// configsAreIdentical returns true if configs a and b
// have fully identical properties.
func configsAreIdentical(a config, b config) bool {
    if a.APIKey       != b.APIKey ||
        a.Debug       != b.Debug ||
        a.DatabaseURL != b.DatabaseURL {
            return false
        }
    
    if len(a.CommandConfigs) != len(b.CommandConfigs) {
        return false
    }
    
    for index, command := range a.CommandConfigs {
        if b.CommandConfigs[index].Name != command.Name ||
           !stringSlicesAreEqual(a.CommandConfigs[index].Alias, b.CommandConfigs[index].Alias) || 
           !stringSlicesAreEqual(a.CommandConfigs[index].Reply, b.CommandConfigs[index].Reply) {
            return false
        }
    }
    
    return true
}

func TestConfigureFromWorkingFile(t *testing.T) {
    actualResult, err := configureFromFile("tests/working_config.json")
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    
    expectedResult := config{
        APIKey:      "TestKey123",
        Debug:       false,
        DatabaseURL: "Poirot.txt",
        CommandConfigs: map[string]commandConfig{
        "start": commandConfig{"start", []string{"/s"}, []string{"Greetings friend!"}},
        "wisdom": commandConfig{"wisdom", []string{"/w"}, []string{""}},
        "horoscope": commandConfig{"horoscope", []string{"/h"}, []string{""}},
        },
    }
    
    if !configsAreIdentical(actualResult, expectedResult) {
        t.Fatalf("Unmarshaling a valid and working json file produced unexpected values")
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
        APIKey:      "TestKey123",
        Debug:       false,
        DatabaseURL: "Poirot.txt",
        CommandConfigs: map[string]commandConfig{
        "start": commandConfig{"start", []string{"/s"}, []string{"Greetings friend!"}},
        "wisdom": commandConfig{"wisdom", []string{"/w"}, []string{""}},
        "horoscope": commandConfig{"horoscope", []string{"/h"}, []string{""}},
        },
    }
    
    if !configsAreIdentical(actualResult, expectedResult) {
        t.Fatalf("Unmarshaling json produced unexpected values")
    }
}
