package main

import (
    "log"
    
    "github.com/ruoskija/juhannusbot/jbot"
)

func main() {

    err := jbot.Start()
    if err != nil {
        log.Panic(err)
    }
    
}
