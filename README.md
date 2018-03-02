# email
Library to send email messages, written in go. This library is using the standard library to make up all the moving parts in sending email messages. No external dependencies required.

## Use
To use the library you need to construct three parts:
* **SMTPServer** The smtp server the message is being submitted to
* **User** The authorized email user sending the message
* **Message** The message that is to be constructed and sent

```go
package main

import (
	"log"
	"net/mail"

	"github.com/ppreeper/email"
)

func main() {

	d := email.Attachment{
		Filename: "test.txt",
		Data:     []byte("oh yeah"),
		Inline:   false,
	}

	//create message
	m := email.Message{
		From:     mail.Address{"Mailer", "mailer@example.com"},
		To:       []mail.Address{{"Standard User", "standard.user@example.com"}},
		Subject:  "test subject",
		Body:     "this is the email body",
		MimeType: "text/html",
	}
	m.Attachments = make(map[string]*email.Attachment)
	m.Attachments["test.txt"] = &d
  
	// Connect to the remote SMTP server.
	user := email.User{Username: "mailer@example.com", Password: "pa55w0rd"}
	server := email.SMTPServer{Host: "mail.example.com", Port: "25", STARTTLS: false}

	err := server.Send(&user, &m)
	check(err)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
```

## Status
This library is a work in progress but has been tested against a MS Exchange server and a Postfix server configured with starttls.

## Issues
* bcc not functioning as expected
* testing on exchange did not send messages outside of organization, but postfix worked sending to internal/external recipients
