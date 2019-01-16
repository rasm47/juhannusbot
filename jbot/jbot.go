package jbot

import (
    "log"
    "bytes"
    "strings"
    "net/http"
    "math/rand"
    "io/ioutil"
    "encoding/json"
    
    "gopkg.in/telegram-bot-api.v4"
    "github.com/ruoskija/juhannusbot/util"
    "github.com/ruoskija/juhannusbot/cfg"
)

type horoscopeMeta struct {
    Intensity string `json:"intensity"`
    Keywords  string `json:"keywords"`
    Mood      string `json:"mood"`
}

type horoscopeResponse struct {
    Date      string `json:"date"`
    Sunsign   string `json:"sunsign"`
    Horoscope string `json:"horoscope"`
    Meta      horoscopeMeta `json:"meta"`
}

func Start() (*tgbotapi.BotAPI, []string, tgbotapi.UpdatesChannel, error) {
    
    cfg, err := config.Configure()
    if err != nil {
        log.Printf("Could not find config.txt")
        return nil, nil, nil, err 
    }
    
    bot, err := tgbotapi.NewBotAPI(cfg.ApiKey) 
    if err != nil {
        log.Printf("API key authentication failed. Try to double check if the key is valid.")
        return nil, nil, nil, err
    }

    bot.Debug = cfg.Debug

    log.Printf("%s authenticated", bot.Self.UserName)
    
    bible, err := util.FileToLines(cfg.BibleFilename)
    if err != nil {
        log.Printf("Bible not found. Have you made a bible.txt?")
        return nil, nil, nil, err
    }
    
    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)
    if err != nil {
        return nil, nil, nil, err
    }
    
    return bot, bible, updates, nil
}

func HandleUpdate(bot *tgbotapi.BotAPI, bible []string, update tgbotapi.Update) (err error) {
    
    gotmsg := update.Message.Text
    log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

    if strings.HasPrefix(gotmsg, "/hello"){
        tosend := "world!"
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, tosend)
        bot.Send(msg)
            
    } else if strings.HasPrefix(gotmsg, "/horos"){
        horoscopeSign := parseHoroscopeMessage(gotmsg)
        if horoscopeSign == "" {
            return
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

                var hresponse horoscopeResponse
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
        tosend := getBibleLine(bible)
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, tosend)
        bot.Send(msg)
    }
    
    return
}

func parseHoroscopeMessage(originalMessage string) string {
    msg := strings.ToLower(originalMessage)
    if strings.Contains(msg, "oina"){
        return "aries"
    } else if strings.Contains(msg, "härk"){
        return "taurus"
    } else if strings.Contains(msg, "kaks"){
        return "gemini"
    } else if strings.Contains(msg, "rap"){
        return "cancer"
    } else if strings.Contains(msg, "leij"){
        return "leo"
    } else if strings.Contains(msg, "neit"){
        return "virgo"
    } else if strings.Contains(msg, "vaa"){
        return "libra"
    } else if strings.Contains(msg, "skor"){
        return "scorpio"
    } else if strings.Contains(msg, "jous"){
        return "sagittarius"
    } else if strings.Contains(msg, "vesi"){
        return "capricorn"
    } else if strings.Contains(msg, "kaur"){
        return "aquarius"
    } else if strings.Contains(msg, "kal"){
        return "pisces"
    } else {
        return ""
    }
}

func getBibleLine(bible []string) string {
  return bible[rand.Intn(len(bible))]
}

