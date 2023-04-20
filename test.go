package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
)

type EndpointEntries struct {
	endpoint string
	entries  int
}

func main() {
	runtime.GOMAXPROCS(MaxThreads)
	endpointSlice := []string{"https://go.dev/tour/moretypes/7", "https://golangforall.com/ru/post/golang-data-handling-concurrent-programs.html", "https://freshman.tech/snippets/go/iterating-over-slices/"}

	var ent, err = CalculateEndpointKeywordEntries(endpointSlice)
	if err != nil {
		fmt.Println(err)
	}

	sum := 0
	for _, v := range ent {
		fmt.Printf("Count for %v :%v\n", v.endpoint, v.entries)
		sum += v.entries
	}
	fmt.Printf("Total: %v\n", sum)

}

func CalculateEndpointKeywordEntries(endpointSlice []string) ([]EndpointEntries, error) {
	if len(endpointSlice) == 0 {
		return nil, errors.New("slice is empty")
	}
	var entries = []EndpointEntries{}
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(endpointSlice))
	syncMutes := sync.Mutex{}
	for _, v := range endpointSlice {
		go func(endpoint string) {
			stringBody, err := GetStringResponse(endpoint)
			KeywordEntries := 0
			if err == nil {
				KeywordEntries = strings.Count(stringBody, substringGO)
			}
			ent := EndpointEntries{endpoint: endpoint, entries: KeywordEntries}

			syncMutes.Lock()
			entries = append(entries, ent)
			syncMutes.Unlock()

			waitGroup.Done()
		}(v)
	}
	waitGroup.Wait()
	return entries, nil
}

func GetStringResponse(endpoint string) (string, error) {
	if len(endpoint) == 0 {
		return "", errors.New("Endpoint is empty")
	}
	resp, err := http.Get(endpoint)
	if err != nil {
		return "", err
	}
	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}
	return string(body), nil
}
