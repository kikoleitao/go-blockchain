package blockchain

// BlockChain represents a sequence of validated Blocks.
type BlockChain struct {
	Blocks []*Block
}

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

// AddBlock adds a new Block to the Blockchain using the provided data.
func (chain *BlockChain) AddBlock(data string) {
	prevBlock := chain.Blocks[len(chain.Blocks)-1]
	new := CreateBlock(data, prevBlock.Hash)
	chain.Blocks = append(chain.Blocks, new)
}

// Genesis creates the first block in the blockchain, known as the "Genesis Block".
func Genesis() *Block {
	return CreateBlock("Genesis", []byte{})
}

// InitBlockChain initializes a new Blockchain with the Genesis Block.
func InitBlockChain() *BlockChain {
	return &BlockChain{[]*Block{Genesis()}}
}
