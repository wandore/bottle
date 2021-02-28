package hash

import (
	"hash/crc32"
	"sort"
	"strconv"
)

type ConHash struct {
	hashFunc func([]byte) uint32
	backup   int
	keys     []int
	hashMap  map[int]string
}

func New(backup int, f func([]byte) uint32) *ConHash {
	c := &ConHash{
		hashFunc: f,
		backup:   backup,
		keys:     make([]int, 0),
		hashMap:  make(map[int]string, 0),
	}
	if c.hashFunc == nil {
		c.hashFunc = crc32.ChecksumIEEE
	}
	return c
}

func (c *ConHash) Add(keys ...string) {
	for _, key := range keys {
		for i := 0; i < c.backup; i++ {
			hashValue := int(c.hashFunc([]byte(strconv.Itoa(i) + "-" + key)))
			c.keys = append(c.keys, hashValue)
			c.hashMap[hashValue] = key
		}
	}
	sort.Ints(c.keys)
}

func (c *ConHash) Get(key string) string {
	if len(c.keys) == 0 {
		return ""
	}

	hashValue := int(c.hashFunc([]byte(key)))

	index := sort.Search(len(c.keys), func(i int) bool {
		return c.keys[i] >= hashValue
	})

	return c.hashMap[c.keys[index%len(c.keys)]]
}
