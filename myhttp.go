package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	parPtr := flag.Int("parallel", 10, "the limit of the number of parallel requests")
	flag.Parse()
	for r := range workerPool(&http.Client{Timeout: 10 * time.Second}, *parPtr, flag.Args()) {
		if r.err == nil {
			fmt.Println(r.url, r.md5)
		}
	}
}

type client interface {
	Do(req *http.Request) (*http.Response, error)
}

func workerPool(client client, numWorkers int, urls []string) <-chan result {
	jobs := make(chan string, len(urls))
	results := make(chan result, len(urls))

	// This starts up numWorkers workers, initially blocked
	// because there are no jobs yet.
	var wg sync.WaitGroup
	wg.Add(len(urls))
	for w := 0; w < numWorkers; w++ {
		go worker(&wg, client, jobs, results)
	}

	// Here we send numJobs `jobs` and then `close` that
	// channel to indicate that's all the work we have.
	for _, j := range urls {
		jobs <- j
	}
	close(jobs)
	go func() {
		wg.Wait()
		close(results)
	}()
	return results
}

func worker(wg *sync.WaitGroup, client client, jobs <-chan string, results chan<- result) {
	for j := range jobs {
		results <- do(client, j)
		wg.Done()
	}
}

func do(client client, url string) result {
	url = "http://" + url
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result{err: err}
	}

	res, err := client.Do(req)
	if err != nil {
		return result{err: err}
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result{err: err}
	}
	return result{url: url, md5: fmt.Sprintf("%x", md5.Sum(body))}

}

type result struct {
	url string
	md5 string
	err error
}
