package tools

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/capitan-beto/vale-backend/models"
	"github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func CreateConnection() (*sql.DB, error) {
	cfg := mysql.Config{
		User:                 os.Getenv("DB_USER"),
		Passwd:               os.Getenv("DB_PASS"),
		Net:                  "tcp",
		Addr:                 os.Getenv("DB_ADDR"),
		DBName:               os.Getenv("DB_NAME"),
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Error(err)
		return nil, err
	}

	db.SetMaxOpenConns(5)

	pingErr := db.Ping()
	if pingErr != nil {
		// log.Error(pingErr)
		return nil, err
	}

	fmt.Println("Connected!")
	return db, nil
}

func NewContestant(ct *models.ContestantData, db *sql.DB) (int64, error) {
	name := ct.Name
	phone := ct.Phone
	createdAt := time.Now().Format("2006-01-02 15:04:05")
	confirmed := 0
	extRef := ct.ExtRef
	num := ct.Number

	fmt.Println("hasta aca ok")
	res, err := db.Exec(`INSERT INTO raffle (name, phone, created_at, confirmed, ext_ref, number)
		VALUES (?, ? ,? ,? ,?, ?)`, name, phone, createdAt, confirmed, extRef, num)
	if err != nil {
		log.Error(err)
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		log.Error(err)
	}

	db.Close()
	fmt.Println(id)
	return id, nil
}

func CheckPayment(ref string, db *sql.DB) error {
	row := db.QueryRow("UPDATE raffle SET confirmed = ? WHERE ext_ref = ?", 1, ref)
	if err := row.Err(); err != nil {
		log.Error(err)
		return err
	}

	db.Close()
	return nil
}
