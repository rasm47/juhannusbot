// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"errors"
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gjson "github.com/tidwall/gjson"
)

type wisdom struct {
	triggerWords []string
}

func (w *wisdom) String() string {
	return "wisdom"
}

func (w *wisdom) init(bot *jbot) error {

	if !connected(bot.database) {
		return errors.New("no database connection")
	}

	var tableExists bool
	err := bot.database.QueryRow("SELECT EXISTS (SELECT * FROM book)").Scan(&tableExists)
	if err != nil {
		return fmt.Errorf("database error: %v", err)
	}
	if !tableExists {
		return errors.New("table book missing from database")
	}

	jsonConfig := gjson.GetBytes(bot.cfg.Features, "wisdom.aliases")
	if !jsonConfig.Exists() {
		return errors.New("missing configs")
	}

	for _, jsonWord := range jsonConfig.Array() {
		w.triggerWords = append(w.triggerWords, jsonWord.String())
	}

	return nil
}

func (w *wisdom) triggers(bot *jbot, u tgbotapi.Update) bool {
	if u.Message == nil {
		return false
	}

	return stringHasAnyPrefix(u.Message.Text, w.triggerWords)
}

func (w *wisdom) execute(bot *jbot, u tgbotapi.Update) error {

	text := createBookResposeString(bot, u.Message.Text)
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	bot.botAPI.Send(msg)

	return nil
}

// createBookResposeString creates a string containing the appropriate
// response to a bookline related command.
func createBookResposeString(bot *jbot, message string) string {
	words := strings.Split(message, " ")
	if len(words) >= 3 {

		line, _ := getBookLine(bot.database, strings.Replace(strings.ToLower(words[1]), ".", "", -1), words[2])
		if line != "" {
			return line
		}
	}

	response := ""
	response, _ = getRandomBookLine(bot.database) // TODO: error is ignored here
	return response
}
