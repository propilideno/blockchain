package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SmartContract represents a smart contract in the blockchain
type SmartContract struct {
	ContractID string `json:"contract_id"`
	Wallet     string `json:"wallet"`
	Code       string `json:"code"`
	Status     string `json:"status"`
}

// Execute executes the action if the contract is completed
func (sc *SmartContract) Execute(blockchain *Blockchain) {
	if sc.Status == "completed" {
		fmt.Printf("Contract %s executed\n", sc.ContractID)
		// Add custom execution logic here
	}
}

// Validate checks if the smart contract is valid
func (sc *SmartContract) Validate(blockchain *Blockchain) bool {
	return true
}

func (sc *SmartContract) periodicCheck(blockchain *Blockchain) {
	for {
		if sc.Status == "completed" {
			sc.Execute(blockchain)
			break
		}

		time.Sleep(10 * time.Second) // Adjust the interval as needed

		resp, err := http.Get(fmt.Sprintf("http://localhost:3000/acme?wallet=%s", sc.Wallet))
		if err != nil {
			fmt.Println("Error checking contract condition:", err)
			continue
		}
		defer resp.Body.Close()

		var result string
		if _, err := fmt.Fscan(resp.Body, &result); err == nil && result == sc.calculateDigest() {
			sc.Status = "completed"
			fmt.Printf("Contract %s condition met\n", sc.ContractID)
		}
	}
}

func (sc *SmartContract) calculateDigest() string {
	data, _ := json.Marshal(sc)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
