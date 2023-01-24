package email

import (
	"fmt"
	"net/mail"
	"testing"
)

func BuildMessage(t *testing.T) {

}

func TestAttach(t *testing.T) {

}

func TestServerName(t *testing.T) {
	host := "test.email.server"
	port := "25"
	s := SMTPServer{Host: host, Port: port, STARTTLS: false}
	if s.ServerName() != host+":"+port {
		t.Errorf("\nSMTPServer.ServerName() failed")
	}
}

func TestSend(t *testing.T) {

}

var addressToStringTests = []struct {
	addresses []mail.Address
	expected  string
}{
	{
		[]mail.Address{
			{Name: "Test User", Address: "test.user@example.com"},
		},
		"Test User <test.user@example.com>",
	},
	{
		[]mail.Address{
			{Name: "TestUser", Address: "testuser@example.com"},
		},
		"TestUser <testuser@example.com>",
	},
	{
		[]mail.Address{
			{Name: "Test User", Address: "test.user@example.com"},
			{Name: "Test User2", Address: "test.user2@example.com"},
		},
		"Test User <test.user@example.com>,Test User2 <test.user2@example.com>",
	},
	{
		[]mail.Address{
			{Name: "Test User", Address: "test.user@example.com"},
			{Name: "Test User2", Address: "test.user2@example.com"},
			{Name: "", Address: "noname@example.com"},
		},
		"Test User <test.user@example.com>,Test User2 <test.user2@example.com>, <noname@example.com>",
	},
}

func TestAddressList(t *testing.T) {
	for _, mt := range addressToStringTests {
		addressString := addressesToString(mt.addresses)
		if addressString != mt.expected {
			t.Errorf("\addressesToString expected: %s got: %s", mt.expected, addressString)
		}
	}
}

var messageTests = []struct {
	message  Message
	expected string
}{
	{
		Message{
			From:    mail.Address{Name: "Test From", Address: "testfrom@example.com"},
			To:      []mail.Address{{Name: "Test To", Address: "testto@example.com"}},
			Subject: "Test Message",
			Body:    "Hello, World!",
		},
		"Test User <test.user@example.com>",
	},
}

func TestHeaders(t *testing.T) {
	for _, mt := range messageTests {
		fmt.Println(mt)
		// fmt.Println(string{mt.message.Headers()})
		// addressString := addressesToString(mt.addresses)
		// if addressString != mt.expected {
		// 	t.Errorf("\addressesToString expected: %s got: %s", mt.expected, addressString)
		// }
	}

}
