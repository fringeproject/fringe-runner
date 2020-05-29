package common

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"
)

type ModuleContext struct {
	Asset     Asset
	NewAssets []Asset
	NewTags   []string
	config    *FringeConfig
}

func NewModuleContext(asset Asset, config *FringeConfig) (*ModuleContext, error) {
	ctx := ModuleContext{
		Asset:     asset,
		NewAssets: make([]Asset, 0),
		NewTags:   make([]string, 0),
		config:    config,
	}

	return &ctx, nil
}

// Get a configuration variable for the module
func (ctx *ModuleContext) GetConfigurationValue(key string) (string, error) {
	value, ok := ctx.config.ModuleConfiguration[key]
	if !ok {
		return "", fmt.Errorf("Configuration variable %s is not set.", key)
	}

	return value, nil
}

// Get the path for a resource file
func (ctx *ModuleContext) GetRessourceFile(filename string) (string, error) {
	return GetRessourceFile(ctx.config, filename)
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

// Add a tag to the current asset
func (ctx *ModuleContext) AddTag(tag string) error {
	if len(tag) == 0 {
		return fmt.Errorf("Tag cannot be an empty string")
	}

	ctx.NewTags = AppendIfMissing(ctx.NewTags, tag)

	return nil
}

// Create a new asset from the module execution
func (ctx *ModuleContext) createNewAsset(assetValue string, assetType AssetType) error {
	asset := Asset{
		Value: assetValue,
		Type:  assetType,
	}
	ctx.NewAssets = append(ctx.NewAssets, asset)

	return nil
}

// Create a raw asset
func (ctx *ModuleContext) CreateNewAssetAsRaw(raw string) error {
	return ctx.createNewAsset(raw, AssetTypes["raw"])
}

// Create a hostname from the current string without format verification
func (ctx *ModuleContext) CreateNewAssetAsHostname(hostname string) error {
	if len(hostname) == 0 {
		return fmt.Errorf("Hostname cannot be an empty string.")
	}

	return ctx.createNewAsset(hostname, AssetTypes["hostname"])
}

// Create an IP from the current string without format verification
func (ctx *ModuleContext) CreateNewAssetAsIP(ip string) error {
	if len(ip) == 0 {
		return fmt.Errorf("IP cannot be an empty string")
	}

	return ctx.createNewAsset(ip, AssetTypes["ip"])
}

// Create a URL from the current string without format verification
func (ctx *ModuleContext) CreateNewAssetAsURL(url string) error {
	if len(url) == 0 {
		return fmt.Errorf("URL cannot be an empty string")
	}

	return ctx.createNewAsset(url, AssetTypes["url"])
}

func (ctx *ModuleContext) GetDefaultHTTPOptions() *HTTPOptions {
	proxy, _ := ctx.GetConfigurationValue("HTTP_PROXY")
	verify, _ := ctx.GetConfigurationValue("VERIFY_CERT")
	verifyCert, err := strconv.ParseBool(verify)
	if err != nil {
		verifyCert = true
	}

	opts := HTTPOptions{
		Proxy:          proxy,
		Timeout:        time.Second * 4,
		FollowRedirect: true,
		VerifyCert:     verifyCert,
		Headers:        make([]HTTPHeader, 0),
		WhiteListIP:    make([]string, 0),
	}

	return &opts
}

func (ctx *ModuleContext) HttpRequest(method string, target string, data io.Reader, opts *HTTPOptions) (*int, *[]byte, *http.Header, error) {
	// If the HTTPOptions is not set, then use the default one
	if opts == nil {
		opts = ctx.GetDefaultHTTPOptions()
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
	// If the HTTPOptions is not set, then use the default one
	if opts == nil {
		opts = ctx.GetDefaultHTTPOptions()
	}

	// Send the request so we need the create the client then do the request
	httpClient, err := NewHTTPClient(context.Background(), opts)
	if err != nil {
		return nil, nil, nil, err
	}

	return httpClient.DoJson(method, target, "", "", request, response)
}
