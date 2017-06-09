package dht

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

// Node info struct
type Node struct {
	id   *BitMap
	ip   string
	port int
}

var chars = "0123456789abcdef"

// IDToHex ..
func IDToHex(id *BitMap) string {
	var buf bytes.Buffer
	for _, b := range id.data {
		buf.WriteByte(chars[b>>4])
		buf.WriteByte(chars[b&0x0f])
	}
	return buf.String()
}

// HexToID ..
func HexToID(hex string) *BitMap {
	return nil
}

// RandomID get
func RandomID() (*BitMap, error) {
	var (
		id  *BitMap
		err error
	)

	return id, err
}

// EnCodeNodes struct
func EnCodeNodes(nodes []Node) ([]byte, error) {
	var (
		b   []byte
		err error
	)
	for _, node := range nodes {
		fmt.Println(node)
	}
	return b, err
}

// DecodeNodes struct
func DecodeNodes(s string) ([]Node, error) {
	var (
		nodes   []Node
		err     error
		nodeLen = 26
		b       = []byte(s)
		start   = 0
		lenB    = len(b)
	)
	if lenB < nodeLen {
		return nodes, errors.New("can't parse nodes from empty string")
	}
	if lenB%nodeLen != 0 {
		return nodes, errors.New("nodes string length not match")
	}

	for {
		s := start * nodeLen
		ipStart := s + 20
		portStart := s + 24
		e := (start + 1) * nodeLen
		if e > lenB {
			break
		}

		id := NewBitMapfromBytes(b[s:ipStart])
		var ip []string
		for _, ipByte := range b[ipStart:portStart] {
			ip = append(ip, strconv.Itoa(int(ipByte)))
		}
		port := int(binary.BigEndian.Uint16(b[portStart:e]))
		node := Node{id, strings.Join(ip, "."), port}
		nodes = append(nodes, node)
		start++
	}

	return nodes, err
}
