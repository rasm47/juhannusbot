package jbot

import (
    "testing"
)

func TestHoroscopeSignStringMethod(t *testing.T) {
    
    signs := [13]horoscopeSign{
        horoscopeSignNone,
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
        horoscopeSignAquarius,
        horoscopeSignPisces,
    }
    
    signStrings := [13]string{
        "",
        "aries",
        "taurus",
        "gemini",
        "cancer",
        "leo",
        "virgo",
        "libra",
        "scorpio",
        "sagittrius",
        "capricorn",
        "aquarius",
        "pisces",
    }
    
    for index, sign := range signs{
        if sign.String() != signStrings[index] {
            t.Fatalf("Got %s, was expecting %s", sign.String(), signStrings[index])
        }
    }
}

func TestStart(t *testing.T) {
    // To test this function properly, a mock telegram bot API is needed.
    t.Skipf("Start test skipped for now")
}

func TestHandleUpdate(t *testing.T) {
    // To test this function properly, a mock telegram bot API is needed.
    t.Skipf("handleUpdate test skipped for now")
}

func TestNewBotInstruction(t *testing.T) {
    // To test this function properly, a mock telegram bot API is needed.
    t.Skipf("newBotInstruction test skipped for now")
}

func TestNewCommand(t *testing.T) {
    
    commandConfigs := commandConfigList{
        commandConfig{[]string{"/start", "/begin"}, "start message"},
        commandConfig{[]string{"/wisdom", "/wisewords"}, ""},
        commandConfig{[]string{"!horoscope"}, ""},
    }
    
    testMessages := [9]string{
        "/start the bot please",        // start command alias with text
        "/begin",                       // start command alias
        "/wisdom for me please",        // wisdom command alias with text
        "/wisewords",                   // wisdom command alias
        "!horoscope",                   // horoscope command alias
        "!horoscope aries",             // horoscope command alias with text
        "Lorem Ipsum",                  // no command aliases
        "please /begin",                // start command alias but not as a prefix
        "/start /wisewords !horoscope", // multiple aliases but start command as prefix
    }
    
    expectedBotCommands := [9]botCommand{
        botCommandStart,
        botCommandStart,
        botCommandWisdom,
        botCommandWisdom,
        botCommandHoroscope,
        botCommandHoroscope,
        botCommandNone,
        botCommandNone,
        botCommandStart,
    }
    
    for index, message := range testMessages {
        if newCommand(commandConfigs, message) != expectedBotCommands[index] {
            t.Fatalf("message %s did not produce the expected command", message)
        }
    }
    
}

func TestExecuteInstruction(t *testing.T) {
    // To test this function properly, a mock telegram bot API is needed.
    t.Skipf("executeInstruction test skipped for now")
}

func TestConvertEmojiToHoroscopeSign(t *testing.T) {
    
    testStrings := [20]string{
        "♒",
        "♓",
        "♈",
        "♉",
        "♊",
        "♋",
        "♌",
        "♍",
        "♎",
        "♏",
        "♐",
        "♑",
        "♑♑",
        "♑   ♑    ♑",
        "♑♎",
        "♎ and ♓",
        " ",
        "",
        "\n",
        "Lorem Ipsum",
    }
    
    expectedHoroscopeSigns := [20]horoscopeSign{
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
        horoscopeSignNone,
        horoscopeSignNone,
        horoscopeSignNone,
        horoscopeSignNone,
        horoscopeSignNone,
        horoscopeSignNone,
        horoscopeSignNone,
        horoscopeSignNone,
    }
    
    for index, testString := range testStrings {
        if convertEmojiToHoroscopeSign(testString) != expectedHoroscopeSigns[index] {
            t.Fatalf("emoji string %s did not produce the expected sign", testString)
        }
    }
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
