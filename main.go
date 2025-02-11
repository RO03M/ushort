package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type RequestData struct {
	LongUrl string `json:"longUrl"`
}

var urlMap = make(map[string]string)

func main() {
	http.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("teste")

		io.WriteString(w, "Hello there!")
	})

	http.HandleFunc("/{shorten}", func(w http.ResponseWriter, r *http.Request) {
		shortenUrl := r.PathValue("shorten")
		fmt.Println(shortenUrl)
		io.WriteString(w, "should redirect")
		io.WriteString(w, fmt.Sprintf("hello %s", shortenUrl))
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
		fmt.Println(hex.EncodeToString(hash[:]))
	})

	fmt.Println("Listening on http://localhost:8000")
	http.ListenAndServe(":8000", nil)

	fmt.Println("Finished ushort")
}
