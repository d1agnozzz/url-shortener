package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	fmt.Println("URL Shortener starting...")

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL Shortener is running!")
	})

	addr := ":3000"
	fmt.Printf("Server listening on %s\n", addr)

	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatal(err)
	}
}
