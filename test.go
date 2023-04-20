package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"
)

type EndpointEntries struct {
	endpoint string
	entries  int
}

func main() {
	runtime.GOMAXPROCS(MaxMachines)
	testSlice := []string{"https://go.dev/tour/moretypes/7", "https://golangforall.com/ru/post/golang-data-handling-concurrent-programs.html", "https://freshman.tech/snippets/go/iterating-over-slices/"}
	var ent, err = CalculateEndpointKeywordEntries(testSlice)
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
	if endpointSlice == nil {
		return nil, errors.New("slice is empty")
	}
	var entries = []EndpointEntries{}
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(endpointSlice))
	syncMutes := sync.Mutex{}
	for _, v := range endpointSlice {
		go func(endpoint string) {
			fmt.Println(time.Now())
			resp, err := http.Get(endpoint)
			if err != nil {
				return
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return
			}
			sb := string(body)

			ent := EndpointEntries{endpoint: endpoint, entries: strings.Count(sb, substringGO)}
			syncMutes.Lock()
			entries = append(entries, ent)
			syncMutes.Unlock()
			waitGroup.Done()
		}(v)
	}
	waitGroup.Wait()
	return entries, nil
}
