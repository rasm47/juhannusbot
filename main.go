package main

import (
    "log"
    "gopkg.in/telegram-bot-api.v4"
    "strings"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "bytes"
    "github.com/ruoskija/juhannusbot/jbot"
)

func main() {

    apikey, err := jbot.FileToLines("apikey.txt")
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
    
    bible, err := jbot.FileToLines("bible.txt")
    if err != nil {
        log.Printf("Bible not found. Have you made a bible.txt?")
        log.Panic(err)
    }

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
            
        } else if strings.HasPrefix(gotmsg, "/horos"){
            horoscopeSign := jbot.ParseHoroscopeMessage(gotmsg)
            if horoscopeSign == "" {
                continue
            }

            response, err := http.Get("http://theastrologer-api.herokuapp.com/api/horoscope/" + horoscopeSign + "/today")
            if err != nil {
                log.Fatal(err)
            } else {
                defer response.Body.Close()
                bodyBytes, err := ioutil.ReadAll(response.Body)
                if err != nil {
                    log.Fatal(err)
                } else {
                    bodyStr := string(bodyBytes)
                    log.Printf(bodyStr)

                    //reset the response body to the original unread state
                    response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

                    var hresponse jbot.HoroscopeResponse
                    err := json.Unmarshal(bodyBytes, &hresponse)
                    if err != nil {
                        log.Panic(err)
                    } else {
                        tosend := "Enkelit välittävät horoskooppinne:\n" +
                        hresponse.Horoscope + "\n\nAvainsanat: " +
                        hresponse.Meta.Keywords + "\nTunnetila: " +
                        hresponse.Meta.Mood  + "\n\nHoroskooppi väittyi energiatasolla " +
                        hresponse.Meta.Intensity + "."
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, tosend)
                        bot.Send(msg)
                    }
                }
            }

        } else if strings.HasPrefix(gotmsg, "/raamatt"){
            tosend := jbot.GetBibleLine(bible)
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, tosend)
            bot.Send(msg)
        }
    }
}
