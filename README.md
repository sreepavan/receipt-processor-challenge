# Receipt Processor API

## Overview
This is a RESTful API that processes retail receipts and calculates reward points based on predefined rules.

## Features
- Accepts receipt data and assigns a unique ID.
- Computes reward points based on receipt details.
- Provides an endpoint to retrieve points for a given receipt.
- Supports local execution and containerized deployment using Docker.

---

## API Endpoints

### 1️ Health Check
**Endpoint:** `GET /health`  
**Description:** Verifies that the server is running.  
**Response:**
```
Server is running!
```

---

### 2️ Process Receipt
**Endpoint:** `POST /receipts/process`  
**Description:** Accepts a JSON receipt, assigns a unique ID, and stores it.

**Request Example:**
```json
{
    "retailer": "Walmart",
    "purchaseDate": "2023-03-15",
    "purchaseTime": "14:30",
    "total": "35.75",
    "items": [
        {"shortDescription": "Milk", "price": "3.50"},
        {"shortDescription": "Bread", "price": "2.75"}
    ]
}
```

**Response Example:**
```json
{
    "id": "a1b2c3d4-e5f6-7890-1234-56789abcdef0"
}
```

---

### 3️ Get Receipt Points
**Endpoint:** `GET /receipts/{id}/points`  
**Description:** Retrieves the reward points assigned to a given receipt ID.

**Response Example:**
```json
{
    "points": 55
}
```

---

## Setup Instructions
### **Option 1: Run Locally (Without Docker)**
```bash
# Clone the repository
git clone https://github.com/sreepavan/receipt-processor-challenge.git
cd receipt-processor-challenge

# Run the application
go run main.go
```

### **Option 2: Run with Docker**
```bash
# Build the Docker image
docker build -t receipt-processor .

# Run the container
docker run -p 8080:8080 receipt-processor
```

### **Option 3: Run with Docker Compose**
```bash
# Start the service
docker-compose up --build
```

---

## Testing the API
### **1️  Health Check**
```bash
curl http://localhost:8080/health
```

### **2️  Process a Receipt**
```bash
curl -X POST "http://localhost:8080/receipts/process" \
     -H "Content-Type: application/json" \
     -d '{ "retailer": "Walmart", "purchaseDate": "2023-03-15", "purchaseTime": "14:30", "total": "35.75", "items": [{"shortDescription": "Milk", "price": "3.50"}, {"shortDescription": "Bread", "price": "2.75"}] }'
```

### **3️  Retrieve Points**
```bash
curl -X GET "http://localhost:8080/receipts/{id}/points"
```
Replace `{id}` with the receipt ID from the process receipt response.

---

## Project Structure
```
receipt-processor-challenge/
│── main.go              # Main server file
│── main_test.go         # Unit tests
│── Dockerfile          # Docker setup
│── docker-compose.yml  # Docker Compose setup
│── go.mod              # Go module dependencies
│── README.md           # Project documentation
```

---

## License
This project is for assessment purposes only. All rights reserved.

---

## Author
Developed by **Sreepavan Kumar Appikonda**

