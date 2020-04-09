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

	"github.com/fringeproject/fringe-runner/common/assets"
)

type ModuleContext struct {
	Asset     assets.Asset
	NewAssets []assets.Asset
}

func NewModuleContext(asset string) (*ModuleContext, error) {
	ctx := ModuleContext{
		Asset: assets.Asset{
			Value: asset,
			Type:  "",
		},
		NewAssets: make([]assets.Asset, 0),
	}

	return &ctx, nil
}

// Get the current asset as a raw string
func (ctx *ModuleContext) GetAssetAsRawString() (string, error) {
	return ctx.Asset.Value, nil
}

// Check if the asset is a hostname and return it
func (ctx *ModuleContext) GetAssetAsHostname() (string, error) {
	asset, err := ctx.GetAssetAsRawString()
	if err != nil {
		return "", err
	}

	if IsHostname(asset) {
		return asset, nil
	} else {
		return "", fmt.Errorf("Current data is not a valid hostname.")
	}
}

// Check if the asset is an IP and return it
func (ctx *ModuleContext) GetAssetAsIP() (string, error) {
	asset, err := ctx.GetAssetAsRawString()
	if err != nil {
		return "", err
	}

	if IsIPv4(asset) {
		return asset, nil
	} else {
		return "", fmt.Errorf("Current data is not a valid IPv4 address.")
	}
}

// Check if the asset is a URL and return it
func (ctx *ModuleContext) GetAssetAsURL() (string, error) {
	asset, err := ctx.GetAssetAsRawString()
	if err != nil {
		return "", err
	}

	if IsURL(asset) {
		return asset, nil
	} else {
		return "", fmt.Errorf("Current data is not a valid url.")
	}
}

// Create a new asset from the module execution
func (ctx *ModuleContext) CreateNewAsset(assetValue string, assetType assets.Type) error {
	asset := assets.Asset{
		Value: assetValue,
		Type:  assetType,
	}
	ctx.NewAssets = append(ctx.NewAssets, asset)

	return nil
}

// Create a hostname from the current string without format verification
func (ctx *ModuleContext) CreateNewAssetAsHostname(hostname string) error {
	if len(hostname) == 0 {
		return fmt.Errorf("Hostname cannot be an empty string.")
	}

	return ctx.CreateNewAsset(hostname, assets.AssetTypes["hostname"])
}

// Create an IP from the current string without format verification
func (ctx *ModuleContext) CreateNewAssetAsIP(ip string) error {
	if len(ip) == 0 {
		return fmt.Errorf("IP cannot be an empty string")
	}

	return ctx.CreateNewAsset(ip, assets.AssetTypes["ip"])
}

// Create a URL from the current string without format verification
func (ctx *ModuleContext) CreateNewAssetAsURL(url string) error {
	if len(url) == 0 {
		return fmt.Errorf("URL cannot be an empty string")
	}

	return ctx.CreateNewAsset(url, assets.AssetTypes["url"])
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

func (ctx *ModuleContext) HttpRequest(method string, target string, data io.Reader, opts *HTTPOptions) (*int, *[]byte, *http.Header, error) {
	// If the HTTPOptions is not set, then use the default one
	if opts == nil {
		opts = ctx.getDefaultHTTPOptions()
	}

	// Send the request so we need the create the client then do the request
	httpClient, err := NewHTTPClient(context.Background(), opts)
	if err != nil {
		return nil, nil, nil, err
	}

	statusCode, responseBody, headers, err := httpClient.DoRequest(method, target, "", "", data)
	if err != nil {
		return nil, nil, nil, err
	}

	return statusCode, responseBody, headers, nil
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
