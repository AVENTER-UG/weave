package types

// Config is a struct of the framework configuration
type Config struct {
	APIURL   string
	SSL      bool
	APIToken string
}

// WeaveNetwork is a struct that holds the customer information for weave
type WeaveNetwork struct {
	Customer           string `json:"customer"`               // Link → Customer
	Name               string `json:"name"`                   // auto Name / slug
	OverlayCIDR        string `json:"overlay_cidr,omitempty"` // optional
	EncryptionPassword string `json:"encryption_password"`    // encrypted field
	Status             string `json:"status"`                 // Active / Paused / Archived
	PeerCount          int    `json:"peer_count,omitempty"`   // computed, read-only
	LastActivity       string `json:"last_activity,omitempty"`
}
