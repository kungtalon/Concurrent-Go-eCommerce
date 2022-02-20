package distributed

import (
	"errors"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// implementation of consistent hash function

type units []uint32

func (x units) Len() int {
	return len(x)
}

func (x units) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x units) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

var errEmpty = errors.New("No data in hash ring!")

type Consistent struct {
	// hash ring, key is the hash value, value is the node information
	circle       map[uint32]string
	sortedHashes units
	VirtualNode  int
	// synchronous lock
	sync.RWMutex
}

// NewConsistent initializes consistent hash struct and sets the default node count
func NewConsistent() *Consistent {
	return &Consistent{
		circle:      make(map[uint32]string),
		VirtualNode: 20,
	}
}

func (c *Consistent) generateKey(element string, index int) string {
	return element + strconv.Itoa(index)
}

func (c *Consistent) hashKey(key string) uint32 {
	if len(key) < 64 {
		var scratch [64]byte
		// copy, equivalent to padding
		copy(scratch[:], key)
		// use IEEE polynomial to calculate the CRC-32 check sum
		return crc32.ChecksumIEEE(scratch[:len(key)])
	}
	return crc32.ChecksumIEEE([]byte(key))
}

// updateSortedHashes resort the sortedHases slice, useful for binary Lookup
func (c *Consistent) updateSortedHashes() {
	// init with an empty slice
	hashes := c.sortedHashes[:0]

	// reset if the slice is too big
	if cap(c.sortedHashes)/(c.VirtualNode*4) > len(c.circle) {
		hashes = nil
	}

	// add existing hases to the new slice
	for k := range c.circle {
		hashes = append(hashes, k)
	}

	sort.Sort(hashes)
	c.sortedHashes = hashes
}

func (c *Consistent) Add(element string) {
	c.Lock()
	defer c.Unlock()
	c.add(element)
}

func (c *Consistent) add(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		c.circle[c.hashKey(c.generateKey(element, i))] = element
	}
	// update sorting
	c.updateSortedHashes()
}

func (c *Consistent) remove(element string) {
	for i := 0; i < c.VirtualNode; i++ {
		delete(c.circle, c.hashKey(c.generateKey(element, i)))
	}
	c.updateSortedHashes()
}

func (c *Consistent) Remove(element string) {
	c.Lock()
	defer c.Unlock()
	c.remove(element)
}

// search finds the closest node to the key on the hash ring, clockwise
func (c *Consistent) search(key uint32) int {
	f := func(x int) bool {
		return c.sortedHashes[x] > key
	}
	// call sort.Search to do binary search
	i := sort.Search(len(c.sortedHashes), f)
	if i >= len(c.sortedHashes) {
		// it means there is no element next to the key clockwise before the end
		// so the closest element must be the first element
		i = 0
	}
	return i
}

// Get returns the information of the closest server node
func (c *Consistent) Get(name string) (string, error) {
	c.RLock()
	defer c.RUnlock()
	if len(c.circle) == 0 {
		return "", errEmpty
	}
	key := c.hashKey(name)
	// use the key to find the closest node
	i := c.search(key)
	return c.circle[c.sortedHashes[i]], nil
}
