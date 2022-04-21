package main

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"
)

type result2 struct {
	url     string
	err     error
	latency time.Duration
}

func get2(ctx context.Context, url string, ch chan<- result2) {
	var r result2

	start := time.Now()
	ticker := time.NewTicker(1 * time.Second).C
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r = result2{url, err, 0}
		return
	}
	t := time.Since(start).Round(time.Millisecond)
	r = result2{url, nil, t}
	errorHttp := resp.Body.Close()
	if err != nil {
		log.Panicln(errorHttp)
		return
	}
	for {
		select {
		case ch <- r:
			return
		case <-ticker:
			log.Println("tick", r)

		}
	}
}

func first(ctx context.Context, urls []string) (*result2, error) {
	results := make(chan result2, len(urls))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	for _, url := range urls {
		go get2(ctx, url, results)

	}
	select {
	case r := <-results:
		return &r, nil
	case <-ctx.Done():
		return nil, ctx.Err()

	}

}

func main() {

	list := []string{
		"https://amazon.com",
		"https://google.com",
		"https://uonbi.ac.ke",
		"https://wsj.com",
		"https://nytimes.com",
	}

	r, _ := first(context.Background(), list)

	if r.err != nil {
		log.Printf("%s %s\n", r.url, r.err)
	} else {
		log.Printf("%s %s\n", r.url, r.latency)
	}
	time.Sleep(9 * time.Second)
	log.Println("quit anyway...", runtime.NumGoroutine(), "still running")
}
