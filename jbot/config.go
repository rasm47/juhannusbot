package jbot

import (
    "io/ioutil"
    "encoding/json"
)

const (
    configFileName = "config.json"
)

// config holds configuration data for jbot.
type config struct {
    APIKey        string `json:"apikey"`
    Debug         bool   `json:"debug"`
    BookFilename  string `json:"book"`
}

// configure reads config.json to a Config struct.
func configure() (config, error) {
    return configureFromFile(configFileName)
}

// configureFromFile reads cofig data from a file specified by fileName
func configureFromFile(fileName string) (cfg config, err error) {
        
    rawBytes, err := ioutil.ReadFile(fileName)
    if err != nil {
        return 
    }
    
    err = json.Unmarshal(rawBytes, &cfg)
    return
    
}
