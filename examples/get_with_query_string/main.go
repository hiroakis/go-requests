package main

import (
	"fmt"
	"net/url"

	requests "github.com/hiroakis/go-requests"
)

func main() {
	qs := &url.Values{}
	qs.Add("param1", "value1")
	qs.Add("param2", "value2")
	resp, err := requests.Get("https://httpbin.org/get", qs, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(resp.Text())
}
