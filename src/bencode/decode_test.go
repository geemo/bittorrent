package bencode

import (
	"fmt"
	"testing"
)

func TestSpace(t *testing.T) {
	// var l = [2]Any{"spam", "eggs"}
	successData := map[string]Any{
		"4:spam": "spam",
		"i3e":    3,
		// "l4:spam4:eggse": l,
		// "d3:cow3:moo4:spam4:eggse": 1,
	}
	for str, value := range successData {
		var torrent = NewTorrent(str)
		var v, err = torrent.value()
		fmt.Println(v, value, err)
		if v != value {
			t.Error("torrent 解析错误")
		}
	}
}
