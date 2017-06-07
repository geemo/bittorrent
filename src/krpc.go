package dht

import (
	"fmt"
	"net"
)

// Krpc struct
type Krpc struct {
	conn *net.UDPConn
}

// NewKrpc instance
func NewKrpc() (Krpc, error) {
	var k Krpc
	ips, err := net.LookupIP("router.bittorrent.com")
	if err != nil {
		return k, err
	}
	fmt.Println(ips, "00")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ips[0], Port: 6881}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		return k, err
	}
	k = Krpc{conn}
	return k, nil
}

func (k *Krpc) send(data Map) (Map, error) {
	b, err := Encode(data)
	if err != nil {
		return nil, err
	}
	k.conn.Write(b)
	res := make([]byte, 1024)
	n, err := k.conn.Read(res)
	if err != nil {
		return nil, err
	}
	msg, err := Decode(res[:n])

	if err != nil {
		return nil, err
	}
	return msg.(Map), nil
}

func getTransactionID() {}

func ping() {}

func findNode() {}

func getPeers() {}

func announcePeer() {}
