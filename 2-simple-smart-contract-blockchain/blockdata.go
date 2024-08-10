package main

// BlockData contains all types of data that can be part of a block
type BlockData struct {
	ContractExecutionHistory []ContractExecution `json:"contract_execution_history"`
	Contracts                []SmartContract            `json:"contracts"`
	Transactions             []Transaction              `json:"transactions"`
}

// BlockReward represents the reward given to the miner
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
	// Ensure the transaction is not a self-transfer
	if t.From == t.To {
		return false
	}

	// Ensure that the transaction amount is positive
	if t.Amount <= 0 {
		return false
	}

	// Ensure that the transaction amount does not exceed the sender's balance
	if t.From != "0" { // "0" indicates a mining reward
		balance := blockchain.getBalance(t.From)
		return balance >= t.Amount
	}
	return true
}
