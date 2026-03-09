package lightningService

type InfoLndNode struct {
	Alias               string   `json:"alias"`
	Color               string   `json:"color"`
	NumPendingChannels  uint32   `json:"numPendingChannels"`
	NumActiveChannels   uint32   `json:"numActiveChannels"`
	NumInactiveChannels uint32   `json:"numInactiveChannels"`
	NumPeers            uint32   `json:"numPeers"`
	BlockHeight         uint32   `json:"blockHeight"`
	Network             string   `json:"network"`
	Uris                []string `json:"uris"`
	SyncedToChain       bool     `json:"syncedToChain"`
	SyncedToGraph       bool     `json:"syncedToGraph"`
}

type NodeBasicInfo struct {
	Id    string `json:"id"`
	Alias string `json:"alias"`
	Color string `json:"color"`
}
