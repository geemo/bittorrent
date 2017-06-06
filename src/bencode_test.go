package dht

import "testing"

func TestDecodeBasic(t *testing.T) {
	testData := Map{
		"4:spam": "spam",
		"i3e":    3,
	}
	for str, value := range testData {
		var v, _ = Decode(str)
		if v != value {
			t.Error("decode torrent error")
		}
	}
}

func TestDecodeList(t *testing.T) {
	var l = List{"spam", "eggs"}
	var value, _ = Decode("l4:spam4:eggse")
	if expect, ok := value.(List); ok {
		for index, v := range l {
			if v != expect[index] {
				t.Error("decode torrent list error")
			}
		}
	}
}

func TestDecodeMap(t *testing.T) {
	var m = Map{"cow": "moo", "spam": "eggs"}
	var value, _ = Decode("d3:cow3:moo4:spam4:eggse")
	if expect, ok := value.(Map); ok {
		for k, v := range m {
			if v != expect[k] {
				t.Error("decode torrent map error")
			}
		}
	}
}

func TestEncodeBasic(t *testing.T) {
	testData := Map{
		"4:spam": "spam",
		"i3e":    3,
	}
	for str, value := range testData {
		var s, _ = Encode(value)
		if s != str {
			t.Error("encode basic error")
		}
	}
}

func TestEncodeList(t *testing.T) {
	var l = List{"spam", "eggs"}
	var str = "l4:spam4:eggse"
	var s, _ = Encode(l)
	if s != str {
		t.Error("encode list error")
	}
}

func TestEncodeMap(t *testing.T) {
	var m = Map{"cow": "moo", "spam": "eggs"}
	var str = "d3:cow3:moo4:spam4:eggse"
	var s, _ = Encode(m)
	if s != str {
		t.Error("encode map error")
	}
}
