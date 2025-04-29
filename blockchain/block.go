package blockchain

import (
	"bytes"
	"encoding/gob"
	"log"
)

// Block represents each 'item' in the blockchain.
type Block struct {
	Hash     []byte // Unique identifier for the block (result of PoW)
	Data     []byte // Actual data stored in the block (e.g., transactions)
	PrevHash []byte // Hash of the previous block in the chain
	Nonce    int    // Arbitrary number used to generate the valid hash
}

// CreateBlock generates a new Block using provided data and previous block hash.
// It also runs the Proof of Work to find a valid hash and nonce.
func CreateBlock(data string, prevHash []byte) *Block {
	block := &Block{[]byte{}, []byte(data), prevHash, 0}
	pow := NewProof(block)
	nonce, hash := pow.Run()

	block.Hash = hash[:]
	block.Nonce = nonce

	return block
}

// Genesis creates the first block in the blockchain, known as the "Genesis Block".
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

func (b *Block) Serialize() []byte {
	var res bytes.Buffer
	encoder := gob.NewEncoder(&res)

	err := encoder.Encode(b)

	Handle(err)

	return res.Bytes()
}

// Deserialize decodes a byte slice into a Block struct.
func Deserialize(data []byte) *Block {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(data))

	err := decoder.Decode(&block)

	Handle(err)

	return &block
}

// Handle panics if an error is encountered (for simplifying error handling).
func Handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}