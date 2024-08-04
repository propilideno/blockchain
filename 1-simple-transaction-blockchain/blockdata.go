package main

// BlockData is an interface for data that can be stored in a block
type BlockData interface {
	Validate(blockchain *Blockchain) bool
}

// BlockReward represents the mining reward for a block
type BlockReward struct {
	Miner  string  `json:"miner"`
	Amount float64 `json:"amount"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// Validate checks if the transaction is valid
func (t Transaction) Validate(blockchain *Blockchain) bool {
	// Ensure that the transaction amount does not exceed the sender's balance
	if t.From != "0" { // "0" indicates a mining reward
		balance := blockchain.getBalance(t.From)
		return balance >= t.Amount
	}
	return true
}
