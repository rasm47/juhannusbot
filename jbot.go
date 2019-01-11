package main

import (
    "log"
    "gopkg.in/telegram-bot-api.v4"
    "strings"
    "os"
    "bufio"
)

func main() {

    apikey, err := FileToLines("apikey.txt")
    if err != nil {
        log.Printf("Have you put your API key to apikey.txt? See README.md")
        log.Panic(err)
    }
    
    bot, err := tgbotapi.NewBotAPI(apikey[0]) 
    if err != nil {
        log.Printf("API key authentication failed. Try to double check if the key is valid.")
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("%s authenticated", bot.Self.UserName)

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

func FileToLines(filePath string) (lines []string, err error) {
      f, err := os.Open(filePath)
      if err != nil {
              return
      }
      defer f.Close()

      scanner := bufio.NewScanner(f)
      for scanner.Scan() {
              line := scanner.Text()
              if line != "" {
                lines = append(lines, line)
              }
      }
      err = scanner.Err()
      return
}



