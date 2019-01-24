package jbot

import (
    "os"
    "bufio"
    "errors"
)

// readFileToLines opens a text file and collects all lines to an array of strings.
func readFileToLines(filePath string) (lines []string, err error) {
    f, err := os.Open(filePath)
    if err != nil {
        return
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if !(line == "" || line == "\n") {
            lines = append(lines, line)
        }
    }
    err = scanner.Err()
    
    if fileEmpty(lines) {
        err = errors.New("File was found, but it was empty.")
    }
    return
}

// fileEmpty returns true if lines is an empty slice.
func fileEmpty(lines []string) bool {
    if len(lines) == 0 {
        return true
    }
    return false
}
