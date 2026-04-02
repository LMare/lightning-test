package lightninghttp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

/* TODO : keep or delete ?

// MockLightningClient mocks the gRPC Lightning client
type MockLightningClient struct {
	mock.Mock
}

func (m *MockLightningClient) GetInfo(ctx context.Context, req *lnrpc.GetInfoRequest, opts ...interface{}) (*lnrpc.GetInfoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*lnrpc.GetInfoResponse), args.Error(1)
}

func (m *MockLightningClient) ConnectPeer(ctx context.Context, req *lnrpc.ConnectPeerRequest, opts ...interface{}) (*lnrpc.ConnectPeerResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*lnrpc.ConnectPeerResponse), args.Error(1)
}

*/

func TestTruncateUri(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		len      int
		expected string
	}{
		{"not uri lnd", "azertyuiop", 3, "azertyuiop"},
		{"uri lnd trunc 3", "azertyuiop@1.2.3.4:1234", 3, "aze...iop@1.2.3.4:1234"},
		{"uri lnd trunc 1", "azertyuiop@1.2.3.4:1234", 1, "a...p@1.2.3.4:1234"},
		{"uri lnd not trunc", "azertyuiop@1.2.3.4:1234", 10, "azertyuiop@1.2.3.4:1234"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncateUri(tt.uri, tt.len)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name     string
		uri      string
		len      int
		expected string
	}{
		{"string trunc 3", "azertyuiop", 3, "aze...iop"},
		{"uri lnd trunc 3", "azertyuiop@1.2.3.4:1234", 3, "aze...234"},
		{"string trunc 1", "azertyuiop", 1, "a...p"},
		{"string not trunc", "azertyuiop", 5, "azertyuiop"},
		{"string not trunc 20", "azertyuiop", 20, "azertyuiop"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := truncate(tt.uri, tt.len)
			if result != tt.expected {
				t.Errorf("got %s, want %s", result, tt.expected)
			}
		})
	}
}

// Test isValidFormatOfId
func TestIsValidFormatOfId(t *testing.T) {
	tests := []struct {
		name     string
		id       string
		expected bool
	}{
		{"valid_id_with_letters", "node1", true},
		{"valid_id_with_underscore", "lnd_node", true},
		{"valid_id_with_dash", "my-node", true},
		{"valid_id_with_all_chars", "node-1_test", true},
		{"invalid_id_with_space", "node 1", false},
		{"invalid_id_with_special_chars", "node@1", false},
		{"invalid_id_with_dot", "node.1", false},
		{"empty_id", "", false}, // empty id should be invalid
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isVaildFormatOfId(tt.id)
			assert.Equal(t, tt.expected, result)
		})
	}
}
