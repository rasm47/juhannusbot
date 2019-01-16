# juhannusbot
A Telegram bot for entertainment purposes.

# How to deploy
Get your Telegram bot API key from the botfather. Put your key in config.txt. You need to create this file or rename the config_example.txt.

You will need a book.txt. You can write whatever wisdom you want to in there. A bot command will later randomly pull lines from this book. A few lines of text is required.

Run go install and then the resulting program.

# Commands
The bot has three commands. 

"/hello" simply answers "world". 

"/horoskooppi *s*" gives the horoscope of the day for the sign *s*. A Finnish sign is expected. For example: /horoskooppi oinas

"/raamatturivi" gives a random line from your book.txt.
