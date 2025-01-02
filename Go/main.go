package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	shell "github.com/ipfs/go-ipfs-api"
)

func storeInIPFSHandler(w http.ResponseWriter, r *http.Request) {
    // Parse request body
    var data map[string]interface{}
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
        log.Println("Error decoding JSON:", err)
        return
    }

    // Initialize IPFS shell
    sh := shell.NewShell("localhost:5001")

    // Marshal the data to JSON bytes
    payloadBytes, err := json.Marshal(data)
    if err != nil {
        http.Error(w, "Error marshaling JSON data", http.StatusInternalServerError)
        log.Println("Error marshaling JSON:", err)
        return
    }

    // Store the JSON data in IPFS
    cid, err := sh.Add(bytes.NewReader(payloadBytes))
    if err != nil {
        http.Error(w, "Failed to store in IPFS", http.StatusInternalServerError)
        log.Println("IPFS storage error:", err)
        return
    }

    // Respond with the CID
    log.Println("Data stored in IPFS with CID:", cid)
    w.Write([]byte(cid))
}
func fetchFromIPFSHandler(w http.ResponseWriter, r *http.Request) {
	// Get CID from query parameters
	cid := r.URL.Query().Get("cid")
	if cid == "" {
		http.Error(w, "CID is required", http.StatusBadRequest)
		return
	}

	// Initialize IPFS shell
	sh := shell.NewShell("localhost:5001")

	// Fetch data from IPFS
	dataReader, err := sh.Cat(cid)
	if err != nil {
		log.Printf("Error fetching from IPFS: %v\n", err)
		http.Error(w, "Failed to fetch from IPFS", http.StatusInternalServerError)
		return
	}
	defer dataReader.Close()

	// Read the data
	data, err := ioutil.ReadAll(dataReader)
	if err != nil {
		log.Printf("Error reading data: %v\n", err)
		http.Error(w, "Failed to read data", http.StatusInternalServerError)
		return
	}

	// Log the retrieved data
	log.Printf("Data retrieved from IPFS: %s\n", string(data))

	// Return the data to the client
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func main() {
    http.HandleFunc("/storeInIPFS", storeInIPFSHandler)
    http.HandleFunc("/fetchFromIPFS", fetchFromIPFSHandler)

    fmt.Println("Go server is running on http://localhost:8545")
    if err := http.ListenAndServe(":8081", nil); err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
}
