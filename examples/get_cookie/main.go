package main

import (
	"fmt"
	"net/url"

	requests "github.com/hiroakis/go-requests"
)

func main() {
	param := &url.Values{}
	param.Add("k1", "v1")
	param.Add("k2", "v2")
	resp, err := requests.Get("https://httpbin.org/cookies/set", param, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	for k, v := range resp.Cookies() {
		fmt.Printf("------ Cookie %d ------\n", k)
		fmt.Printf("Name: %s\n", v.Name)
		fmt.Printf("Value: %s\n", v.Value)
		fmt.Printf("Path: %s\n", v.Path)
	}
}
