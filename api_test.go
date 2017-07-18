package requests

import (
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"encoding/json"

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
			t.Errorf("Found error: %s", err.Error())
		}
	}()
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
	assert.Equal(t, "text/plain", resp.Headers().Get("Content-Type"), "Content-type should be text/plain")
	// History
	assert.Equal(t, 0, len(resp.History()), "History should be 0")
	// Body
	assert.Equal(t, "Get Test", resp.Text(), "Response body should contain Get Test")
	// Cookie
	assert.Equal(t, 0, len(resp.Cookies()), "len(Cookies()) should be 0")
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
