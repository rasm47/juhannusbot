package config

import (
    "strings"
    
    "github.com/ruoskija/juhannusbot/util"
)

type Config struct {
    ApiKey string
    Debug bool
    BibleFilename string
}

func Configure() (c Config, err error) {
    
    lines, err := util.FileToLines("config.txt")
    if err != nil {
        return 
    }

    c.ApiKey = lines[0]
    if strings.HasPrefix( strings.ToLower(lines[1]), "true") {
        c.Debug = true
    } else {
        c.Debug = false
    }
    c.BibleFilename = lines[2]
    
    return
}
