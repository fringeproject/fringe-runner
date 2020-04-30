package common

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	ContentTypeHeaderName = "Content-Type"
	AcceptHeaderName      = "Accept"
	ApplicationJSON       = "application/json"
)

// Represents options for the HTTP client
// By default a user pass a HTTP option to instanciate a HTTP client. This struct
// must contains every options the user needs to pass to configure the client.
type HTTPOptions struct {
	// We use a default User-Agent but a user may need a custom one to bypass
	// some verification
	UserAgent string
	// Custom headers
	Headers []HTTPHeader
	// Some use need a custom timeout, specially on a brute-force attack
	Timeout time.Duration
	// Do e follow redirection on 301/302 ?
	FollowRedirect bool
	// Set a proxy for the HTTP and HTTPS requests
	Proxy string
	// Do we need to check the SSL/TLS server certificate. Can be usefull when
	// you use a custom proxy like Burp
	VerifyCert bool
	// White list for IP addresses
	WhiteListIP []string
}

// Represents a (custom) HTTP header
type HTTPHeader struct {
	Name  string
	Value string
}

// Represents an HTTP client
type HTTPClient struct {
	// The internal Go HTTP client
	client *http.Client
	// A context to stop the request
	context context.Context
	// Set this here instead of in the headers field
	userAgent string
	// Custom headers
	headers []HTTPHeader
}

func NewHTTPClient(c context.Context, opt *HTTPOptions) (*HTTPClient, error) {
	// Create the variables for the HTTP client
	var client HTTPClient

	// Here some defaut functions
	var proxyURLFunc func(*http.Request) (*url.URL, error)
	var redirectFunc func(req *http.Request, via []*http.Request) error

	// First we check the options before continuing
	if opt == nil {
		return nil, fmt.Errorf("HTTP options cannot be nil")
	}

	// Parse the Proxy URL
	if opt.Proxy != "" {
		proxyURL, err := url.Parse(opt.Proxy)
		if err != nil {
			return nil, fmt.Errorf("Proxy URL is invalid (%v)", err)
		}
		proxyURLFunc = http.ProxyURL(proxyURL)
	} else {
		proxyURLFunc = http.ProxyFromEnvironment
	}

	// Check the redirection
	if !opt.FollowRedirect {
		redirectFunc = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	} else {
		redirectFunc = nil
	}

	// Generate a dialer to white/black list the IP
	if opt.WhiteListIP == nil {
		opt.WhiteListIP = []string{}
	}

	dialer, err := GenerateDialer(opt.WhiteListIP)
	if err != nil {
		return nil, err
	}

	// Intanciate the HTTP client
	client.client = &http.Client{
		Timeout:       opt.Timeout,
		CheckRedirect: redirectFunc,
		Transport: &http.Transport{
			Proxy:               proxyURLFunc,
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 100,
			TLSClientConfig: &tls.Config{
				MinVersion:         tls.VersionTLS10,
				InsecureSkipVerify: !opt.VerifyCert,
			},
			Dial: dialer,
		},
	}
	client.context = c
	client.userAgent = opt.UserAgent
	client.headers = opt.Headers

	return &client, nil
}

func (client *HTTPClient) DoRequest(method, target, host, cookie string, data io.Reader) (*int, *[]byte, *http.Header, error) {
	resp, err := client.makeRequest(method, target, host, cookie, data)
	if err != nil {
		if client.context.Err() == context.Canceled {
			return nil, nil, nil, nil
		}
		return nil, nil, nil, err
	}
	defer resp.Body.Close()

	// Becareful! Even if we don't want the response body then we read it because
	// Go reuse connections.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("could not read body: %v", err)
	}

	return &resp.StatusCode, &body, &resp.Header, nil
}

func (client *HTTPClient) Get(target, host, cookie string, data io.Reader) (*int, *[]byte, *http.Header, error) {
	return client.DoRequest(http.MethodGet, target, host, cookie, data)
}

func (client *HTTPClient) Trace(target, host, cookie string, data io.Reader) (*int, *[]byte, *http.Header, error) {
	return client.DoRequest(http.MethodTrace, target, host, cookie, data)
}

func (client *HTTPClient) Post(target, host, cookie string, data io.Reader) (*int, *[]byte, *http.Header, error) {
	return client.DoRequest(http.MethodPost, target, host, cookie, data)
}

func (client *HTTPClient) Put(target, host, cookie string, data io.Reader) (*int, *[]byte, *http.Header, error) {
	return client.DoRequest(http.MethodPut, target, host, cookie, data)
}

func (client *HTTPClient) DoJson(method, target, host, cookie string, request interface{}, response interface{}) (*int, *[]byte, *http.Header, error) {
	// As we need to add 2 custom headers, first we check the field is not nil
	if client.headers == nil {
		client.headers = make([]HTTPHeader, 0)
	}

	var requestBody io.Reader
	if request != nil {
		requestBodyRequest, err := json.Marshal(request)
		if err != nil {
			return nil, nil, nil, err
		}
		requestBody = bytes.NewReader(requestBodyRequest)

		foundContentType := false
		for _, header := range client.headers {
			if header.Name == ContentTypeHeaderName {
				foundContentType = true
				break
			}
		}

		if !foundContentType {
			client.headers = append(client.headers, HTTPHeader{Name: ContentTypeHeaderName, Value: ApplicationJSON})
		}
	}

	// Check if the "Accept" header is set. We want to receive a JSON payload
	// then we need to set, as default, the value `application/json`
	if response != nil {
		foundAccept := false
		for _, header := range client.headers {
			if header.Name == AcceptHeaderName {
				foundAccept = true
				break
			}
		}

		if !foundAccept {
			client.headers = append(client.headers, HTTPHeader{Name: AcceptHeaderName, Value: ApplicationJSON})
		}
	}

	// Send the payload
	statusCode, responseBody, headers, err := client.DoRequest(method, target, host, cookie, requestBody)
	if err != nil {
		return nil, nil, nil, err
	}

	// We've got a JSON response, let's deserialize it
	if response != nil {
		decoder := json.NewDecoder(bytes.NewReader(*responseBody))
		err = decoder.Decode(response)

		if err != nil {
			// return response info if we want to do something with it even if the
			// decoding failed
			return statusCode, responseBody, headers, err
		}
	}

	return statusCode, responseBody, headers, nil
}

// Make the request and perform it based on the client options and arguments.
func (client *HTTPClient) makeRequest(method, target, host, cookie string, data io.Reader) (*http.Response, error) {
	// Crate the new request. `data` can be nil.
	req, err := http.NewRequest(method, target, data)
	if err != nil {
		return nil, err
	}

	// Add the context to the request so we can cancel it
	req = req.WithContext(client.context)

	// Now we set the HTTP headers:
	// The user can override the `Host` header
	if host != "" {
		req.Host = host
	}

	// Add the `Cookie` header
	if cookie != "" {
		req.Header.Set("Cookie", cookie)
	}

	// The user can set the `User-Agent` header or we use default one
	if client.userAgent != "" {
		req.Header.Set("User-Agent", client.userAgent)
	} else {
		req.Header.Set("User-Agent", DefaultUserAgent())
	}

	// Then the other headers. This call is at the end because the use may want
	// to override the the previous headers (Host, Cookie and User-Agent) using
	// the client `headers` field.
	for _, h := range client.headers {
		req.Header.Set(h.Name, h.Value)
	}

	// Perform the request and check the errors
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
