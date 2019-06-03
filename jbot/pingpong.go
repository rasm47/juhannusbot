// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type pingpong struct {
}

func (p pingpong) String() string {
	return "pingpong"
}

func (p pingpong) init(_ *jbot) error {
	return nil
}

func (p pingpong) triggers(bot *jbot, u tgbotapi.Update) bool {
	// any message will trigger
	return u.Message != nil
}

func (p pingpong) execute(bot *jbot, u tgbotapi.Update) error {

	for _, command := range bot.cfg.CommandConfigs {
		if command.Type == "message" {
			toSend := findPingpongReply(strings.ToLower(u.Message.Text), command)
			if toSend != "" {
				msg := tgbotapi.NewMessage(u.Message.Chat.ID, toSend)
				bot.botAPI.Send(msg)
			}
		}
	}
	return nil
}

func findPingpongReply(text string, command commandConfig) string {
	answer := ""
	for _, keyword := range command.Aliases {
		if command.IsPrefixCommand {
			if strings.HasPrefix(text, keyword) {
				answer = command.ReplyMessages[rand.Intn(len(command.ReplyMessages))]
			}
		} else if strings.Contains(text, keyword) {
			answer = command.ReplyMessages[rand.Intn(len(command.ReplyMessages))]
		}
	}

	if answer == "" {
		return ""
	}

	success := true
	if command.SuccessPropability > 0 && command.SuccessPropability < 1 {
		if rand.Float64() > command.SuccessPropability {
			success = false
		}
	}

	if !success {
		return ""
	}
	return answer
}
