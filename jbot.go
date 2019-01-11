package main

import (
    "log"
    "gopkg.in/telegram-bot-api.v4"
    "strings"
)

func main() {

    bot, err := tgbotapi.NewBotAPI("Your key") //Insert your API key here
    if err != nil {
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("Bot %s authenticated", bot.Self.UserName)

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        gotmsg := update.Message.Text
            log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        if strings.HasPrefix(gotmsg, "/hello"){
            tosend := "world!"
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, tosend)
            bot.Send(msg)
        }
    }
}



