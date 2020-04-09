package assets

// Define a custome type for an asset type
type Type string

// Set the values for the asset types
const (
	IP       Type = "ip"
	HOSTNAME Type = "hostname"
	URL      Type = "url"
	RAW      Type = "raw"
)

var AssetTypes = map[string]Type{
	"ip":       IP,
	"hostname": HOSTNAME,
	"url":      URL,
	"raw":      RAW,
}

// Represents an asset with a value and a type
type Asset struct {
	Value string `json:"value"`
	Type  Type   `json:"type"`
}
