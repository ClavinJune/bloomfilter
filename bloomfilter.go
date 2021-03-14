package bloomfilter

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"hash"
	"math"
	"strings"
)

// BloomFilter is struct is a simple implementation for https://en.wikipedia.org/wiki/Bloom_filter
// This implementation is inspired from https://pypi.org/project/pybloom/
type BloomFilter struct {
	Bitset  []bool
	HashFns []hash.Hash
}

// Add stores sentence to the bitset word by word
func (bf *BloomFilter) Add(sentence string) {
	bf.add(bf.tokenize(sentence)...)
}

// Check checks if sentence exists in the bitset
func (bf *BloomFilter) Check(sentence string) bool {
	for _, w := range bf.tokenize(sentence) {
		if !bf.check(w) {
			return false
		}
	}

	return true
}

// add stores words to the bitset
func (bf *BloomFilter) add(words ...string) {
	m := len(bf.Bitset)

	for _, w := range words {
		hs := bf.hash(w)

		for _, h := range hs {
			bf.Bitset[h%uint64(m)] = true
		}
	}
}

// check checks a single string if it exists in the bitset
func (bf *BloomFilter) check(w string) bool {
	m := len(bf.Bitset)

	hs := bf.hash(w)

	for _, h := range hs {
		if !bf.Bitset[h%uint64(m)] {
			return false
		}
	}

	return true
}

// tokenize splits sentence into words
func (bf *BloomFilter) tokenize(sentence string) []string {
	return strings.Split(strings.ToLower(sentence), " ")
}

// hash hashes a single string using BloomFilter.HashFns
func (bf *BloomFilter) hash(w string) []uint64 {
	var r []uint64

	b := []byte(w)
	for _, fn := range bf.HashFns {
		fn.Write(b)
		u, _ := binary.ReadUvarint(bytes.NewBuffer(fn.Sum(nil)))

		r = append(r, u)
		fn.Reset()
	}

	return r
}

// New creates pointer of BloomFilter by given n as capacity, and e as error rate
func New(n int, e float64) (*BloomFilter, error) {
	if n <= 0 {
		return nil, ErrCapacity
	}

	if e <= 0 || e >= 1 {
		return nil, ErrErrorRate
	}

	k := findK(e)
	m := findM(n, e)

	return &BloomFilter{
		Bitset:  make([]bool, m, m),
		HashFns: createHashFns(k),
	}, nil
}

// findK inspired from pybloom
func findK(e float64) int {
	return int(math.Ceil(math.Abs(math.Log2(1 / e))))
}

// findM inspired from pybloom
// I adjust the x divisor by removing the K multiplication
func findM(n int, e float64) int {
	log2Sqr := math.Pow(math.Log(2), 2)
	logErrAbs := math.Abs(math.Log(e))
	x := (float64(n) * logErrAbs) / log2Sqr
	return int(math.Abs(math.Ceil(x)))
}

// createHashFns creates k hmac function with different key
func createHashFns(k int) []hash.Hash {
	result := make([]hash.Hash, k, k)

	for i := range result {
		result[i] = hmac.New(sha256.New, []byte(fmt.Sprintf("key-%d", i)))
	}

	return result
}
