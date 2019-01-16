package main

import (
    "log"
    
    "github.com/ruoskija/juhannusbot/jbot"
)

func main() {

    bot, bible, updates, err := jbot.Start()
    if err != nil {
        log.Panic(err)
    }

    for update := range updates {
        if update.Message == nil {
            continue
        }

        jbot.HandleUpdate(bot, bible, update)
    }
}
