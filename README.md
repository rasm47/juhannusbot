# juhannusbot
A Telegram bot for entertainment purposes.

# Commands
The bot has two types of commands: message commands and special commands.

A message command is a simple command where the bot reacts to text with text. When the bot sees a message with a certain keyword, it send back a corresponding message. For example, when `/start` is seen, a message containing instructions could be sent back. Any number of such commands can be configured. 

Special commands are commands that perform more complicated tasks. These commands are unique so there is only one of each special command. These commands are "horoscope", "wisdom" and "decide". Horoscope gives the horoscope of the day. Wisdom gives a random line from a book in your database. If you give additional parameters e.g. `Ch1 Verse3`, you can request your favorite lines from the book. Decide randomly picks one of the words coming with the command.

Special commands require a working database connection with a database that has certain tables in it. 

# How to deploy
To deploy this bot, three things are required:
* A distribution of the [go porgramming language](https://golang.org/doc/install)
* A [posgres](https://www.postgresql.org/download/) SQL database (with proper tables created)
* A Telegram bot [API key](https://core.telegram.org/#bot-api)

The bot needs to be configured. This is done by creating a file called `config.json`.
An example of such a file (`example_config.json`) is provided, you can simply rename it to `config.json`.
The configuration file has many filds. These fields are documented later. When running the bot, please have the config file in the current working directory.

Install and run with these commands:
* `go get` the dependencies.
* `go install` this bot.
* `./juhannusbot` in your go paths bin directory.

To stop the bot, use CTRL+C/CMD+C.

# Populating the database
The special commands require a connection to a PostgreSQL database with some book and horoscope data.
The bot expects there to be a table called `book` with rows `chapter`, `verse` and `text`.
The bot also expects  a table called `horoscope` with rows `datestring`, `signstring`, `text`, `intensity`, `keywords` and `mood`.


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


If you want the special commands to work, you need to [insert](https://www.postgresql.org/docs/11/tutorial-populate.html) a few rows to both tables. 

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
* "debug": bool for debug mode, default is false
* "databaseurl": your PostgreSQL connection string
* "commands": list of commands

Each command in the list of commands has the following fields:
* "name": name of command, please give each command a unigue name.
* "type": "message" or "special", depending on the type of the command.
* "alias": list of strings that trigger the command when seen by the bot.
* "reply": list of strings that the command can send back to the user. If there are multiple entries, a random one is chosen.
* "isprefixcommand": bool for whether the alias string needs to be at the start of the recieved message or can be anywhere.
* "isreply": bool, true if the reply message is treated as a telegram reply or a telegram message.
* "successpropability": 0.0-1.0, if less than 1.0, the command has a chance of not sending back anyting.

Here is an example of a command: 
```json
{
    "name": "dice",
    "type": "message",
    "alias": ["!d6","!dice6"],
    "reply": ["1","2","3","4","5","6"],
    "isprefixcommand": true,
    "isreply": true,
    "successpropability": 1.0
}
```
Now, whenever the bot sees a message that starts with "!d6" or "!dice6", it reacts. One of the six replies is sent to the user. The randomly chosen reply is sent as a telegram reply to the message that contained the alias word. 
