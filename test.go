package main

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
)

const substringGO = "Go"

func main() {
	testSlice := []string{"https://go.dev/tour/moretypes/7", "https://golangforall.com/ru/post/golang-data-handling-concurrent-programs.html", "https://freshman.tech/snippets/go/iterating-over-slices/"}
	var ent, _ = CalculateEndpointKeywordEntries(testSlice)
	fmt.Println(ent)

}
func CalculateEndpointKeywordEntries(endpointSlice []string) (int, error) {
	if endpointSlice == nil {
		return 0, errors.New("slice is empty")
	}
	entriesCount := 0
	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(endpointSlice))
	syncMutes := sync.Mutex{}
	for _, v := range endpointSlice {
		go func(endpoint string) {
			resp, err := http.Get(endpoint)
			if err != nil {
				entriesCount += 0
				return
			}
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				entriesCount += 0
				return
			}
			sb := string(body)
			syncMutes.Lock()
			entriesCount += strings.Count(sb, substringGO)
			syncMutes.Unlock()
			waitGroup.Done()
		}(v)
	}
	waitGroup.Wait()
	return entriesCount, nil
}
