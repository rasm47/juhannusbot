// Package util provides utilities to github.com/ruoskija/juhannusbot/jbot
package util

import (
    "os"
    "bufio"
)

// ReadFileToLines opens a text file and collects all lines to an array of strings.
func ReadFileToLines(filePath string) (lines []string, err error) {
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
