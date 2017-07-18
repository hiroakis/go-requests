package main

import (
	"bytes"
	"fmt"

	requests "github.com/hiroakis/go-requests"
)

func main() {

	postData := &bytes.Buffer{}
	postData.WriteString("postdata")

	resp, err := requests.Post("https://httpbin.org/post", nil, &requests.RequestParams{
		Data: postData,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Text())
}
