// Package jbot provides a telegram bot for entertainment purposes.
package jbot

import (
	"database/sql"
	"log"
	"math/rand"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/lib/pq" // blank import to use PostgreSQL
)

// bot is a collection of relevant pointers.
type jbot struct {
	botAPI   *tgbotapi.BotAPI
	database *sql.DB
	cfg      *config
}

// feature is an interface that all of the bots features must satisfy
type feature interface {
	init(*jbot) error
	triggers(*jbot, tgbotapi.Update) bool
	execute(*jbot, tgbotapi.Update) error
	String() string
}

// Start starts and runs the bot.
func Start() error {

	cfg, err := configure()
	if err != nil {
		return err
	}

	botAPI, err := tgbotapi.NewBotAPI(cfg.APIKey)
	if err != nil {
		return err
	}
	log.Printf("Telegram botAPI authenticated for %v", botAPI.Self.UserName)

	botAPI.Debug = cfg.Debug
	botAPIUpdateConfig := tgbotapi.NewUpdate(0)
	botAPIUpdateConfig.Timeout = 60

	updates, err := botAPI.GetUpdatesChan(botAPIUpdateConfig)
	if err != nil {
		return err
	}

	db, err := sql.Open("postgres", cfg.DatabaseURL)
	if err != nil {
		return err
	}
	defer db.Close()

	if connected(db) {
		log.Println("connected to database")
	} else {
		log.Println("no database connection")
	}

	rand.Seed(time.Now().UnixNano())

	mybot := jbot{botAPI, db, &cfg}

	allFeatures := []feature{
		new(decide),
		new(pingpong),
		new(horoscope),
		new(wisdom),
	}
	features := []feature{}

	for _, feat := range allFeatures {
		if err = feat.init(&mybot); err != nil {
			log.Printf("not running %v: %v", feat.String(), err)
		} else {
			features = append(features, feat)
			log.Printf("running %v", feat.String())
		}
	}

	for update := range updates {
		for _, feat := range features {
			if feat.triggers(&mybot, update) {
				feat.execute(&mybot, update)
			}
		}
	}

	return nil
}
