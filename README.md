# go-requests

A HTTP library for Go, inspired by [requests module written in Python](http://docs.python-requests.org/en/master/). The usage is similer to Python requests.

# Installation

```
go get github.com/hiroakis/go-requests
```

# Usage

Basic usage

```
resp, err := requests.<Method>(urlStr, // URL string
    nil, // *url.Values QueryString
    nil) // *requests.RequestParams
```

Supported HTTP request types are GET, POST, PUT, PATCH, DELETE, HEAD and OPTIONS. Followings are simple usage.

```
# GET
resp, err := requests.Get("https://httpbin.org/", nil, nil)

# POST
postData := &bytes.Buffer{}
postData.WriteString("postdata")
resp, err := requests.Post("https://httpbin.org/", nil, &requests.RequestParams{
    Data:    postData,
})

# PUT
resp, err := requests.Put("https://httpbin.org/", nil, nil)

# PATCH
resp, err := requests.Patch("https://httpbin.org/", nil, nil)

# DELETE
resp, err := requests.Delete("https://httpbin.org/", nil, nil)

# HEAD
resp, err := requests.Head("https://httpbin.org/", nil, nil)

# OPTIONS
resp, err := requests.Options("https://httpbin.org/", nil, nil)
```

## Get

```
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
```

## Get with QueryString

```
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
```

## Post

```
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
```

## Post JSON

```
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
```

## More complicated POST requests

```
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
	headers.Add("X-requests", "i-am-go-requests")
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
```

## Async API

It has also asynchronous API.

```
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
```

## Bindata


# TODO

* File uploading
* Client certificate authentication
* Proxy support

# License

MIT