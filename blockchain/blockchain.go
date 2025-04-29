package blockchain

import (
	"fmt"
	"github.com/dgraph-io/badger/v3"
)

const (
	dbPath = "./tmp/blocks" // Path where the BadgerDB will store blockchain data
)

// BlockChain represents the chain as a whole, including a reference to the last block.
type BlockChain struct {
	LastHash []byte       // Hash of the latest block in the chain
	Database *badger.DB   // Handle to the BadgerDB instance
}

// BlockChainIterator helps iterate through the blockchain blocks from newest to oldest.
type BlockChainIterator struct {
	CurrentHash []byte     // The hash of the block currently being accessed
	Database    *badger.DB // Reference to the database to retrieve block data
}

// InitBlockChain initializes the blockchain, creating the genesis block if needed.
func InitBlockChain() *BlockChain {
	var lastHash []byte

	opts := badger.DefaultOptions(dbPath)
	opts.Logger = nil // Disables logging completely 
	db, err := badger.Open(opts)
	Handle(err)

	err = db.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte("lh")) // "lh" = last hash
		if err == badger.ErrKeyNotFound {
			// If not found, create genesis block
			fmt.Println("No existing blockchain found")
			genesis := Genesis()
			fmt.Println("Genesis created")

			err = txn.Set(genesis.Hash, genesis.Serialize())
			Handle(err)
			err = txn.Set([]byte("lh"), genesis.Hash)
			lastHash = genesis.Hash

			return err
		} else {
			// Load existing last hash
			item, err := txn.Get([]byte("lh"))
			Handle(err)
			err = item.Value(func(val []byte) error {
				lastHash = append([]byte{}, val...)
				return nil
			})
			return err
		}
	})
	Handle(err)

	return &BlockChain{LastHash: lastHash, Database: db}
}

// AddBlock creates and stores a new block with the given data in the database.
func (chain *BlockChain) AddBlock(data string) {
	var lastHash []byte

	// Get the last hash to reference the previous block
	err := chain.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("lh"))
		Handle(err)
		err = item.Value(func(val []byte) error {
			lastHash = append([]byte{}, val...)
			return nil
		})
		return err
	})
	Handle(err)

	newBlock := CreateBlock(data, lastHash)

	// Store the new block and update the last hash
	err = chain.Database.Update(func(txn *badger.Txn) error {
		err := txn.Set(newBlock.Hash, newBlock.Serialize())
		Handle(err)
		err = txn.Set([]byte("lh"), newBlock.Hash)
		chain.LastHash = newBlock.Hash

		return err
	})
	Handle(err)
}

// Iterator creates a new blockchain iterator starting from the latest block.
func (chain *BlockChain) Iterator() *BlockChainIterator {
	return &BlockChainIterator{CurrentHash: chain.LastHash, Database: chain.Database}
}

// Next returns the next block in the chain (going backward in time).
func (iter *BlockChainIterator) Next() *Block {
	var block *Block

	err := iter.Database.View(func(txn *badger.Txn) error {
		item, err := txn.Get(iter.CurrentHash)
		Handle(err)
		err = item.Value(func(val []byte) error {
			block = Deserialize(val)
			return nil
		})
		return err
	})
	Handle(err)

	iter.CurrentHash = block.PrevHash
	return block
}
