package glass

/* This file implements the Fowler–Noll–Vo hash function
   (warning: this is not a cryptographic hash)
*/

import (
    "hash"
    "bytes"
    "encoding/binary"
)

const offset_basis = 14695981039346656037
const prime = 1099511628211

type fnv1 struct {
    state uint64
}

func NewFNV1() hash.Hash {
    var this fnv1
    this.Reset()
    return &this
}

func (this *fnv1) BlockSize() int {
    return 1
}

func (this *fnv1) Reset() {
    this.state = offset_basis
}

func (this *fnv1) Size() int {
    return 8
}

func (this *fnv1) Write(p []byte) (n int, err error) {
    state := this.state
    n = len(p)
    for i := 0; i < n; i++ {
		state = state * prime
        state = state ^ uint64(p[i])
    }
    this.state = state
    return n, nil
}

func (this *fnv1) Sum(in []byte) []byte {
    buf := new(bytes.Buffer)
    err := binary.Write(buf, binary.LittleEndian, &this.state)
    if err != nil {
        return nil
    }
    b := buf.Bytes()
    return append(in, b[:]...)
}
