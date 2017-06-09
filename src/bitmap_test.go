package dht

import "testing"

func TestBitMap(t *testing.T) {
	b := []byte{byte(0xfe)}
	s := "11111110"

	cb := []byte{byte(0xde)}
	//11011110

	bitMap := NewBitMapfromBytes(b)
	cBitMap := NewBitMapfromBytes(cb)
	preBitMap := bitMap.Prefix(cBitMap)
	if preBitMap.String() != "11" || preBitMap.Size != 2 {
		t.Error("bit map prefix error")
	}

	if bitMap.String() != s {
		t.Error("bit map from bytes error")
	}

	bitMapFromString := NewBitMapfromString(string(b))
	if bitMapFromString.String() != s {
		t.Error("bit map from string error")
	}

	if bitMap.Bit(7) != 0 {
		t.Error("bit map Bit error")
	}

	bitMap.Set(7)
	if bitMap.Bit(7) != 1 {
		t.Error("bit map Set error")
	}

	bitMap.Unset(7)
	if bitMap.Bit(7) != 0 {
		t.Error("bit map Unset error")
	}

	if bitMap.RawString() != string(b) {
		t.Error("bit map rawString error")
	}
}
