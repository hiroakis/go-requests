package requests

import (
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"encoding/json"

	"bytes"

	"bufio"

	"github.com/stretchr/testify/assert"
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

func TestHead(t *testing.T) {
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Head Test"))

	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := Head(ts.URL, nil, nil)

	assert.Nil(t, err)
	// Url()
	assert.Equal(t, ts.URL, resp.Url().Scheme+"://"+resp.Url().Host, "")
	assert.Empty(t, resp.Url().Path, "Requested path should be empty")
	assert.Empty(t, resp.Url().RawQuery, "Query String should be empty")
	// Status
	assert.Equal(t, 200, resp.StatusCode(), "Response code should be 200")
	assert.Equal(t, "200 OK", resp.Status(), "Response status line should be 200 OK")
	// Header
	assert.Equal(t, "text/plain", resp.Headers().Get("Content-Type"), "Content-type should be text/plain")
	// History
	assert.Equal(t, 0, len(resp.History()), "History should be 0")
	// Body
	assert.Empty(t, resp.Text(), "Response body should be empty")
	// Cookie
	assert.Equal(t, 0, len(resp.Cookies()), "len(Cookies()) should be 0")
}

func TestHeadAsync(t *testing.T) {
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Head Test"))

	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	respCh, errCh := HeadAsync(ts.URL, nil, nil)
	var (
		resp Response
		err  error
	)

	doneCh := make(chan struct{})
	go func() {
		select {
		case resp = <-respCh:
			// Url()
			assert.Equal(t, ts.URL, resp.Url().Scheme+"://"+resp.Url().Host, "")
			assert.Empty(t, resp.Url().Path, "Requested path should be empty")
			assert.Empty(t, resp.Url().RawQuery, "Query String should be empty")
			// Status
			assert.Equal(t, 200, resp.StatusCode(), "Response code should be 200")
			assert.Equal(t, "200 OK", resp.Status(), "Response status line should be 200 OK")
			// Header
			assert.Equal(t, "text/plain", resp.Headers().Get("Content-Type"), "Content-type should be text/plain")
			// History
			assert.Equal(t, 0, len(resp.History()), "History should be 0")
			// Body
			assert.Empty(t, resp.Text(), "Response body should be empty")
			// Cookie
			assert.Equal(t, 0, len(resp.Cookies()), "len(Cookies()) should be 0")
		case err = <-errCh:
			assert.Nil(t, err)
		}
		doneCh <- struct{}{}
	}()
	<-doneCh
	assert.Nil(t, err)

}

func TestGet(t *testing.T) {
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Get Test"))

	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := Get(ts.URL, nil, nil)
	assert.Nil(t, err)
	// Url()
	assert.Equal(t, ts.URL, resp.Url().Scheme+"://"+resp.Url().Host, "")
	assert.Empty(t, resp.Url().Path, "Requested path should be empty")
	assert.Empty(t, resp.Url().RawQuery, "Query String should be empty")
	// Status
	assert.Equal(t, 200, resp.StatusCode(), "Response code should be 200")
	assert.Equal(t, "200 OK", resp.Status(), "Response status line should be 200 OK")
	// Header
	assert.Equal(t, "text/plain", resp.Headers().Get("Content-Type"), "Content-type should be text/plain")
	// History
	assert.Equal(t, 0, len(resp.History()), "History should be 0")
	// Body
	assert.Equal(t, "Get Test", resp.Text(), "Response body should contain Get Test")
	// Cookie
	assert.Equal(t, 0, len(resp.Cookies()), "len(Cookies()) should be 0")
}

func TestGetWithQueryString(t *testing.T) {
	var (
		val1 string
		val2 string
	)
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		val1 = r.FormValue("key1")
		val2 = r.FormValue("key2")
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Get Test"))

	})
	ts := httptest.NewServer(handler)
	defer ts.Close()
	qs := &url.Values{}
	qs.Add("key1", "value1")
	qs.Add("key2", "value2")
	resp, err := Get(ts.URL, qs, nil)
	assert.Nil(t, err)
	// Url()
	assert.Equal(t, ts.URL, resp.Url().Scheme+"://"+resp.Url().Host, "")
	assert.Empty(t, resp.Url().Path, "Requested path should be empty")
	assert.Equal(t, "value1", val1, "Query String should be empty")
	assert.Equal(t, "value2", val2, "Query String should be empty")
	// Status
	assert.Equal(t, 200, resp.StatusCode(), "Response code should be 200")
	assert.Equal(t, "200 OK", resp.Status(), "Response status line should be 200 OK")
	// Header
	assert.Equal(t, "text/plain", resp.Headers().Get("Content-Type"), "Content-type should be text/plain")
	// History
	assert.Equal(t, 0, len(resp.History()), "History should be 0")
	// Body
	assert.Equal(t, "Get Test", resp.Text(), "Response body should contain Get Test")
	// Cookie
	assert.Equal(t, 0, len(resp.Cookies()), "len(Cookies()) should be 0")
}

func TestGetWithSpecifiedUserAgent(t *testing.T) {
	var ua string
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua = r.Header.Get("User-Agent")
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Get Test"))

	})
	ts := httptest.NewServer(handler)
	defer ts.Close()
	headers := make(http.Header)
	headers.Add("User-Agent", "Test-UA")
	_, err := Get(ts.URL, nil, &RequestParams{
		Headers: headers,
	})
	assert.Nil(t, err)
	assert.Equal(t, "Test-UA", ua, "User-Agent should be Test-UA")
}

func TestGetJson(t *testing.T) {
	type Body struct {
		ID            int               `json:"id"`
		Name          string            `json:"name"`
		CreatedAt     time.Time         `json:"created_at"`
		AddtionalInfo map[string]string `json:"additional_info"`
	}
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		body := Body{
			ID:        1,
			Name:      "hiroakis",
			CreatedAt: time.Now(),
			AddtionalInfo: map[string]string{
				"FavoriteFood": "Sushi, Yakitori, Ajillo, Curry and ลาบไก่",
			},
		}
		b, _ := json.Marshal(&body)
		w.Write(b)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := Get(ts.URL, nil, nil)
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)
	// Url()
	assert.Equal(t, ts.URL, resp.Url().Scheme+"://"+resp.Url().Host, "")
	assert.Empty(t, resp.Url().Path, "Requested path should be empty")
	assert.Empty(t, resp.Url().RawQuery, "Query String should be empty")
	// Status
	assert.Equal(t, 200, resp.StatusCode(), "Response code should be 200")
	assert.Equal(t, "200 OK", resp.Status(), "Response status line should be 200 OK")
	// Header
	assert.Equal(t, "application/json", resp.Headers().Get("Content-Type"), "Content-type should be application/json")
	// History
	assert.Equal(t, 0, len(resp.History()), "History should be 0")
	// Body
	assert.Equal(t, true, strings.Contains(resp.Text(), "hiroakis"), "Response body should contain hiroakis")
	var jsonBody Body
	err = resp.Json(&jsonBody)
	assert.Nil(t, err)
	assert.Equal(t, 1, jsonBody.ID, "Response body.id should 1")
	assert.Equal(t, "hiroakis", jsonBody.Name, "Response body.name should hiroakis")
	assert.NotEqual(t, 0, jsonBody.CreatedAt.Unix(), "Response body.created_at should not 0")
	assert.Equal(t, "Sushi, Yakitori, Ajillo, Curry and ลาบไก่", jsonBody.AddtionalInfo["FavoriteFood"], "Response body.additional_info should have key FavoriteFood")
	// Cookie
	assert.Equal(t, 0, len(resp.Cookies()), "len(Cookies()) should be 0")
}

func TestBindata(t *testing.T) {

	var sizeTestPng int64
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		f, err := os.Open("testdata/test.png")
		if err != nil {
			t.Errorf("Open error: %s", err.Error())
		}
		stat, _ := f.Stat()
		sizeTestPng = stat.Size()
		defer f.Close()

		reader := bufio.NewReader(f)
		io.Copy(w, reader)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := Get(ts.URL, nil, nil)
	assert.Nil(t, err)

	out, err := os.Create("testdata/out.png")
	assert.Nil(t, err)
	defer func() {
		out.Close()
		os.Remove("testdata/out.png")
	}()

	io.Copy(out, bufio.NewReader(resp.Raw()))
	f, _ := os.Open("testdata/out.png")
	stat, _ := f.Stat()
	assert.Equal(t, sizeTestPng, stat.Size(), "Should be equal")
}

func TestGetCookie(t *testing.T) {
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.SetCookie(w, &http.Cookie{
			Name:  "cookie1",
			Value: "value1",
			Path:  "/",
		})
		http.SetCookie(w, &http.Cookie{
			Name:  "cookie2",
			Value: "value2",
			Path:  "/cookie2_path",
		})
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := Get(ts.URL, nil, nil)
	if err != nil {
		fmt.Println(err)
	}
	assert.Nil(t, err)
	// Status
	assert.Equal(t, 200, resp.StatusCode(), "Response code should be 200")
	// Cookie
	assert.Equal(t, "cookie1", resp.Cookies()[0].Name, "Should have cookie1")
	assert.Equal(t, "value1", resp.Cookies()[0].Value, "Cookie1 Should have value1")
	assert.Equal(t, "/", resp.Cookies()[0].Path, "Cookie1 Should have path /")
	assert.Equal(t, "cookie2", resp.Cookies()[1].Name, "Should have cookie2")
	assert.Equal(t, "value2", resp.Cookies()[1].Value, "Cookie2 should have value2")
	assert.Equal(t, "/cookie2_path", resp.Cookies()[1].Path, "Cookie1 Should have path /cookie2_path")
}

func TestRedirectAllow(t *testing.T) {

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Found"))
	})
	ts := httptest.NewServer(handler)
	red2Ts := httptest.NewServer(http.RedirectHandler(ts.URL, 302))
	red1Ts := httptest.NewServer(http.RedirectHandler(red2Ts.URL, 301))
	defer func() {
		ts.Close()
		red1Ts.Close()
		red2Ts.Close()
	}()

	resp, err := Get(red1Ts.URL, nil, nil) // Default: &RequestParams{AllowRedirects: Redirect().Allow(),}
	assert.Nil(t, err)
	// status
	assert.Equal(t, 200, resp.StatusCode(), "Response code should be 200")
	assert.Empty(t, resp.Headers().Get("Location"), "Should not have Location Header")
	assert.Equal(t, 2, len(resp.History()), "History should be 2")
	assert.Equal(t, red1Ts.URL, "http://"+resp.History()[0].URL.Host, "History should be 0")
	assert.Equal(t, red2Ts.URL, "http://"+resp.History()[1].URL.Host, "History should be 0")
}

func TestRedirectNowAllow(t *testing.T) {

	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Found"))
	})
	ts := httptest.NewServer(handler)
	red2Ts := httptest.NewServer(http.RedirectHandler(ts.URL, 302))
	red1Ts := httptest.NewServer(http.RedirectHandler(red2Ts.URL, 301))
	defer func() {
		ts.Close()
		red1Ts.Close()
		red2Ts.Close()
	}()

	resp, err := Get(red1Ts.URL, nil, &RequestParams{
		AllowRedirects: Redirect().NotAllow(),
	})
	assert.Nil(t, err)
	// Status
	assert.Equal(t, 301, resp.StatusCode(), "Response code should be 301")
	assert.NotEmpty(t, resp.Headers().Get("Location"), "Should have Location Header")
	assert.Equal(t, 0, len(resp.History()), "History should be 0")
}

func TestSendCookie(t *testing.T) {
	var receivedCookies []*http.Cookie
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedCookies = r.Cookies()
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(ts.URL)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:  "cookie_name",
		Value: "value",
	}
	cookies = append(cookies, cookie)
	jar.SetCookies(u, cookies)
	r := &RequestParams{
		Cookies: jar,
	}
	_, err := Get(ts.URL, nil, r)
	assert.Nil(t, err)
	assert.Equal(t, "cookie_name", receivedCookies[0].Name, "")
	assert.Equal(t, "value", receivedCookies[0].Value, "")
}

func TestReadTimeout(t *testing.T) {
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1000 * time.Millisecond)
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	r := &RequestParams{
		Timeout: &Timeout{
			Read:    time.Duration(500) * time.Millisecond,
			Connect: time.Duration(5000) * time.Millisecond,
		},
	}
	resp, err := Get(ts.URL, nil, r)
	assert.NotNil(t, err)
	assert.Empty(t, resp.Raw())
}

func TestConnectTimeout(t *testing.T) {
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(1000 * time.Millisecond)
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	r := &RequestParams{
		Timeout: &Timeout{
			Read:    time.Duration(5000) * time.Millisecond,
			Connect: time.Duration(500) * time.Millisecond,
		},
	}
	resp, err := Get(ts.URL, nil, r)
	assert.NotNil(t, err)
	assert.Empty(t, resp.Raw())
}

func TestGetAsync(t *testing.T) {
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/plain")
		w.Write([]byte("Get Test"))

	})
	ts := httptest.NewServer(handler)
	defer ts.Close()

	resp, err := GetAsync(ts.URL, nil, nil)

	doneCh := make(chan struct{})
	var (
		code  int
		cType string
		body  string
	)
	go func() {
		select {
		case r := <-resp:
			code = r.StatusCode()
			cType = r.Headers().Get("Content-Type")
			body = r.Text()
		case e := <-err:
			assert.Nil(t, e)
		}
		doneCh <- struct{}{}
	}()
	<-doneCh
	assert.Equal(t, 200, code, "Response code should be 200")
	assert.Equal(t, "text/plain", cType, "Content-type should be text/plain")
	// Body
	assert.Equal(t, "Get Test", body, "Response body should contain Get Test")
}

// POST
func TestPost(t *testing.T) {
	var buf *bytes.Buffer
	var cType string
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body := r.Body
		cType = r.Header.Get("Content-Type")
		defer body.Close()
		buf = new(bytes.Buffer)
		io.Copy(buf, body)
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	headers := make(http.Header)
	headers.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := Post(ts.URL, nil, &RequestParams{
		Data:    bytes.NewBufferString("a=b"),
		Headers: headers,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode(), "")
	assert.Equal(t, "a=b", buf.String(), "")
	assert.Equal(t, "application/x-www-form-urlencoded", cType, "")
}

func TestPostJson(t *testing.T) {
	type PostData struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	}
	var postData PostData
	var cType string
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cType = r.Header.Get("Content-Type")
		json.NewDecoder(r.Body).Decode(&postData)
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	headers := make(http.Header)
	headers.Add("Content-Type", "application/json")
	resp, err := Post(ts.URL, nil, &RequestParams{
		Json: PostData{
			ID:   1,
			Name: "hiroakis",
		},
		Headers: headers,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode(), "Should be 200")
	assert.Equal(t, 1, postData.ID, "Should be 1")
	assert.Equal(t, "hiroakis", postData.Name, "Should be hiroakis")
	assert.Equal(t, "application/json", cType, "Should be application/json")
}

func TestPut(t *testing.T) {
	var method string
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	resp, err := Put(ts.URL, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode(), "Should be 200")
	assert.Equal(t, "PUT", method, "method should be PUT")
}

func TestPutAsync(t *testing.T) {
	var method string
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	resp, err := PutAsync(ts.URL, nil, nil)

	doneCh := make(chan struct{})
	var code int
	var e error
	go func() {
		select {
		case r := <-resp:
			code = r.StatusCode()
		case e = <-err:
		}
		doneCh <- struct{}{}
	}()
	<-doneCh
	assert.Nil(t, e)
	assert.Equal(t, 200, code, "Should be 200")
	assert.Equal(t, "PUT", method, "Method should be PUT")
}

func TestPatch(t *testing.T) {
	var method string
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	resp, err := Patch(ts.URL, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode(), "Should be 200")
	assert.Equal(t, "PATCH", method, "method should be PATCH")
}

func TestPatchAsync(t *testing.T) {
	var method string
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	resp, err := PatchAsync(ts.URL, nil, nil)

	doneCh := make(chan struct{})
	var code int
	var e error
	go func() {
		select {
		case r := <-resp:
			code = r.StatusCode()
		case e = <-err:
		}
		doneCh <- struct{}{}
	}()
	<-doneCh
	assert.Nil(t, e)
	assert.Equal(t, 200, code, "Should be 200")
	assert.Equal(t, "PATCH", method, "Method should be PATCH")
}

func TestDelete(t *testing.T) {
	var method string
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	resp, err := Delete(ts.URL, nil, nil)
	assert.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode(), "")
	assert.Equal(t, "DELETE", method, "method should be DELETE")
}

func TestDeleteAsync(t *testing.T) {
	var method string
	var handler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		method = r.Method
		w.WriteHeader(200)
	})
	ts := httptest.NewServer(handler)

	resp, err := DeleteAsync(ts.URL, nil, nil)

	doneCh := make(chan struct{})
	var code int
	var e error
	go func() {
		select {
		case r := <-resp:
			code = r.StatusCode()
		case e = <-err:
		}
		doneCh <- struct{}{}
	}()
	<-doneCh
	assert.Nil(t, e)
	assert.Equal(t, 200, code, "Should be 200")
	assert.Equal(t, "DELETE", method, "method should be DELETE")
}

func TestOptions(t *testing.T) {
	resp, err := Options("https://httpbin.org", nil, nil)
	assert.Nil(t, err, "")
	assert.Equal(t, "GET, POST, PUT, DELETE, PATCH, OPTIONS",
		resp.Headers().Get("Access-Control-Allow-Methods"), "")
}
