package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func CreateConnection() *sql.DB {
	db, err := sql.Open("postgres", "host=localhost port=5432 user=admin password=admin dbname=ushort sslmode=disable")

	if err != nil {
		panic(err)
	}

	return db
}

type Url struct {
	Id    int    `json:"id"`
	Url   string `json:"url"`
	Alias string `json:"alias"`
}

func GetUrls() ([]Url, error) {
	db := CreateConnection()
	defer db.Close()

	var urls []Url

	rows, err := db.Query(`SELECT id, url, alias FROM urls`)

	if err != nil {
		log.Fatalf("Unable to execute the query. %v", err)
	}

	defer rows.Close()

	for rows.Next() {
		var url Url

		err = rows.Scan(&url.Id, &url.Url, &url.Alias)

		if err != nil {
			log.Fatalf("Unable to scan the row. %v", err)
		}

		urls = append(urls, url)
	}

	return urls, err
}
