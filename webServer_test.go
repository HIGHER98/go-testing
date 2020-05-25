package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func create_table() {
	db := dbConnect()
	if db == nil {
		log.Fatal("Failed to connect to DB")
		return
	}
	const query = `
		CREATE TABLE users (
			id int  AUTO_INCREMENT,
			first_name varchar(20),
			last_name varchar(20),
			PRIMARY KEY (id)
		)
	`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Close()
}

func drop_table() {
	db := dbConnect()
	if db == nil {
		log.Fatal("Failed to connect to DB")
	}
	const query = `
		DROP TABLE IF EXISTS users
	`
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
		return
	}
	db.Close()
}

func insert_record(query string) {
	db := dbConnect()
	if db == nil {
		log.Fatal("Failed to connect to DB")
	}
	_, err := db.Exec(query)
	if err != nil {
		fmt.Println(err)
	}
	db.Close()
}

func Test_count(t *testing.T) {
	var count int
	create_table()
	insert_record("INSERT INTO users (first_name, last_name) VALUES ('gary', 'hyer')")
	insert_record("INSERT INTO users (first_name, last_name) VALUES ('john', 'jones')")
	insert_record("INSERT INTO users (first_name, last_name) VALUES ('flip', 'flop')")
	insert_record("INSERT INTO users (first_name, last_name) VALUES ('test', 'icle')")

	db := dbConnect()
	if db == nil {
		t.Errorf("Failed to connect to database.")
	}
	row := db.QueryRow("SELECT COUNT(*) FROM users")
	err := row.Scan(&count)
	if err != nil {
		t.Error(err)
	}
	if count != 4 {
		t.Errorf("Select query returned count of: %d", count)
	}
	db.Close()
	drop_table()
}

func Test_queryDB(t *testing.T) {
	create_table()
	db := dbConnect()
	if db == nil {
		t.Errorf("Failed to connect to database.")
	}
	query := "INSERT INTO users (first_name, last_name) VALUES ('random', '1234')"
	insert_record(query)
	rows, err := db.Query("SELECT * FROM users WHERE last_name=?", `1234`)
	if err != nil {
		t.Error("Failed to execute query: ", err)
	}
	var col1 int
	var col2 string
	var col3 string

	for rows.Next() {
		rows.Scan(&col1, &col2, &col3)
	}
	if col2 != "random" {
		t.Errorf("First name returned: %s", col2)
	}
	if col3 != "1234" {
		t.Errorf("Last name returned: %s", col3)
	}
	db.Close()
	drop_table()
}

func Test_record(t *testing.T) {
	create_table()
	insert_record("INSERT INTO users (first_name, last_name) VALUES ('John', 'Doe')")
	req, err := http.NewRequest("GET", "/getdata", nil)
	if err != nil {
		t.Errorf("Failed to execute query: %v", err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(getData)
	handler.ServeHTTP(rr, req)

	status := rr.Code
	if status != http.StatusOK {
		t.Errorf("Handler returned %v", status)
	}

	if rr.Body.String() != "<h3 align=\"center\">1, John, Doe</h3>" {
		t.Errorf("Received response:\n%v\n-------", rr.Body.String())
	}
	drop_table()
}
