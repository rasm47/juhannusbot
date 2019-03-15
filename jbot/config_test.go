package jbot

import (
    "math"
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
           b.CommandConfigs[index].Type != command.Type ||
           b.CommandConfigs[index].IsPrefixCommand != command.IsPrefixCommand ||
           b.CommandConfigs[index].IsReply != command.IsReply ||
           math.Abs(b.CommandConfigs[index].SuccessPropability - command.SuccessPropability) > 0.000001 ||
           !stringSlicesAreEqual(a.CommandConfigs[index].Aliases, b.CommandConfigs[index].Aliases) || 
           !stringSlicesAreEqual(a.CommandConfigs[index].ReplyMessages, b.CommandConfigs[index].ReplyMessages) {
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
        DatabaseURL: "Poirot",
        CommandConfigs: map[string]commandConfig{
        "start": commandConfig{"start", "message", []string{"/s"}, true, false, []string{"Greetings friend!"}, 1.0},
        "wisdom": commandConfig{"wisdom", "special", []string{"/w"}, true, false, []string{""}, 1.0},
        "horoscope": commandConfig{"horoscope", "special", []string{"/h"}, true, false, []string{""}, 1.0},
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
        DatabaseURL: "Poirot",
        CommandConfigs: map[string]commandConfig{
        "start": commandConfig{"start", "message", []string{"/s"}, true, false, []string{"Greetings friend!"}, 1.0},
        "wisdom": commandConfig{"wisdom", "special", []string{"/w"}, true, false, []string{""}, 1.0},
        "horoscope": commandConfig{"horoscope", "special", []string{"/h"}, true, false, []string{""}, 1.0},
        },
    }
    
    if !configsAreIdentical(actualResult, expectedResult) {
        t.Fatalf("Unmarshaling json produced unexpected values")
    }
}
