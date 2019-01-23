// Package config provides configurations for github.com/ruoskija/juhannusbot/jbot
package config

import (
    "io/ioutil"
    "encoding/json"
)

// Config holds configuration data for jbot.
type Config struct {
    APIKey        string `json:"apikey"`
    Debug         bool   `json:"debug"`
    BookFilename  string `json:"book"`
}

// Configure reads config.json to a Config struct.
func Configure() (cfg Config, err error) {
        
    rawBytes, err := ioutil.ReadFile("config.json")
    if err != nil {
        return 
    }
    
    err = json.Unmarshal(rawBytes, &cfg)
    return
}
