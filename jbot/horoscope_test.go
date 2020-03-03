package jbot

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

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

func TestGetHoroscopeData(t *testing.T) {
	// create mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT .*").WithArgs("sagittarius").WillReturnRows(sqlmock.NewRows([]string{"datestring", "signstring", "text", "intensity", "keywords", "mood"}).AddRow("1.1.1980", "sagittarius", "Good fortune for your friend but not you", "5 percent", "keyword1, keyword2", "neutral"))

	expectedContents := horoscopeData{"1.1.1980", "sagittarius", "Good fortune for your friend but not you", horoscopeMeta{"5 percent", "keyword1, keyword2", "neutral"}}

	contents := getHoroscopeData(db, horoscopeSignSagittarius)
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	if contents != expectedContents {
		t.Fatalf("contents of mock database were fetched incorrectly")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
