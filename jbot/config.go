package jbot

import (
    "errors"
    "io/ioutil"
    "encoding/json"
)

const configFileName = "config.json"


// config holds the configuration data for jbot.
type config struct {
    APIKey         string
    Debug          bool
    DatabaseURL    string
    CommandConfigs map[string]commandConfig
}

// configData is a struct used for unmarshaling config.json
type configData struct {
    APIKey         string            `json:"apikey"`
    Debug          bool              `json:"debug"`
    DatabaseURL    string            `json:"databaseurl"`
    CommandConfigs []commandConfig   `json:"commands"`
}

// commandConfig holds configurations for a single command
type commandConfig struct {
    Name  string   `json:"name"`  // name of the command
    Type  string   `json:"type"`  // type of the command (message/special)
    Aliases []string `json:"alias"` // list of words to trigger the command
    IsPrefixCommand bool `json:"ispreifxcommand"` // if true, command triggers only if it is a prefix
    IsReply bool `json:"isreply"` // if true, the telegram message is replying to the command (feature in telegram)
    ReplyMessages []string `json:"reply"` // list of possible answers to command, random one will be sent
    SuccessPropability float64 `json:"successpropability"` // 0.0-1.0 propability, used to make the command randomly fail
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
    
    var data configData
    
    err = json.Unmarshal(rawBytes, &data)
    if err != nil {
        return
    }
    
    cfg = buildConfigFromData(data)
    
    if !verifyConfig(cfg) {
        err = errors.New("config.json was found and opened succesfully. " + 
            "Some of the fields in the file are missing or mistyped.")
    }
    
    return 
    
}

// buildConfigFromData creates a config from a configData
func buildConfigFromData(data configData) (cfg config) {
    cfg.APIKey = data.APIKey
    cfg.Debug = data.Debug
    cfg.DatabaseURL = data.DatabaseURL
    
    // put all commands from the data to a map
    cfg.CommandConfigs = make(map[string]commandConfig)
    for _, command := range data.CommandConfigs {
        cfg.CommandConfigs[command.Name] = command
    }
    
    return
}

// verifyConfig checks that a config has non-empty fields.
// APIKey, DatabaseURL and at lest one command has to be there.
func verifyConfig(cfg config) bool {
    
    if cfg.APIKey == "" || cfg.DatabaseURL == "" || len(cfg.CommandConfigs) == 0 {
        return false
    }
    return true
}
