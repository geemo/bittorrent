package dht

import (
	"fmt"
	"strings"
)

// BitMap struct
type BitMap struct {
	data []byte
	Size int
}

// NewBitMap func
func NewBitMap(Size int) *BitMap {
	div, mod := Size/8, Size%8
	if mod > 0 {
		div++
	}
	return &BitMap{make([]byte, div), Size}
}

// NewBitMapfromString func
func NewBitMapfromString(s string) *BitMap {
	return NewBitMapfromBytes([]byte(s))
}

// NewBitMapfromBytes func
func NewBitMapfromBytes(b []byte) *BitMap {
	Size := len(b) * 8
	return &BitMap{b, Size}
}

func (bitMap *BitMap) set(i int, bit int) {
	if i >= bitMap.Size {
		panic("set bitMap out of index")
	}

	div, mod := i/8, i%8
	shift := byte(1 << uint(7-mod))
	bitMap.data[div] &= ^shift
	if bit > 0 {
		bitMap.data[div] |= shift
	}
}

// Bit return bit of bitMap
func (bitMap *BitMap) Bit(i int) int {
	if i >= bitMap.Size {
		panic("index out of range")
	}

	div, mod := i/8, i%8
	return int((uint(bitMap.data[div]) & (1 << uint(7-mod))) >> uint(7-mod))
}

// Set bit 1
func (bitMap *BitMap) Set(index int) {
	bitMap.set(index, 1)
}

// Unset set bit 0
func (bitMap *BitMap) Unset(index int) {
	bitMap.set(index, 0)
}

// Prefix two bitmap
func (bitMap *BitMap) Prefix(compareBitMap *BitMap) *BitMap {
	if bitMap.Size != compareBitMap.Size {
		panic("two bit map have differece size, can't get prefix")
	}
	prefixBitMap := NewBitMap(bitMap.Size)
	for index := 0; index < bitMap.Size; index++ {
		bit := bitMap.Bit(index)
		if bit == compareBitMap.Bit(index) {
			prefixBitMap.set(index, bit)
		} else {
			prefixBitMap.Size = index
			return prefixBitMap
		}
	}
	return prefixBitMap
}

// String []byte to binary string
func (bitMap *BitMap) String() string {
	div, mod := bitMap.Size/8, bitMap.Size%8
	buff := make([]string, div+mod)

	for i := 0; i < div; i++ {
		buff[i] = fmt.Sprintf("%08b", bitMap.data[i])
	}

	for i := div; i < div+mod; i++ {
		buff[i] = fmt.Sprintf("%1b", bitMap.Bit(div*8+(i-div)))
	}

	return strings.Join(buff, "")
}

// RawString []byte to string
func (bitMap *BitMap) RawString() string {
	return string(bitMap.data)
}
