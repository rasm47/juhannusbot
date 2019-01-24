package jbot

import (
    "testing"
)

func TestStart(t *testing.T) {
    t.Skipf("Start test skipped for now")
}

func TestHandleUpdate(t *testing.T) {
    t.Skipf("HandleUpdate test skipped for now")
}

func TestParseHoroscopeMessageAries(t *testing.T) {
    originalMessage := "oinas"
    targetOutput    := "aries"
    output          := parseHoroscopeMessage(originalMessage)
    if output != targetOutput {
        t.Fatalf("input %s produced output %s instead of desired output %s", originalMessage, output, targetOutput)
    }
}

func TestParseHoroscopeMessageNoSign(t *testing.T) {
    originalMessage := "No signs here"
    targetOutput    := ""
    output          := parseHoroscopeMessage(originalMessage)
    if output != targetOutput {
        t.Fatalf("input %s produced output %s instead of desired output %s", originalMessage, output, targetOutput)
    }
}

func TestResolveHoroscope(t *testing.T) {
    t.Skipf("resolveHoroscope test skipped for now")
}

func TestGetBookLineEmptyBook(t *testing.T) {
    emptyBook := []string{}
    if getBookLine(emptyBook) != "" {
        t.Fatalf("Empty book did not return an empty line")
    }
}

func TestGetBookLineNormalBook(t *testing.T) {
    book := []string{"line1", "line2", "line3", "line4", "line5"}
    line := getBookLine(book)
    
    lineIsInBook := false
    for i := range book {
        if book[i] == line {
            lineIsInBook = true
            break
        }
    }
    
    if !lineIsInBook {
        t.Fatalf("getBookLine did not return any of the lines in the book")
    }
}

func TestSendMessage(t *testing.T) {
    t.Skipf("sendMessage test skipped for now")
}

func TestCreateResponseCommandHello(t *testing.T) {
    desiredResponse := "world!"
    message := "/hello friend!"
    book := []string{"Test", "Book "}
    response, err := createResponse(message, book)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    
    if response != "world!" {
        t.Fatalf("\"%s\" resulted in \"%s\" instead of \"%s\"", message, response, desiredResponse)
    }
}

func TestCreateResponseCommandHelloWithEmptyBook(t *testing.T) {
    desiredResponse := "world!"
    message := "/hello"
    book := []string{}
    response, err := createResponse(message, book)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    
    if response != desiredResponse {
        t.Fatalf("\"%s\" resulted in \"%s\" instead of \"%s\"", message, response, desiredResponse)
    }
}

func TestCreateResponseCommandRaamatturivi(t *testing.T) {
    desiredResponse1 := "Test"
    desiredResponse2 := "Book "
    message := "/raamatturivi"
    book := []string{"Test", "Book "}
    response, err := createResponse(message, book)
    if err != nil {
        t.Error(err)
        t.Fail()
    }
    
    if !(response == desiredResponse1 || response == desiredResponse2) {
        t.Fatalf("\"%s\" resulted in \"%s\" instead of \"%s\" or \"%s\"", message, response, desiredResponse1, desiredResponse2)
    }
}
