package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// This crap is rate limited to somehing weak like 9 reqs per minute.
const perMinute = 9

var (
	cnt       = int64(0)
	caches    = make(map[string]bytes.Buffer)
	cacheLock sync.Mutex
)

type h struct{}

func (h) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	cacheKey := fmt.Sprintf("%s:%s", r.URL.Path, r.URL.Query().Get("guid"))

	w.Header().Add("Access-Control-Allow-Origin", "*")

	if atomic.LoadInt64(&cnt) >= perMinute {
		log.Printf("serving cached response for %s", cacheKey)
		cacheLock.Lock()
		buf := caches[cacheKey]
		io.Copy(w, &buf)
		cacheLock.Unlock()
		return
	}

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

	cacheLock.Lock()
	buf := caches[cacheKey]
	buf.Reset()
	tee := io.TeeReader(resp.Body, &buf)
	io.Copy(w, tee)
	caches[cacheKey] = buf
	cacheLock.Unlock()
	atomic.AddInt64(&cnt, 1)
	log.Printf("freshly served %s", cacheKey)
}

func main() {
	go func() {
		for {
			time.Sleep(time.Minute)
			atomic.StoreInt64(&cnt, 0)
		}
	}()

	log.Fatal(http.ListenAndServe(":5353", h{}))
}
