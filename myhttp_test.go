package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestDo(t *testing.T) {
	testCases := map[string]struct {
		url      string
		client   client
		expected result
	}{
		"dummy": {
			url:      "google.com",
			client:   &dummyClient{},
			expected: result{url: "google.com", md5: fmt.Sprintf("%x", md5.Sum([]byte("http://google.com")))},
		},
		"error": {
			client:   &errClient{},
			expected: result{err: errors.New("")},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			ok := assertRespValues(tc.expected, do(tc.client, tc.url))
			if !ok {
				t.Fatal()
			}

		})
	}
}

func assertRespValues(expected, actual result) bool {
	return expected.url == actual.url && expected.md5 == actual.md5 &&
		(expected.err == nil || expected.err.Error() == actual.err.Error())
}

func TestWorker(t *testing.T) {
	urls := []string{"google.com", "facebook.com", "netflix.com"}
	jobs := make(chan string, len(urls))
	for _, u := range urls {
		jobs <- u
	}
	close(jobs)
	results := make(chan result, len(urls))
	var wg sync.WaitGroup
	wg.Add(len(urls))
	worker(&wg, &dummyClient{}, jobs, results)
	go func() {
		<-time.After(5 * time.Second)
		t.Fatal()
	}()
	wg.Wait()
}

func TestWorkerPool(t *testing.T) {
	urls := []string{"google.com", "facebook.com", "netflix.com"}
	go func() {
		<-time.After(1 * time.Second)
		t.Fatal()
	}()
	for _ = range workerPool(&slowClient{}, len(urls), urls) {
	}
}

type dummyClient struct {
}

func (d *dummyClient) Do(req *http.Request) (*http.Response, error) {
	return &http.Response{Body: ioutil.NopCloser(strings.NewReader(req.URL.String()))}, nil
}

type slowClient struct {
}

func (s *slowClient) Do(req *http.Request) (*http.Response, error) {
	time.Sleep(750 * time.Millisecond)
	return &http.Response{Body: ioutil.NopCloser(strings.NewReader(req.URL.String()))}, nil
}

type errClient struct {
}

func (e *errClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("")
}
