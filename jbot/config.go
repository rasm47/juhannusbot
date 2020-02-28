package jbot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

const configFileName = "config.json"

// config holds the configuration data for jbot.
type config struct {
	APIKey      string          `json:"apikey"`
	DatabaseURL string          `json:"databaseurl"`
	Features    json.RawMessage `json:"features"`
}

// configure reads config.json to a config struct.
func configure() (config, error) {

	rawBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		errorMessage := "Failed to open \"" + configFileName +
			"\". Check that your current working directory has " +
			"a file called \"" + configFileName + "\"."
		err = errors.New(errorMessage)
		return config{}, err
	}

	var cfg config

	err = json.Unmarshal(rawBytes, &cfg)
	if err != nil {
		return config{}, err
	}

	if cfg.APIKey == "" {
		err = errors.New("Could not find apikey in " + configFileName)
		return config{}, err
	}

	if cfg.DatabaseURL == "" {
		err = errors.New("Could not find databaseurl in" + configFileName)
		return config{}, err
	}

	return cfg, nil

}
