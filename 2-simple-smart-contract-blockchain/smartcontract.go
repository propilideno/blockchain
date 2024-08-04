package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// SmartContract represents a smart contract in the blockchain
type SmartContract struct {
	ContractID    string `json:"contract_id"`
	CreatorWallet string `json:"creator_wallet"`
	Condition     string `json:"condition"`
	Action        string `json:"action"`
	Status        string `json:"status"`
	Result        string `json:"result"`
}

// GenerateToken generates a random token for the contract
func GenerateToken() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
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

		resp, err := http.Get(fmt.Sprintf("http://localhost:3000/acme?wallet=%s", sc.CreatorWallet))
		if err != nil {
			fmt.Println("Error checking contract condition:", err)
			continue
		}
		defer resp.Body.Close()

		var result string
		if _, err := fmt.Fscan(resp.Body, &result); err == nil && result == sc.calculateDigest() {
			sc.Status = "completed"
			sc.Result = result
			fmt.Printf("Contract %s condition met\n", sc.ContractID)
		}
	}
}

func (sc *SmartContract) calculateDigest() string {
	data, _ := json.Marshal(sc)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
