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
// certain words are never picked.
// Other words are biased to be picked more often.
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
	return stringHasAnyPrefix(u.Message.Text, d.triggerWords)
}

// execute sends the chosen option back to the user
func (d *decide) execute(bot *jbot, u tgbotapi.Update) error {
	message := u.Message.Text

	// compress whitespace to single spaces
	spaceRegexp := regexp.MustCompile(`\s+`)
	trimmedMessage := spaceRegexp.ReplaceAllString(message, " ")

	// split incoming message string to words
	inputWords := strings.Split(trimmedMessage, " ")
	if len(inputWords) < 3 {
		return nil
	}

	// remove the command word (e.g. !decide)
	inputWords = inputWords[1:]

	// maps lowercase inputs to original inputs
	originalInputs := make(map[string]string)

	// inputWords to lower case
	for i, inputword := range inputWords {
		inputWords[i] = strings.ToLower(inputword)
		originalInputs[inputWords[i]] = inputword
	}

	// filter out common filler words (don't want to return e.g. 'or')
	WordsToSkip := []string{"or", "vai", "tai", "vaiko"}
	inputWords = filterWords(inputWords, WordsToSkip)
	if len(inputWords) < 2 {
		return nil
	}

	// double the chance for drinking realted words to get chosen
	preferredWords := []string{"kalja", "beer", "olut", "bisse", "kaljaa", "viina"}
	inputWords = duplicateWords(inputWords, preferredWords)

	chosenWord := inputWords[rand.Intn(len(inputWords))]
	chosenWord = originalInputs[chosenWord]
	msg := tgbotapi.NewMessage(u.Message.Chat.ID, chosenWord)
	bot.botAPI.Send(msg)
	return nil
}

func filterWords(words, wordsToRemove []string) []string {
	result := []string{}
	for _, word := range words {
		skipThisWord := false
		for _, badWord := range wordsToRemove {
			if word == badWord {
				skipThisWord = true
				break
			}
		}
		if !skipThisWord {
			result = append(result, word)
		}
	}
	return result
}

func duplicateWords(words, wordsToDuplicate []string) []string {
	wordsLen := len(words)
	for i := 0; i < wordsLen; i++ {
		word := words[i]
		for _, goodWord := range wordsToDuplicate {
			if word == goodWord {
				words = append(words, goodWord)
				break
			}
		}
	}
	return words
}
