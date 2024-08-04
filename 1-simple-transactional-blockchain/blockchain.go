package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// Block represents each 'item' in the blockchain
type Block struct {
	Data         []BlockData
	Reward       BlockReward
	PreviousHash string
	Hash         string
	Timestamp    time.Time
	Nonce        int
}

// Blockchain represents the entire chain
type Blockchain struct {
	GenesisBlock   Block
	Chain          []Block
	MemoryPool     []BlockData
	Difficulty     int
	RewardPerBlock float64
	MaxCoins       float64
}

// calculateHash calculates the hash of a block
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.Data)
	reward, _ := json.Marshal(b.Reward)
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

// addBlockData adds new data to the memory pool
func (b *Blockchain) addBlockData(data BlockData) {
	if data.Validate(b) {
		b.MemoryPool = append(b.MemoryPool, data)
	}
}

// mine mines a new block containing data from the memory pool
func (b *Blockchain) mine(miner string) (Block, error) {
	if len(b.Chain) == 1 && b.Chain[0].Hash == b.Chain[0].calculateHash() {
		// Genesis block should be mined first
		b.Chain[0].mine(b.Difficulty)
		b.Chain[0].Hash = b.Chain[0].calculateHash()
		b.Chain[0].Reward = BlockReward{Miner: miner, Amount: b.RewardPerBlock}
		return b.Chain[0], nil
	}

	lastBlock := b.Chain[len(b.Chain)-1]
	reward := BlockReward{
		Miner:  miner,
		Amount: b.RewardPerBlock,
	}
	newBlock := Block{
		Data:         b.MemoryPool,
		Reward:       reward,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}

	// Check if adding the reward exceeds the maximum coins limit
	if b.getMinedCoins()+b.RewardPerBlock > b.MaxCoins {
		return Block{}, fmt.Errorf("max coins limit reached")
	}

	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)

	// Clear the memory pool
	b.MemoryPool = nil

	return newBlock, nil
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
func (b Blockchain) getBalance(address string) float64 {
	balance := 0.0
	for _, block := range b.Chain {
		for _, data := range block.Data {
			if tx, ok := data.(Transaction); ok {
				if tx.From == address {
					balance -= tx.Amount
				}
				if tx.To == address {
					balance += tx.Amount
				}
			}
		}
		if block.Reward.Miner == address {
			balance += block.Reward.Amount
		}
	}
	return balance
}

// getMinedCoins calculates the total mined coins
func (b Blockchain) getMinedCoins() float64 {
	totalMined := 0.0
	for _, block := range b.Chain {
		totalMined += block.Reward.Amount
	}
	return totalMined
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
	app.Get("/mine", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		miner := c.Query("wallet")
		if miner == "" {
			return c.Status(fiber.StatusBadRequest).SendString("Missing miner wallet")
		}
		block, err := blockchain.mine(miner)
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

	// Add new block data (transaction)
	app.Post("/data/new", func(c *fiber.Ctx) error {
		var data Transaction
		if err := c.BodyParser(&data); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		blockchain := c.Locals("blockchain").(*Blockchain)
		blockchain.addBlockData(data)

		response := fiber.Map{"message": "Data added to the memory pool"}
		return c.Status(fiber.StatusCreated).JSON(response)
	})

	// Get the full blockchain
	app.Get("/chain", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		response := fiber.Map{
			"chain":     blockchain.Chain,
			"length":    len(blockchain.Chain),
			"isValid":   blockchain.isValid(),
			"minedCoins": blockchain.getMinedCoins(),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	// Get data from the memory pool
	app.Get("/memorypool", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		response := fiber.Map{
			"memorypool": blockchain.MemoryPool,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Listen(":7000")
}
