# Simple Blockchain implementation

## Routes
- GET /chain
- GET /memorypool
- GET /mine
- POST /transactions/new
    > body: `{ "from": "Lucas", "to": "Filipe", "amount": 10 }`

## Lacks of
- Transaction validation
- Persistence
- Miner Reward
- Descentralization
    - P2P Network
    - Node discovery
