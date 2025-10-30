package slip2go

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// Slip2GoClient handles communication with Slip2Go API
type Slip2GoClient struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	MockMode   bool // If true, skip real API call and return mock response
}

// VerifySlipRequest contains parameters for slip verification
type VerifySlipRequest struct {
	FileData      io.Reader
	Filename      string
	AccountType   string  // e.g., "01004" for SCB
	AccountName   string  // Account holder name in Thai
	AccountNumber string  // Bank account number
	Amount        float64 // Expected amount
	PaymentDate   string  // Optional: ISO date (YYYY-MM-DD)
}

// Slip2GoResponse represents the response from Slip2Go API
type Slip2GoResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  struct {
		Verified  bool   `json:"verified"`
		SlipID    string `json:"slip_id"`
		Amount    string `json:"amount"`
		Sender    string `json:"sender"`
		Receiver  string `json:"receiver"`
		Timestamp string `json:"timestamp"`
	} `json:"result"`
}

// NewSlip2GoClient creates a new Slip2Go API client
func NewSlip2GoClient() *Slip2GoClient {
	apiKey := os.Getenv("SLIP2GO_SECRET_KEY")
	mockMode := os.Getenv("MOCK_SLIP2GO") == "true"

	return &Slip2GoClient{
		APIKey:  apiKey,
		BaseURL: "https://connect.slip2go.com/api/verify-slip/qr-image/info",
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		MockMode: mockMode,
	}
}

// VerifySlip sends a slip image to Slip2Go for verification
func (c *Slip2GoClient) VerifySlip(req VerifySlipRequest) (*Slip2GoResponse, error) {
	// 🧪 Mock Mode - Skip real API call (for development/testing)
	if c.MockMode {
		return c.mockVerifySlip(req), nil
	}

	// 🔑 Check API Key
	if c.APIKey == "" {
		return nil, fmt.Errorf("SLIP2GO_SECRET_KEY environment variable is not set")
	}

	// 📦 Prepare multipart/form-data
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 1. Add file to form-data
	part, err := writer.CreateFormFile("file", filepath.Base(req.Filename))
	if err != nil {
		return nil, fmt.Errorf("failed to create form file: %w", err)
	}

	if _, err := io.Copy(part, req.FileData); err != nil {
		return nil, fmt.Errorf("failed to copy file data: %w", err)
	}

	// 2. Add payload for verification rules
	checkPayload := map[string]interface{}{
		"checkDuplicate": true, // Prevent duplicate slips
	}

	// Check receiver (bank account)
	if req.AccountType != "" && req.AccountName != "" && req.AccountNumber != "" {
		checkPayload["checkReceiver"] = []map[string]string{{
			"accountType":   req.AccountType,
			"accountNameTH": req.AccountName,
			"accountNumber": req.AccountNumber,
		}}
	}

	// Check amount
	if req.Amount > 0 {
		checkPayload["checkAmount"] = map[string]interface{}{
			"type":   "eq", // Equal to expected amount
			"amount": fmt.Sprintf("%.2f", req.Amount),
		}
	}

	// Check date (optional)
	if req.PaymentDate != "" {
		checkPayload["checkDate"] = map[string]interface{}{
			"type": "eq",
			"date": req.PaymentDate,
		}
	}

	checkPayloadJSON, err := json.Marshal(checkPayload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal check payload: %w", err)
	}

	writer.WriteField("payload", string(checkPayloadJSON))
	writer.Close()

	// 🌐 Send HTTP Request to Slip2Go
	httpReq, err := http.NewRequest("POST", c.BaseURL, &buf)
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	httpReq.Header.Set("Content-Type", writer.FormDataContentType())
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	// Execute request
	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Slip2Go API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Parse JSON response
	var slipResp Slip2GoResponse
	if err := json.Unmarshal(body, &slipResp); err != nil {
		return nil, fmt.Errorf("failed to parse Slip2Go response: %w (body: %s)", err, string(body))
	}

	return &slipResp, nil
}

// mockVerifySlip returns a mock response for development/testing
func (c *Slip2GoClient) mockVerifySlip(req VerifySlipRequest) *Slip2GoResponse {
	// 🎭 Mock successful verification
	return &Slip2GoResponse{
		Status:  "success",
		Message: "Verification successful (MOCK MODE)",
		Result: struct {
			Verified  bool   `json:"verified"`
			SlipID    string `json:"slip_id"`
			Amount    string `json:"amount"`
			Sender    string `json:"sender"`
			Receiver  string `json:"receiver"`
			Timestamp string `json:"timestamp"`
		}{
			Verified:  true,
			SlipID:    "MOCK_" + time.Now().Format("20060102150405"),
			Amount:    fmt.Sprintf("%.2f", req.Amount),
			Sender:    "Mock Sender",
			Receiver:  req.AccountName,
			Timestamp: time.Now().Format(time.RFC3339),
		},
	}
}
