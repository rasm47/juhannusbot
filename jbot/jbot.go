// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
    "log"
    "regexp"
    "strings"
    "math/rand"
    "database/sql"
    
    _ "github.com/lib/pq" // blank import to use PostgreSQL
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
    botActionNone                  botAction = iota
    botActionSendMessage           
    botActionSendReplyMessage      
    botActionCallbackReply         
    botActionSendHoroscopeKeyboard 
)

// botInstruction tells the bot what to do and contains
// all of the necessary data for performing that botAction.
type botInstruction struct {
    Action botAction
    ChatID int64
    MessageID int
    CallbackQueryID string
    Text string
}

// botCommand represents a supported command type for the bot.
// Command types are none, message and special.
// Special is a specifically programmed command.
// Message is a general type of command that sends back a message.
type botCommand int
const (
    botCommandNone    botCommand = iota
    botCommandMessage 
    botCommandSpecial 
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
    horoscopeSignNone        horoscopeSign = iota
    horoscopeSignAries       
    horoscopeSignTaurus      
    horoscopeSignGemini     
    horoscopeSignCancer   
    horoscopeSignLeo       
    horoscopeSignVirgo       
    horoscopeSignLibra       
    horoscopeSignScorpio     
    horoscopeSignSagittarius 
    horoscopeSignCapricorn   
    horoscopeSignAquarius    
    horoscopeSignPisces      
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
    
    // log message if it causes an action
    if update.Message != nil && instruction.Action != botActionNone {
        if update.Message.Text != "" {
            log.Println(update.Message.From.FirstName + " " + update.Message.Text)
        }
    }
    
    executeInstruction(jbot,instruction)
    return nil
    
}

// newBotInstruction creates the appropriate botInstruction based
// on the data contained in update.
func newBotInstruction(jbot *bot, update tgbotapi.Update) (botInstruction, error) {
    
    // A non-nil CallbackQuery means that someone pressed a button on
    // the inline keybord.
    if update.CallbackQuery != nil {
        return newCallbackInstruction(jbot, update)
        
    // A non-nil Message means the bot recieved a message
    } else if update.Message != nil {
        
        commandType, commandName := findCommand(jbot.cfg.CommandConfigs, update.Message.Text)
        configs := jbot.cfg.CommandConfigs[commandName]
        
        switch commandType {
            case botCommandMessage:
                return newMessageInstruction(configs, update), nil
            case botCommandSpecial:
                return newSpecialInstruction(jbot, configs, commandName, update)
            default:
                return botInstruction{}, nil
        }
    }
    return botInstruction{}, nil
}

// newCallbackInstruction creates a new botInstruction from a callback query
func newCallbackInstruction(jbot *bot, update tgbotapi.Update) (bi botInstruction, err error) {
    bi.Action = botActionCallbackReply
    bi.ChatID = update.CallbackQuery.Message.Chat.ID
    bi.CallbackQueryID = update.CallbackQuery.ID
    bi.Text, err = resolveHoroscope(convertEmojiToHoroscopeSign(update.CallbackQuery.Data),jbot.database)
    return
}

// newMessageInstruction creates a new botInstruction for a message command
func newMessageInstruction(cc commandConfig, update tgbotapi.Update) (bi botInstruction) {
    
    // see if SuccessPropability is properly configured for this command
    if cc.SuccessPropability < 1.0 && 
        cc.SuccessPropability > 0.0 {
                    
        // see if the command fails
        if rand.Float64() > cc.SuccessPropability {
            // failed the random check
            bi.Action = botActionNone
            return
        }
    }
            
    if cc.IsReply {
        bi.Action = botActionSendReplyMessage
        bi.MessageID = update.Message.MessageID
    } else {
        bi.Action = botActionSendMessage
    }
    bi.Text = cc.ReplyMessages[rand.Intn(len(cc.ReplyMessages))]
    bi.ChatID = update.Message.Chat.ID
    return
}

// newSpecialInstruction creates a new botInstruction for a special command
func newSpecialInstruction(jbot *bot, cc commandConfig, name string, update tgbotapi.Update) (bi botInstruction, err error) {
    
    if name == "wisdom" {
        bi = newWisdomInstruction(jbot, update)
    } else if name == "horoscope" {
        bi, err = newHoroscopeInstruction(jbot, update)
    } else if name == "decide" {
        bi = newDecideInstruction(cc, update)
    } 
    return
}

// newWisdomInstruction creates a new botInstruction for the wisom command
func newWisdomInstruction(jbot *bot, update tgbotapi.Update) (bi botInstruction) {
    bi.ChatID = update.Message.Chat.ID
    bi.Action = botActionSendMessage
    bi.Text = createBookResposeString(jbot, update.Message.Text)
    return
}

// newHoroscopeInstruction creates a new botInstruction for the horoscope command
func newHoroscopeInstruction(jbot *bot, update tgbotapi.Update) (bi botInstruction, err error) {
    bi.ChatID = update.Message.Chat.ID
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
    return
}

// newDecideInstruction creates a new botInstruction for the decide command
func newDecideInstruction(cc commandConfig, update tgbotapi.Update) (bi botInstruction) {
    if cc.IsReply {
        bi.Action = botActionSendReplyMessage
        bi.MessageID = update.Message.MessageID
    } else {
        bi.Action = botActionSendMessage
    }
    bi.ChatID = update.Message.Chat.ID
    bi.Text = createDecideString(update.Message.Text)
    return
}

// findCommand searches if message contains any of 
// the command aliases from the commandConfigs and
// returns a corresponding botCommand and its name.
func findCommand(commandConfigs map[string]commandConfig, message string) (commandType botCommand, commandName string) {
    
    // force commands to be case insensitive 
    messageLower := strings.ToLower(message)
    
    for _, command := range commandConfigs {
        for _, alias := range command.Aliases {
            if strings.Contains(messageLower, alias) {
                if command.IsPrefixCommand && !strings.HasPrefix(messageLower, alias){
                    continue
                }
                if command.Type == "special" {
                    return botCommandSpecial, command.Name
                } else if command.Type == "message" {
                    return botCommandMessage, command.Name
                }
            }
        }
    }
    
    return botCommandNone, ""
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
            
    case botActionSendReplyMessage:
        msg := tgbotapi.NewMessage(instructions.ChatID, instructions.Text)
        msg.ReplyToMessageID = instructions.MessageID
        jbot.botAPI.Send(msg)
            
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

    emojiToHoroscopeMap := map[string]horoscopeSign{
        "♒": horoscopeSignAquarius,
        "♓": horoscopeSignPisces,
        "♈": horoscopeSignAries,
        "♉": horoscopeSignTaurus,
        "♊": horoscopeSignGemini,
        "♋": horoscopeSignCancer,
        "♌": horoscopeSignLeo,
        "♍": horoscopeSignVirgo,
        "♎": horoscopeSignLibra,
        "♏": horoscopeSignScorpio,
        "♐": horoscopeSignSagittarius,
        "♑": horoscopeSignCapricorn,
    }

    return emojiToHoroscopeMap[emoji]
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

// createDecideString creates a message to send for the decide command
func createDecideString(message string) string {
    
    spaceRegexp := regexp.MustCompile(`\s+`)
    trimmedMessage := spaceRegexp.ReplaceAllString(message, " ")
    inputWords := strings.Split(trimmedMessage, " ")

    // remove the command (e.g. !decide)
    inputWords = inputWords[1:]
    
    skippedWords := []string{"or", "vai", "tai", "vaiko"}
    preferredWords := []string{"kalja", "beer", "olut", "bisse", "kotiin"}
    outputWords := []string{}

    if len(inputWords) < 2 {
        return ""
    }
    
    var lowercaseWord string
    for _, inputWord := range inputWords {
        
        lowercaseWord = strings.ToLower(inputWord)
        
        for _, preferredWord := range preferredWords {
            if lowercaseWord == preferredWord {
                return inputWord
            }
        }

        skip := false
        for _, skippedWord := range skippedWords {
            if lowercaseWord == skippedWord {
                skip = true
            }
        }

        if !skip {
            outputWords = append(outputWords, inputWord)
        }
    }

    var chosenWord string
    if len(outputWords) == 0 {
        chosenWord = ""
    } else {
        chosenWord = outputWords[rand.Intn(len(outputWords))]
    }

    return chosenWord
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
