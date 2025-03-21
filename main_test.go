package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test -> (POST /receipts/process)
func TestProcessReceipt(t *testing.T) {
	requestBody := `{
		"retailer": "Target",
		"purchaseDate": "2023-03-10",
		"purchaseTime": "14:15",
		"total": "25.50",
		"items": [{"shortDescription": "Item A", "price": "5.00"}]
	}`

	req, err := http.NewRequest("POST", "/receipts/process", bytes.NewBuffer([]byte(requestBody)))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(processReceiptHandler)
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("Expected status 200, got %d", status)
	}
}

// Test -> points calculation
func TestComputePoints(t *testing.T) {
	receipt := Receipt{
		Retailer:     "Amazon",
		PurchaseDate: "2023-03-11",
		PurchaseTime: "15:30",
		Total:        "45.75",
		Items: []Item{
			{"Book", "12.00"},
			{"Pen", "3.25"},
		},
	}

	points := computeReceiptPoints(receipt)

	if points <= 0 {
		t.Errorf("Expected positive points, got %d", points)
	}
}
