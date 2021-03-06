// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"database/sql"
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

func (w *wisdom) triggers(u tgbotapi.Update) bool {
	if u.Message == nil {
		return false
	}

	return stringHasAnyPrefix(u.Message.Text, w.triggerWords)
}

func (w *wisdom) execute(bot *jbot, u tgbotapi.Update) error {

	text, err := createBookResposeString(bot, u.Message.Text)
	if err != nil {
		return err
	}
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	bot.botAPI.Send(msg)

	return nil
}

// createBookResposeString creates a string containing the appropriate
// response to a bookline related command.
func createBookResposeString(bot *jbot, message string) (string, error) {
	words := strings.Split(message, " ")
	if len(words) >= 3 {
		// try a specific line
		line, _ := getBookLine(bot.database, strings.Replace(strings.ToLower(words[1]), ".", "", -1), words[2])
		if line != "" {
			return line, nil
		}
	}

	response := ""
	response, err := getRandomBookLine(bot.database)
	if err != nil {
		return "", fmt.Errorf("database error: %v", err)
	}
	return response, nil
}

// getBookLine fetches a particular bookline from a database.
func getBookLine(database *sql.DB, chapter string, verse string) (string, error) {
	var text string
	err := database.QueryRow("SELECT text FROM book WHERE chapter = $1 and verse = $2", chapter, verse).Scan(&text)
	return text, err
}

// getBookLine fetches and formats a random bookline from a database.
func getRandomBookLine(database *sql.DB) (string, error) {

	var chapter, verse, text string
	rows, err := database.Query("SELECT chapter, verse, text FROM book ORDER BY RANDOM() LIMIT 1")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&chapter, &verse, &text)
		if err != nil {
			return "", err
		}
	}
	err = rows.Err()
	if err != nil {
		return "", err
	}

	return strings.ToUpper(chapter) + ". " + verse + " " + text, nil
}
