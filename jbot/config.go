package jbot

import (
    "errors"
    "io/ioutil"
    "encoding/json"
)

const (
    configFileName = "config.json"
)

// config holds configuration data for jbot.
type config struct {
    apiKey      string      `json:"apikey"`
    debug       bool        `json:"debug"`
    databaseURL string      `json:"databaseurl"`
    commands    commandlist `json:"commands"`
}

type commandlist struct {
    hello     command `json:"hello"`
    start     command `json:"start"`
    wisdom    command `json:"wisdom"`
    horoscope command `json:"horoscope"`
}

type command struct {
    alias []string `json:"alias"`
    reply string   `json:"reply"`
}

// configure reads config.json to a config struct.
func configure() (config, error) {
    return configureFromFile(configFileName)
}

// configureFromFile reads cofig data from a file specified by fileName
func configureFromFile(fileName string) (cfg config, err error) {
        
    rawBytes, err := ioutil.ReadFile(fileName)
    if err != nil {
        err = errors.New("Failed to open \"config.json\". Check that your current working " + 
            "directory has a file \"config.json\".")
        return 
    }
    
    err = json.Unmarshal(rawBytes, &cfg)
    
    if !verifyConfig(cfg) {
        err = errors.New("config.json was found and opened succesfully. " + 
            "Some of the fields in the file are missing or mistyped.")
    }
    
    return
    
}

// verifyConfig checks that a config has non-empty fields.
// Unmarshaling a file with missing entries leaves the fields empty.
func verifyConfig(cfg config) bool {
    // Debug mode defaults to false in case of errors and needs not be verified
    if cfg.apiKey == "" || cfg.databaseURL == "" {
        return false
    }
    return true
}
