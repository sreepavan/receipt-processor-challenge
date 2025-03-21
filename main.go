package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Receipt -> the structure of a receipt.
type Receipt struct {
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Total        string `json:"total"`
	Items        []Item `json:"items"`
}

// Item -> each item in a receipt.
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// This holds receipts.
type InMemoryStore struct {
	sync.RWMutex
	receipts map[string]Receipt
}

var receiptStore = InMemoryStore{
	receipts: make(map[string]Receipt),
}

// verify the server is running.
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Server is running!")
}

// processReceiptHandler handles POST /receipts/process endpoint.
func processReceiptHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}

	var receiptData Receipt
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&receiptData)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	uniqueReceiptID := uuid.New().String()

	receiptStore.Lock()
	receiptStore.receipts[uniqueReceiptID] = receiptData
	receiptStore.Unlock()

	responsePayload := map[string]string{
		"id": uniqueReceiptID,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

// getPointsHandler -> GET /receipts/{id}/points endpoint.
func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET allowed", http.StatusMethodNotAllowed)
		return
	}

	pathParts := strings.Split(r.URL.Path, "/")
	if len(pathParts) != 4 || pathParts[1] != "receipts" || pathParts[3] != "points" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}
	receiptID := pathParts[2]

	receiptStore.RLock()
	receipt, exists := receiptStore.receipts[receiptID]
	receiptStore.RUnlock()

	if !exists {
		http.Error(w, "Receipt not found", http.StatusNotFound)
		return
	}

	totalPoints := computeReceiptPoints(receipt)

	responsePayload := map[string]int{
		"points": totalPoints,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responsePayload)
}

// computeReceiptPoints -> calculates points for a given receipt.
func computeReceiptPoints(receipt Receipt) int {
	points := 0

	// 1. alphanumeric character.
	alphaNumericRegex := regexp.MustCompile(`[a-zA-Z0-9]`)
	points += len(alphaNumericRegex.FindAllString(receipt.Retailer, -1))

	totalAmount, err := strconv.ParseFloat(receipt.Total, 64)
	if err != nil {
		log.Printf("Error parsing total amount: %v", err)
	}

	// 2. total is a round dollar amount with no cents.
	if totalAmount == math.Trunc(totalAmount) {
		points += 50
	}

	// 3. total is a multiple of 0.25.
	if math.Mod(totalAmount, 0.25) == 0 {
		points += 25
	}

	// 4. 5 points for every two items on the receipt.
	itemCountPoints := (len(receipt.Items) / 2) * 5
	points += itemCountPoints

	// 5. Multiply the price by 0.2 and round up.
	for _, item := range receipt.Items {
		trimmedDescription := strings.TrimSpace(item.ShortDescription)
		if len(trimmedDescription) > 0 && len(trimmedDescription)%3 == 0 {
			itemPrice, err := strconv.ParseFloat(item.Price, 64)
			if err != nil {
				log.Printf("Error parsing item price: %v", err)
				continue
			}
			calculatedPoints := int(math.Ceil(itemPrice * 0.2))
			points += calculatedPoints
		}
	}

	// 6. 5 points if total > 10.00.
	if totalAmount > 10.00 {
		points += 5
	}

	// 7. if the day in the purchase date is odd.
	dateParts := strings.Split(receipt.PurchaseDate, "-")
	if len(dateParts) == 3 {
		day, err := strconv.Atoi(dateParts[2])
		if err == nil && day%2 == 1 {
			points += 6
		}
	}

	// 8. if the time of purchase is after 2:00pm and before 4:00pm.
	parsedTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil {
		twoPM, _ := time.Parse("15:04", "14:00")
		fourPM, _ := time.Parse("15:04", "16:00")
		if parsedTime.After(twoPM) && parsedTime.Before(fourPM) {
			points += 10
		}
	}

	return points
}

func main() {
	// Setup HTTP routes.
	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/receipts/process", processReceiptHandler)
	http.HandleFunc("/receipts/", getPointsHandler)

	// Start the HTTP server.
	serverAddress := ":8080"
	fmt.Printf("Starting server on %s...\n", serverAddress)
	if err := http.ListenAndServe(serverAddress, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
