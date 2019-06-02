// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"math/rand"
	"regexp"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// decide is a feature of jbot
// it responds to "keyword option1 option2 ..."
// an option is randomly chosen.
// certain words are biased to not get picked.
// Other words are biased to be picked evry time.
type decide struct {
}

func (d decide) String() string {
	return "decide"
}

// init for decide always works
func (d decide) init(_ *jbot) error {
	return nil
}

// triggers when one of the configured keywords is seen
// as a prefix of a message seen by the bot
func (d decide) triggers(bot *jbot, u tgbotapi.Update) bool {
	if u.Message == nil {
		return false
	}

	triggeringPrefixes := bot.cfg.CommandConfigs["decide"].Aliases
	return stringHasAnyPrefix(u.Message.Text, triggeringPrefixes)
}

// execute sends the somewhat randomly chosen option back to the user
func (d decide) execute(bot *jbot, u tgbotapi.Update) error {
	message := u.Message.Text
	spaceRegexp := regexp.MustCompile(`\s+`)
	trimmedMessage := spaceRegexp.ReplaceAllString(message, " ")
	inputWords := strings.Split(trimmedMessage, " ")

	// remove the command word (e.g. !decide)
	inputWords = inputWords[1:]

	skippedWords := []string{"or", "vai", "tai", "vaiko"}
	preferredWords := []string{"kalja", "beer", "olut", "bisse", "kaljaa"}
	outputWords := []string{}

	if len(inputWords) < 2 {
		return nil
	}

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

	var chosenWord string
	if len(outputWords) == 0 {
		chosenWord = ""
	} else {
		chosenWord = outputWords[rand.Intn(len(outputWords))]
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, chosenWord)

	bot.botAPI.Send(msg)
	return nil
}
