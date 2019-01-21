package util

import (
    "os"
    "bufio"
)

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
