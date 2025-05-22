package blockchain

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"encoding/hex"
	"fmt"
	"log"
)

// Transaction groups inputs and outputs, tracked by an ID.
type Transaction struct {
	ID      []byte
	Inputs  []TxInput
	Outputs []TxOutput
}

// TxOutput represents coins that can be claimed using a specific key.
type TxOutput struct {
	Value  int
	PubKey string
}

// TxInput refers to a previous output being used in a new transaction.
type TxInput struct {
	ID  []byte
	Out int
	Sig string
}

// SetID hashes the entire transaction structure to set a unique ID.
func (tx *Transaction) SetID() {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(tx)
	Handle(err)

	hash := sha256.Sum256(buffer.Bytes())
	tx.ID = hash[:]
}

// CoinbaseTx creates the first transaction in a block, awarding coins.
func CoinbaseTx(recipient, customData string) *Transaction {
	if customData == "" {
		customData = fmt.Sprintf("Reward to %s", recipient)
	}

	input := TxInput{ID: []byte{}, Out: -1, Sig: customData}
	output := TxOutput{Value: 100, PubKey: recipient}

	coinbase := Transaction{
		ID:      nil,
		Inputs:  []TxInput{input},
		Outputs: []TxOutput{output},
	}
	coinbase.SetID()

	return &coinbase
}

// NewTransaction creates a standard transaction with inputs and outputs.
func NewTransaction(sender, recipient string, amount int, bc *BlockChain) *Transaction {
	var inputs []TxInput
	var outputs []TxOutput

	total, validOuts := bc.FindSpendableOutputs(sender, amount)
	if total < amount {
		log.Panic("Insufficient funds")
	}

	for txIDStr, outs := range validOuts {
		txID, err := hex.DecodeString(txIDStr)
		Handle(err)

		for _, outIdx := range outs {
			inputs = append(inputs, TxInput{ID: txID, Out: outIdx, Sig: sender})
		}
	}

	outputs = append(outputs, TxOutput{Value: amount, PubKey: recipient})
	if total > amount {
		outputs = append(outputs, TxOutput{Value: total - amount, PubKey: sender})
	}

	tx := Transaction{nil, inputs, outputs}
	tx.SetID()

	return &tx
}

// IsCoinbase returns true if the transaction is a coinbase (mining reward).
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].ID) == 0 && tx.Inputs[0].Out == -1
}

// CanUnlock determines if the input signature matches the unlocking key.
func (in *TxInput) CanUnlock(key string) bool {
	return in.Sig == key
}

// CanBeUnlocked checks if the provided key can unlock the output.
func (out *TxOutput) CanBeUnlocked(key string) bool {
	return out.PubKey == key
}
