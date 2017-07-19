package main

import (
	"fmt"

	"net/url"

	requests "github.com/hiroakis/go-requests"
)

type HttpBin struct {
	Args    map[string]string      `json:"args"`
	Data    string                 `json:"data"`
	Files   []string               `json:"files"`
	Form    []string               `json:"form"`
	Headers map[string]string      `json:"headers"`
	Json    map[string]interface{} `json:"json"`
	Origin  string                 `json:"origin"`
	Url     string                 `json:"url"`
}

func main() {
	params := &url.Values{}
	params.Add("args1", "argVal1")
	params.Add("args2", "argVal2")
	resp, err := requests.Get("https://httpbin.org/get", params, nil)
	if err != nil {
		fmt.Println(err)
		return
	}

	var respJson HttpBin
	if err := resp.Json(&respJson); err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(respJson)
	fmt.Println("----- Query string -----")
	fmt.Println(respJson.Args["args1"])
	fmt.Println(respJson.Args["args2"])
	fmt.Println("----- Response headers -----")
	for k := range respJson.Headers {
		fmt.Printf("%s: %s\n", k, respJson.Headers[k])
	}
	fmt.Println("----- Origin -----")
	fmt.Println(respJson.Origin)
	fmt.Println("----- URL -----")
	fmt.Println(respJson.Url)
}
