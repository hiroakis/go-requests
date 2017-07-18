package requests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type Redirection bool

func Redirect() *Redirection {
	return new(Redirection)
}

func (r *Redirection) Allow() *Redirection {
	allow := new(Redirection)
	*allow = true
	return allow
}

func (r *Redirection) NotAllow() *Redirection {
	allow := new(Redirection)
	*allow = false
	return allow
}

type (
	RequestParams struct {
		Data    *bytes.Buffer
		Json    interface{}
		Headers http.Header
		Cookies *cookiejar.Jar
		// files     string
		Auth           *Auth
		Timeout        *Timeout
		AllowRedirects *Redirection // redirects bool
		// proxies   string
		// verify    string
		// cert      SSLClientCert
	}
	Timeout struct {
		Connect time.Duration
		Read    time.Duration
	}
	Auth struct {
		Username string
		Password string
	}
	SSLClientCert struct {
		Crt string
		Key string
	}
)

type Response struct {
	_url          *url.URL
	status        string
	statusCode    int
	contentLength int64
	history       []http.Request
	body          *bytes.Buffer
	cookies       []*http.Cookie
	headers       http.Header
}

func redirectPolicyFunc(r *RequestParams) func(*http.Request, []*http.Request) error {

	var f func(*http.Request, []*http.Request) error
	f = func(req *http.Request, via []*http.Request) error {
		return errors.New("go-requests handles redirect")
	}

	if r == nil || r.AllowRedirects == nil {
		return f
	}

	if *r.AllowRedirects {
		f = func(req *http.Request, via []*http.Request) error {
			return errors.New("go-requests handles redirect")
		}
	} else {
		f = func(req *http.Request, via []*http.Request) error {
			// net/http: use last response
			return http.ErrUseLastResponse
		}
	}
	return f
}

const (
	defaultReadTimeout    = time.Duration(0)
	defaultConnectTimeout = time.Duration(0)
)

func timeout(r *RequestParams) (time.Duration, time.Duration) {
	if r == nil || r.Timeout == nil {
		return defaultReadTimeout, defaultConnectTimeout
	}
	return r.Timeout.Read, r.Timeout.Connect
}

func setCookie(r *RequestParams) http.CookieJar {
	if r == nil || r.Cookies == nil {
		return nil
	}
	return r.Cookies
}

func send(method, urlStr string, queryString *url.Values, r *RequestParams) (Response, error) {

	// redirect
	httpClient.client.CheckRedirect = redirectPolicyFunc(r)

	// cookie
	httpClient.client.Jar = setCookie(r)

	// timeout
	readTimeout, connTimeout := timeout(r)
	httpClient.client.Timeout = readTimeout
	var (
		ctx    context.Context
		cancel context.CancelFunc
	)
	if connTimeout == 0 {
		ctx, cancel = context.WithCancel(context.Background())
	} else {
		ctx, cancel = context.WithTimeout(context.Background(), connTimeout)
	}
	defer cancel()

	req, err := httpClient.newRequest(method, urlStr, queryString, r)
	if err != nil {
		return Response{}, err
	}
	req = req.WithContext(ctx)

	resp, err := httpClient.do(req)
	if err != nil {
		return Response{}, err
	}

	return resp, nil
}

func sendAsync(method, urlStr string, queryString *url.Values, r *RequestParams) (chan Response, chan error) {
	respCh := make(chan Response)
	errCh := make(chan error)
	go func() {
		defer func() {
			close(respCh)
			close(errCh)
		}()
		resp, err := send(method, urlStr, queryString, r)
		if err != nil {
			errCh <- err
			return
		}
		respCh <- resp
	}()

	return respCh, errCh
}

// Head makes HTTP(s) HEAD request with given urlStr, queryString and RequestParams
func Head(urlStr string, queryString *url.Values, r *RequestParams) (Response, error) {
	return send(http.MethodHead, urlStr, queryString, r)
}

// HeadAsync makes asynchronous HTTP(s) HEAD request with given urlStr, queryString and RequestParams
func HeadAsync(urlStr string, queryString *url.Values, r *RequestParams) (chan Response, chan error) {
	return sendAsync(http.MethodHead, urlStr, queryString, r)
}

// Get makes HTTP(s) GET request with given urlStr, queryString and RequestParams
func Get(urlStr string, queryString *url.Values, r *RequestParams) (Response, error) {
	return send(http.MethodGet, urlStr, queryString, r)
}

// GetAsync makes asynchronous HTTP(s) GET request with given urlStr, queryString and RequestParams
func GetAsync(urlStr string, queryString *url.Values, r *RequestParams) (chan Response, chan error) {
	return sendAsync(http.MethodGet, urlStr, queryString, r)
}

// Post makes HTTP(s) POST request with given urlStr, queryString and RequestParams
func Post(urlStr string, queryString *url.Values, r *RequestParams) (Response, error) {
	return send(http.MethodPost, urlStr, queryString, r)
}

// PostAsync makes asynchronous HTTP(s) POST request with given urlStr, queryString and RequestParams
func PostAsync(urlStr string, queryString *url.Values, r *RequestParams) (chan Response, chan error) {
	return sendAsync(http.MethodPost, urlStr, queryString, r)
}

// Put makes HTTP(s) PUT request with given urlStr, queryString and RequestParams
func Put(urlStr string, queryString *url.Values, r *RequestParams) (Response, error) {
	return send(http.MethodPut, urlStr, queryString, r)
}

// PutAsync makes asynchronous HTTP(s) PUT request with given urlStr, queryString and RequestParams
func PutAsync(urlStr string, queryString *url.Values, r *RequestParams) (chan Response, chan error) {
	return sendAsync(http.MethodPut, urlStr, queryString, r)
}

// Patch makes HTTP(s) PATCH request with given urlStr, queryString and RequestParams
func Patch(urlStr string, queryString *url.Values, r *RequestParams) (Response, error) {
	return send(http.MethodPatch, urlStr, queryString, r)
}

// PatchAsync makes asynchronous HTTP(s) PATCH request with given urlStr, queryString and RequestParams
func PatchAsync(urlStr string, queryString *url.Values, r *RequestParams) (chan Response, chan error) {
	return sendAsync(http.MethodDelete, urlStr, queryString, r)
}

// Delete makes HTTP(s) DELETE request with given urlStr, queryString and RequestParams
func Delete(urlStr string, queryString *url.Values, r *RequestParams) (Response, error) {
	return send(http.MethodPatch, urlStr, queryString, r)
}

// DeleteAsync makes asynchronous HTTP(s) DELETE request with given urlStr, queryString and RequestParams
func DeleteAsync(urlStr string, queryString *url.Values, r *RequestParams) (chan Response, chan error) {
	return sendAsync(http.MethodDelete, urlStr, queryString, r)
}

// Options makes HTTP(s) OPTIONS request with given urlStr, queryString and RequestParams
func Options(urlStr string, queryString *url.Values, r *RequestParams) (Response, error) {
	return send(http.MethodOptions, urlStr, queryString, r)
}

// OptionsAsync makes asynchronous HTTP(s) OPTIONS request with given urlStr, queryString and RequestParams
func OptionsAsync(urlStr string, queryString *url.Values, r *RequestParams) (chan Response, chan error) {
	return sendAsync(http.MethodOptions, urlStr, queryString, r)
}

func (resp Response) Url() *url.URL { return resp._url }

func (resp Response) StatusCode() int { return resp.statusCode }

func (resp Response) Status() string { return resp.status }

func (resp Response) Headers() http.Header { return resp.headers }

func (resp *Response) History() []http.Request { return resp.history }

func (resp Response) Text() string { return resp.body.String() }

func (resp Response) Content() []byte { return resp.body.Bytes() }

func (resp Response) Raw() *bytes.Buffer { return resp.body }

func (resp Response) Json(dst interface{}) error {
	return json.Unmarshal(resp.Content(), dst)
}

func (resp Response) Len() int64 { return resp.contentLength }

func (resp Response) Cookies() []*http.Cookie { return resp.cookies } // .cookies['requests-is']
