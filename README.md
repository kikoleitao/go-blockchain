# Simple Blockchain in Go — Part 1: Basic Prototype

## Overview

This is the first part of a simple blockchain implementation in Go.  
In this stage, we have:

- A `Block` structure containing data, its hash, and the hash of the previous block.
- A `BlockChain` structure managing a list of blocks.
- A genesis block that starts the chain.
- Basic functionality to add new blocks linked through hashes.

Each block's hash is generated using SHA-256 over its `Data` and `PrevHash`.

## How It Works

- The blockchain is initialized with a **genesis block** (`"Genesis"` data).
- New blocks are added manually by calling `AddBlock(data)`, which:
  - References the previous block's hash.
  - Computes its own hash.
  - Appends the block to the chain.
- When run, the program prints the data, previous hash, and hash of each block in the chain.

## How to Run

To run this prototype locally:

```bash
# Initialize the Go module
go mod init github.com/kikoleitao/go-blockchain

# Run the main program
go run main.go
```

### Example Output
```bash
Previous Hash: 
Data in Block: Genesis
Hash: 81ddc8d248b2dccdd3fdd5e84f0cad62b08f2d10b57f9a831c13451e5c5c80a5
Previous Hash: 81ddc8d248b2dccdd3fdd5e84f0cad62b08f2d10b57f9a831c13451e5c5c80a5
Data in Block: First Block after Genesis
Hash: 50493b76a2b7bec8d33620d6310d5578b1dda079684405ed5e6bd55510146daf
Previous Hash: 50493b76a2b7bec8d33620d6310d5578b1dda079684405ed5e6bd55510146daf
Data in Block: Second Block after Genesis
Hash: 213e91a4ae1be45a651695ede0e75cba50818dce027dd4f0fe35742dc90158e1
Previous Hash: 213e91a4ae1be45a651695ede0e75cba50818dce027dd4f0fe35742dc90158e1
Data in Block: Third Block after Genesis
Hash: e22b76962d23ed3e327b9ababac19270b56c4d70d8878446609b13fa72ebc0e1

```
