package tools

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/capitan-beto/vale-backend/models"
)

var td = models.ContestantData{
	Id:        1,
	Name:      "test1",
	Phone:     "1234 12-1234",
	CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	Confirmed: false,
	ExtRef:    "text-ext-ref",
	Number:    1,
}

func TestNewContestant(t *testing.T) {
	stmt := "INSERT INTO raffle (name, phone, created_at, confirmed, ext_ref, number) VALUES (?, ? ,? ,? ,?, ?)"
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	mock.ExpectExec(regexp.QuoteMeta(stmt)).WithArgs(td.Name, td.Phone,
		td.CreatedAt, 0, td.ExtRef, td.Number).WillReturnResult(sqlmock.NewResult(1, 1))

	if _, err = NewContestant(&td, db); err != nil {
		t.Fatalf("error! %v", err)
	}

}

func TestCheckPayment(t *testing.T) {
	refMock := "1234abcd"

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to create sqlmock: %v", err)
	}

	rows := sqlmock.NewRows([]string{"confirmed"}).AddRow("ok")
	mock.ExpectQuery("UPDATE raffle (.*)").WillReturnRows(rows)

	if err := CheckPayment(refMock, db); err != nil {
		t.Fatalf("error! expected no error, got %v", err)
	}
}
