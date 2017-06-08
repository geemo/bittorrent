package dht

import (
	"fmt"
	"net"
	"testing"
)

func TestSend(t *testing.T) {
	ips, err := net.LookupIP("router.bittorrent.com")
	if err != nil {
		fmt.Println(err)
	}
	ip := ips[0]
	port := 6881
	fmt.Println(ips, "00")
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: ip, Port: port}
	conn, err := net.DialUDP("udp", srcAddr, dstAddr)
	if err != nil {
		fmt.Println(err, "111")
	}

	defer conn.Close()
	// var m = Map{"t": "aa", "y": "q", "q": "ping", "a": Map{"id": "abcdefghij0123456789"}}
	var m = Map{"t": "aa", "y": "q", "q": "find_node", "a": Map{"id": "abcdefghij0123456789", "target": "mnopqrstuvwxyz123456"}}
	// var m = Map{"t": "aa", "y": "q", "q": "get_peers", "a": Map{"id": "abcdefghij0123456789", "info_hash": "mnopqrstuvwxyz123456"}}
	b, err := Encode(m)
	fmt.Println(string(b))
	if err != nil {
		fmt.Println(err)
	}
	conn.Write(b)
	data := make([]byte, 1024)
	n, err := conn.Read(data)
	if err != nil {
		fmt.Println(err, "22")
	}
	fmt.Println(string(data[:n]))
	r, err := Decode(data[:n])
	if err != nil {
		fmt.Println(err)
	}
	v := Value{r}
	s := v.key("r").key("nodes").toString()
	nodes, err := DecodeNodes(s)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println([]byte("mnopqrstuvwxyz123456"))
	for _, node := range nodes {
		fmt.Println(node)
	}
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("read333 %s from <%s>\n", r, conn.RemoteAddr())

}
