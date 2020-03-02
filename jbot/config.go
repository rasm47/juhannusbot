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
	return configureFromFile(configFileName)
}

func configureFromFile(fileName string) (config, error) {
	rawBytes, err := ioutil.ReadFile(fileName)
	if err != nil {
		errorMessage := "Failed to open \"" + fileName +
			"\". Check that your current working directory has " +
			"a file called \"" + fileName + "\"."
		return config{}, errors.New(errorMessage)
	}

	var cfg config
	err = json.Unmarshal(rawBytes, &cfg)
	if err != nil {
		return config{}, err
	}

	if cfg.APIKey == "" {
		err = errors.New("Could not find apikey in " + fileName)
		return config{}, err
	}

	if cfg.DatabaseURL == "" {
		err = errors.New("Could not find databaseurl in" + fileName)
		return config{}, err
	}

	return cfg, nil
}
