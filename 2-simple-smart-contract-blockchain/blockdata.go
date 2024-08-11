package main

// BlockData contains all types of data that can be part of a block
type BlockData struct {
	ContractExecutionHistory []ContractExecution `json:"contract_execution_history"`
	Contracts                []SmartContract     `json:"contracts"`
	Transactions             []Transaction       `json:"transactions"`
}

// Transaction represents a blockchain transaction
type Transaction struct {
	From   string  `json:"from"`
	To     string  `json:"to"`
	Amount float64 `json:"amount"`
}

// Validate checks if the transaction is valid
func (t Transaction) Validate(blockchain *Blockchain) bool {
	if t.From == t.To {
		return false
	} else if  t.From == BLOCK_REWARD_WALLET || t.To == BLOCK_REWARD_WALLET {
		return false
	} else if t.Amount <= 0 {
		return false
	} else {
		balance := blockchain.getBalance(t.From)
		return balance >= t.Amount
	}
}
