package common

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math/rand"
	"unsafe"
)

// deep clone object from a to b
func Clone(src, dst interface{}) error {
	buff := new(bytes.Buffer)
	enc := gob.NewEncoder(buff)
	dec := gob.NewDecoder(buff)
	if err := enc.Encode(src); err != nil {
		return err
	}
	if err := dec.Decode(dst); err != nil {
		return err
	}
	return nil
}

// PaddedAppend appends the src byte slice to dst, returning the new slice.
// If the length of the source is smaller than the passed size, leading zero
// bytes are appended to the dst slice before appending src.
func PaddedAppend(size uint, src []byte) []byte {
	dst := make([]byte, int(size)-len(src), int(size)-len(src))
	return append(dst, src...)
}

func String(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Rcopy(dst []byte, src []byte) int {
	dstLen := len(dst)
	srcLen := len(src)
	if srcLen > dstLen {
		src = src[srcLen-dstLen:]
	}
	return copy(dst[dstLen-srcLen:], src)
}

func Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

func Expand(b []byte) []byte {
	bLen := len(b)
	expand := make([]byte, bLen*2)
	for i := range b {
		expand[i*2] = b[i] / 16
		expand[i*2+1] = b[i] % 16
	}
	for i := range expand {
		if expand[i] < 10 {
			expand[i] += '0'
		} else {
			expand[i] += 'a'
		}
	}
	return expand
}

func Compress(b []byte) []byte {
	bLen := len(b)
	compress := make([]byte, bLen/2)
	for i := range b {
		if b[i] < 10+'0' {
			b[i] -= '0'
		} else {
			b[i] -= 'a'
		}
	}
	for i := range compress {
		compress[i] = b[i*2]*16 + b[i*2+1]
	}
	return compress
}

func PrettyPrintJson(b []byte) []byte {
	var out bytes.Buffer
	err := json.Indent(&out, b, "", "  ")
	if err != nil {
		fmt.Println(err)
	}
	return out.Bytes()
}

func RandString(ulen uint8) string {
	slen := int(ulen) / 2
	strBytes := make([]byte, slen)
	for i := 0; i < slen; i += 1 {
		strBytes[i] = byte(rand.Int() % 256)
	}
	randStr := string(Expand(strBytes))
	return randStr
}
