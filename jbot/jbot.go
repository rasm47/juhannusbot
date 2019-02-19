// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
    "log"
    "strings"
    "strconv"
    "net/http"
    "io/ioutil"
    "database/sql"
    "encoding/json"
    
    _ "github.com/lib/pq"
    "github.com/go-telegram-bot-api/telegram-bot-api"
)

// bot is a collection of relevant pointers.
type bot struct {
    botAPI   *tgbotapi.BotAPI
    Updates  *tgbotapi.UpdatesChannel
    database *sql.DB
    cfg      *config
}

// Start starts and runs the bot.
func Start() error {
    
    cfg, err := configure()
    if err != nil {
        return err 
    }
    
    botAPI, err := tgbotapi.NewBotAPI(cfg.apiKey) 
    if err != nil {
        return err
    }
    
    botAPI.Debug = cfg.debug
    botAPIUpdateConfig := tgbotapi.NewUpdate(0)
    botAPIUpdateConfig.Timeout = 60

    updates, err := botAPI.GetUpdatesChan(botAPIUpdateConfig)
    if err != nil {
        return err
    }
    log.Printf("%s authenticated", botAPI.Self.UserName)
    
    db, err := sql.Open("postgres", cfg.databaseURL)
    if err != nil {
        return err
    }
    defer db.Close()
    
    // Ping the database to check if the db connection is there.
    err = db.Ping()
    if err != nil {
        return err
    }
    log.Printf("Database connection established")
    
    var mybot bot
    mybot.botAPI   = botAPI
    mybot.Updates  = &updates
    mybot.database = db
    mybot.cfg      = &cfg
    
    for update := range *mybot.Updates {
        if update.Message == nil {
            continue
        }

        err = handleUpdate(&mybot, update)
        if err != nil {
            log.Panic(err)
        }
    }
    
    return nil
}

// handleUpdate processes an update from the channel provided by tgbotapi. 
func handleUpdate(jbot *bot, update tgbotapi.Update) (err error) {
    
    log.Printf("Recieved message: [%s %s %s] %s",
               strconv.Itoa(update.Message.From.ID), 
               update.Message.From.UserName, 
               update.Message.From.FirstName, 
               update.Message.Text)
    
    response, err := createResponse(jbot, update.Message.Text)
    if err != nil {
        return
    }
    
    sendMessage(jbot.botAPI, update.Message.Chat.ID, response)
    log.Printf("Message sent: %s", response)
    
    return
}

// horoscopeResponse contains data from
// a particular REST API json response 
type horoscopeResponse struct {
    Date      string `json:"date"`
    Sunsign   string `json:"sunsign"`
    Horoscope string `json:"horoscope"`
    Meta      horoscopeMeta `json:"meta"`
}

// horoscopeResponseMetaData contains data from 
// a particular REST API json response
type horoscopeMeta struct {
    Intensity string `json:"intensity"`
    Keywords  string `json:"keywords"`
    Mood      string `json:"mood"`
}

// horoscopeSign represents a particular horoscope sign.
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
    
    //set "out of range" to horoscopeSignNone
    if sign < 0 || sign > 12 {
        return ""
    }
    
    signs := [13]string{
        "",
        "aries",
        "taurus",
        "gemini",
        "cancer",
        "leo",
        "virgo",
        "libra",
        "scorpio",
        "sagittrius",
        "capricorn",
        "aquarius",
        "pisces",
        }
        
        return signs[sign]
}

// parseHoroscopeMessage searches originalMessage for certain
// key phrases and returns a corresponding horoscopeSign if one is found.
func parseHoroscopeMessage(originalMessage string) horoscopeSign {
    msg := strings.ToLower(originalMessage)
    if strings.Contains(msg,        "oina") || strings.Contains(msg, "aries")      {
        return horoscopeSignAries
    } else if strings.Contains(msg, "härk") || strings.Contains(msg, "taurus")     {
        return horoscopeSignTaurus
    } else if strings.Contains(msg, "kaks") || strings.Contains(msg, "gemini")     {
        return horoscopeSignGemini
    } else if strings.Contains(msg, "rap")  || strings.Contains(msg, "cancer")     {
        return horoscopeSignCancer
    } else if strings.Contains(msg, "leij") || strings.Contains(msg, "leo")        {
        return horoscopeSignLeo
    } else if strings.Contains(msg, "neit") || strings.Contains(msg, "virgo")      {
        return horoscopeSignVirgo
    } else if strings.Contains(msg, "vaa")  || strings.Contains(msg, "libra")      {
        return horoscopeSignLibra
    } else if strings.Contains(msg, "skor") || strings.Contains(msg, "scorpio")    {
        return horoscopeSignScorpio
    } else if strings.Contains(msg, "jous") || strings.Contains(msg, "sagittrius") {
        return horoscopeSignSagittarius
    } else if strings.Contains(msg, "vesi") || strings.Contains(msg, "aquarius")   {
        return horoscopeSignCapricorn
    } else if strings.Contains(msg, "kaur") || strings.Contains(msg, "capricorn")  {
        return horoscopeSignAquarius
    } else if strings.Contains(msg, "kal")  || strings.Contains(msg, "pisces")     {
        return horoscopeSignPisces
    } else {
        return horoscopeSignNone
    }
}

// resolveHoroscope provides a message string to send based on a horoscopeSign
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
    
    var hresponse horoscopeResponse
    err = json.Unmarshal(bodyBytes, &hresponse)
    if err != nil {
        return
    }
    
    reply = horoscopeReply(hresponse)
    return
}

// horoscopeReply builds a reply string from a horoscopeResponse
func horoscopeReply(hresponse horoscopeResponse) (reply string) {
    reply = "The Angels transfer your horoscope:\n👼👼👼\n" +
    hresponse.Horoscope + "\n👼👼 👼 \n\nKeywords: " +
    hresponse.Meta.Keywords + "\n\nMood: " +
    hresponse.Meta.Mood  + "\n\nEnergy level of transfer: " +
    hresponse.Meta.Intensity + "."
    
    return
}

// sendMessage sends message to the chat specified by chatID.
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
    msg := tgbotapi.NewMessage(chatID, message)
    bot.Send(msg)    
}

// createResponse generates a response string based on the recieved message string.
func createResponse(jbot *bot, message string) (response string, err error) {
    messageLower := strings.ToLower(message)
    
    for _, alias := range jbot.cfg.commands.hello.alias {
        if strings.HasPrefix(messageLower, alias){
            response = jbot.cfg.commands.hello.reply
            return
        }
    }
    for _, alias := range jbot.cfg.commands.start.alias {
        if strings.HasPrefix(messageLower, alias){
            response = jbot.cfg.commands.start.reply
            return
        }
    }
    for _, alias := range jbot.cfg.commands.wisdom.alias {
        if strings.HasPrefix(messageLower, alias){
            response = createBookResposeString(jbot, message)
            return
        }
    }
    for _, alias := range jbot.cfg.commands.horoscope.alias {
        if strings.HasPrefix(messageLower, alias){
            sign := parseHoroscopeMessage(message)
            if sign == horoscopeSignNone {
                response = horoscopeSignNone.String()
                return
            }
            response, err = resolveHoroscope(sign)
            if err != nil {
                // If getting the horoscope fails, log the error and move on.
                log.Println(err)
                return "", nil
            }
            return
        }
    }
    
    response = ""
    return
}

// createBookResposeString creates a string containing the appropriate
// response to a bookline related command. 
func createBookResposeString(jbot *bot, message string) string {
    words := strings.Split(message, " ")
    if len(words) >= 3 {
        
        line, _ := getBookLine(jbot.database, strings.ToLower(words[1]), words[2])
        if line != "" {
            return line
        }
    }
    
    response := ""
    response, _ = getRandomBookLine(jbot.database)
    return response
}

// createStartMessage generates a reply string for the command start.
func createStartMessage() string {
    return "Greetings traveler!\n\n" + 
    "This bot supports two commands:\n" + 
    "/horoscope SIGN\nSIGN is your horoscope sign (e.g. aries).\n" + 
    "/wisdom\nThis provides wisdom for you.\n\n" + 
    "An advanced user can request a particular set of wise words " + 
    "by specifying a chapter and a verse (e.g. 1Moos 1:1)\n\n" + 
    "Komennot myös suomeksi:\n" +
    "/horoskooppi vesimies\n" + 
    "/raamatturivi 1Moos 1:1"
}
