package main

import (
	"fmt"
	"strconv"

	"github.com/kikoleitao/go-blockchain/blockchain"
)

func main() {
	// Initialize a new blockchain with the Genesis Block.
	chain := blockchain.InitBlockChain()

	// Add new blocks to the chain.
	chain.AddBlock("First Block after Genesis")
	chain.AddBlock("Second Block after Genesis")
	chain.AddBlock("Third Block after Genesis")

	// Iterate over each block and print its details.
	for _, block := range chain.Blocks {
		fmt.Printf("Data in Block: %s\n", block.Data)
		fmt.Printf("Previous Hash: %x\n", block.PrevHash)
		fmt.Printf("Hash: %x\n", block.Hash)

		// Validate the block's Proof of Work
		pow := blockchain.NewProof(block)
		fmt.Printf("PoW Valid: %s\n", strconv.FormatBool(pow.Validate()))
		fmt.Println()
	}
}
