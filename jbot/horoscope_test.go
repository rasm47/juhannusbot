package jbot

import "testing"

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

	for i, testString := range testStrings {
		if convertEmojiToHoroscopeSign(testString) != expectedHoroscopeSigns[i] {
			t.Fatalf("emoji string %s did not produce the expected sign", testString)
		}
	}
}

func TestParseHoroscopeMessageAries(t *testing.T) {
	originalMessage := "oinas"
	targetOutput := horoscopeSignAries
	output := parseHoroscopeMessage(originalMessage)
	if output != targetOutput {
		t.Fatalf("input %s produced output %s instead of desired output %s", originalMessage, output, targetOutput)
	}
}

func TestParseHoroscopeMessageNoSign(t *testing.T) {
	originalMessage := "No signs here"
	targetOutput := horoscopeSignNone
	output := parseHoroscopeMessage(originalMessage)
	if output != targetOutput {
		t.Fatalf("input %s produced output %s instead of desired output %s", originalMessage, output, targetOutput)
	}
}

func TestResolveHoroscope(t *testing.T) {
	// Needs a mock of the horoscope API
	t.Skipf("resolveHoroscope test skipped for now")
}

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
		"sagittarius",
		"capricorn",
		"aquarius",
		"pisces",
	}

	for index, sign := range signs {
		if sign.String() != signStrings[index] {
			t.Fatalf("Got %s, epected %s", sign.String(), signStrings[index])
		}
	}
}
