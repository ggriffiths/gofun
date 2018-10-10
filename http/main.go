package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}

}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	respBody := make(map[interface{}]interface{})
	if err := json.NewDecoder(r.Body).Decode(&respBody); err != nil {
		w.Write([]byte("Decode failed"))
		return
	}

	fmt.Println("Hello")
}
