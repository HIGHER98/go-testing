package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
	"strconv"
	"time"
)

var (
	dbdriver = "mysql"
	dbuser   = "Golang"
	dbpass   = "uberDuberSecret"
	dbname   = "golangTestDB"
)

func dbConnect() *sql.DB {
	connStr := fmt.Sprintf("%s:%s@/%s", dbuser, dbpass, dbname)
	db, err := sql.Open(dbdriver, connStr)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return db
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Serving %s\n", r.URL.Path)
	fmt.Printf("Served %s\n", r.Host)
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(time.RFC1123)
	Body := "The current time is "
	fmt.Fprintf(w, "<h1 align=\"center\">%s</h1>", Body)
	fmt.Fprintf(w, "<h2 align=\"center\">%s</h2>", t)
	fmt.Fprintf(w, "Serving %s\n", r.URL.Path)
	fmt.Printf("Served time for %s\n", r.Host)
}

func getData(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Serving %s\n", r.URL.Path)
	fmt.Printf("Served %s\n", r.Host)

	//CREATE USER 'Golang'@'localhost' IDENTIFIED BY 'uberDuberSecret'
	//user:password@/dbname
	connStr := fmt.Sprintf("%s:%s@/%s", dbuser, dbpass, dbname)
	db, err := sql.Open(dbdriver, connStr)
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Failed to connect to database. Please try again later. ")
		return
	}

	rows, err := db.Query("SELECT * FROM users")
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(w, "Failed to execute query. Please try again later")
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var firstName string
		var lastName string
		err = rows.Scan(&id, &firstName, &lastName)
		if err != nil {
			fmt.Println(err)
			fmt.Fprintf(w, "Failed to scan data")
			return
		}
		fmt.Fprintf(w, "<h3 align=\"center\">%d, %s, %s</h3>", id, firstName, lastName)
	}
	err = rows.Err()
	if err != nil {
		fmt.Fprintf(w, "<h3 align=\"center\">%s</h3>", err)
		fmt.Println(err)
		return
	}
}

func main() {
	portNo := flag.Int("p", 8081, "Port to host web server on")
	flag.Parse()
	portstr := strconv.Itoa(*portNo)
	port := ":" + portstr
	fmt.Printf("Using port %s", port)

	http.HandleFunc("/time", timeHandler)
	http.HandleFunc("/getdata", getData)
	http.HandleFunc("/", myHandler)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
}
