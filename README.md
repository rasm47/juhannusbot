# juhannusbot
A Telegram bot for entertainment purposes.

# Features
The bot comes with four features: pingpong, wisdom, decide and horoscope. Adding your own features is also supported.

Pingpong is a feature that is triggered by a phrase called ping and responds with a phrase called pong. For example, when recieving `/ping`, the bot can be configured to respond with "pong". This feature can be customized with any amount of "pings" and "pongs". 

Wisdom is a feature that send back a random or specified line from a book. The book has to be in a database.

Decide randomly picks one of the specified words. For example, `/decide eat sleep drink` picks either eat, sleep or drink with a 33% chance to land on any of them.

Horoscope gives the daily forecast of your life based on your zodiac sign. Requires a database with a horoscope table.

Your custom feature must satisfy the feature interface:
```go
type feature interface {
    init(*jbot) error
    triggers(*jbot, tgbotapi.Update) bool
    execute(*jbot, tgbotapi.Update) error
    String() string
}
```
`init` is called when the bot starts, `triggers` checks if an update triggers your feature, `execute` is called if your feature was triggered and `string` returns the name of your feature.

# How to deploy
To deploy this bot, two things are required:
* A distribution of the [go porgramming language](https://golang.org/doc/install)
* A Telegram bot [API key](https://core.telegram.org/#bot-api)

Optional:
* A [posgres](https://www.postgresql.org/download/) SQL database (with proper tables created)


The bot needs to be configured. This is done by creating a file called `config.json`.
An example of such a file (`example_config.json`) is provided, you can simply rename it to `config.json`.
The configuration file has many filds. These fields are documented later. When running the bot, please have the config file in the current working directory.

Install and run with these commands:
* `go get` the dependencies.
* `go install` this bot.
* `./juhannusbot` in your go paths bin directory.

To stop the bot, use CTRL+C/CMD+C.

# Populating the database
Some features require a PostgreSQL database connection. You can still run the bot without a database connection, the database related features will simply be disabled.

The bot expects a table called `book` with rows `chapter`, `verse` and `text`. The bot also expects a table called `horoscope` with rows `datestring`, `signstring`, `text`, `intensity`, `keywords` and `mood`.

You can create these tables with:
```sql
CREATE TABLE book (
chapter varchar(7),
verse   varchar(7),
text    varchar(4096)
);
```

```sql
CREATE TABLE horoscope (
    datestring varchar(20),
    signstring varchar(20),
    text varchar(1000),
    intensity varchar(100),
    keywords varchar(100),
    mood varchar(100)
);
```

For some of the features to work, you need to [insert](https://www.postgresql.org/docs/11/tutorial-populate.html) a few rows to both tables. 

Place some rows to your book with a statement such as:
```sql
COPY book(chapter,verse,text) FROM 'YOUR_PATH/book_data_example.txt' WITH DELIMITER '%'; 
```

You can add some rows to the horoscope table with:
```sql
INSERT INTO horoscope (datestring, signstring, text, intensity, keywords, mood)
    VALUES ('1.1.1980', 'aries', 'good luck', '1%', 'a, b, c', 'happy'),
           ('1.1.1980', 'taurus', 'no luck', '2%', 'a, b, c', 'happy'),
           ('1.1.1980', 'gemini', 'bad luck', '3%', 'a, b, c', 'happy'),
           ('1.1.1980', 'cancer', 'perfect luck', '4%', 'a, b, c', 'happy'),
           ('1.1.1980', 'leo', 'slight luck', '5%', 'a, b, c', 'happy'),
           ('1.1.1980', 'virgo', 'some luck', '6%', 'a, b, c', 'happy'),
           ('1.1.1980', 'libra', 'steady luck', '7%', 'a, b, c', 'happy'),
           ('1.1.1980', 'scorpio', 'very good luck', '8%', 'a, b, c', 'happy'),
           ('1.1.1980', 'sagittarius', 'lucky', '9%', 'a, b, c', 'happy'),
           ('1.1.1980', 'capricorn', 'unlucky weather', '10%', 'a, b, c', 'happy'),
           ('1.1.1980', 'aquarius', 'alternating luck', '11%', 'a, b, c', 'happy'),
           ('1.1.1980', 'pisces', 'unforseeable luck', '12%', 'a, b, c', 'happy');
```

# Configuring the bot

The bot is configured by editing `confg.json`. An example of a config file is given in the file `example_config.json`. 

The following fields need to be configured: 
* "apikey": put your telegram APIkey here
* "databaseurl": your PostgreSQL connection string

Some of the features can be customized by further editing of `config.json`.

The pingpong config has a list of features.
Each pingpong feature in the list has the following fields:
* "pings": list of strings that trigger the command when seen by the bot.
* "pongs": list of strings that the command can send back to the user. If there are multiple entries, a random one is chosen.
* "isprefixcommand": bool for whether the ping string needs to be at the start of the recieved message (false means it can be anywhere).
* "isreply": bool, true if the reply message is treated as a telegram reply.
* "successpropability": 0.0-1.0, if less than 1.0, the command has a chance of not sending back anyting.

Here is an example: 
```json
{
    "pings": ["!d6","!dice6"],
    "pongs": ["1","2","3","4","5","6"],
    "isprefixcommand": true,
    "isreply": true,
    "successpropability": 1.0
}
```
Now, whenever the bot sees a message that starts with "!d6" or "!dice6", it reacts. One of the six pongs is sent to the user. The randomly chosen pong is sent as a telegram reply to the message that contained the pinging word. 

Or:
```json
{
    "pings": ["/start","/info"],
    "pongs": ["Hello, try my features:\n/decide\n/wisdom\n!d6"],
    "isprefixcommand": true,
    "isreply": false,
    "successpropability": 1.0
}
```
Now, a message starting with "/start" or "/info" will promt the bot to answer with some information. The information is sen as a normal telegram message.

Most of the feures can be configured similarly to pingpong. You can experiment with them or use the defaults.
