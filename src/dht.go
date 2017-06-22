package dht

import "net"

// DHT struct
type DHT struct {
	selfNode      *Node
	selfAddr      *net.UDPAddr
	rootNode      *BitMap
	k             *Krpc
	transactionID int
}

// ips, err := net.LookupIP("router.bittorrent.com")
var maxTansactionID = 0xffffffff

// NewDHT struct
func NewDHT() *DHT {
	return &DHT{}
}

func (dht *DHT) getTransactionID() int {
	if dht.transactionID >= maxTansactionID {
		dht.transactionID = 0
		return dht.transactionID
	}
	dht.transactionID++
	return dht.transactionID
}
