package main

import (
	"fmt"

	requests "github.com/hiroakis/go-requests"
)

func main() {
	respCh, errCh := requests.GetAsync("https://httpbin.org/get", nil, nil)

	doneCh := make(chan struct{})
	var respBody string
	go func() {
		select {
		case resp := <-respCh:
			respBody = resp.Text()
		case err := <-errCh:
			if err != nil {
				fmt.Println(err)
			}
		}
		doneCh <- struct{}{}
	}()

	fmt.Println("do something while waiting for the response")
	<-doneCh

	fmt.Println(respBody)
}
