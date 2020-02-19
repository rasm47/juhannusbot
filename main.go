package main

import (
	"log"

	"github.com/rasm47/juhannusbot/jbot"
)

func main() {

	err := jbot.Start()
	if err != nil {
		log.Printf("Closing bot due to error: %v", err)
	}

	return
}
