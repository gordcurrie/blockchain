package main

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"strconv"
	"strings"
	"time"
)

// Block is the basic building block
type Block struct {
	data         map[string]interface{}
	hash         string
	previousHash string
	timestamp    time.Time
	pow          int
}

// Blockchain represents the blockchain
type Blockchain struct {
	genesisBlock Block
	chain        []Block
	difficulty   int
}

func (b Block) clacluateHash() string {
	data, err := json.Marshal(b.data)
	if err != nil {
		log.Fatal(err.Error())
	}
	blockData := b.previousHash + string(data) + b.timestamp.String() + strconv.Itoa(b.pow)
	blockhash := sha256.Sum256([]byte(blockData))

	return fmt.Sprintf("%x", blockhash)
}

func (b *Block) mine(difficulty int) {
	for !strings.HasPrefix(b.hash, strings.Repeat("0", difficulty)) {
		b.pow++
		b.hash = b.clacluateHash()
	}
}

// CreateBlockchain Creates a blockchain with given difficulty
func CreateBlockchain(difficulty int) Blockchain {
	genesisBlock := Block{
		hash:      "0",
		timestamp: time.Now(),
	}

	return Blockchain{
		genesisBlock: genesisBlock,
		chain:        []Block{genesisBlock},
		difficulty:   difficulty,
	}
}

func (b *Blockchain) addBlock(from, to string, amount float64) {
	blockData := map[string]interface{}{
		"from":   from,
		"to":     to,
		"amount": amount,
	}

	lastBlock := b.chain[len(b.chain)-1]
	newBlock := Block{
		data:         blockData,
		previousHash: lastBlock.hash,
		timestamp:    time.Now(),
	}
	newBlock.mine(b.difficulty)
	b.chain = append(b.chain, newBlock)
}

func (b Blockchain) isValid() bool {
	for i := range b.chain[1:] {
		previousBlock := b.chain[i]
		currentBlock := b.chain[i+1]
		if currentBlock.hash != currentBlock.clacluateHash() || currentBlock.previousHash != previousBlock.hash {
			return false
		}
	}

	return true
}

func main() {
	go func() {
		log.Println(http.ListenAndServe("localhost:6060", nil))
	}()
	start := time.Now()
	fmt.Println(start)
	blockchain := CreateBlockchain(5)

	blockchain.addBlock("Alice", "Bob", 5)
	blockchain.addBlock("John", "Bob", 2)

	fmt.Printf("%#v", blockchain)
	fmt.Println(blockchain.isValid())
	fmt.Println(time.Now().Sub(start))
}
