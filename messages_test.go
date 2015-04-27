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
	0x08, 0x6D, 0x79, 0x64, 0x6F, 0x6D, 0x61, 0x69, 0x6E, 0x03, 0x63, 0x6F, 0x6D, 0x00, // mydomain.com
	0x00, 0x05,
	0x00, 0x01,
	0x00, 0x00, 0x00, 0x00,
	0x00, 0x02,
	0x01, 0x01,
}

func TestParseLabel(t *testing.T) {
	input := []byte{0x08, 0x6D, 0x79, 0x64, 0x6F, 0x6D, 0x61, 0x69, 0x6E, 0x03, 0x63, 0x6F, 0x6D, 0x00}

	output := parseLabel(input)
	if output != "mydomain.com" {
		t.Error("Expected mydomain.com, got %s", output)
	}
}

func TestParseLabelComplex(t *testing.T) {
	input := []byte{0x03, 0x77, 0x77, 0x77, 0x08, 0x6D, 0x79, 0x64, 0x6F, 0x6D, 0x61, 0x69, 0x6E, 0x03, 0x63, 0x6F, 0x6D, 0x00}

	output := parseLabel(input)
	if output != "www.mydomain.com" {
		t.Error("Expected www.mydomain.com, got %s", output)
	}
}

func TestNewHeader(t *testing.T) {
	offset := 0
	header := newHeader(message, &offset)

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
	offset := 12
	questions := newQuestions(message, 1, &offset)

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

func TestParseLabelWithPtrPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Expected to panic since the operation is not supported.")
		}
	}()
	parseLabelWithPtr([]byte{0xFF, 0xFF})
}

func TestParseLabelWithPtr(t *testing.T) {
	input := []byte{0x08, 0x6D, 0x79, 0x64, 0x6F, 0x6D, 0x61, 0x69, 0x6E, 0x03, 0x63, 0x6F, 0x6D, 0x00}
	output := parseLabelWithPtr(input)

	if output != "mydomain.com" {
		t.Error("Expected mydomain.com, got ", output)
	}
}

func TestParseLabelWithPtrWithEmptyArray(t *testing.T) {
	output := parseLabelWithPtr([]byte{})
	if output != "" {
		t.Error("Expected an empty string, got ", output)
	}
}

func TestNewResourceRecord(t *testing.T) {
	offset := 30
	rrs := newResourceRecords(message, 1, &offset)

	if rrs[0].Rrname != "mydomain.com" {
		t.Error("Expected mydomain.com, got ", rrs[0].Rrname)
	}

	if rrs[0].Rrtype != 5 {
		t.Error("Expected 5, got ", rrs[0].Rrtype)
	}

	if rrs[0].Rrclass != 1 {
		t.Error("Expected 1, got ", rrs[0].Rrclass)
	}

	if rrs[0].Rrttl != 0 {
		t.Error("Expected 0, got ", rrs[0].Rrttl)
	}

	if rrs[0].Rrdlength != 2 {
		t.Error("Expected 2, got ", rrs[0].Rrdlength)
	}
}
