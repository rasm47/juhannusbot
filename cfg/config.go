package config

import (
    "io/ioutil"
    "encoding/json"
)

type Config struct {
    ApiKey        string `json:"apikey"`
    Debug         bool   `json:"debug"`
    BibleFilename string `json:"book"`
}

func Configure() (cfg Config, err error) {
        
    rawBytes, err := ioutil.ReadFile("config.json")
    if err != nil {
        return 
    }
    
    err = json.Unmarshal(rawBytes, &cfg)   
    return
}
