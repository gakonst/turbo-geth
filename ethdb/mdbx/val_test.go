package mdbx

import (
	"bytes"
	"testing"
)

func TestValBytes(t *testing.T) {
	ptr, n := valBytes(nil)
	if len(ptr) == 0 {
		t.Errorf("unexpected unaddressable slice")
	}
	if n != 0 {
		t.Errorf("unexpected length: %d (expected 0)", n)
	}

	b := []byte("abc")
	ptr, n = valBytes(b)
	if len(ptr) == 0 {
		t.Errorf("unexpected unaddressable slice")
	}
	if n != 3 {
		t.Errorf("unexpected length: %d (expected %d)", n, len(b))
	}
}

func TestVal(t *testing.T) {
	orig := []byte("hey hey")
	val := wrapVal(orig)

	p := getBytes(val)
	if !bytes.Equal(p, orig) {
		t.Errorf("getBytes() not the same as original data: %q", p)
	}
	if &p[0] != &orig[0] {
		t.Errorf("getBytes() is not the same slice as original")
	}

	p = getBytesCopy(val)
	if !bytes.Equal(p, orig) {
		t.Errorf("getBytesCopy() not the same as original data: %q", p)
	}
	if &p[0] == &orig[0] {
		t.Errorf("getBytesCopy() overlaps with orignal slice")
	}
}
