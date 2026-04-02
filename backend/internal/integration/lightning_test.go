package integration

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Lmare/lightning-playground/backend/internal/testsupport"
	"github.com/stretchr/testify/assert"
)

// Test handleShowUri with invalid ID
func TestHandleShowUri_InvalidID(t *testing.T) {

	h := testsupport.NewTestHarness()

	// Create a test request with invalid ID
	req := httptest.NewRequest("GET", "/lightning/uri?id=invalid@id#", nil)

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleShowUri with missing ID
func TestHandleShowUri_MissingID(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request without ID
	req := httptest.NewRequest("GET", "/lightning/uri", nil)

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleNodeInfo with invalid ID
func TestHandleNodeInfo_InvalidID(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request with invalid ID
	req := httptest.NewRequest("GET", "/lightning/nodeInfo?id=node!invalid", nil)

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleAddPeer with invalid ID
func TestHandleAddPeer_InvalidID(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request with invalid ID
	body := strings.NewReader("id=invalid@id&uri=02abc@host:port")
	req := httptest.NewRequest("POST", "/lightning/peer", body)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleOpenChannel with invalid ID
func TestHandleOpenChannel_InvalidID(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request with invalid ID
	body := strings.NewReader("id=invalid#id&pubKey=02abc123&amount=100000")
	req := httptest.NewRequest("POST", "/lightning/channel", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleOpenChannel with invalid amount
func TestHandleOpenChannel_InvalidAmount(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request with invalid amount
	body := strings.NewReader("id=node1&pubKey=02abc123&amount=not_a_number")
	req := httptest.NewRequest("POST", "/lightning/channel", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleCreateInvoice with invalid ID
func TestHandleCreateInvoice_InvalidID(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request with invalid ID
	body := strings.NewReader("id=bad$id&memo=test&amount=1000")
	req := httptest.NewRequest("POST", "/lightning/invoice", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleCreateInvoice with invalid amount
func TestHandleCreateInvoice_InvalidAmount(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request with invalid amount
	body := strings.NewReader("id=node1&memo=test&amount=abc")
	req := httptest.NewRequest("POST", "/lightning/invoice", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleMakePayment with invalid ID
func TestHandleMakePayment_InvalidID(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request with invalid ID
	body := strings.NewReader("id=invalid.id&paymentRequest=ln_invoice_123")
	req := httptest.NewRequest("POST", "/lightning/paiment", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleMakePayment with missing payment request
func TestHandleMakePayment_MissingRequest(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request without paymentRequest
	body := strings.NewReader("id=node1")
	req := httptest.NewRequest("POST", "/lightning/paiment", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleUpdateNodeAlias with invalid ID
func TestHandleUpdateNodeAlias_InvalidID(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request with invalid ID
	body := strings.NewReader("id=invalid!id&alias=myAlias&color=#FF0000")
	req := httptest.NewRequest("PUT", "/lightning/alias", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

// Test handleUpdateNodeAlias with missing ID
func TestHandleUpdateNodeAlias_MissingID(t *testing.T) {
	h := testsupport.NewTestHarness()
	// Create a test request without ID
	body := strings.NewReader("alias=myAlias&color=#FF0000")
	req := httptest.NewRequest("PUT", "/lightning/alias", body)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Create a response writer
	rr := httptest.NewRecorder()

	// Call the handler
	h.Router.ServeHTTP(rr, req)

	// Check the status code - should be 400 (Bad Request)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
