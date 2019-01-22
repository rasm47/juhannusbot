package jbot

import (
    "log"
    "strings"
    "strconv"
    "net/http"
    "math/rand"
    "io/ioutil"
    "encoding/json"
    
    "github.com/go-telegram-bot-api/telegram-bot-api"
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
    
    updateConfig := tgbotapi.NewUpdate(0)
    updateConfig.Timeout = 60

    updates, err := bot.GetUpdatesChan(updateConfig)
    if err != nil {
        return nil, nil, nil, err
    }
    
    return bot, bible, updates, nil
}

func HandleUpdate(bot *tgbotapi.BotAPI, bible []string, update tgbotapi.Update) (err error) {
    
    log.Printf("[%s %s %s] %s", strconv.Itoa(update.Message.From.ID), update.Message.From.UserName, update.Message.From.FirstName, update.Message.Text)
    
    response, err := createResponse(update.Message.Text, bible)
    if err != nil {
        return
    }
    sendMessage(bot, update.Message.Chat.ID, response)
    
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

func resolveHoroscope(sign string) (reply string, err error) {
    
    response, err := http.Get("http://theastrologer-api.herokuapp.com/api/horoscope/" + sign + "/today")
    if err != nil {
        log.Fatal(err)
        return
    }
    defer response.Body.Close()
    
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        log.Fatal(err)
        return
    }
    
    var hresponse horoscopeResponse
    err = json.Unmarshal(bodyBytes, &hresponse)
    if err != nil {
        log.Panic(err)
        return
    }
    
    reply = "Enkelit välittävät horoskooppinne:\n" +
    hresponse.Horoscope + "\n\nAvainsanat: " +
    hresponse.Meta.Keywords + "\nTunnetila: " +
    hresponse.Meta.Mood  + "\n\nHoroskooppi väittyi energiatasolla " +
    hresponse.Meta.Intensity + "."
    
    return
}

func getBibleLine(bible []string) string {
  return bible[rand.Intn(len(bible))]
}

func sendMessage(bot *tgbotapi.BotAPI, chatId int64, message string) {
    msg := tgbotapi.NewMessage(chatId, message)
    bot.Send(msg)    
}

func createResponse(message string, bible []string) (response string, err error) {
    
    if strings.HasPrefix(message, "/hello"){
        response = "world!"
        
    } else if strings.HasPrefix(message, "/horos"){
        horoscopeSign := parseHoroscopeMessage(message)
        if horoscopeSign == "" {
            response = ""
            return
        }
        response, err = resolveHoroscope(horoscopeSign)
        if err != nil {
            return
        }
        
    } else if strings.HasPrefix(message, "/raamatt"){
        response = getBibleLine(bible)
    } else {
        response = ""
    }
    
    return
}
