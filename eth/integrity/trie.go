package integrity

import (
	"encoding/binary"
	"fmt"
	"math/bits"

	"github.com/ledgerwatch/turbo-geth/common"
	"github.com/ledgerwatch/turbo-geth/common/dbutils"
	"github.com/ledgerwatch/turbo-geth/ethdb"
)

// AssertSubset a & b == a - checks whether a is subset of b
func AssertSubset(a, b uint16) {
	if (a & b) != a {
		panic(fmt.Errorf("invariant 'is subset' failed: %b, %b", a, b))
	}
}

func Trie(tx ethdb.Tx) {
	{
		c := tx.Cursor(dbutils.TrieOfAccountsBucket)
		defer c.Close()
		parentC := tx.Cursor(dbutils.TrieOfAccountsBucket)
		defer parentC.Close()
		for k, v, err := c.First(); k != nil; k, v, err = c.Next() {
			if err != nil {
				panic(err)
			}
			if len(k) == 1 {
				continue
			}
			hasState := binary.BigEndian.Uint16(v)
			hasBranch := binary.BigEndian.Uint16(v[2:])
			hasHash := binary.BigEndian.Uint16(v[4:])
			AssertSubset(hasBranch, hasState)
			AssertSubset(hasHash, hasState)
			if bits.OnesCount16(hasHash) != len(v[6:])/common.HashLength {
				panic(fmt.Errorf("invariant bits.OnesCount16(hasHash) == len(hashes) failed: %d, %d", bits.OnesCount16(hasHash), len(v[6:])/common.HashLength))
			}
			found := false
			var parentK []byte
			for i := len(k) - 1; i > 0; i-- {
				parentK = k[:i]
				kParent, vParent, err := parentC.SeekExact(parentK)
				if err != nil {
					panic(err)
				}
				if kParent == nil {
					continue
				}
				found = true
				parentHasBranch := binary.BigEndian.Uint16(vParent[2:])
				parentHasBit := uint16(1)<<uint16(k[len(parentK)])&parentHasBranch != 0
				if !parentHasBit {
					panic(fmt.Errorf("for %x found parent %x, but it has no branchBit: %016b", k, parentK, parentHasBranch))
				}
			}
			if !found {
				panic(fmt.Errorf("trie hash %x has no parent", k))
			}
		}
	}
	{
		c := tx.Cursor(dbutils.TrieOfStorageBucket)
		defer c.Close()
		parentC := tx.Cursor(dbutils.TrieOfAccountsBucket)
		defer parentC.Close()
		for k, v, err := c.First(); k != nil; k, v, err = c.Next() {
			if err != nil {
				panic(err)
			}
			if len(k) == 40 {
				continue
			}
			hasState := binary.BigEndian.Uint16(v)
			hasBranch := binary.BigEndian.Uint16(v[2:])
			hasHash := binary.BigEndian.Uint16(v[4:])
			AssertSubset(hasBranch, hasState)
			AssertSubset(hasHash, hasState)
			if bits.OnesCount16(hasHash) != len(v[6:])/common.HashLength {
				panic(fmt.Errorf("invariant bits.OnesCount16(hasHash) == len(hashes) failed: %d, %d", bits.OnesCount16(hasHash), len(v[6:])/common.HashLength))
			}

			found := false
			var parentK []byte
			for i := len(k) - 1; i >= 40; i-- {
				parentK = k[:i]
				kParent, vParent, err := parentC.SeekExact(parentK)
				fmt.Printf("qaa: %x,%x\n", kParent, vParent)
				if err != nil {
					panic(err)
				}
				if kParent == nil {
					continue
				}
				found = true
				parentBranches := binary.BigEndian.Uint16(vParent[2:])
				parentHasBit := uint16(1)<<uint16(k[len(parentK)])&parentBranches != 0
				if !parentHasBit {
					panic(fmt.Errorf("for %x found parent %x, but it has no branchBit for child: %016b", k, parentK, parentBranches))
				}
			}
			if !found {
				panic(fmt.Errorf("trie hash %x has no parent. Last checked: %x", k, parentK))
			}
		}
	}
}
