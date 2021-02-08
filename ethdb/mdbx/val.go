package mdbx

/*
#include <stdlib.h>
#include <stdio.h>
#include "mdbxgo.h"
#include "dist/mdbx.h"
*/
import "C"

import (
	"unsafe"

	"github.com/ledgerwatch/turbo-geth/ethdb/mdbx/internal/lmdbarch"
)

// Just for docs:
//struct MDBX_val {
//	void *iov_base; /**< pointer to some data */
//	size_t iov_len; /**< the length of data in bytes */
//};

// valSizeBits is the number of bits which constraining the length of the
// single values in an LMDB database, either 32 or 31 depending on the
// platform.  valMaxSize is the largest data size allowed based.  See runtime
// source file malloc.go and the compiler typecheck.go for more information
// about memory limits and array bound limits.
//
//		https://github.com/golang/go/blob/a03bdc3e6bea34abd5077205371e6fb9ef354481/src/runtime/malloc.go#L151-L164
//		https://github.com/golang/go/blob/36a80c5941ec36d9c44d6f3c068d13201e023b5f/src/cmd/compile/internal/gc/typecheck.go#L383
//
// On 64-bit systems, luckily, the value 2^32-1 coincides with the maximum data
// size for LMDB (MAXDATASIZE).
const (
	valSizeBits = lmdbarch.Width64*32 + (1-lmdbarch.Width64)*31
	valMaxSize  = 1<<valSizeBits - 1
)

var eb = []byte{0}

func valBytes(b []byte) ([]byte, int) {
	if len(b) == 0 {
		return eb, 0
	}
	return b, len(b)
}

func wrapVal(b []byte) *C.MDBX_val {
	p, n := valBytes(b)
	return &C.MDBX_val{
		iov_base: unsafe.Pointer(&p[0]),
		iov_len:  C.size_t(n),
	}
}

func getBytes(val *C.MDBX_val) []byte {
	return (*[valMaxSize]byte)(val.iov_base)[:val.iov_len:val.iov_len]
}

func getBytesCopy(val *C.MDBX_val) []byte {
	return C.GoBytes(val.iov_base, C.int(val.iov_len))
}
