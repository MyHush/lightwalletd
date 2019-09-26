package common

import (
	"bytes"
	"sync"

	"github.com/adityapk00/lightwalletd/walletrpc"
	"github.com/pkg/errors"
)

type BlockCache struct {
	MaxEntries int

	FirstBlock int
	LastBlock  int

	m map[int]*walletrpc.CompactBlock

	mutex sync.RWMutex
}

func NewBlockCache(maxEntries int) *BlockCache {
	return &BlockCache{
		MaxEntries: maxEntries,
		FirstBlock: -1,
		LastBlock:  -1,
		m:          make(map[int]*walletrpc.CompactBlock),
	}
}

func (c *BlockCache) Add(height int, block *walletrpc.CompactBlock) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	//println("Cache add", height)
	if c.FirstBlock == -1 && c.LastBlock == -1 {
		// If this is the first block, prep the data structure
		c.FirstBlock = height
		c.LastBlock = height - 1
	}

	// Don't allow out-of-order blocks. This is more of a sanity check than anything
	// If there is a reorg, then the ingestor needs to handle it.
	if c.m[height-1] != nil && !bytes.Equal(block.PrevHash, c.m[height-1].Hash) {
		return errors.New("Prev hash of the block didn't match")
	}

	// Add the entry and update the counters
	c.m[height] = block

	c.LastBlock = height

	// If the cache is full, remove the oldest block
	if c.LastBlock-c.FirstBlock+1 > c.MaxEntries {
		//println("Deleteing at height", c.FirstBlock)
		delete(c.m, c.FirstBlock)
		c.FirstBlock = c.FirstBlock + 1
	}

	//println("Cache size is ", len(c.m))
	return nil
}

func (c *BlockCache) Get(height int) *walletrpc.CompactBlock {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	//println("Cache get", height)
	if c.LastBlock == -1 || c.FirstBlock == -1 {
		return nil
	}

	if height < c.FirstBlock || height > c.LastBlock {
		//println("Cache miss: index out of range")
		return nil
	}

	//println("Cache returned")
	return c.m[height]
}

func (c *BlockCache) GetLatestBlock() int {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	return c.LastBlock
}
