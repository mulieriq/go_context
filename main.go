package main

import (
	"context"
	"log"
	"net/http"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(ctx context.Context, url string, ch chan<- result) {

	start := time.Now()
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		ch <- result{url, err, 0}
		return
	}
	t := time.Since(start).Round(time.Millisecond)
	ch <- result{url, nil, t}
	errorHttp := resp.Body.Close()
	if err != nil {
		log.Panicln(errorHttp)
		return
	}

}

func main() {

	result := make(chan result)

	list := []string{
		"https://amazon.com",
		"https://google.com",
		"https://uonbi.ac.ke",
		"https://wsj.com",
		"https://nytimes.com",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)

	defer cancel()

	for _, url := range list {
		go get(ctx, url, result)

	}
	for range list {
		r := <-result

		if r.err != nil {
			log.Printf("%s %s\n", r.url, r.err)
		} else {
			log.Printf("%s %s\n", r.url, r.latency)
		}
	}

}
