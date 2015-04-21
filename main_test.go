package bdns

import "testing"

var message = []byte{
	// header
	0xFF, 0x00,
	0x00, 0x00,
	0x00, 0x01,
	0x00, 0x00,
	0x00, 0x00,
	0x00, 0x00,

	// questions
	0x08, 0x6D, 0x79, 0x64, 0x6F, 0x6D, 0x61, 0x69, 0x6E, 0x03, 0x63, 0x6F, 0x6D, 0x00, // mydomain.com
	0x00, 0x01,
	0x00, 0x01,

	// resource records
}

func TestParseLabel(t *testing.T) {
	var input = []byte{0x08, 0x6D, 0x79, 0x64, 0x6F, 0x6D, 0x61, 0x69, 0x6E, 0x03, 0x63, 0x6F, 0x6D, 0x00}

	output := parseLabel(input)
	if output != "mydomain.com" {
		t.Error("Expected mydomain.com, got %s", output)
	}
}

func TestParseLabelComplex(t *testing.T) {
	var input = []byte{0x03, 0x77, 0x77, 0x77, 0x08, 0x6D, 0x79, 0x64, 0x6F, 0x6D, 0x61, 0x69, 0x6E, 0x03, 0x63, 0x6F, 0x6D, 0x00}

	output := parseLabel(input)
	if output != "www.mydomain.com" {
		t.Error("Expected www.mydomain.com, got %s", output)
	}
}

func TestNewHeader(t *testing.T) {
	header := newHeader(message)

	if header.Id != 65280 {
		t.Error("Expected Id == 65280, got ", header.Id)
	}

	if header.Flags != 0 {
		t.Error("Expected Flags == 0, got ", header.Flags)
	}

	if header.Qdcount != 1 {
		t.Error("Expected Qdcount == 1, got ", header.Qdcount)
	}

	if header.Ancount != 0 {
		t.Error("Expected Ancount == 0, got ", header.Ancount)
	}

	if header.Nscount != 0 {
		t.Error("Expected Nscount == 0, got ", header.Nscount)
	}

	if header.Arcount != 0 {
		t.Error("Expected Arcount == 0, got ", header.Arcount)
	}
}

func TestNewQuestion(t *testing.T) {
	questions := newQuestions(message, 1)

	if nbQuestions := len(questions); nbQuestions != 1 {
		t.Error("Expected 1 question, got", nbQuestions)
	}

	if questions[0].Qname != "mydomain.com" {
		t.Error("Expected mydomain.com, got ", questions[0].Qname)
	}

	if questions[0].Qtype != 1 {
		t.Error("Expected Qtype == 1, got ", questions[0].Qtype)
	}

	if questions[0].Qclass != 1 {
		t.Error("Expected Qclass == 1, got ", questions[0].Qclass)
	}
}

func TestNewMessage(t *testing.T) {
	msg := NewMessage(message)

	if msg.Header.Qdcount != 1 {
		t.Error("Expected Qdcount == 1, got ", msg.Header.Qdcount)
	}

	if msg.Questions[0].Qname != "mydomain.com" {
		t.Error("Expected Qname == mydomain.com, got", msg.Questions[0].Qname)
	}
}
