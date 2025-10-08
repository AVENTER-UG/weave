package types

import (
	"time"
)

// Config is a struct of the framework configuration
type Config struct {
	APIHost  string
	SSL      bool
	APIToken string
}

// WeaveNetwork is a struct that holds the customer information for weave
type WeaveNetwork struct {
	Customer           string      `json:"customer"`
	Name               string      `json:"name"`
	OverlayCIDR        string      `json:"overlay_cidr,omitempty"`
	EncryptionPassword string      `json:"encryption_password"`
	Status             string      `json:"status"`
	PeerCount          int         `json:"peer_count,omitempty"`
	LastActivity       string      `json:"last_activity,omitempty"`
	Peers              []WeavePeer `json:"peers"`
}

type WeavePeer struct {
	Hostname string     `json:"hostname"`
	Status   string     `json:"status,omitempty"`
	LastSeen *time.Time `json:"last_seen,omitempty"`
}

type WeaveAllowedIP struct {
	CIDR string `json:"cidr"`
}
