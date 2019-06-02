// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type pingpong struct {
}

func (p pingpong) init(_ *jbot) error {
	return nil
}

func (p pingpong) triggers(bot *jbot, u tgbotapi.Update) bool {
	// any message will trigger
	return u.Message != nil
}

func (p pingpong) execute(bot *jbot, u tgbotapi.Update) error {

	text := findPingpongReply(strings.ToLower(u.Message.Text), bot.cfg)

	if text == "" {
		return nil
	}

	msg := tgbotapi.NewMessage(u.Message.Chat.ID, text)
	bot.botAPI.Send(msg)

	return nil
}

func findPingpongReply(message string, conf *config) string {

	for _, command := range conf.CommandConfigs {
		if command.Type == "message" {
			for _, keyword := range command.Aliases {
				if command.IsPrefixCommand {
					if strings.HasPrefix(message, keyword) {
						return command.ReplyMessages[rand.Intn(len(command.ReplyMessages))]
					}
				} else if strings.Contains(message, keyword) {
					return command.ReplyMessages[rand.Intn(len(command.ReplyMessages))]
				}
			}
		}
	}
	return ""
}
