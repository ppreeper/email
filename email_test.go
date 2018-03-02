package email_test

import (
	"testing"

	"github.com/ppreeper/email"
)

func BuildMessage(t *testing.T) {

}

func TestAttach(t *testing.T) {

}

func TestServerName(t *testing.T) {
	host := "test.email.server"
	port := "25"
	s := email.SMTPServer{Host: host, Port: port, STARTTLS: false}
	if s.ServerName() != host+":"+port {
		t.Errorf("\nSMTPServer.ServerName() failed")
	}
}

func TestSend(t *testing.T) {

}
