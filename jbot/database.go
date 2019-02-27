package jbot

import (
    "log"
    "time"
    "strings"
    "net/http"
    "io/ioutil"
    "database/sql"
    "encoding/json"
    
    _ "github.com/lib/pq"
)

// getBookLine fetches a particular bookline from a database.
func getBookLine(database *sql.DB, chapter string, verse string) (string, error) {
    var text string
    err := database.QueryRow("SELECT text FROM book WHERE chapter = $1 and verse = $2", chapter, verse).Scan(&text)
    if err != nil {
        return "", err
    }
    return text, nil
}

// getBookLine fetches and formats a random bookline from a database.
func getRandomBookLine(database *sql.DB) (string, error) {
    var chapter string
    var verse   string
    var text    string
    
    rows, err := database.Query("SELECT chapter, verse, text FROM book ORDER BY RANDOM() LIMIT 1")
    if err != nil {
        return "", err
    }
    defer rows.Close()
    for rows.Next() {
        err := rows.Scan(&chapter, &verse, &text)
        if err != nil {
            return "", err
        }
    }
    err = rows.Err()
    if err != nil {
        return "", err
    }
    
    return strings.ToUpper(chapter) + ". " + verse + " " + text, nil
}

// getHoroscopeData queries the database for the data of a particular sign
func getHoroscopeData(database *sql.DB, sign horoscopeSign) (data horoscopeData) {
    
    rows, err := database.Query("SELECT datestring, signstring, text, intensity, keywords, mood FROM horoscope WHERE signstring = $1", sign.String())
    if err != nil {
        return
    }
    defer rows.Close()
    for rows.Next() {
        err := rows.Scan(&data.Date, &data.Sunsign, &data.Text, &data.Meta.Intensity, &data.Meta.Keywords, &data.Meta.Mood)
        if err != nil {
            return
        }
    }
    err = rows.Err()
    if err != nil {
        return
    }
    return
    
}

// startHoroscopeUpdater starts a process that updates all horoscopes 
// in the database daily at 04:00 / 4am
func startHoroscopeUpdater(database *sql.DB) {
    
    durationToNextFourAm := time.Duration(24 + 4 - time.Now().Hour()) * time.Hour + 
        time.Duration(time.Now().Minute()) * time.Minute + 
        time.Duration(time.Now().Second()) * time.Second
        
    time.AfterFunc(durationToNextFourAm, func(){updateHoroscopeDaily(database)})
    log.Println(durationToNextFourAm)
}

// updateHoroscopeDaily starts a repeating goroutine
// that updates all horoscopes once every 24 hours.
func updateHoroscopeDaily(database *sql.DB) {
    
    log.Println("Horoscopes are updating for the first time: starting the daily updater")
    // call updateHoroscopeData one time
    updateAllHoroscopeData(database)
    
    // start an endless anonymous go routine of daily updating
    go func() {
        for _ = range time.NewTicker(24 * time.Hour).C {
            log.Println("Attempting to fetch new horoscopes...")
            updateAllHoroscopeData(database)
        }
    }()
    
}

// updateAllHoroscopeData updates the database rows for
// all of the horoscopes
func updateAllHoroscopeData(database *sql.DB) {
    
    signs := [12]horoscopeSign{
        horoscopeSignAquarius,
        horoscopeSignPisces,
        horoscopeSignAries,
        horoscopeSignTaurus,
        horoscopeSignGemini,
        horoscopeSignCancer,
        horoscopeSignLeo,
        horoscopeSignVirgo,
        horoscopeSignLibra,
        horoscopeSignScorpio,
        horoscopeSignSagittarius,
        horoscopeSignCapricorn,
    }
    
    for _, sign := range signs {
        updateHoroscopeData(database, sign)
    }
    return
}

// updateHoroscopeData fetches the new horoscope of the day for a 
// partucular horoscopeSign and updates that data to the database.
func updateHoroscopeData(database *sql.DB, sign horoscopeSign) {
    
    data, err := httpGetHoroscopeData(sign)
    if err != nil {
        log.Println("Failed to get new horoscopes from the web, database not updated")
        return
    }
        
    rows, err := database.Query("UPDATE horoscope SET (datestring, text, intensity, keywords, mood) = ($1, $2, $3, $4, $5) WHERE signstring = $6", data.Date, data.Text, data.Meta.Intensity, data.Meta.Keywords, data.Meta.Mood, strings.ToLower(data.Sunsign))
    if err != nil {
        log.Println("Error with the database, database not updated")
        return
    }
    defer rows.Close()
    
    log.Println("Updated data to the database,", sign.String())
    return
}

// httpGetHoroscopeData fetches the new horoscopeData for the day for
// a particular horoscopeSign. The data comes from a REST API whose
// url is hard coded inside this function. 
func httpGetHoroscopeData(sign horoscopeSign) (data horoscopeData, err error) {
    
    response, err := http.Get("http://theastrologer-api.herokuapp.com/api/horoscope/" + sign.String() + "/today")
    if err != nil {
        return
    }
    defer response.Body.Close()
    
    bodyBytes, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return
    }
    
    err = json.Unmarshal(bodyBytes, &data)
    return
}
