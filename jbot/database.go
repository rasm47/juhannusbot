package jbot

import (
    "strings"
    "database/sql"
    
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
    
    // This query might not scale when using HUGE databases
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
