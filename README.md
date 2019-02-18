# juhannusbot
A Telegram bot for entertainment purposes.

# Commands
The bot has three commands. 

`/hello` simply answers "world". 

`/horoskooppi *s*` gives the horoscope of the day for the sign *s*. For example: `/horoskooppi aries`.

`/raamatturivi` gives a random line from a book in your database. If you give additional parameters e.g. `/raamatturivi Ch1 Verse2`, you can request your favorite lines.

# How to deploy
To deploy this bot, three things are required:
* A distribution of the [go porgramming language](https://golang.org/doc/install)
* A [posgres](https://www.postgresql.org/download/) SQL database (preferably populated appropriately)
* A Telegram bot [API key](https://core.telegram.org/#bot-api)

The bot needs to be configured. This is done by creating a file called `config.json`.
An example of such a file (`example_config.json`) is provided, you can simply rename it to `config.json`.
The configuration file has three filds, one for the Telegram bot API key, one for debug mode (defaults to false) and one for a database connection string. 
When running the bot, have the config file in the current working directory.

Install and run with these commands:
* `go get` the dependencies.
* `go install` this bot.
* `./juhannusbot` in your go paths bin directory.

To stop the bot, use CTRL+C/CMD+C.

# Populating the database
The command `/raamatturivi` requires a connection to a SQL database with some book data.
The bot expects there to be a table called `book` with rows `chapter`, `verse` and `text`.

You can create this table with:
```sql
CREATE TABLE book (
chapter varchar(7),
verse   varchar(7),
text    varchar(4096)
);
```

If you want the `/raamatturivi` command to work, you need to [insert](https://www.postgresql.org/docs/11/tutorial-populate.html) a few rows that have a chapter, a verse and a text.
