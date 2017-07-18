package requests

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	defaultUserAgent  = "go-requests/" + version
	maxRedirectCounts = 5 // http://www.ietf.org/rfc/rfc1945.txt
)

type client struct {
	client *http.Client
}

var httpClient = newClient()

func newClient() *client {
	return &client{
		client: &http.Client{},
	}
}

func (c *client) requestData(r *RequestParams) (io.ReadWriter, error) {
	if r == nil {
		return nil, nil
	}

	var (
		buf io.ReadWriter
		err error
	)

	if r.Data != nil {
		buf = r.Data
	} else if r.Json != nil {
		buf = &bytes.Buffer{}
		err = json.NewEncoder(buf).Encode(r.Json)
		if err != nil {
			return nil, err
		}
	}
	return buf, err
}

func (c *client) newRequest(method, urlStr string, queryString *url.Values, r *RequestParams) (*http.Request, error) {

	u, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}
	if queryString != nil {
		u.RawQuery = queryString.Encode()
	}

	payload, err := c.requestData(r)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, u.String(), payload)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", defaultUserAgent)
	if r != nil {
		if r.Headers != nil {
			req.Header = r.Headers
		}
		if r.Auth != nil {
			req.SetBasicAuth(r.Auth.Username, r.Auth.Password)
		}
	}

	return req, nil
}

func (c *client) do(req *http.Request) (Response, error) {
	var (
		resp    *http.Response
		history []http.Request
		cookies []*http.Cookie
		err     error
	)

	for x := 0; x < maxRedirectCounts; x++ {
		resp, err = c.client.Do(req)
		if err != nil {
			if strings.Contains(err.Error(), "go-requests handles redirect") {
				loc := resp.Header.Get("Location")
				u, err := req.URL.Parse(loc)
				if err != nil {
					return Response{}, err
				}
				history = append(history, *req)
				req.URL = u
				cookies = append(cookies, resp.Cookies()...)
				defer resp.Body.Close()
				continue
			} else {
				return Response{}, err
			}
		}
		cookies = append(cookies, resp.Cookies()...)
		defer resp.Body.Close()
		break
	}

	buf := &bytes.Buffer{}
	io.Copy(buf, resp.Body)
	response := Response{
		_url:          req.URL,
		headers:       resp.Header,
		status:        resp.Status,
		statusCode:    resp.StatusCode,
		contentLength: resp.ContentLength,
		history:       history,
		body:          buf,
		cookies:       cookies,
	}
	return response, nil
}
