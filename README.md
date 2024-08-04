# Blockchain
The best way to learn blockchain is by building one. This repository guides you from creating a basic blockchain to understanding advanced concepts, providing hands-on experience along the way. Dive in and start your blockchain journey today!

## [Simple Blockchain](./0-simple-blockchain/README.md)
```
docker run -p 7000:7000 propilideno/simple-blockchain
```

```mermaid
classDiagram
    direction LR

    class Blockchain {
        +Block[] Chain
        +int Difficulty
        +Transaction[] MemoryPool
        +Block mine()
        +void addTransaction(transaction Transaction)
        +bool isValid()
    }

    class Block {
        +Transaction[] Transactions
        +string PreviousHash
        +string Hash
        +time.Time Timestamp
        +int Nonce
        +string calculateHash()
        +void mine(difficulty int)
    }

    Blockchain "1" --> "*" Block : contains

    Block <|-- Block1 : PreviousHash
    Block1 <|-- Block2 : PreviousHash
    Block2 <|-- Block3 : PreviousHash
    Block3 <|-- Block4 : PreviousHash

```
#### Routes
- GET /chain
- GET /memorypool
- GET /mine
- POST /transactions/new
    - body: `{ "from": "Lucas", "to": "Filipe", "amount": 10 }`
#### Lacks of
- Transaction validation
- Persistence
- Miner Reward
- Descentralization
    - P2P Network
    - Node discovery

## [Simple Transactional Blockchain](./1-simple-transactional-blockchain/README.md)
```
docker run -p 7000:7000 propilideno/simple-transactional-blockchain
```
```mermaid
classDiagram
    direction LR

    class Blockchain {
        +Block[] Chain
        +int Difficulty
        +float64 RewardPerBlock
        +float64 MaxCoins
        +float64 getMinedCoins()
        +float64 getBalance(address string)
        +bool isValid()
        +Block mine(miner string) Block
        +void addBlockData(data BlockData)
    }

    class Block {
        +BlockData[] Data
        +BlockReward Reward
        +string PreviousHash
        +string Hash
        +time.Time Timestamp
        +int Nonce
        +string calculateHash()
        +void mine(difficulty int)
    }

    Blockchain "1" --> "*" Block : contains

    Block <|-- Block1 : PreviousHash
    Block1 <|-- Block2 : PreviousHash
    Block2 <|-- Block3 : PreviousHash
    Block3 <|-- Block4 : PreviousHash
```
#### Routes
- GET /chain
- GET /memorypool
- GET /mine?wallet=**base64_encoded_public_key**
- POST /data/new
    - body: `{ "from": "Lucas", "to": "Filipe", "amount": 10 }`

#### Lacks of
- Persistence
- Descentralization
    - P2P Network
    - Node discovery
