package bencode

import "testing"

func TestSpace(t *testing.T) {
	successData := map[string]Any{
		"4:spam": "spam",
		"i3e":    3,
	}
	for str, value := range successData {
		var torrent = NewTorrent(str)
		var v, _ = torrent.value()
		if v != value {
			t.Error("parse torrent 解析错误")
		}
	}
}

func TestList(t *testing.T) {
	var l = [2]string{"spam", "eggs"}
	var torrent = NewTorrent("l4:spam4:eggse")
	var value, _ = torrent.value()
	if expect, ok := value.([]Any); ok {
		for index, v := range l {
			if v != expect[index] {
				t.Error("parse torrent list 解析错误")
			}
		}
	}
}

func TestMap(t *testing.T) {
	var m = map[string]Any{"cow": "moo", "spam": "eggs"}
	var torrent = NewTorrent("d3:cow3:moo4:spam4:eggse")
	var value, _ = torrent.value()
	if expect, ok := value.(map[string]Any); ok {
		for k, v := range m {
			if v != expect[k] {
				t.Error("parse torrent map 解析错误")
			}
		}
	}
}
