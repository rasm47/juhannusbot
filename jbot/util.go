package jbot

import (
    "os"
    "bufio"
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
        if line != "" {
            lines = append(lines, line)
        }
    }
    err = scanner.Err()
    return
}
