package port

// NodeRepository is the outbound port listing LND node identifiers.
type NodeRepository interface {
	GetNodesIds() ([]string, error)
}
