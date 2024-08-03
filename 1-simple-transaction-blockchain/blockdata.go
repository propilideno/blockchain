package main

// BlockData is an interface for data that can be stored in a block
type BlockData interface {
	Validate() bool
}

// Transaction represents a blockchain transaction
type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// Validate checks if the transaction is valid
func (t Transaction) Validate() bool {
	return t.From != "" && t.To != "" && t.Amount > 0
}
