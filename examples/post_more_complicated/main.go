package main

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"

	"time"

	"bytes"

	requests "github.com/hiroakis/go-requests"
)

func main() {
	// make http headers
	var headers http.Header
	headers = make(http.Header)
	headers.Add("X-Requests", "i-am-go-requests")
	headers.Add("Content-Type", "application/x-www-form-urlencoded")

	// make post data
	postData := &bytes.Buffer{}
	postData.WriteString("postdata=valvalval")

	// make cookie
	cookieJar, _ := cookiejar.New(nil)
	u, _ := url.Parse("https://httpbin.org/post")
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:     "cookie",
		Value:    "value",
		Path:     "/",
		HttpOnly: true,
	})
	cookieJar.SetCookies(u, cookies)

	resp, err := requests.Post("https://httpbin.org/post", nil, &requests.RequestParams{
		Data:    postData,
		Headers: headers,
		Auth: &requests.Auth{
			Username: "hiroakis",
			Password: "password",
		},
		Timeout: &requests.Timeout{
			Connect: 5 * time.Second,
			Read:    10 * time.Second,
		},
		Cookies:        cookieJar,
		AllowRedirects: requests.Redirect().NotAllow(),
	})

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(resp.Text())
}
