package types

import (
	"tinychain/common"
	"math/big"
	"sync/atomic"
	"encoding/json"
)

// BNonce is a 64-bit hash which proves that a sufficient amount of
// computation has been carried out on a block
type BNonce [8]byte

type Header struct {
	ParentHash common.Hash    `json:"parent_hash"` // Hash of parent block
	Height     *big.Int       `json:"height"`      // Block height
	Difficulty *big.Int       `json:"difficulty"`  // Difficulty of miner
	Root       common.Hash    `json:"root"`        // Root of merkle tree
	Miner      common.Address `json:"miner"`       // Miner address who receives reward of this block
	Extra      []byte         `json:"extra"`       // Extra data
	Nonce      BNonce         `json:"nonce"`
	Time       int64          `json:"time"` // Timestamp
}

func NewHeader(
	parentHash common.Hash,
	height *big.Int,
	difficulty *big.Int,
	root common.Hash,
	miner common.Address,
	extra []byte,
	nonce BNonce,
	tm int64,
) *Header {
	header := &Header{
		parentHash,
		height,
		difficulty,
		root,
		miner,
		extra,
		nonce,
		tm,
	}
	return header
}

func (hd *Header) Hash() common.Hash {
	data, _ := json.Marshal(hd)
	hash := common.Sha256(data)
	return hash
}

func (hd *Header) Serialize() ([]byte, error) { return json.Marshal(hd) }

func (hd *Header) Desrialize(d []byte) error { return json.Unmarshal(d, hd) }

type Block struct {
	Header       *Header        `json:"header"`
	Transactions []*Transaction `json:"transactions"`

	hash atomic.Value // Header hash

	// Total difficulty, to avoid hard fork
	// Tiny will accept the block  with the largest difficulty
	// and link it to the main chain
	TD *big.Int
}

func NewBlock(header *Header, txs []*Transaction, td *big.Int) *Block {
	block := &Block{
		Header:       header,
		Transactions: txs,
	}
	block.hash.Store(header.Hash())
	return block
}

func (bl *Block) Hash() common.Hash { return bl.hash.Load().(common.Hash) }

func (bl *Block) Serialize() ([]byte, error) { return json.Marshal(bl) }

func (bl *Block) Desrialize(d []byte) error { return json.Unmarshal(d, bl) }