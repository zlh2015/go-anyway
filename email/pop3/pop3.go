package pop3

import (
	"bufio"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/textproto"
	"strconv"
	"strings"
)

// Client for POP3.
type Client struct {
	stls bool
	text *textproto.Conn
	conn net.Conn
	bin  *bufio.Reader
}

// Dial creates an unsecured connection to the POP3 server at the given address
// and returns the corresponding Client.
func Dial(addr string) (*Client, error) {
	// conn, err := net.Dial("tcp", "pop-mail.outlook.com:995")
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}
	return NewClient(conn)
}

// DialTLS creates a TLS-secured connection to the POP3 server at the given
// address and returns the corresponding Client.
func DialTLS(addr string) (*Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		return nil, err
	}
	return NewClient(conn)
}

// NewClient returns a new Client object using an existing connection.
func NewClient(conn net.Conn) (*Client, error) {
	client := &Client{
		stls: false,
		text: nil,
		bin:  bufio.NewReader(conn),
		conn: conn,
	}
	// send dud command, to read a line
	resp, err := client.Cmd("")
	if err != nil {
		return nil, err
	}
	fmt.Println(resp)
	return client, nil
}

// STLS sends the STLS command and encrypts all further communication.
func (c *Client) STLS(config *tls.Config) (error) {
	// if err := c.hello(); err != nil {
	// return err
	// }
	_, err := c.Cmd("STLS\r\n", nil)
	if err != nil {
		return  err
	}
	c.conn = tls.Client(c.conn, config)
	c.text = textproto.NewConn(c.conn)
	c.stls = true
	return nil 
}

// STLSCmd is a convenience function that sends a command and returns the response
func (c *Client) STLSCmd(expectCode int, format string, args ...interface{}) (int, string, error) {
	id, err := c.text.Cmd(format, args...)
	if err != nil {
		return 0, "", err
	}
	c.text.StartResponse(id)
	defer c.text.EndResponse(id)
	code, msg, err := c.text.ReadResponse(expectCode)
	return code, msg, err
}

// Cmd sent command and receive the first line. the left reponse lines must be retrieved via readLines.
func (c *Client) Cmd(format string, args ...interface{}) (string, error) {
	fmt.Println(format, args)
	fmt.Fprintf(c.conn, format, args...)
	line, _, err := c.bin.ReadLine()
	if err != nil {
		return "", err
	}
	l := string(line)
	if l[0:3] != "+OK" {
		err = errors.New(l[5:])
	}
	if len(l) >= 4 {
		return l[4:], err
	}
	return "", err
}

// ReadLines get all lines from the response io
func (c *Client) ReadLines() (lines []string, err error) {
	lines = make([]string, 0)
	l, _, err := c.bin.ReadLine()
	line := string(l)
	for err == nil && line != "." {
		if len(line) > 0 && line[0] == '.' {
			line = line[1:]
		}
		lines = append(lines, line)
		l, _, err = c.bin.ReadLine()
		line = string(l)
	}
	return
}

// User sends the given username to the server. Generally, there is no reason
// not to use the Auth convenience method.
func (c *Client) User(username string) (err error) {
	_, err = c.Cmd("USER %s\r\n", username)
	return
}

// Pass sends the given password to the server. The password is sent
// unencrypted unless the connection is already secured by TLS (via DialTLS or
// some other mechanism). Generally, there is no reason not to use the Auth
// convenience method.
func (c *Client) Pass(password string) (err error) {
	_, err = c.Cmd("PASS %s\r\n", password)
	return
}

// Auth sends the given username and password to the server, calling the User
// and Pass methods as appropriate.
func (c *Client) Auth(username, password string) (err error) {
	err = c.User(username)
	if err != nil {
		return
	}
	err = c.Pass(password)
	return
}

// Stat retrieves a drop listing for the current maildrop, consisting of the
// number of messages and the total size (in octets) of the maildrop.
// Information provided besides the number of messages and the size of the
// maildrop is ignored. In the event of an error, all returned numeric values
// will be 0.
func (c *Client) Stat() (count, size int, err error) {
	l, err := c.Cmd("STAT\r\n")
	if err != nil {
		return 0, 0, err
	}
	parts := strings.Fields(l)
	count, err = strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, errors.New("Invalid server response")
	}
	size, err = strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, errors.New("Invalid server response")
	}
	return
}

// List returns the size of the given message, if it exists. If the message
// does not exist, or another error is encountered, the returned size will be
// 0.
func (c *Client) List(msg int) (size int, err error) {
	l, err := c.Cmd("LIST %d\r\n", msg)
	if err != nil {
		return 0, err
	}
	size, err = strconv.Atoi(strings.Fields(l)[1])
	if err != nil {
		return 0, errors.New("Invalid server response")
	}
	return size, nil
}

// ListAll returns a list of all messages and their sizes.
func (c *Client) ListAll() (msgs []int, sizes []int, err error) {
	_, err = c.Cmd("LIST\r\n")
	if err != nil {
		return
	}
	lines, err := c.ReadLines()
	if err != nil {
		return
	}
	msgs = make([]int, len(lines), len(lines))
	sizes = make([]int, len(lines), len(lines))
	for i, l := range lines {
		var m, s int
		fs := strings.Fields(l)
		m, err = strconv.Atoi(fs[0])
		if err != nil {
			return
		}
		s, err = strconv.Atoi(fs[1])
		if err != nil {
			return
		}
		msgs[i] = m
		sizes[i] = s
	}
	return
}
