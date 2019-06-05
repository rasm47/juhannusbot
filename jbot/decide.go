// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"errors"
	"math/rand"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gjson "github.com/tidwall/gjson"
)

// decide is a feature of jbot
// it responds to "keyword option1 option2 ..."
// an option is randomly chosen.
// certain words are biased to not get picked.
// Other words are biased to be picked evry time.
type decide struct {
	triggerWords []string
}

func (d *decide) String() string {
	return "decide"
}

func (d *decide) init(bot *jbot) error {

	jsonConfig := gjson.GetBytes(bot.cfg.Features, "decide.aliases")
	if !jsonConfig.Exists() {
		return errors.New("missing configs")
	}

	for _, jsonWord := range jsonConfig.Array() {
		d.triggerWords = append(d.triggerWords, jsonWord.String())
	}

	return nil
}

// triggers when one of the configured keywords is seen
// as a prefix of a message seen by the bot
func (d *decide) triggers(bot *jbot, u tgbotapi.Update) bool {
	if u.Message == nil {
		return false
	}

	return stringHasAnyPrefix(u.Message.Text, d.triggerWords)
}

// execute sends the chosen option back to the user
func (d *decide) execute(bot *jbot, u tgbotapi.Update) error {
	message := u.Message.Text
	spaceRegexp := regexp.MustCompile(`\s+`)
	trimmedMessage := spaceRegexp.ReplaceAllString(message, " ")
	inputWords := strings.Split(trimmedMessage, " ")

	// remove the command word (e.g. !decide)
	inputWords = inputWords[1:]
	if len(inputWords) < 2 {
		return nil
	}

	skippedWords := []string{"or", "vai", "tai", "vaiko"}
	preferredWords := []string{"kalja", "beer", "olut", "bisse", "kaljaa"}
	outputWords := []string{}

	var lowercaseWord string
	for _, inputWord := range inputWords {
		lowercaseWord = strings.ToLower(inputWord)

		for _, preferredWord := range preferredWords {
			if lowercaseWord == preferredWord {
				msg := tgbotapi.NewMessage(u.Message.Chat.ID, inputWord)
				bot.botAPI.Send(msg)
				return nil
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

	if len(outputWords) == 0 {
		return nil
	}

	chosenWord := outputWords[rand.Intn(len(outputWords))]
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, chosenWord)
	bot.botAPI.Send(msg)
	return nil
}
