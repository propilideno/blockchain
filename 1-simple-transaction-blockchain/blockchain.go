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
	PreviousHash string
	Hash         string
	Timestamp    time.Time
	Nonce        int
}

// Blockchain represents the entire chain
type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	MemoryPool   []BlockData
	Difficulty   int
}

// calculateHash calculates the hash of a block
func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.Data)
	blockData := b.PreviousHash + string(data) + b.Timestamp.String() + strconv.Itoa(b.Nonce)
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
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		Hash:      "0",
		Timestamp: time.Now(),
	}
	return Blockchain{
		GenesisBlock: genesisBlock,
		Chain:        []Block{genesisBlock},
		Difficulty:   difficulty,
	}
}

// addBlockData adds new data to the memory pool
func (b *Blockchain) addBlockData(data BlockData) {
	if data.Validate() {
		b.MemoryPool = append(b.MemoryPool, data)
	}
}

// mine mines a new block containing data from the memory pool
func (b *Blockchain) mine() Block {
	if len(b.MemoryPool) == 0 {
		return Block{}
	}

	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Data:         b.MemoryPool,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}
	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)

	// Clear the memory pool
	b.MemoryPool = nil

	return newBlock
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

// main sets up the server and routes
func main() {
	app := fiber.New()

	// Initialize the blockchain with a difficulty of 2
	blockchain := CreateBlockchain(2)

	// Middleware to set blockchain in context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("blockchain", &blockchain)
		return c.Next()
	})

	// Mine a new block
	app.Get("/mine", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)
		block := blockchain.mine()

		if block.Hash == "" {
			return c.Status(fiber.StatusForbidden).SendString("No data to mine")
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
			"chain":  blockchain.Chain,
			"length": len(blockchain.Chain),
			"isValid": blockchain.isValid(),
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