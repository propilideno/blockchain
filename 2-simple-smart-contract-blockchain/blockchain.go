package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Data         BlockData   `json:"data"`
	PreviousHash string      `json:"previous_hash"`
	Hash         string      `json:"hash"`
	Timestamp    time.Time   `json:"timestamp"`
	Nonce        int         `json:"nonce"`
	BlockReward  BlockReward `json:"block_reward"`
}

// Blockchain represents the entire chain
type Blockchain struct {
	GenesisBlock          Block
	Chain                 []Block
	TransactionPool       []Transaction
	ContractExecutionPool []ContractExecution
	Difficulty            int
	RewardPerBlock        float64
	MaxCoins              float64
}

// NewBlock creates a new block with the given parameters and calculates its hash
func NewBlock(previousHash string) Block {
	block := Block{
		PreviousHash: previousHash,
		Timestamp:    time.Now(),
	}
	block.Hash = block.calculateHash()
	return block
}

// calculateHash calculates the hash of a block
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.Data)
	reward, _ := json.Marshal(b.BlockReward)
	blockData := b.PreviousHash + string(data) + string(reward) + b.Timestamp.String() + strconv.Itoa(b.Nonce)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

// mine mines a block
func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.Nonce++
		b.Hash = b.calculateHash()
	}
}

// CreateBlockchain creates a new blockchain with a genesis block
func CreateBlockchain(difficulty int, rewardPerBlock float64, maxCoins float64) Blockchain {
	genesisBlock := Block{
		Timestamp: time.Now(),
	}
	genesisBlock.Hash = genesisBlock.calculateHash() // Set initial hash without mining
	return Blockchain{
		GenesisBlock:   genesisBlock,
		Chain:          []Block{genesisBlock},
		Difficulty:     difficulty,
		RewardPerBlock: rewardPerBlock,
		MaxCoins:       maxCoins,
	}
}

func (b *Blockchain) findContractByID(contractID string) *SmartContract {
	fmt.Printf("Blockchain state: %+v\n", b.Chain)
	fmt.Printf("Looking for contract ID: %s\n", contractID)

	for blockIndex, block := range b.Chain {
		fmt.Printf("Inspecting block %d with hash: %s\n", blockIndex, block.Hash)
		for i := range block.Data.Contracts {
			fmt.Printf("Checking contract ID: %s in block %d\n", block.Data.Contracts[i].ContractID, blockIndex)
			if block.Data.Contracts[i].ContractID == contractID {
				fmt.Println("Contract found!")
				return &block.Data.Contracts[i]
			}
		}
	}

	fmt.Println("Contract not found.")
	return nil
}

// addContract adds a smart contract directly to the current block
func (b *Blockchain) addContract(contract SmartContract) {
	lastBlock := &b.Chain[len(b.Chain)-1]
	lastBlock.Data.Contracts = append(lastBlock.Data.Contracts, contract)
}

// addTransaction adds a transaction to the transaction pool after validating it
func (b *Blockchain) addTransaction(tx Transaction) error {
	// Validate the transaction
	if !tx.Validate(b) {
		return fmt.Errorf("transaction validation failed: insufficient balance or invalid transaction")
	}

	// If valid, add the transaction to the pool
	b.TransactionPool = append(b.TransactionPool, tx)
	return nil
}

// mineTransaction mines transactions from the transaction pool into the current block
func (b *Blockchain) mineTransaction() error {
	if len(b.TransactionPool) == 0 {
		return fmt.Errorf("no transactions to mine")
	}

	// Get the current block (Last block in the chain)
	lastBlock := &b.Chain[len(b.Chain)-1]

	// Process the first transaction in the pool (FIFO)
	transaction := b.TransactionPool[0]
	lastBlock.Data.Transactions = append(lastBlock.Data.Transactions, transaction)

	// Remove the processed transaction from the pool
	b.TransactionPool = b.TransactionPool[1:]

	return nil
}

// mineContractExecution mines contract executions from the execution pool into the current block
func (b *Blockchain) mineContractExecution(miner string) float64 {
	lastBlock := &b.Chain[len(b.Chain)-1]

	if len(b.ContractExecutionPool) > 0 {
		// Process the first contract execution in the pool (FIFO)
		execpool := b.ContractExecutionPool[0]

		// Execute the contract
		contract := b.findContractByID(execpool.ContractID)
		if contract != nil {
			contract.Execute(b)
			execpool.Miner = miner
			lastBlock.Data.ContractExecutionHistory = append(lastBlock.Data.ContractExecutionHistory, execpool)
		}

		// Remove the processed contract execution from the pool
		b.ContractExecutionPool = b.ContractExecutionPool[1:]
		return execpool.ConsumedGas
	}
	return 0
}


func (b *Blockchain) mineBlock(miner string) (Block, error) {
	// Get the current block (Last block in the chain)
	currentBlock := &b.Chain[len(b.Chain)-1]

	// Determine the block reward based on the maximum coins limit
	if b.getMinedCoins()+b.RewardPerBlock > b.MaxCoins {
		currentBlock.BlockReward = BlockReward{Miner: miner, Amount: 0}
	} else {
		currentBlock.BlockReward = BlockReward{Miner: miner, Amount: b.RewardPerBlock}
	}

	// Mine the current block
	currentBlock.mine(b.Difficulty)

	// Create a new empty block and append it to the chain using the constructor
	b.Chain = append(b.Chain, NewBlock(currentBlock.Hash))

	// Return the mined block
	return *currentBlock, nil
}



// isValid checks if the blockchain is valid
func (b Blockchain) isValid() bool {
	for i := range b.Chain[1:] {
		previousBlock := b.Chain[i]
		currentBlock := b.Chain[i+1]
		if currentBlock.Hash != currentBlock.calculateHash() || currentBlock.PreviousHash != previousBlock.Hash {
			return false
		}
	}
	return true
}

// getBalance calculates the balance of a specific address
func (b *Blockchain) getBalance(address string) float64 {
	balance := 0.0
	for _, block := range b.Chain {
		for _, data := range block.Data.Transactions {
			if tx := data; tx.From == address {
				balance -= tx.Amount
			} else if tx.To == address {
				balance += tx.Amount
			}
		}
		if block.BlockReward.Miner == address {
			balance += block.BlockReward.Amount
		}
	}

	for _, history := range b.ContractExecutionPool {
		if history.ContractID == address {
			balance -= history.ConsumedGas
		}
	}

	return balance
}

// getMinedCoins calculates the total mined coins
func (b Blockchain) getMinedCoins() float64 {
	totalMined := 0.0
	for _, block := range b.Chain {
		totalMined += block.BlockReward.Amount
	}
	return totalMined
}

// generateRandomID generates a random 16-byte hex string
func generateRandomID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// extract the public key from a certificate
func extractPublicKeyFromCertificate(certificate string) string {
	return "TODO"
}

// main sets up the server and routes
func main() {
	app := fiber.New()

	// Initialize the blockchain with a difficulty of 2, reward of 10 coins per block, and a maximum of 1000 coins
	blockchain := CreateBlockchain(2, 10, 1000)

	// Middleware to set blockchain in context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("blockchain", &blockchain)
		return c.Next()
	})

	// Mine a new block
	app.Get("/mine/block", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		miner := c.Query("wallet")
		if miner == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing miner wallet")
		}
		block, err := blockchain.mineBlock(miner)
		if err != nil {
			return c.Status(fiber.StatusForbidden).SendString(err.Error())
		}

		response := fiber.Map{
			"message": "New Block Forged",
			"index":   len(blockchain.Chain) - 1,
			"block":   block,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Get("/mine/transaction", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		miner := c.Query("wallet")
		if miner == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing miner wallet")
		}

		// Mine the transaction
		err := blockchain.mineTransaction()
		if err != nil {
			return c.Status(fiber.StatusNoContent).SendString("No transactions to mine")
		}

		response := fiber.Map{
			"message": "Transaction mined successfully",
		}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// Mine contract executions
	app.Get("/mine/contract", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		miner := c.Query("wallet")
		if miner == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing miner wallet")
		}

		// Mine and process the contract executions
		gas := blockchain.mineContractExecution(miner)

		if gas != 0 {
			response := fiber.Map{
				"message": "Contract Executed Successfully",
				"gas":     gas,
			}
			return c.Status(fiber.StatusOK).JSON(response)
		} else {
			return c.Status(fiber.StatusNoContent).SendString("No contracts to mine")
		}
	})

	// Add new block data (transaction)
	app.Post("/transaction/new", func(c *fiber.Ctx) error {
		var tx Transaction
		if err := c.BodyParser(&tx); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		blockchain := c.Locals("blockchain").(*Blockchain)
		blockchain.addTransaction(tx)

		response := fiber.Map{"message": "Transaction added to the pool"}
		return c.Status(fiber.StatusCreated).JSON(response)
	})

	// Add new smart contract
	app.Post("/certificate/request", func(c *fiber.Ctx) error {
		var request struct {
			Certificate string `json:"certificate"`
		}

		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		contractID, err := generateRandomID()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString("Could not generate contract ID")
		}

		smartContract := SmartContract{
			ContractID:    contractID,
			Wallet:        extractPublicKeyFromCertificate(request.Certificate),
			Type:          "certificate",
			Specification: request.Certificate,
			Code:          &CertificateRequest{Certificate: request.Certificate},
		}

		blockchain := c.Locals("blockchain").(*Blockchain)
		blockchain.addContract(smartContract)

		response := fiber.Map{
			"message":    "Smart contract added to the current block",
			"contractID": contractID,
			"wallet":     smartContract.Wallet,
		}
		return c.Status(fiber.StatusCreated).JSON(response)
	})

	// Execute a smart contract (add to execution pool)
	app.Post("/contract/execute", func(c *fiber.Ctx) error {
		// Define a struct to parse the request body
		var request struct {
			ContractID string `json:"contractId"`
		}

		// Parse the request body
		if err := c.BodyParser(&request); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		// Validate the contract ID
		if request.ContractID == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing contract ID")
		}

		fmt.Printf("Received request to execute contract ID: %s\n", request.ContractID)

		blockchain := c.Locals("blockchain").(*Blockchain)

		// Find the contract in the blockchain
		contract := blockchain.findContractByID(request.ContractID)
		if contract == nil {
			return c.Status(fiber.StatusNotFound).SendString("Contract not found")
		}

		// Add the contract execution request to the ContractExecutionPool
		execution := ContractExecution{
			ContractID:  request.ContractID,
			ConsumedGas: 0.1, // Fixed gas fee
			Result:      "",  // Result will be set when mined
			Miner:       "",  // Miner will be set when mined
			Timestamp:   time.Now(),
		}

		blockchain.ContractExecutionPool = append(blockchain.ContractExecutionPool, execution)

		response := fiber.Map{
			"message": "Contract execution added to the pool",
		}
		return c.Status(fiber.StatusCreated).JSON(response)
	})

	// Get the full blockchain
	app.Get("/chain", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		response := fiber.Map{
			"chain":      blockchain.Chain,
			"length":     len(blockchain.Chain),
			"isValid":    blockchain.isValid(),
			"minedCoins": blockchain.getMinedCoins(),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// Get data from the transaction pool
	app.Get("/memorypool", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		response := fiber.Map{
			"transactionpool":       blockchain.TransactionPool,
			"contractexecutionpool": blockchain.ContractExecutionPool,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// Get information of a wallet
	app.Get("/info", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		wallet := c.Query("wallet")
		response := fiber.Map{
			"balance": blockchain.getBalance(wallet),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Listen(":7000")
}
