package main

import (
	"fmt"
	"net/http"
	"time"

	requests "github.com/hiroakis/go-requests"
)

func main() {

	jsonData := struct {
		ID        int               `json:"id"`
		Name      string            `json:"name"`
		CreatedAt time.Time         `json:"created_at"`
		Info      map[string]string `json:"info"`
	}{
		ID:        1,
		Name:      "hiroakis",
		CreatedAt: time.Now(),
		Info: map[string]string{
			"info1": "xxxx",
			"info2": "yyyy",
		},
	}

	h := make(http.Header)
	h.Add("Content-Type", "application/json")

	resp, err := requests.Post("https://httpbin.org/post", nil, &requests.RequestParams{
		Json:    jsonData,
		Headers: h,
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Text())
}
