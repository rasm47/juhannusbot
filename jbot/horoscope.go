// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gjson "github.com/tidwall/gjson"
)

type horoscope struct {
	triggerWords []string
}

func (h *horoscope) String() string {
	return "horoscope"
}

func (h *horoscope) init(bot *jbot) error {

	if !connected(bot.database) {
		return errors.New("no database connection")
	}

	var tableExists bool
	err := bot.database.QueryRow("SELECT EXISTS (SELECT * FROM horoscope)").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	if !tableExists {
		return errors.New("table horoscope missing from database")
	}

	jsonConfig := gjson.GetBytes(bot.cfg.Features, "horoscope.aliases")
	if !jsonConfig.Exists() {
		return errors.New("missing configs")
	}

	for _, jsonWord := range jsonConfig.Array() {
		h.triggerWords = append(h.triggerWords, jsonWord.String())
	}

	return nil
}

func (h *horoscope) triggers(u tgbotapi.Update) bool {
	if u.Message != nil {
		return stringHasAnyPrefix(u.Message.Text, h.triggerWords)
	} else if u.CallbackQuery != nil {
		return true
	}

	return false
}

func (h *horoscope) execute(bot *jbot, u tgbotapi.Update) error {

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

// getHoroscopeData queries the database for the data of a particular sign
func getHoroscopeData(database *sql.DB, sign horoscopeSign) (data horoscopeData) {

	rows, err := database.Query("SELECT datestring, signstring, text, intensity, keywords, mood FROM horoscope WHERE signstring = $1", sign.String())
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&data.Date, &data.Sunsign, &data.Text, &data.Meta.Intensity, &data.Meta.Keywords, &data.Meta.Mood)
		if err != nil {
			return
		}
	}
	err = rows.Err()
	if err != nil {
		return
	}

	return
}

// startHoroscopeUpdater starts a process that updates all horoscopes
// in the database daily at roughly 04:00 / 4am
func startHoroscopeUpdater(database *sql.DB) {

	durationToNextFourAm := time.Duration(24+4-time.Now().Hour())*time.Hour +
		time.Duration(time.Now().Minute())*time.Minute +
		time.Duration(time.Now().Second())*time.Second

	time.AfterFunc(durationToNextFourAm, func() { updateHoroscopeDaily(database) })
}

// updateHoroscopeDaily starts a repeating goroutine
// that updates all horoscopes once every 24 hours.
func updateHoroscopeDaily(database *sql.DB) {

	log.Println("Horoscopes are updating for the first time: starting the daily updater")
	// call updateHoroscopeData one time
	updateAllHoroscopeData(database)

	// start an endless anonymous go routine of daily updating
	go func() {
		for range time.NewTicker(24 * time.Hour).C {
			log.Println("Attempting to fetch new horoscopes...")
			updateAllHoroscopeData(database)
		}
	}()

}

// updateAllHoroscopeData updates the database rows for
// all of the horoscopes
func updateAllHoroscopeData(database *sql.DB) {

	signs := [12]horoscopeSign{
		horoscopeSignAquarius,
		horoscopeSignPisces,
		horoscopeSignAries,
		horoscopeSignTaurus,
		horoscopeSignGemini,
		horoscopeSignCancer,
		horoscopeSignLeo,
		horoscopeSignVirgo,
		horoscopeSignLibra,
		horoscopeSignScorpio,
		horoscopeSignSagittarius,
		horoscopeSignCapricorn,
	}

	for _, sign := range signs {
		updateHoroscopeData(database, sign)
	}
	return
}

// updateHoroscopeData fetches the new horoscope of the day for a
// partucular horoscopeSign and updates that data to the database.
func updateHoroscopeData(database *sql.DB, sign horoscopeSign) {

	data, err := httpGetHoroscopeData(sign)
	if err != nil {
		log.Println("Failed to get new horoscopes from the web, database not updated")
		return
	}

	rows, err := database.Query("UPDATE horoscope SET (datestring, text, intensity, keywords, mood) = ($1, $2, $3, $4, $5) WHERE signstring = $6", data.Date, data.Text, data.Meta.Intensity, data.Meta.Keywords, data.Meta.Mood, strings.ToLower(data.Sunsign))
	if err != nil {
		log.Println("Error with the database, database not updated")
		return
	}
	defer rows.Close()

	log.Println("Updated data to the database,", sign.String())
	return
}

// httpGetHoroscopeData fetches the new horoscopeData for the day for
// a particular horoscopeSign. The data comes from a REST API whose
// url is hard coded inside this function.
func httpGetHoroscopeData(sign horoscopeSign) (data horoscopeData, err error) {

	response, err := http.Get("http://theastrologer-api.herokuapp.com/api/horoscope/" + sign.String() + "/today")
	if err != nil {
		return
	}
	defer response.Body.Close()

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(bodyBytes, &data)
	return
}
