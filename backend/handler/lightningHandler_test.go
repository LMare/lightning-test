package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	lnrpc "github.com/Lmare/lightning-playground/backend/gRPC/github.com/lightningnetwork/lnd/lnrpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

func TestTruncateUri(t *testing.T) {
	tests := []struct {
        name     string
        uri		 string
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
        uri		 string
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

// Test handleShowUri with invalid ID
func TestHandleShowUri_InvalidID(t *testing.T) {
	// Create a test request with invalid ID
	req, err := http.NewRequest("GET", "/uri?id=invalid@id#", nil)
	assert.NoError(t, err)

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleShowUri(rr, req)

	// Check the status code - should be 500 (fail)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleShowUri with missing ID
func TestHandleShowUri_MissingID(t *testing.T) {
	// Create a test request without ID
	req, err := http.NewRequest("GET", "/uri", nil)
	assert.NoError(t, err)

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleShowUri(rr, req)

	// Check the status code - should be 500 (fail)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleNodeInfo with invalid ID
func TestHandleNodeInfo_InvalidID(t *testing.T) {
	// Create a test request with invalid ID
	req, err := http.NewRequest("GET", "/info?id=node!invalid", nil)
	assert.NoError(t, err)

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleNodeInfo(rr, req)

	// Check the status code - should be 500 (fail)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleAddPeer with invalid ID
func TestHandleAddPeer_InvalidID(t *testing.T) {
	// Create a test request with invalid ID
	body := strings.NewReader("id=invalid@id&uri=02abc@host:port")
	req, err := http.NewRequest("POST", "/addpeer", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleAddPeer(rr, req)

	// Check the status code - should be 500 (fail)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleOpenChannel with invalid ID
func TestHandleOpenChannel_InvalidID(t *testing.T) {
	// Create a test request with invalid ID
	body := strings.NewReader("id=invalid#id&pubKey=02abc123&amount=100000")
	req, err := http.NewRequest("POST", "/openchannel", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleOpenChannel(rr, req)

	// Check the status code - should be 500 (fail)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleOpenChannel with invalid amount
func TestHandleOpenChannel_InvalidAmount(t *testing.T) {
	// Create a test request with invalid amount
	body := strings.NewReader("id=node1&pubKey=02abc123&amount=not_a_number")
	req, err := http.NewRequest("POST", "/openchannel", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleOpenChannel(rr, req)

	// Check the status code
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleCreateInvoice with invalid ID
func TestHandleCreateInvoice_InvalidID(t *testing.T) {
	// Create a test request with invalid ID
	body := strings.NewReader("id=bad$id&memo=test&amount=1000")
	req, err := http.NewRequest("POST", "/invoice", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleCreateInvoice(rr, req)

	// Check the status code - should be 500 (fail)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleCreateInvoice with invalid amount
func TestHandleCreateInvoice_InvalidAmount(t *testing.T) {
	// Create a test request with invalid amount
	body := strings.NewReader("id=node1&memo=test&amount=abc")
	req, err := http.NewRequest("POST", "/invoice", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleCreateInvoice(rr, req)

	// Check the status code - should be 500
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleMakePayment with invalid ID
func TestHandleMakePayment_InvalidID(t *testing.T) {
	// Create a test request with invalid ID
	body := strings.NewReader("id=invalid.id&paymentRequest=ln_invoice_123")
	req, err := http.NewRequest("POST", "/pay", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleMakePaiment(rr, req)

	// Check the status code - should be 500 (fail)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleMakePayment with missing payment request
func TestHandleMakePayment_MissingRequest(t *testing.T) {
	// Create a test request without paymentRequest
	body := strings.NewReader("id=node1")
	req, err := http.NewRequest("POST", "/pay", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleMakePaiment(rr, req)

	// Check the status code - should be 500
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleUpdateNodeAlias with invalid ID
func TestHandleUpdateNodeAlias_InvalidID(t *testing.T) {
	// Create a test request with invalid ID
	body := strings.NewReader("id=invalid!id&alias=myAlias&color=#FF0000")
	req, err := http.NewRequest("POST", "/updatealias", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleUpdateNodeAlias(rr, req)

	// Check the status code - should be 500 (fail)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}

// Test handleUpdateNodeAlias with missing ID
func TestHandleUpdateNodeAlias_MissingID(t *testing.T) {
	// Create a test request without ID
	body := strings.NewReader("alias=myAlias&color=#FF0000")
	req, err := http.NewRequest("POST", "/updatealias", body)
	assert.NoError(t, err)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	handleUpdateNodeAlias(rr, req)

	// Check the status code - should be 500 (fail)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)
}
