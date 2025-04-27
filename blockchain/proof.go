package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"log"
	"math"
	"math/big"
)

// Proof-of-Work Algorithm:
// - Take the block's data
// - Create a counter (nonce) starting at 0
// - Hash the data combined with the nonce
// - Check if the hash meets the difficulty target (leading zeros)
// - Repeat until a valid hash is found

// Difficulty defines how hard the proof-of-work puzzle is (number of leading zeros).
const Difficulty = 18

// ProofOfWork holds the block and the target needed for mining.
type ProofOfWork struct {
	Block  *Block   // Block to be mined
	Target *big.Int // Target threshold (hash must be lower than this)
}

// NewProof creates a new ProofOfWork instance for a given block.
func NewProof(b *Block) *ProofOfWork {
	target := big.NewInt(1)
	// Left shift to set the target based on the difficulty
	target.Lsh(target, uint(256-Difficulty))

	pow := &ProofOfWork{b, target}

	return pow
}

// InitData prepares the block data combined with the nonce to hash.
func (pow *ProofOfWork) InitData(nonce int) []byte {
	data := bytes.Join(
		[][]byte{
			pow.Block.PrevHash,
			pow.Block.Data,
			ToHex(int64(nonce)),
			ToHex(int64(Difficulty)),
		},
		[]byte{},
	)

	return data
}

// Run executes the proof-of-work algorithm to find a valid nonce and hash.
func (pow *ProofOfWork) Run() (int, []byte) {
	var intHash big.Int
	var hash [32]byte

	nonce := 0

	// Try different nonce values until a valid hash is found
	for nonce < math.MaxInt64 {
		data := pow.InitData(nonce)
		hash = sha256.Sum256(data)

		fmt.Printf("\r%x", hash) // Print current hash (debug)
		intHash.SetBytes(hash[:])

		// Check if the hash is below the target
		if intHash.Cmp(pow.Target) == -1 {
			break // Success
		} else {
			nonce++ // Try next nonce
		}
	}
	fmt.Println()

	return nonce, hash[:]
}

// Validate verifies that the block's proof-of-work is correct.
func (pow *ProofOfWork) Validate() bool {
	var intHash big.Int

	data := pow.InitData(pow.Block.Nonce)
	hash := sha256.Sum256(data)
	intHash.SetBytes(hash[:])

	// Return true if hash is below target
	return intHash.Cmp(pow.Target) == -1
}

// ToHex converts an int64 to a byte slice.
func ToHex(num int64) []byte {
	buff := new(bytes.Buffer)
	err := binary.Write(buff, binary.BigEndian, num)
	if err != nil {
		log.Panic(err)
	}

	return buff.Bytes()
}
