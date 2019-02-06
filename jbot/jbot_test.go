package jbot

import (
    "testing"
)

func TestStart(t *testing.T) {
    // To test this function properly, a mock telegram bot API is needed.
    t.Skipf("Start test skipped for now")
}

func TestHandleUpdate(t *testing.T) {
    // To test this function properly, a mock telegram bot API is needed.
    t.Skipf("HandleUpdate test skipped for now")
}

func TestParseHoroscopeMessageAries(t *testing.T) {
    originalMessage := "oinas"
    targetOutput    := horoscopeSignAries
    output          := parseHoroscopeMessage(originalMessage)
    if output != targetOutput {
        t.Fatalf("input %s produced output %s instead of desired output %s", originalMessage, output, targetOutput)
    }
}

func TestParseHoroscopeMessageNoSign(t *testing.T) {
    originalMessage := "No signs here"
    targetOutput    := horoscopeSignNone
    output          := parseHoroscopeMessage(originalMessage)
    if output != targetOutput {
        t.Fatalf("input %s produced output %s instead of desired output %s", originalMessage, output, targetOutput)
    }
}

func TestResolveHoroscope(t *testing.T) {
    // Needs a mock of the horoscope API
    t.Skipf("resolveHoroscope test skipped for now")
}

func TestHoroscopeReply(t *testing.T) {
    t.Skipf("asd test skipped for now")
}

func TestSendMessage(t *testing.T) {
    // To test this function properly, a mock telegram bot API is needed.
    t.Skipf("sendMessage test skipped for now")
}

func TestCreateResponse(t *testing.T) {
    t.Skipf("createResponse test skipped for now")
}

func TestCreateBookResponseString(t *testing.T) {
    t.Skipf("createBookResponseString test skipped for now")
}
