// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	gjson "github.com/tidwall/gjson"
)

type pingpong struct {
	features []pingpongFeature
}

type pingpongFeature struct {
	Pings              []string `json:"pings"`              // list of words to trigger the command
	IsPrefixCommand    bool     `json:"isprefixcommand"`    // if true, command triggers only if it is a prefix
	IsReply            bool     `json:"isreply"`            // if true, the telegram message is replying to the command (replying is a feature in telegram)
	Pongs              []string `json:"pongs"`              // list of possible answers to command, random one will be sent
	SuccessPropability float64  `json:"successpropability"` // 0.0-1.0 propability, used to make the command randomly fail
}

func (p *pingpong) String() string {
	return "pingpong"
}

func (p *pingpong) init(bot *jbot) error {

	jsonConfig := gjson.GetBytes(bot.cfg.Features, "pingpong")
	if !jsonConfig.Exists() {
		return errors.New("missing configs")
	}

	var err error
	ppFeat := new(pingpongFeature)
	for _, jsonWord := range jsonConfig.Array() {
		ppFeat = new(pingpongFeature)
		err = json.Unmarshal([]byte(jsonWord.Raw), ppFeat)
		if err != nil {
			log.Println(err)
		} else {
			p.features = append(p.features, *ppFeat)
		}
	}

	return nil
}

func (p *pingpong) triggers(bot *jbot, u tgbotapi.Update) bool {
	// any message will trigger
	return u.Message != nil
}

func (p *pingpong) execute(bot *jbot, u tgbotapi.Update) error {

	for _, feat := range p.features {

		toSend := findPingpongReply(strings.ToLower(u.Message.Text), feat)
		if toSend != "" {

			msg := tgbotapi.NewMessage(u.Message.Chat.ID, toSend)
			if feat.IsReply {
				msg.ReplyToMessageID = u.Message.MessageID
			}
			bot.botAPI.Send(msg)

		}

	}
	return nil
}

func findPingpongReply(text string, feature pingpongFeature) string {
	answer := ""

	for _, keyword := range feature.Pings {

		if feature.IsPrefixCommand {
			if strings.HasPrefix(text, keyword) {
				answer = feature.Pongs[rand.Intn(len(feature.Pongs))]
				break
			}
		} else if strings.Contains(text, keyword) {
			answer = feature.Pongs[rand.Intn(len(feature.Pongs))]
			break
		}
	}

	if answer == "" {
		return ""
	}

	success := true
	if feature.SuccessPropability > 0 && feature.SuccessPropability < 1 {
		if rand.Float64() > feature.SuccessPropability {
			success = false
		}
	}

	if !success {
		return ""
	}
	return answer
}
