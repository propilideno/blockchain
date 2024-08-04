# Simple Blockchain implementation

```bash
# Generate key pairs for testing in transactions and mining
genkeypair(){
    WALLETS="./keys/wallets.txt"
    KEY_DIR="./keys/$1"
    PVT_KEY="$KEY_DIR/$1.key"
    PUB_KEY="$KEY_DIR/$1.pub"
    mkdir -p $KEY_DIR && openssl genpkey -algorithm RSA -out $PVT_KEY -pkeyopt rsa_keygen_bits:1024 && openssl rsa -pubout -in $PVT_KEY -out $PUB_KEY && echo "$1: $(base64 -w0 $PUB_KEY)" >> $WALLETS
}
```

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
