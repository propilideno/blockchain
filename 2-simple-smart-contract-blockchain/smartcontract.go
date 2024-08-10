package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
)

// SmartContract represents a smart contract in the blockchain
type SmartContract struct {
	ContractID    string `json:"contract_id"`
	Wallet        string `json:"wallet"`
	Type          string `json:"type"`
	Specification string `json:"specification"`
	Code          Code   `json:"-"`
}

// Code interface defines the methods for a smart contract
type Code interface {
	Execute(blockchain *Blockchain) error
	Validate(blockchain *Blockchain) bool
}

type ContractExecution struct {
	ContractID  string  `json:"contract_id"`
	ConsumedGas float64 `json:"consumed_gas"`
	Result      string  `json:"result"`
	Timestamp   time.Time `json:"timestamp"`
	Miner       string  `json:"miner"`
}

// CertificateRequest implements the Code interface for certificate requests
type CertificateRequest struct {
	Certificate string `json:"certificate"`
}

func (cr *CertificateRequest) Execute(blockchain *Blockchain) error {
	// Add logic to process the certificate request
	fmt.Println("Executing certificate request...")
	return nil
}

func (cr *CertificateRequest) Validate(blockchain *Blockchain) bool {
	// Add validation logic for the certificate request
	fmt.Println("Validating certificate request...")
	return true
}

// Execute calls the Execute method of the Code interface
func (sc *SmartContract) Execute(blockchain *Blockchain) error {
	return sc.Code.Execute(blockchain)
}

// Validate calls the Validate method of the Code interface
func (sc *SmartContract) Validate(blockchain *Blockchain) bool {
	return sc.Code.Validate(blockchain)
}

// calculateDigest generates a SHA256 digest of the contract data
func (sc *SmartContract) calculateDigest() string {
	data, _ := json.Marshal(sc)
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}
