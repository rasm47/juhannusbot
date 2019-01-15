package jbot

import (
    "os"
    "bufio"
    "strings"
    "math/rand"
)

type HoroscopeMeta struct {
    Intensity string `json:"intensity"`
    Keywords  string `json:"keywords"`
    Mood      string `json:"mood"`
}

type HoroscopeResponse struct {
    Date      string `json:"date"`
    Sunsign   string `json:"sunsign"`
    Horoscope string `json:"horoscope"`
    Meta      HoroscopeMeta `json:"meta"`
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

func GetBibleLine(bible []string) string {
  return bible[rand.Intn(len(bible))]
}

