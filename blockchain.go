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

type Block struct {
	Data         map[string]interface{}
	Hash         string
	PreviousHash string
	Timestamp    time.Time
	PoW          int
}

type Blockchain struct {
	GenesisBlock Block
	Chain        []Block
	Difficulty   int
}

func (b Block) calculateHash() string {
	data, _ := json.Marshal(b.Data)
	blockData := b.PreviousHash + string(data) + b.Timestamp.String() + strconv.Itoa(b.PoW)
	blockHash := sha256.Sum256([]byte(blockData))
	return fmt.Sprintf("%x", blockHash)
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.Hash, strings.Repeat("0", difficulty)) {
		b.PoW++
		b.Hash = b.calculateHash()
	}
}

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

func (b *Blockchain) addBlock(data map[string]interface{}) {
	lastBlock := b.Chain[len(b.Chain)-1]
	newBlock := Block{
		Data:         data,
		PreviousHash: lastBlock.Hash,
		Timestamp:    time.Now(),
	}
	newBlock.mine(b.Difficulty)
	b.Chain = append(b.Chain, newBlock)
}

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

func main() {
	app := fiber.New()

	// Initialize the blockchain with a difficulty of 2
	blockchain := CreateBlockchain(2)

	// Middleware to set blockchain in context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("blockchain", &blockchain)
		return c.Next()
	})

	app.Get("/mine", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)

		// Mine a new block with a reward transaction
		transaction := map[string]interface{}{
			"from":   "0",
			"to":     "miner-address",
			"amount": 1,
		}
		blockchain.addBlock(transaction)

		response := fiber.Map{
			"message": "New Block Forged",
			"index":   len(blockchain.Chain) - 1,
		}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Post("/transactions/new", func(c *fiber.Ctx) error {
		var transaction map[string]interface{}
		if err := c.BodyParser(&transaction); err != nil {
			return c.Status(fiber.StatusBadRequest).SendString("Invalid input")
		}

		// Ensure transaction data is correct
		if transaction["from"] == nil || transaction["to"] == nil || transaction["amount"] == nil {
			return c.Status(fiber.StatusBadRequest).SendString("Missing transaction data")
		}

		blockchain := c.Locals("blockchain").(*Blockchain)
		blockchain.addBlock(transaction)

		response := fiber.Map{"message": "Transaction added to the blockchain"}
		return c.Status(fiber.StatusCreated).JSON(response)
	})

	app.Get("/chain", func(c *fiber.Ctx) error {
		blockchain := c.Locals("blockchain").(*Blockchain)

		response := fiber.Map{
			"chain":  blockchain.Chain,
			"length": len(blockchain.Chain),
		}
		return c.Status(fiber.StatusOK).JSON(response)
	})

	app.Listen(":3000")
}
