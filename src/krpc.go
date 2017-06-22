package dht

import (
	"io/ioutil"
	"net"
	"strconv"
)

// Krpc struct
type Krpc struct {
	conn *net.UDPConn
}

// NewKrpc instance
func NewKrpc(ip string, port int) (Krpc, error) {
	var k Krpc
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		return k, err
	}
	k = Krpc{conn}
	return k, nil
}

func (k *Krpc) send(data Map, addr *net.UDPAddr) (Map, error) {
	b, err := Encode(data)
	if err != nil {
		return nil, err
	}
	k.conn.WriteToUDP(b, addr)
	res, err := ioutil.ReadAll(k.conn)
	if err != nil {
		return nil, err
	}

	msg, err := Decode(res)
	if err != nil {
		return nil, err
	}
	return msg.(Map), nil
}

func (k *Krpc) ping(dht *DHT, addr *net.UDPAddr) (Map, error) {
	query := Map{
		"t": strconv.Itoa(dht.getTransactionID()),
		"y": "q",
		"q": "ping",
		"a": Map{
			"id": dht.selfNode.id.RawString(),
		},
	}
	return k.send(query, addr)
}

func (k *Krpc) findNode(dht *DHT, addr *net.UDPAddr, targetNode *Node) (Map, error) {
	query := Map{
		"t": strconv.Itoa(dht.getTransactionID()),
		"y": "q",
		"q": "find_node",
		"a": Map{
			"id":     dht.selfNode.id.RawString(),
			"target": targetNode.id.RawString(),
		},
	}
	return k.send(query, addr)
}

func (k *Krpc) getPeers(dht *DHT, addr *net.UDPAddr, infoHash *BitMap) (Map, error) {
	query := Map{
		"t": strconv.Itoa(dht.getTransactionID()),
		"y": "q",
		"q": "get_peers",
		"a": Map{
			"id":        dht.selfNode.id.RawString(),
			"info_hash": infoHash.RawString(),
		},
	}
	return k.send(query, addr)
}

func (k *Krpc) announcePeer(dht *DHT, addr *net.UDPAddr, token string, infoHash *BitMap) (Map, error) {
	query := Map{
		"t": strconv.Itoa(dht.getTransactionID()),
		"y": "q",
		"q": "announce_peer",
		"a": Map{
			"implied_port": 1,
			"port":         12231,
			"id":           dht.selfNode.id.RawString(),
			"info_hash":    infoHash.RawString(),
			"token":        token,
		},
	}
	return k.send(query, addr)
}
