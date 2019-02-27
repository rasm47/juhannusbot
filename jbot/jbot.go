// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
    "log"
    "strings"
    "database/sql"
    
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

// botAction represents something the bot can do.
type botAction int
const (
    botActionNone                  botAction = 0
    botActionSendMessage           botAction = 1
    botActionCallbackReply         botAction = 2
    botActionSendHoroscopeKeyboard botAction = 3
)

// botInstruction tells the bot what to do and contains
// all of the necessary data for performing that botAction.
type botInstruction struct {
    Action botAction
    ChatID int64
    MessageID int64
    CallbackQueryID string
    Text string
}

// botCommand represents a supported comand for the bot.
type botCommand int
const (
    botCommandNone      botCommand = 0
    botCommandStart     botCommand = 1
    botCommandWisdom    botCommand = 2
    botCommandHoroscope botCommand = 3
)

// horoscopeData contains data from
// a particular APIs json response 
type horoscopeData struct {
    Date      string `json:"date"`
    Sunsign   string `json:"sunsign"`
    Text      string `json:"horoscope"`
    Meta      horoscopeMeta `json:"meta"`
}

// horoscopeMeta contains data from 
// a particular APIs json response
type horoscopeMeta struct {
    Intensity string `json:"intensity"`
    Keywords  string `json:"keywords"`
    Mood      string `json:"mood"`
}

// horoscopeSign represents a particular horoscope sign.
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
        "sagittarius",
        "capricorn",
        "aquarius",
        "pisces",
        }
        
        return signs[sign]
}

// Start starts and runs the bot.
func Start() error {
    
    cfg, err := configure()
    if err != nil {
        return err 
    }
    
    botAPI, err := tgbotapi.NewBotAPI(cfg.APIKey) 
    if err != nil {
        return err
    }
    
    botAPI.Debug = cfg.Debug
    botAPIUpdateConfig := tgbotapi.NewUpdate(0)
    botAPIUpdateConfig.Timeout = 60

    updates, err := botAPI.GetUpdatesChan(botAPIUpdateConfig)
    if err != nil {
        return err
    }
    log.Printf("%s authenticated", botAPI.Self.UserName)
    
    db, err := sql.Open("postgres", cfg.DatabaseURL)
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
    
    startHoroscopeUpdater(db)
    
    mybot := bot{botAPI,&updates,db,&cfg}
    for update := range *mybot.Updates {

        err = handleUpdate(&mybot, update)
        if err != nil {
            log.Panic(err)
        }
    }
    
    return nil
}

// handleUpdate processes an update from the channel provided by tgbotapi. 
func handleUpdate(jbot *bot, update tgbotapi.Update) error {
    
    instruction, err := newBotInstruction(jbot, update)
    if err != nil {
        return err
    }
       
    executeInstruction(jbot,instruction)
    return nil
}

// newBotInstruction creates the appropriate botInstruction based
// on the data contained in update.
func newBotInstruction(jbot *bot, update tgbotapi.Update) (bi botInstruction, err error) {
    
    // A non-nil CallbackQuery means that someone pressed a button on
    // the inline keybord.
    if update.CallbackQuery != nil {
        bi.Action = botActionCallbackReply
        bi.ChatID = update.CallbackQuery.Message.Chat.ID
        bi.CallbackQueryID = update.CallbackQuery.ID
        bi.Text, err = resolveHoroscope(convertEmojiToHoroscopeSign(update.CallbackQuery.Data),jbot.database)
        
    // A non-nil Message means the bot recieved a message
    } else if update.Message != nil {
        
        bi.ChatID = update.Message.Chat.ID
        
        switch command := newCommand(jbot.cfg.CommandConfigs, update.Message.Text); command {
        case botCommandStart:
            bi.Action = botActionSendMessage
            bi.Text = jbot.cfg.CommandConfigs.Start.Reply
        case botCommandWisdom:
            bi.Action = botActionSendMessage
            bi.Text = createBookResposeString(jbot, update.Message.Text)
        case botCommandHoroscope:
            sign := parseHoroscopeMessage(update.Message.Text)
            
            if sign == horoscopeSignNone {
                bi.Action = botActionSendHoroscopeKeyboard
                bi.Text = "Try a button"
            } else {
                bi.Action = botActionSendMessage
                messageToSend, err := resolveHoroscope(sign, jbot.database)
                if err != nil {
                    bi.Text = "Horoscope failed"
                } else {
                    bi.Text = messageToSend
                }
            }
        default:
            bi.Action = botActionNone
        }
    }
    return
}

// newCommand searches if message contains any of 
// the command aliases from the commandConfigs and
// returns a corresponding botCommand.
func newCommand(commandConfigs commandConfigList, message string) botCommand {
    
    messageLower := strings.ToLower(message)
    
    for _, alias := range commandConfigs.Start.Alias {
        if strings.HasPrefix(messageLower, alias) {
            return botCommandStart
        }
    }
    
    for _, alias := range commandConfigs.Wisdom.Alias {
        if strings.HasPrefix(messageLower, alias) {
            return botCommandWisdom
        }
    }
    
    for _, alias := range commandConfigs.Horoscope.Alias {
        if strings.HasPrefix(messageLower, alias) {
            return botCommandHoroscope
        }
    }
    
    return botCommandNone
}

// executeInstruction makes jbot act according to the instructions.
func executeInstruction(jbot *bot, instructions botInstruction) {
    
    switch instructions.Action {
        
        case botActionCallbackReply:
            jbot.botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(
                instructions.CallbackQueryID, "Fortune delivered"))
            jbot.botAPI.Send(tgbotapi.NewMessage(
                instructions.ChatID, instructions.Text))
    
        case botActionSendMessage:
            jbot.botAPI.Send(tgbotapi.NewMessage(instructions.ChatID, instructions.Text))
            
        case botActionSendHoroscopeKeyboard:
            msg := tgbotapi.NewMessage(instructions.ChatID, instructions.Text)
            msg.ReplyMarkup = getSignKeyboard()
            jbot.botAPI.Send(msg)
            
        default:
            return
    }
    return
}

// convertEmojiToHoroscopeSign matches the string emoji
// to the horoscope emojis and returns a horoscopeSign
// that matches that emoji. Returns horoscopeSignNone if
// no match was found.
func convertEmojiToHoroscopeSign(emoji string) (sign horoscopeSign) {
    switch emoji {
        case "♒": 
            sign = horoscopeSignAquarius
        case "♓": 
            sign = horoscopeSignPisces
        case "♈": 
            sign = horoscopeSignAries
        case "♉": 
            sign = horoscopeSignTaurus
        case "♊": 
            sign = horoscopeSignGemini
        case "♋": 
            sign = horoscopeSignCancer
        case "♌": 
            sign = horoscopeSignLeo
        case "♍": 
            sign = horoscopeSignVirgo
        case "♎": 
            sign = horoscopeSignLibra
        case "♏": 
            sign = horoscopeSignScorpio
        case "♐": 
            sign = horoscopeSignSagittarius
        case "♑": 
            sign = horoscopeSignCapricorn
        default:
            sign = horoscopeSignNone
    }
    return sign
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

// resolveHoroscope provides a string to send to the user
// based on a horoscopeSign.
func resolveHoroscope(sign horoscopeSign, database *sql.DB) (reply string, err error) {
    
    hresponse := getHoroscopeData(database, sign)
    reply = horoscopeReply(hresponse)
    return
}

// horoscopeReply builds a reply string from horoscopeData
func horoscopeReply(hresponse horoscopeData) (reply string) {
    
    reply = "The Angels transfer your horoscope:\n👼👼👼\n" +
    hresponse.Text + "\n👼👼 👼 \n\nKeywords: " +
    hresponse.Meta.Keywords + "\n\nMood: " +
    hresponse.Meta.Mood  + "\n\nEnergy level of transfer: " +
    hresponse.Meta.Intensity + "."
    
    return
}

// createBookResposeString creates a string containing the appropriate
// response to a bookline related command. 
func createBookResposeString(jbot *bot, message string) string {
    words := strings.Split(message, " ")
    if len(words) >= 3 {
        
        line, _ := getBookLine(jbot.database, strings.Replace(strings.ToLower(words[1]), ".", "", -1), words[2])
        if line != "" {
            return line
        }
    }
    
    response := ""
    response, _ = getRandomBookLine(jbot.database)
    return response
}

// getSignKeyboard returns an inline keyboard with buttons for
// all horoscope signs.
func getSignKeyboard() tgbotapi.InlineKeyboardMarkup {

    return tgbotapi.NewInlineKeyboardMarkup(
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("♒","♒"),
            tgbotapi.NewInlineKeyboardButtonData("♓","♓"),
            tgbotapi.NewInlineKeyboardButtonData("♈","♈"),
            tgbotapi.NewInlineKeyboardButtonData("♉","♉"),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("♊","♊"),
            tgbotapi.NewInlineKeyboardButtonData("♋","♋"),
            tgbotapi.NewInlineKeyboardButtonData("♌","♌"),
            tgbotapi.NewInlineKeyboardButtonData("♍","♍"),
        ),
        tgbotapi.NewInlineKeyboardRow(
            tgbotapi.NewInlineKeyboardButtonData("♎","♎"),
            tgbotapi.NewInlineKeyboardButtonData("♏","♏"),
            tgbotapi.NewInlineKeyboardButtonData("♐","♐"),
            tgbotapi.NewInlineKeyboardButtonData("♑","♑"),
        ),
    )
    
}
