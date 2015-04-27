package bdns

import (
	"bytes"
	"encoding/binary"
	"log"
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

func NewMessage(buffer []byte) *Message {
	var offset int
	message := new(Message)
	header := newHeader(buffer, &offset)
	message.Header = *header
	message.Questions = newQuestions(buffer, header.Qdcount, &offset)
	message.Answers = newResourceRecords(buffer, header.Ancount, &offset)
	message.Authorities = newResourceRecords(buffer, header.Nscount, &offset)
	message.Additionnals = newResourceRecords(buffer, header.Arcount, &offset)

	return message
}

func newHeader(buffer []byte, offset *int) *Header {
	reader := bytes.NewReader(buffer)
	header := new(Header)

	err := binary.Read(reader, binary.BigEndian, header)
	checkErr(err)
	*offset = binary.Size(header)

	return header
}

func newQuestions(buffer []byte, count uint16, offset *int) []Question {
	questions := make([]Question, count)
	for i, _ := range questions {
		name := parseLabel(buffer[*offset:])
		questions[i].Qname = name
		*offset += len(name) + 2

		questions[i].Qtype = binary.BigEndian.Uint16(buffer[*offset : *offset+2])
		*offset += 2

		questions[i].Qclass = binary.BigEndian.Uint16(buffer[*offset : *offset+2])
		*offset += 2
	}
	return questions
}

func newResourceRecords(buffer []byte, count uint16, offset *int) []ResourceRecord {
	rrs := make([]ResourceRecord, count)
	for i, _ := range rrs {
		name := parseLabelWithPtr(buffer)
		rrs[i].Rrname = name
		*offset += len(name) + 2

		rrs[i].Rrtype = binary.BigEndian.Uint16(buffer[*offset : *offset+2])
		*offset += 2

		rrs[i].Rrclass = binary.BigEndian.Uint16(buffer[*offset : *offset+2])
		*offset += 2

		rrs[i].Rrttl = binary.BigEndian.Uint32(buffer[*offset : *offset+4])
		*offset += 4

		rrs[i].Rrdlength = binary.BigEndian.Uint16(buffer[*offset : *offset+2])
		*offset += 2

		//rrs[i].Rrdata = buffer[*offset : *offset+int(rrs[i].Rrdlength)-2]
	}
	return rrs
}

func checkErr(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func parseLabel(label []byte) string {
	var buffer bytes.Buffer
	i, j := 0, 0
	for ; i < len(label) && label[i] != 0x00; i++ {
		if i-j == 0 {
			j += int(label[i]) + 1
			if i > 0 {
				buffer.WriteString(".")
			}
			continue
		}

		buffer.WriteString(string(label[i]))
	}
	return buffer.String()
}

func parseLabelWithPtr(label []byte) string {
	if len(label) == 2 && binary.BigEndian.Uint16(label[0:2])&0xC000 != 0 {
		panic("Need to find docs on this.")
	} else {
		return parseLabel(label)
	}
}
