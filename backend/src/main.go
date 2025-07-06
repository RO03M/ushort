package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
	"urlshort/src/db"
	"urlshort/src/modules/urls"

	"github.com/redis/go-redis/v9"

	_ "github.com/lib/pq"
)

type RequestData struct {
	LongUrl string `json:"longUrl"`
}

var urlMap = make(map[string]string)

var client = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // No password set
	DB:       0,  // Use default DB
	Protocol: 2,  // Connection protocol
})

func main() {

	http.HandleFunc("GET /{shorten}", func(w http.ResponseWriter, r *http.Request) {
		shortenUrl := r.PathValue("shorten")
		start := time.Now()

		ctx := context.Background()
		cachedUrl, err := client.Get(ctx, fmt.Sprintf("url:alias:%s", shortenUrl)).Result()

		if err == nil {
			var url urls.Url
			err := json.Unmarshal([]byte(cachedUrl), &url)

			if err == nil && url.Url != "" {
				http.Redirect(w, r, url.Url, http.StatusMovedPermanently)
				fmt.Println("elapsed cache", time.Since(start))
				return
			}
		}

		url := urls.FindUrlByAlias(shortenUrl)

		if url == nil {
			io.WriteString(w, "No ushort found :(")
			return
		}

		http.Redirect(w, r, url.Url, http.StatusMovedPermanently)
		fmt.Println("elapsed total", time.Since(start))
	})
	http.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("teste")

		ctx := context.Background()
		err := client.Set(ctx, "foo", "bar", 0).Err()

		if err != nil {
			panic(err)
		}

		val, err := client.Get(ctx, "foo").Result()

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

		response, _ := json.Marshal(urls)

		io.Writer.Write(w, response)
	})

	http.HandleFunc("POST /create-shorten", func(w http.ResponseWriter, r *http.Request) {
		var data RequestData
		err := json.NewDecoder(r.Body).Decode(&data)

		if err != nil {
			fmt.Println(err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		url := urls.CreateUrl(data.LongUrl)
		jsonUrl, _ := json.Marshal(url)

		ctx := context.Background()
		client.Set(ctx, fmt.Sprintf("url:alias:%s", url.Alias), string(jsonUrl), time.Minute*100)

		io.Writer.Write(w, jsonUrl)
	})

	fmt.Println("Listening on http://localhost:8000")
	http.ListenAndServe(":8000", nil)

	fmt.Println("Finished ushort")
}
