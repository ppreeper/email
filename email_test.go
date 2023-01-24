package email

import (
	"net/mail"
	"testing"
)

func TestServerName(t *testing.T) {
	host := "mail.example.com"
	port := "25"
	s := SMTPServer{Host: host, Port: port, STARTTLS: false}
	if s.ServerName() != host+":"+port {
		t.Errorf("\nSMTPServer.ServerName() failed")
	}
}

// func TestSend(t *testing.T) {

// }

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
			{Name: "", Address: "test.user@example.com"},
		},
		"test.user@example.com",
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
		"Test User <test.user@example.com>,Test User2 <test.user2@example.com>,noname@example.com",
	},
}

func TestAddressList(t *testing.T) {
	for _, mt := range addressToStringTests {
		addressString := addressesToString(mt.addresses)
		if addressString != mt.expected {
			t.Errorf("\naddressesToString expected: %s got: %s", mt.expected, addressString)
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
	{
		Message{
			From:    mail.Address{Name: "Test From", Address: "testfrom@example.com"},
			To:      []mail.Address{{Name: "", Address: "testto@example.com"}},
			Subject: "Test Message",
			Body:    "Hello, World!",
		},
		"Test User <test.user@example.com>",
	},
	{
		Message{
			From:    mail.Address{Name: "Test From", Address: "testfrom@example.com"},
			To:      []mail.Address{{Name: "Test To", Address: "testto@example.com"}},
			Cc:      []mail.Address{{Name: "Test Cc", Address: "testcc@example.com"}},
			Subject: "Test Message CC",
			Body:    "Hello, World!",
		},
		"Test User <test.user@example.com>",
	},
	{
		Message{
			From:    mail.Address{Name: "Test From", Address: "testfrom@example.com"},
			To:      []mail.Address{{Name: "Test To", Address: "testto@example.com"}},
			Cc:      []mail.Address{{Name: "Test Cc", Address: "testcc@example.com"}},
			Bcc:     []mail.Address{{Name: "Test Bcc", Address: "testbcc@example.com"}},
			Subject: "Test Message BCC",
			Body:    "Hello, World!",
		},
		"Test User <test.user@example.com>",
	},
	{
		Message{
			From:    mail.Address{Name: "Test From", Address: "testfrom@example.com"},
			To:      []mail.Address{{Name: "Test To", Address: "testto@example.com"}},
			Cc:      []mail.Address{{Name: "Test Cc", Address: "testcc@example.com"}},
			Bcc:     []mail.Address{{Name: "Test Bcc", Address: "testbcc@example.com"}},
			Subject: "Test Message Attachment",
			Body:    "Hello, World!",
		},
		"Test User <test.user@example.com>",
	},
}

// func TestGenHeader(t *testing.T) {
// 	for _, mt := range messageTests {
// 		r := mt.message.GenHeader()
// 		fmt.Println(string(r))
// fmt.Println(string{mt.message.Headers()})
// addressString := addressesToString(mt.addresses)
// if addressString != mt.expected {
// 	t.Errorf("\addressesToString expected: %s got: %s", mt.expected, addressString)
// }
// 	}
// }

// func TestGenBody(t *testing.T) {
// 	for _, mt := range messageTests {
// 		r := mt.message.GenBody()
// 		fmt.Println(string(r))
// fmt.Println(string{mt.message.Headers()})
// addressString := addressesToString(mt.addresses)
// if addressString != mt.expected {
// 	t.Errorf("\addressesToString expected: %s got: %s", mt.expected, addressString)
// }
// 	}
// }

// func TestGenMixedMessage(t *testing.T) {
// 	for _, mt := range messageTests {
// 		fmt.Println(mt)
// 		r := mt.message.GenMixedMessage("boundary")
// 		fmt.Println(string(r))
// fmt.Println(string{mt.message.Headers()})
// addressString := addressesToString(mt.addresses)
// if addressString != mt.expected {
// 	t.Errorf("\addressesToString expected: %s got: %s", mt.expected, addressString)
// }
// 	}
// }

// func TestGenAttachment(t *testing.T) {
// 	for _, mt := range messageTests {
// 		r := mt.message.GenAttachment("boundary")
// 		fmt.Println(string(r))
// fmt.Println(string{mt.message.Headers()})
// addressString := addressesToString(mt.addresses)
// if addressString != mt.expected {
// 	t.Errorf("\addressesToString expected: %s got: %s", mt.expected, addressString)
// }
// 	}
// }

// func TestBuildMessage(t *testing.T) {
// 	for _, mt := range messageTests {
// 		r := mt.message.BuildMessage()
// 		fmt.Println(string(r))
// fmt.Println(string{mt.message.Headers()})
// addressString := addressesToString(mt.addresses)
// if addressString != mt.expected {
// 	t.Errorf("\addressesToString expected: %s got: %s", mt.expected, addressString)
// }
// 	}
// }

func BenchmarkBuildMessageBase(b *testing.B) {
	var m = Message{
		From:    mail.Address{Name: "Test From", Address: "testfrom@example.com"},
		To:      []mail.Address{{Name: "Test To", Address: "testto@example.com"}},
		Cc:      []mail.Address{{Name: "Test Cc", Address: "testcc@example.com"}},
		Bcc:     []mail.Address{{Name: "Test Bcc", Address: "testbcc@example.com"}},
		Subject: "Test Message Attachment",
		Body:    "Hello, World!",
	}

	for n := 0; n < b.N; n++ {
		m.BuildMessage()
	}
}
func BenchmarkBuildMessageCC(b *testing.B) {
	var m = Message{
		From:    mail.Address{Name: "Test From", Address: "testfrom@example.com"},
		To:      []mail.Address{{Name: "Test To", Address: "testto@example.com"}},
		Cc:      []mail.Address{{Name: "Test Cc", Address: "testcc@example.com"}},
		Subject: "Test Message Attachment",
		Body:    "Hello, World!",
	}

	for n := 0; n < b.N; n++ {
		m.BuildMessage()
	}
}

func BenchmarkBuildMessageBCC(b *testing.B) {
	var m = Message{
		From:    mail.Address{Name: "Test From", Address: "testfrom@example.com"},
		To:      []mail.Address{{Name: "Test To", Address: "testto@example.com"}},
		Cc:      []mail.Address{{Name: "Test Cc", Address: "testcc@example.com"}},
		Bcc:     []mail.Address{{Name: "Test Bcc", Address: "testbcc@example.com"}},
		Subject: "Test Message Attachment",
		Body:    "Hello, World!",
	}

	for n := 0; n < b.N; n++ {
		m.BuildMessage()
	}
}
func BenchmarkBuildMessageAttachment(b *testing.B) {
	var a = Attachment{
		Filename: "test.txt",
		Data: []byte("Sunt reprehenderit ad Lorem sunt ea velit qui consequat mollit ut duis dolore ullamco adipisicing. Velit laboris laborum culpa voluptate Lorem laboris esse mollit cupidatat et. Irure excepteur dolor mollit consectetur ut reprehenderit." +
			"\nQui eu excepteur non ea eu velit. Ut elit ea ullamco nostrud laborum ullamco. Non dolor esse dolor est enim do labore officia. Dolor aliquip id nostrud exercitation. Ut do culpa ad irure fugiat id incididunt nulla quis dolor velit qui mollit dolor." +
			"\nNon sint reprehenderit irure enim esse ad enim est mollit exercitation eu veniam sit nostrud. Aliquip ad incididunt anim consectetur veniam laborum nisi minim irure Lorem pariatur. Ut sit ut do deserunt cupidatat pariatur ad quis veniam ullamco laboris non quis incididunt." +
			"\nSunt ut qui voluptate nulla. Excepteur reprehenderit aliqua eiusmod exercitation magna consectetur excepteur enim eu officia pariatur eiusmod anim. Ad tempor aliquip aliqua labore excepteur eu tempor aliqua. Cupidatat nostrud ullamco qui amet ut commodo occaecat velit ullamco in in esse proident. Eiusmod culpa dolore sit est incididunt quis in commodo commodo. Tempor sit elit laboris ullamco laboris labore ex. Mollit duis amet laboris in laboris Lorem voluptate minim laborum aute dolor est non." +
			"\nEx elit quis id id ut elit commodo elit elit amet veniam magna nostrud aute. Et aliqua dolor id adipisicing minim et exercitation anim nisi consequat. Reprehenderit deserunt et aliqua veniam nulla in magna. Non proident veniam incididunt anim reprehenderit enim adipisicing mollit."),
	}
	var m = Message{
		From:        mail.Address{Name: "Test From", Address: "testfrom@example.com"},
		To:          []mail.Address{{Name: "Test To", Address: "testto@example.com"}},
		Cc:          []mail.Address{{Name: "Test Cc", Address: "testcc@example.com"}},
		Bcc:         []mail.Address{{Name: "Test Bcc", Address: "testbcc@example.com"}},
		Subject:     "Test Message Attachment",
		Body:        "Hello, World!",
		Attachments: map[string]*Attachment{"test.txt": &a},
	}

	for n := 0; n < b.N; n++ {
		m.BuildMessage()
	}
}
