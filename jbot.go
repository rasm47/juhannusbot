package main

import (
    "log"
    "gopkg.in/telegram-bot-api.v4"
    "strings"
    "os"
    "bufio"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "bytes"
    "math/rand"
)

type horoscopeMeta struct {
    Intensity string `json:"intensity"`
    Keywords  string `json:"keywords"`
    Mood      string `json:"mood"`
}

type horoscopeResponse struct {
    Date      string `json:"date"`
    Sunsign   string `json:"sunsign"`
    Horoscope string `json:"horoscope"`
    Meta      horoscopeMeta `json:"meta"`
}

func main() {

    apikey, err := FileToLines("apikey.txt")
    if err != nil {
        log.Printf("Have you put your API key to apikey.txt? See README.md")
        log.Panic(err)
    }
    
    bot, err := tgbotapi.NewBotAPI(apikey[0]) 
    if err != nil {
        log.Printf("API key authentication failed. Try to double check if the key is valid.")
        log.Panic(err)
    }

    bot.Debug = true

    log.Printf("%s authenticated", bot.Self.UserName)
    
    bible, err := FileToLines("bible.txt")
    if err != nil {
        log.Printf("Bible not found. Have you made a bible.txt?")
        log.Panic(err)
    }

    u := tgbotapi.NewUpdate(0)
    u.Timeout = 60

    updates, err := bot.GetUpdatesChan(u)

    for update := range updates {
        if update.Message == nil {
            continue
        }

        gotmsg := update.Message.Text
            log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

        if strings.HasPrefix(gotmsg, "/hello"){
            tosend := "world!"
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, tosend)
            bot.Send(msg)
            
        } else if strings.HasPrefix(gotmsg, "/horos"){
            horoscopeSign := ParseHoroscopeMessage(gotmsg)
            if horoscopeSign == "" {
                continue
            }

            response, err := http.Get("http://theastrologer-api.herokuapp.com/api/horoscope/" + horoscopeSign + "/today")
            if err != nil {
                log.Fatal(err)
            } else {
                defer response.Body.Close()
                bodyBytes, err := ioutil.ReadAll(response.Body)
                if err != nil {
                    log.Fatal(err)
                } else {
                    bodyStr := string(bodyBytes)
                    log.Printf(bodyStr)

                    //reset the response body to the original unread state
                    response.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

                    var hresponse horoscopeResponse
                    err := json.Unmarshal(bodyBytes, &hresponse)
                    if err != nil {
                        log.Panic(err)
                    } else {
                        tosend := "Enkelit välittävät horoskooppinne:\n" +
                        hresponse.Horoscope + "\n\nAvainsanat: " +
                        hresponse.Meta.Keywords + "\nTunnetila: " +
                        hresponse.Meta.Mood  + "\n\nHoroskooppi väittyi energiatasolla " +
                        hresponse.Meta.Intensity + "."
                        msg := tgbotapi.NewMessage(update.Message.Chat.ID, tosend)
                        bot.Send(msg)
                    }
                }
            }

        } else if strings.HasPrefix(gotmsg, "/raamatt"){
            tosend := GetBibleLine(bible)
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, tosend)
            bot.Send(msg)
        }
    }
}

func FileToLines(filePath string) (lines []string, err error) {
      f, err := os.Open(filePath)
      if err != nil {
              return
      }
      defer f.Close()

      scanner := bufio.NewScanner(f)
      for scanner.Scan() {
              line := scanner.Text()
              if line != "" {
                lines = append(lines, line)
              }
      }
      err = scanner.Err()
      return
}

func ParseHoroscopeMessage(originalMessage string) string {
    msg := strings.ToLower(originalMessage)
    if strings.Contains(msg, "oina"){
        return "aries"
    } else if strings.Contains(msg, "härk"){
        return "taurus"
    } else if strings.Contains(msg, "kaks"){
        return "gemini"
    } else if strings.Contains(msg, "rap"){
        return "cancer"
    } else if strings.Contains(msg, "leij"){
        return "leo"
    } else if strings.Contains(msg, "neit"){
        return "virgo"
    } else if strings.Contains(msg, "vaa"){
        return "libra"
    } else if strings.Contains(msg, "skor"){
        return "scorpio"
    } else if strings.Contains(msg, "jous"){
        return "sagittarius"
    } else if strings.Contains(msg, "vesi"){
        return "capricorn"
    } else if strings.Contains(msg, "kaur"){
        return "aquarius"
    } else if strings.Contains(msg, "kal"){
        return "pisces"
    } else {
        return ""
    }
}

func GetBibleLine(bible []string) string {
  return bible[rand.Intn(len(bible))]
}

