package email

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"mime"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var coder = base64.StdEncoding

// Attachment for email attachment
type Attachment struct {
	Filename string
	Data     []byte
	Inline   bool
}

// Message SMTP Message
type Message struct {
	From        mail.Address
	To          []mail.Address
	Cc          []mail.Address
	Bcc         []mail.Address
	ReplyTo     mail.Address
	Subject     string
	Body        string
	MimeType    string
	Attachments map[string]*Attachment
}

// Attach file to SMTP Message
func (m *Message) Attach(file string, inline bool) error {
	_, filename := filepath.Split(file)
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	m.Attachments[filename] = &Attachment{
		Filename: filename,
		Data:     data,
		Inline:   inline,
	}

	return nil
}

// BuildMessage returns byte ready email file
func (m *Message) BuildMessage() []byte {
	t := time.Now()
	buf := bytes.NewBuffer(nil)

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	datarand := strconv.Itoa(r.Intn(1000000))
	digest := sha1.Sum([]byte(datarand))
	boundary := coder.EncodeToString(digest[:])

	// Setup headers
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", m.From.Name, m.From.Address)
	headers["Date"] = fmt.Sprintf("%s", t.Format(time.RFC1123Z))
	var hdr string
	for i := 0; i < len(m.To); i++ {
		if i == len(m.To)-1 {
			hdr += fmt.Sprintf(m.To[i].Name + " <" + m.To[i].Address + ">")
		} else {
			hdr += fmt.Sprintf(m.To[i].Name + " <" + m.To[i].Address + ">,")
		}
	}
	headers["To"] = hdr

	if len(m.Cc) > 0 {
		hdr = ""
		for i := 0; i < len(m.Cc); i++ {
			if i == len(m.Cc)-1 {
				hdr += fmt.Sprintf(m.Cc[i].Name + " <" + m.Cc[i].Address + ">")
			} else {
				hdr += fmt.Sprintf(m.Cc[i].Name + " <" + m.Cc[i].Address + ">,")
			}
		}
		headers["Cc"] = hdr
	}
	if len(m.Bcc) > 0 {
		hdr = ""
		for i := 0; i < len(m.Bcc); i++ {
			if i == len(m.Bcc)-1 {
				hdr += fmt.Sprintf(m.Bcc[i].Name + " <" + m.Bcc[i].Address + ">")
			} else {
				hdr += fmt.Sprintf(m.Bcc[i].Name + " <" + m.Bcc[i].Address + ">,")
			}
		}
		headers["Bcc"] = hdr
	}
	if m.ReplyTo.Address != "" {
		headers["Reply-To"] = fmt.Sprintf("%s <%s>", m.ReplyTo.Name, m.ReplyTo.Address)
	}
	headers["Subject"] = "=?UTF-8?B?" + coder.EncodeToString([]byte(m.Subject)) + "?="

	// Setup message = headers + body
	for k, v := range headers {
		buf.WriteString(k + ": " + v + "\r\n")
	}
	buf.WriteString("MIME-Version: 1.0\r\n")

	if len(m.Attachments) > 0 {
		buf.WriteString("Content-Type: multipart/mixed; boundary=" + boundary + "\r\n")
		buf.WriteString("\r\n--" + boundary + "\r\n")
	}
	// Body text
	buf.WriteString("Content-Type: " + m.MimeType + "; charset=utf-8\r\n\r\n")
	buf.WriteString(m.Body + "\r\n")

	if len(m.Attachments) > 0 {
		for _, f := range m.Attachments {
			buf.WriteString("\r\n\r\n--" + boundary)
			if f.Inline {
				buf.WriteString("\r\nContent-Type: message/rfc822\r\n")
				buf.WriteString("Content-Disposition: inline; filename=\"" + f.Filename + "\"\r\n\r\n")
				buf.Write(f.Data)
				buf.WriteString("\r\n--" + boundary)
			} else {
				mimetype := mime.TypeByExtension(filepath.Ext(f.Filename))
				if mimetype != "" {
					buf.WriteString("\r\nContent-Type: " + mimetype + "\r\n")
				} else {
					buf.WriteString("\r\nContent-Type: application/octet-stream\r\n")
				}
				buf.WriteString("Content-Transfer-Encoding: base64\r\n")
				buf.WriteString("Content-Disposition: attachment; filename=\"" + "=?UTF-8?B?" + coder.EncodeToString([]byte(f.Filename)) + "?=" + "\"\r\n\r\n")
				b := make([]byte, base64.StdEncoding.EncodedLen(len(f.Data)))
				base64.StdEncoding.Encode(b, f.Data)
				for i, l := 0, len(b); i < l; i++ {
					buf.WriteByte(b[i])
					if (i+1)%76 == 0 {
						buf.WriteString("\r\n")
					}
				}
			}
		}
		buf.WriteString("\r\n--" + boundary)
		buf.WriteString("--")
	}

	return buf.Bytes()
}

// User client
type User struct {
	Username string
	Password string
	Auth     bool
}

// SMTPServer host setup
type SMTPServer struct {
	Host     string
	Port     string
	STARTTLS bool
}

// ServerName return host port combo
func (s *SMTPServer) ServerName() string {
	return string(s.Host + ":" + s.Port)
}

// Send smtp message
func (s *SMTPServer) Send(u *User, m *Message) error {
	if err := validateLine(m.From.String()); err != nil {
		return err
	}
	for _, recp := range m.To {
		if err := validateLine(recp.String()); err != nil {
			return err
		}
	}
	for _, recp := range m.Cc {
		if err := validateLine(recp.String()); err != nil {
			return err
		}
	}
	for _, recp := range m.Bcc {
		if err := validateLine(recp.String()); err != nil {
			return err
		}
	}
	var to, cc, bcc []string
	for _, v := range m.To {
		to = append(to, v.Address)
	}
	for _, v := range m.Cc {
		cc = append(cc, v.Address)
	}
	for _, v := range m.Bcc {
		bcc = append(bcc, v.Address)
	}
	// var auth Auth
	// if u.Auth {
	// 	auth = PlainAuth("", u.Username, u.Password, s.Host)
	// } else {
	// 	auth = nil
	// }
	if len(to) > 0 {
		err := s.SendMail(s.ServerName(), m.From.String(), to, m.BuildMessage())
		if err != nil {
			return err
		}
	}
	if len(cc) > 0 {
		err := s.SendMail(s.ServerName(), m.From.String(), to, m.BuildMessage())
		if err != nil {
			log.Printf("err: %v", err)
		}
	}
	if len(bcc) > 0 {
		err := s.SendMail(s.ServerName(), m.From.Address, bcc, m.BuildMessage())
		if err != nil {
			log.Printf("err: %v", err)
		}
	}

	return nil
}

// SendMail requires no tls
func (s *SMTPServer) SendMail(serverName, from string, to []string, msg []byte) error {
	c, err := smtp.Dial(serverName)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// validateLine checks to see if a line has CR or LF as per RFC 5321
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}
