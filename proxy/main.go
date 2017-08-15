package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

type h struct{}

func (h) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://realmofthemadgodhrd.appspot.com%s", r.URL.String()), nil)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("failed: %s", err)
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("failed: %s", err)
		return
	}
	defer resp.Body.Close()
	io.Copy(w, resp.Body)
}

func main() {
	log.Fatal(http.ListenAndServe(":5353", h{}))
}
