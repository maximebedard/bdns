package bdns

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net"
)

type Message struct {
	Header       Header
	Questions    []Question
	Answers      []ResourceRecord
	Authorities  []ResourceRecord
	Additionnals []ResourceRecord
}

type Header struct {
	Id      uint16
	Flags   uint16
	Qdcount uint16
	Ancount uint16
	Nscount uint16
	Arcount uint16
}

type Question struct {
	Qname  string
	Qtype  uint16
	Qclass uint16
}

type ResourceRecord struct {
	Rrname    string
	Rrtype    uint16
	Rrclass   uint16
	Rrttl     uint32
	Rrdlength uint16
	Rrdata    []byte
}

func main() {
	var localAddr = flag.String("host", ":8008", "Local address to listen for messages")

	addr, err := net.ResolveUDPAddr("udp", *localAddr)
	checkErr(err)

	conn, err := net.ListenUDP("udp", addr)
	checkErr(err)

	defer conn.Close()

	go handleConnection(conn)

	var input string
	fmt.Scanln(&input)
}

func handleConnection(conn *net.UDPConn) {
	buffer := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		checkErr(err)

		fmt.Println(hex.Dump(buffer[0:n]))
		fmt.Println("Received ", string(buffer[0:n]), " from ", addr)
		newHeader(buffer[0:n])
	}
}

func NewMessage(buffer []byte) *Message {
	message := new(Message)
	header := newHeader(buffer)
	message.Header = *header
	message.Questions = newQuestions(buffer, header.Qdcount)
	message.Answers = newResourceRecords(buffer, header.Ancount, binary.Size(message))
	message.Authorities = newResourceRecords(buffer, header.Nscount, binary.Size(message))
	message.Additionnals = newResourceRecords(buffer, header.Arcount, binary.Size(message))

	return message
}

func newHeader(buffer []byte) *Header {
	reader := bytes.NewReader(buffer)
	header := new(Header)

	err := binary.Read(reader, binary.BigEndian, header)
	checkErr(err)

	return header
}

func newQuestions(buffer []byte, count uint16) []Question {
	questions := make([]Question, count)
	offset := 12
	for _, question := range questions {
		//question
		label, length := parseLabel(buffer[offset:])
		question.Qname = label
		offset += length - 10

		question.Qtype = binary.BigEndian.Uint16(buffer[offset:2])
		offset += 2

		question.Qclass = binary.BigEndian.Uint16(buffer[offset:2])
		offset += 2
	}
	return questions
}

func newResourceRecords(buffer []byte, count uint16, offset int) []ResourceRecord {
	return nil
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func parseLabel(label []byte) (string, int) {
	var buffer bytes.Buffer
	i, j := 0, 0
	for ; label[i] != 0x00; i++ {
		if i-j == 0 {
			j += int(label[i]) + 1
			if i > 0 {
				buffer.WriteString(".")
			}
			continue
		}

		buffer.WriteString(string(label[i]))
	}
	return buffer.String(), i
}
