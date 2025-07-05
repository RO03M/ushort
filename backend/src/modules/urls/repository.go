package urls

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"urlshort/src/db"
)

func CreateUrl(url string) int {
	db := db.CreateConnection()
	defer db.Close()

	var id int

	hash := sha256.Sum256([]byte(url))
	var shortenUrl string = hex.EncodeToString(hash[:])[:7]

	err := db.QueryRow(`INSERT INTO urls (url, alias) VALUES ($1, $2) RETURNING id`, url, shortenUrl).Scan(&id)

	if err != nil {
		fmt.Printf("Failed to insert url. %v\n", err)
	}

	return id
}
