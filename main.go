package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	dsn := "mysql_username:CHANGEME@tcp(localhost:3306)/dbname"

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer sql.Close()

	rows, err := db.Query("select id, first_name from user limit 10")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		var username string
		if err := rows.Scan(&id, &username); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%d-%s\n", id, username)
	}
}
