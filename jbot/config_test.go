package jbot

import (
    "testing"
)

// configsAreIdentical returns true if configs a and b
// have fully identical porperties.
func configsAreIdentical(a config, b config) bool {
    if a.APIKey                          == b.APIKey &&
        a.Debug                          == b.Debug &&
        a.DatabaseURL                    == b.DatabaseURL &&
        a.CommandConfigs.Start.Reply     == b.CommandConfigs.Start.Reply &&
        a.CommandConfigs.Wisdom.Reply    == b.CommandConfigs.Wisdom.Reply &&
        a.CommandConfigs.Horoscope.Reply == b.CommandConfigs.Horoscope.Reply &&
        stringSlicesAreEqual(a.CommandConfigs.Start.Alias, b.CommandConfigs.Start.Alias) &&
        stringSlicesAreEqual(a.CommandConfigs.Wisdom.Alias, b.CommandConfigs.Wisdom.Alias) &&
        stringSlicesAreEqual(a.CommandConfigs.Horoscope.Alias, b.CommandConfigs.Horoscope.Alias){
            return true
        }
        
    return false
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
        CommandConfigs : commandConfigList{
        commandConfig{[]string{"/s"}, "Greetings friend!"},
        commandConfig{[]string{"/w"}, ""},
        commandConfig{[]string{"/h"}, ""},
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
        CommandConfigs : commandConfigList{
        commandConfig{[]string{"/s"}, "Greetings friend!"},
        commandConfig{[]string{"/w"}, ""},
        commandConfig{[]string{"/h"}, ""},
        },
    }
    
    if !configsAreIdentical(actualResult, expectedResult) {
        t.Fatalf("Unmarshaling json produced unexpected values")
    }
}
