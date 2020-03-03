package jbot

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestGetBookLine(t *testing.T) {
	// create mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT .*").WithArgs("test", "1:1").WillReturnRows(sqlmock.NewRows([]string{"text"}).AddRow("contents of mock database"))

	contents, err := getBookLine(db, "test", "1:1")
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	if contents != "contents of mock database" {
		t.Fatalf("contents of mock database were fetched incorrectly")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}

func TestGetRandomBookLine(t *testing.T) {
	// create mock database connection
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("^SELECT .*").WillReturnRows(sqlmock.NewRows([]string{"chapter", "verse", "text"}).AddRow("test", "123", "lorem ipsum"))

	contents, err := getRandomBookLine(db)
	if err != nil {
		t.Errorf("error was not expected: %s", err)
	}

	if contents != "TEST. 123 lorem ipsum" {
		t.Fatalf("contents of mock database were fetched incorrectly")
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}

}
