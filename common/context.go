package common

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type ModuleContext struct {
	Asset string
}

func NewModuleContext(asset string) (*ModuleContext, error) {
	ctx := ModuleContext{
		Asset: asset,
	}

	return &ctx, nil
}

func (ctx *ModuleContext) GetAssetAsRawString() (string, error) {
	return ctx.Asset, nil
}

func (ctx *ModuleContext) GetAssetAsHostname() (string, error) {
	asset := ctx.Asset

	if IsHostname(asset) {
		return asset, nil
	} else {
		return "", fmt.Errorf("Current data is not a valid hostname.")
	}
}

func (ctx *ModuleContext) GetAssetAsIP() (string, error) {
	asset := ctx.Asset

	if IsIPv4(asset) {
		return asset, nil
	} else {
		return "", fmt.Errorf("Current data is not a valid IPv4 address.")
	}
}

func (ctx *ModuleContext) GetAssetAsURL() (string, error) {
	asset := ctx.Asset

	if IsURL(asset) {
		return asset, nil
	} else {
		return "", fmt.Errorf("Current data is not a valid url.")
	}
}

func (ctx *ModuleContext) getDefaultHTTPOptions() *HTTPOptions {
	opts := HTTPOptions{
		Proxy:          "",
		Timeout:        time.Second * 4,
		FollowRedirect: true,
		VerifyCert:     true,
		Headers:        make([]HTTPHeader, 0),
	}

	return &opts
}

func (ctx *ModuleContext) HTTPRequestJson(method string, target string, request interface{}, response interface{}, opts *HTTPOptions) (*int, *[]byte, *http.Header, error) {
	var requestBody io.Reader

	if request != nil {
		requestBodyRequest, err := json.Marshal(request)
		if err != nil {
			logrus.Debug(err)
			return nil, nil, nil, err
		}
		requestBody = bytes.NewReader(requestBodyRequest)
	}

	// If the HTTPOptions is not set, then use the default one
	if opts == nil {
		opts = ctx.getDefaultHTTPOptions()
	}

	// Set the headers
	if opts.Headers == nil {
		opts.Headers = make([]HTTPHeader, 0)
	}

	// When we receive a body, then we need a longer timeout
	opts.Timeout = time.Second * 20

	// Check if the "Accept" header is set. We want to receive a JSON payload
	// then we need to set, as default, the value `application/json`
	if response != nil {
		foundAccept := false
		for _, header := range opts.Headers {
			if header.Name == "Accept" {
				foundAccept = true
				break
			}
		}

		if !foundAccept {
			opts.Headers = append(opts.Headers, HTTPHeader{Name: "Accept", Value: "application/json"})
		}
	}

	// Send the request so we need the create the client then do the request
	httpClient, err := NewHTTPClient(context.Background(), opts)
	if err != nil {
		logrus.Debug(err)
		return nil, nil, nil, err
	}

	statusCode, responseBody, headers, err := httpClient.DoRequest(method, target, "", "", requestBody)
	if err != nil {
		logrus.Debug(err)
		return nil, nil, nil, err
	}

	if response != nil {
		decoder := json.NewDecoder(bytes.NewReader(*responseBody))
		err = decoder.Decode(response)

		if err != nil {
			logrus.Debug(err)
			// return response info if we want to do something with it even if the
			// decoding failed
			return statusCode, responseBody, headers, err
		}
	}

	return statusCode, responseBody, headers, nil
}
