package chain

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Block is a block in the blockchain
type Block struct {
	Index        int       `json:"index"`
	PreviousHash string    `json:"peviousHash"`
	Timestamp    time.Time `json:"timestamp"`
	Data         string    `json:"data"`
	Hash         string    `json:"hash"`
	Nonce        int       `json:"nonce"`
}

// Blockchain holds the blockchain
type Blockchain struct {
	Blocks            []Block
	TotalTransactions int       `json:"totalTransactions"`
	Inception         time.Time `json:"inception"`
}

// BlockData data to store on blockchain
type BlockData struct {
	I string `json:"i"`
	B int    `json:"b"`
}

// AddBlock to Blockchain struct
func (chain *Blockchain) AddBlock(block Block) []Block {
	chain.Blocks = append(chain.Blocks, block)
	return chain.Blocks
}

// IsPrime is # prime
func IsPrime(num int) bool {
	for i := 2; i < num; i++ {
		if num%i == 0 {
			return false
		}
	}
	return num != 1
}

func calculateHashForBlock(index int, previousHash string, timestamp time.Time, data BlockData, nonce int, chain Blockchain) string {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	h := sha256.New()
	h.Write([]byte(strconv.Itoa(index) + previousHash + timestamp.String() + string(bytes) + strconv.Itoa(nonce)))
	chain.TotalTransactions++
	hash := hex.EncodeToString(h.Sum(nil))
	return hash
}

// GetGenesis gets a genesis block
func GetGenesis() Block {
	block := Block{0, "0", time.Now().UTC(), "{ name: 'Genesis Block', value: 0 }", "1d79e9eef321cac0aa8f73d1245a5604a8a665e6daacf64d1b9843e2ab98fa29", 745}
	return block
}

// IsValidHashDifficulty check if hash is valid difficult
func IsValidHashDifficulty(hash string) bool {
	totalCount := 0
	charACount := 0

	for i := 0; i < len(hash)-1; i++ {
		if string([]rune(hash)[i]) == "a" {
			charACount++
		}
		hashItem := string(hash[i])
		numberFromHash, err := strconv.Atoi(hashItem)
		if err == nil {
			totalCount += numberFromHash
		}
	}
	isValid := IsPrime(totalCount) && charACount >= 10
	return isValid
}

//GenerateNextBlock gen block of chain
func GenerateNextBlock(chain Blockchain, data BlockData) Block {
	latestBlock := chain.Blocks[len(chain.Blocks)-1]
	nextIndex := latestBlock.Index
	previousHash := latestBlock.PreviousHash
	timestamp := time.Now().UTC()
	nonce := 0
	nextHash := calculateHashForBlock(nextIndex, previousHash, timestamp, data, nonce, chain)

	for !IsValidHashDifficulty(nextHash) {
		nonce++
		timestamp = time.Now().UTC()
		nextHash = calculateHashForBlock(nextIndex, previousHash, timestamp, data, nonce, chain)
	}

	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	return Block{nextIndex, previousHash, timestamp, string(bytes), nextHash, nonce}
}

// Mine new blocks
func Mine(chain Blockchain, data BlockData) {
	newBlock := GenerateNextBlock(chain, data)
	// fmt.Println("newBlock", newBlock)
	chain.AddBlock(newBlock)
}

// Main intiailizes blockchain
func Main() {
	start := time.Now().UTC()
	chainToMine := []BlockData{
		{"test 1", 1},
		{"test 2", 2},
		{"test 3", 3},
		{"test 4", 4},
		{"test 5", 5},
		{"test 6", 6},
		{"test 13", 13},
	}
	initialBlock := GetGenesis()
	blocks := []Block{initialBlock}
	Chain := Blockchain{blocks, 0, time.Now().UTC()}

	for i := 0; i < len(chainToMine); i++ {
		Mine(Chain, chainToMine[i])
	}
	duration := time.Since(start)
	fmt.Println(duration.Seconds())
}
