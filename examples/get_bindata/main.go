package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	requests "github.com/hiroakis/go-requests"
)

func main() {
	resp, err := requests.Get("https://httpbin.org/image/png", nil, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	img := resp.Raw()

	out, err := os.Create("image.png")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer out.Close()

	io.Copy(out, bufio.NewReader(img))
}
