package main

import (
	"fmt"

	"github.com/hiroakis/go-requests"
)

func main() {

	resp, err := requests.Get("https://httpbin.org/", nil, nil)

	if err != nil {
		fmt.Println(err)
		return
	}
	// Response body
	fmt.Println(resp.Text())
	// status code
	fmt.Printf("Code: %d\n", resp.StatusCode())
	for k, v := range resp.Headers() {
		// Response Headers
		fmt.Printf("%s: %s\n", k, v)
	}
}
