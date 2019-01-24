// Package jbot provides a telegram bot for entertainment purposes.
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
)

// Start starts the bot.
func Start() (*tgbotapi.BotAPI, []string, tgbotapi.UpdatesChannel, error) {
    
    cfg, err := configure()
    if err != nil {
        return nil, nil, nil, err 
    }
    
    bot, err := tgbotapi.NewBotAPI(cfg.APIKey) 
    if err != nil {
        return nil, nil, nil, err
    }
    
    bot.Debug = cfg.Debug
    log.Printf("%s authenticated", bot.Self.UserName)
    
    book, err := readFileToLines(cfg.BookFilename)
    if err != nil {
        return nil, nil, nil, err
    }
    
    updateConfig := tgbotapi.NewUpdate(0)
    updateConfig.Timeout = 60

    updates, err := bot.GetUpdatesChan(updateConfig)
    if err != nil {
        return nil, nil, nil, err
    }
    
    return bot, book, updates, nil
}

// HandleUpdate processes an update from the channel created with Start. 
func HandleUpdate(bot *tgbotapi.BotAPI, book []string, update tgbotapi.Update) (err error) {
    
    log.Printf("[%s %s %s] %s", strconv.Itoa(update.Message.From.ID), update.Message.From.UserName, update.Message.From.FirstName, update.Message.Text)
    
    response, err := createResponse(update.Message.Text, book)
    if err != nil {
        return
    }
    sendMessage(bot, update.Message.Chat.ID, response)
    
    return
}

// horoscopeResponseData contains data from
// a certain REST API json response 
type horoscopeResponseData struct {
    Date      string `json:"date"`
    Sunsign   string `json:"sunsign"`
    Horoscope string `json:"horoscope"`
    Meta      horoscopeResponseMetaData `json:"meta"`
}

// horoscopeResponseMetaData contains data from 
// a certain REST API json response
type horoscopeResponseMetaData struct {
    Intensity string `json:"intensity"`
    Keywords  string `json:"keywords"`
    Mood      string `json:"mood"`
}


// golang has no native support for enums,
// so each horoscope is associated with a number.
type horoscopeSign int 
const (
    horoscopeSignNone        horoscopeSign = 0
    horoscopeSignAries       horoscopeSign = 1
    horoscopeSignTaurus      horoscopeSign = 2
    horoscopeSignGemini      horoscopeSign = 3
    horoscopeSignCancer      horoscopeSign = 4
    horoscopeSignLeo         horoscopeSign = 5
    horoscopeSignVirgo       horoscopeSign = 6
    horoscopeSignLibra       horoscopeSign = 7
    horoscopeSignScorpio     horoscopeSign = 8
    horoscopeSignSagittarius horoscopeSign = 9
    horoscopeSignCapricorn   horoscopeSign = 10
    horoscopeSignAquarius    horoscopeSign = 11
    horoscopeSignPisces      horoscopeSign = 12
)

// String method for type horoscopeSign
func (sign horoscopeSign) String() string {
    //set out of range to horoscopeSignNone
    if sign < 0 || sign > 12 {
        return ""
    }
    
    signs := [13]string{
        "",
        "Aries",
        "Taurus",
        "Gemini",
        "Cancer",
        "Leo",
        "Virigo",
        "Libra",
        "Scorpio",
        "Sagittrius",
        "Capricorn",
        "Aquarius",
        "Pisces",
        }
        
        return signs[sign]
}

// parseHoroscopeMessage searches originalMessage for certain
// key phrases and returns a horoscopeSign if one is found.
func parseHoroscopeMessage(originalMessage string) horoscopeSign {
    msg := strings.ToLower(originalMessage)
    if strings.Contains(msg, "oina"){
        return horoscopeSignAries
    } else if strings.Contains(msg, "härk"){
        return horoscopeSignTaurus
    } else if strings.Contains(msg, "kaks"){
        return horoscopeSignGemini
    } else if strings.Contains(msg, "rap"){
        return horoscopeSignCancer
    } else if strings.Contains(msg, "leij"){
        return horoscopeSignLeo
    } else if strings.Contains(msg, "neit"){
        return horoscopeSignVirgo
    } else if strings.Contains(msg, "vaa"){
        return horoscopeSignLibra
    } else if strings.Contains(msg, "skor"){
        return horoscopeSignScorpio
    } else if strings.Contains(msg, "jous"){
        return horoscopeSignSagittarius
    } else if strings.Contains(msg, "vesi"){
        return horoscopeSignCapricorn
    } else if strings.Contains(msg, "kaur"){
        return horoscopeSignAquarius
    } else if strings.Contains(msg, "kal"){
        return horoscopeSignPisces
    } else {
        return horoscopeSignNone
    }
}

func resolveHoroscope(sign horoscopeSign) (reply string, err error) {
    
    response, err := http.Get("http://theastrologer-api.herokuapp.com/api/horoscope/" + sign.String() + "/today")
    if err != nil {
        return
    }
    defer response.Body.Close()
    
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return
    }
    
    var hresponse horoscopeResponseData
    err = json.Unmarshal(bodyBytes, &hresponse)
    if err != nil {
        return
    }
    
    reply = "Enkelit välittävät horoskooppinne:\n" +
    hresponse.Horoscope + "\n\nAvainsanat: " +
    hresponse.Meta.Keywords + "\nTunnetila: " +
    hresponse.Meta.Mood  + "\n\nHoroskooppi väittyi energiatasolla " +
    hresponse.Meta.Intensity + "."
    
    return
}

func getBookLine(book []string) string {
    if len(book) != 0 {
        return book[rand.Intn(len(book))]
    }
    return ""    
}

func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
    msg := tgbotapi.NewMessage(chatID, message)
    bot.Send(msg)    
}

func createResponse(message string, book []string) (response string, err error) {
    
    if strings.HasPrefix(message, "/hello"){
        response = "world!"
        
    } else if strings.HasPrefix(message, "/horos"){
        sign := parseHoroscopeMessage(message)
        if sign == horoscopeSignNone {
            response = horoscopeSignNone.String()
            return
        }
        response, err = resolveHoroscope(sign)
        if err != nil {
            return
        }
        
    } else if strings.HasPrefix(message, "/raamatt"){
        response = getBookLine(book)
    } else {
        response = ""
    }
    
    return
}
