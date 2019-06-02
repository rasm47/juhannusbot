// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type wisdom struct {
}

func (w wisdom) init(bot *jbot) error {

	if err := bot.database.Ping(); err != nil {
		return err
	}

	var tableExists bool
	err := bot.database.QueryRow("SELECT EXISTS (SELECT * FROM book)").Scan(&tableExists)
	if err != nil {
		return err
	}
	if !tableExists {
		return nil // todo make errors
	}
	return nil
}

func (w wisdom) triggers(bot *jbot, u tgbotapi.Update) bool {
	if u.Message != nil {
		for _, triggerWord := range bot.cfg.CommandConfigs["wisdom"].Aliases {
			if strings.HasPrefix(u.Message.Text, triggerWord) {
				return true
			}
		}
	}
	return false
}

func (w wisdom) execute(bot *jbot, u tgbotapi.Update) error {

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
