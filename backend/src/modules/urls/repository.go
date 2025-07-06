package urls

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"urlshort/src/db"
)

func FindUrlByAlias(alias string) *Url {
	db := db.CreateConnection()
	defer db.Close()

	var url Url

	var err = db.QueryRow(`SELECT id, url, alias FROM urls WHERE alias = $1`, alias).Scan(&url.Id, &url.Url, &url.Alias)

	if err != nil {
		fmt.Printf("Failed to find url by alias. %v", err)
	}

	return &url
}

func CreateUrl(url string) *Url {
	db := db.CreateConnection()
	defer db.Close()

	var id int

	hash := sha256.Sum256([]byte(url))
	var alias string = hex.EncodeToString(hash[:])[:7]

	urlInstance := FindUrlByAlias(alias)

	fmt.Println(urlInstance)
	if urlInstance != nil {
		return urlInstance
	}

	err := db.QueryRow(`INSERT INTO urls (url, alias) VALUES ($1, $2) RETURNING id`, url, alias).Scan(&id)

	if err != nil {
		fmt.Printf("Failed to insert url. %v\n", err)
	}

	urlInstance = &Url{
		Id:    id,
		Url:   url,
		Alias: alias,
	}

	return urlInstance
}
