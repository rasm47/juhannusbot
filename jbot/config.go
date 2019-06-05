package jbot

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

const configFileName = "testconfig.json"

// config holds the configuration data for jbot.
type config struct {
	APIKey      string          `json:"apikey"`
	Debug       bool            `json:"debug"`
	DatabaseURL string          `json:"databaseurl"`
	Features    json.RawMessage `json:"features"`
}

// configure reads config.json to a config struct.
func configure() (config, error) {

	rawBytes, err := ioutil.ReadFile(configFileName)
	if err != nil {
		err = errors.New("Failed to open \"config.json\". Check that your current working " +
			"directory has a file \"config.json\".")
		return config{}, err
	}

	var cfg config

	err = json.Unmarshal(rawBytes, &cfg)
	if err != nil {
		return config{}, err
	}

	if cfg.APIKey == "" {
		err = errors.New("Could not find apikey in config.json")
		return config{}, err
	}

	if cfg.DatabaseURL == "" {
		err = errors.New("Could not find apikey in config.json")
		return config{}, err
	}

	return cfg, nil

}
