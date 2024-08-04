# Blockchain

Our goal is start simple and making the blockchain much complex

## [Simple Blockchain](./0-simple-blockchain/README.md)
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
