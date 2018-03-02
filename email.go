package email

import (
	"bytes"
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"math/rand"
	"mime"
	"net/mail"
	"net/smtp"
	"path/filepath"
	"strconv"
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
	if s.STARTTLS {
		auth := smtp.PlainAuth("", u.Username, u.Password, s.Host)
		// config
		// tlsconfig := &tls.Config{
		// 	InsecureSkipVerify: true,
		// 	ServerName:         host,
		// }

		// c.StartTLS(tlsconfig)

		// Auth
		// err = c.Auth(auth)
		smtp.SendMail(s.ServerName(), auth, m.From.Address, to, m.BuildMessage())
		smtp.SendMail(s.ServerName(), auth, m.From.Address, cc, m.BuildMessage())
		smtp.SendMail(s.ServerName(), auth, m.From.Address, bcc, m.BuildMessage())
	} else {
		smtp.SendMail(s.ServerName(), nil, m.From.Address, to, m.BuildMessage())
		smtp.SendMail(s.ServerName(), nil, m.From.Address, cc, m.BuildMessage())
		smtp.SendMail(s.ServerName(), nil, m.From.Address, bcc, m.BuildMessage())
	}

	return nil
}
