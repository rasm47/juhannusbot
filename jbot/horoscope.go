// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"database/sql"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type horoscope struct {
}

func (h horoscope) init(bot *jbot) error {

	if err := bot.database.Ping(); err != nil {
		return err
	}

	var exists bool
	err := bot.database.QueryRow("SELECT EXISTS (SELECT * FROM horoscope)").Scan(&exists)
	if err != nil {
		return err
	}

	if !exists {
		return nil // TODO: make new error
	}

	return nil

}

func (h horoscope) triggers(bot *jbot, u tgbotapi.Update) bool {
	if u.Message != nil {
		for _, triggerWord := range bot.cfg.CommandConfigs["horoscope"].Aliases {
			if strings.HasPrefix(u.Message.Text, triggerWord) {
				return true
			}
		}
	} else if u.CallbackQuery != nil {
		return true
	}

	return false
}

func (h horoscope) execute(bot *jbot, u tgbotapi.Update) error {

	text := ""

	if u.CallbackQuery != nil {

		text, err := resolveHoroscope(convertEmojiToHoroscopeSign(u.CallbackQuery.Data), bot.database)
		if err != nil {
			return err
		}

		bot.botAPI.AnswerCallbackQuery(tgbotapi.NewCallback(u.CallbackQuery.ID, "Fortune delivered"))
		bot.botAPI.Send(tgbotapi.NewMessage(u.CallbackQuery.Message.Chat.ID, text))
		return nil

	}

	chatID := u.Message.Chat.ID
	sign := parseHoroscopeMessage(u.Message.Text)

	if sign == horoscopeSignNone {
		text = "Try a button"
		msg := tgbotapi.NewMessage(chatID, text)
		msg.ReplyMarkup = getSignKeyboard()
		bot.botAPI.Send(msg)
	} else {
		text, err := resolveHoroscope(sign, bot.database)
		if err != nil {
			text = "Horoscope failed"
		}

		msg := tgbotapi.NewMessage(chatID, text)
		bot.botAPI.Send(msg)
	}
	return nil

}

// horoscopeData contains data from
// a particular APIs json response
type horoscopeData struct {
	Date    string        `json:"date"`
	Sunsign string        `json:"sunsign"`
	Text    string        `json:"horoscope"`
	Meta    horoscopeMeta `json:"meta"`
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
	horoscopeSignNone horoscopeSign = iota
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

// horoscopeReply builds a reply string from horoscopeData
func horoscopeReply(hresponse horoscopeData) (reply string) {

	reply = "The Angels transfer your horoscope:\n👼👼👼\n" +
		hresponse.Text + "\n👼👼 👼 \n\nKeywords: " +
		hresponse.Meta.Keywords + "\n\nMood: " +
		hresponse.Meta.Mood + "\n\nEnergy level of transfer: " +
		hresponse.Meta.Intensity + "."

	return
}

// getSignKeyboard returns an inline keyboard with buttons for
// all horoscope signs.
func getSignKeyboard() tgbotapi.InlineKeyboardMarkup {

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("♒", "♒"),
			tgbotapi.NewInlineKeyboardButtonData("♓", "♓"),
			tgbotapi.NewInlineKeyboardButtonData("♈", "♈"),
			tgbotapi.NewInlineKeyboardButtonData("♉", "♉"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("♊", "♊"),
			tgbotapi.NewInlineKeyboardButtonData("♋", "♋"),
			tgbotapi.NewInlineKeyboardButtonData("♌", "♌"),
			tgbotapi.NewInlineKeyboardButtonData("♍", "♍"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("♎", "♎"),
			tgbotapi.NewInlineKeyboardButtonData("♏", "♏"),
			tgbotapi.NewInlineKeyboardButtonData("♐", "♐"),
			tgbotapi.NewInlineKeyboardButtonData("♑", "♑"),
		),
	)

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
	if strings.Contains(msg, "oina") || strings.Contains(msg, "aries") {
		return horoscopeSignAries
	} else if strings.Contains(msg, "härk") || strings.Contains(msg, "taurus") {
		return horoscopeSignTaurus
	} else if strings.Contains(msg, "kaks") || strings.Contains(msg, "gemini") {
		return horoscopeSignGemini
	} else if strings.Contains(msg, "rap") || strings.Contains(msg, "cancer") {
		return horoscopeSignCancer
	} else if strings.Contains(msg, "leij") || strings.Contains(msg, "leo") {
		return horoscopeSignLeo
	} else if strings.Contains(msg, "neit") || strings.Contains(msg, "virgo") {
		return horoscopeSignVirgo
	} else if strings.Contains(msg, "vaa") || strings.Contains(msg, "libra") {
		return horoscopeSignLibra
	} else if strings.Contains(msg, "skor") || strings.Contains(msg, "scorpio") {
		return horoscopeSignScorpio
	} else if strings.Contains(msg, "jous") || strings.Contains(msg, "sagittrius") {
		return horoscopeSignSagittarius
	} else if strings.Contains(msg, "vesi") || strings.Contains(msg, "aquarius") {
		return horoscopeSignCapricorn
	} else if strings.Contains(msg, "kaur") || strings.Contains(msg, "capricorn") {
		return horoscopeSignAquarius
	} else if strings.Contains(msg, "kal") || strings.Contains(msg, "pisces") {
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
