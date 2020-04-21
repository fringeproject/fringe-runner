package common

// Define a custome type for an asset type
type AssetType string

// Set the values for the asset types
const (
	ASSET_IP       AssetType = "ip"
	ASSET_HOSTNAME AssetType = "hostname"
	ASSET_URL      AssetType = "url"
	ASSET_RAW      AssetType = "raw"
)

var AssetTypes = map[string]AssetType{
	"ip":       ASSET_IP,
	"hostname": ASSET_HOSTNAME,
	"url":      ASSET_URL,
	"raw":      ASSET_RAW,
}

// Represents an asset with a value and a type
type Asset struct {
	ID    string    `json:"id"`
	Value string    `json:"value"`
	Type  AssetType `json:"type"`
}

// Represents a job: a module to execute with an asset
type Job struct {
	ID     string `json:"id"`
	Module string `json:"module"`
	Asset  Asset  `json:"data"`
}
