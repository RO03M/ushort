package main

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"urlshort/src/db"

	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

type RequestData struct {
	LongUrl string `json:"longUrl"`
}

var urlMap = make(map[string]string)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("teste")
		client := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "", // No password set
			DB:       0,  // Use default DB
			Protocol: 2,  // Connection protocol
		})

		context := context.Background()
		err := client.Set(context, "foo", "bar", 0).Err()

		if err != nil {
			panic(err)
		}

		val, err := client.Get(context, "foo").Result()

		io.WriteString(w, val)
	})

	http.HandleFunc("GET /pg", func(w http.ResponseWriter, r *http.Request) {
		db, err := sql.Open("postgres", "host=localhost port=5432 user=admin password=admin dbname=ushort sslmode=disable")

		if err != nil {
			io.WriteString(w, fmt.Sprintf("Failed to connect to the database :( %s", err))
			return
		}

		err = db.Ping()
		defer db.Close()
		fmt.Println(err)

		io.WriteString(w, "return")
	})

	http.HandleFunc("GET /urls", func(w http.ResponseWriter, r *http.Request) {
		urls, err := db.GetUrls()

		if err != nil {
			io.WriteString(w, "There was an error")
			return
		}

		response, err := json.Marshal(urls)

		io.Writer.Write(w, response)
	})

	http.HandleFunc("/{shorten}", func(w http.ResponseWriter, r *http.Request) {
		shortenUrl := r.PathValue("shorten")
		fmt.Println(shortenUrl)

		longUrl := urlMap[shortenUrl]

		if longUrl == "" {
			return
		}
		fmt.Println(longUrl)
		http.Redirect(w, r, longUrl, http.StatusFound)
	})

	http.HandleFunc("/create-shorten", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		var data RequestData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		hash := sha256.Sum256([]byte(data.LongUrl))
		var shortenUrl string = hex.EncodeToString(hash[:])[:7]

		urlMap[shortenUrl] = data.LongUrl

		fmt.Println(shortenUrl)
	})

	fmt.Println("Listening on http://localhost:8000")
	http.ListenAndServe(":8000", nil)

	fmt.Println("Finished ushort")
}
