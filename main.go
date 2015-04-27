package bdns

import (
	"encoding/hex"
	"flag"
	"fmt"
	"net"
)

func main() {
	var localAddr = flag.String("host", ":8008", "Local address to listen for messages")

	addr, err := net.ResolveUDPAddr("udp", *localAddr)
	checkErr(err)

	conn, err := net.ListenUDP("udp", addr)
	checkErr(err)

	defer conn.Close()

	go handleRequests(conn)

	var input string
	fmt.Scanln(&input)
}

func handleRequests(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		checkErr(err)

		fmt.Println(hex.Dump(buffer[0:n]))
		fmt.Println("Received ", string(buffer[0:n]), " from ", addr)
		NewMessage(buffer[0:n])
	}
}
