package sonic

import (
	"bufio"
	"errors"
	"log"
	"net"
	"strings"
	"sync"
)

// Connection represents a connection to the Sonic server
type Connection struct {
	Password string
	async    map[string]chan *Message
	tcp      net.Conn
	buf      int
	protocol int
	rcv      chan *Message
	mux      sync.Mutex
}

// Connect establishes a connection to the Sonic server
func Connect(addr, pass string) (conn *Connection, err error) {
	tcp, err := net.Dial("tcp", addr)
	if err != nil {
		return
	}

	conn = &Connection{
		Password: pass,
		tcp:      tcp,
		rcv:      make(chan *Message),
		async:    make(map[string]chan *Message),
		mux:      sync.Mutex{},
	}
	go conn.listen()

	msg := <-conn.rcv
	if msg.Name != "CONNECTED" {
		err = ErrUnexpectedResponse
	}

	return
}

// Send sends a message and returns the response
func (c *Connection) Send(m *Message) (res *Message, err error) {
	c.mux.Lock()
	defer c.mux.Unlock()

	for _, msg := range m.Split(c.buf) {
		b := []byte(msg.String() + "\n")
		_, err = c.tcp.Write(b)
		if err != nil {
			return
		}

		res = <-c.rcv
		if res.Name == "ERR" {
			err = errors.New(strings.Join(res.Args, " "))
		}
	}

	return
}

// SendAsync sends an async message to Sonic
func (c *Connection) SendAsync(m *Message) (res *Message, err error) {
	res, err = c.Send(m)
	if err != nil {
		return
	}

	if res.Name != "PENDING" {
		err = ErrUnexpectedResponse
		return
	}

	id := res.Args[0]
	ch, del := c.ensureAsyncChan(id)
	defer del()

	res = <-ch
	return
}

// Search establishes a connection in search mode
func (c *Connection) Search() (s *Search, err error) {
	err = c.handshake("search")
	if err != nil {
		return
	}

	return &Search{Conn: c}, nil
}

// Ingest establishes the connection in ingest mode
func (c *Connection) Ingest() (i *Ingest, err error) {
	err = c.handshake("ingest")
	if err != nil {
		return
	}

	return &Ingest{Conn: c}, nil
}

func (c *Connection) handshake(mode string) error {
	msg := &Message{Name: "START", Args: []string{mode, c.Password}}
	msg, err := c.Send(msg)
	if err != nil {
		return err
	}

	if msg.Name != "STARTED" {
		return ErrUnexpectedResponse
	}

	_, c.protocol = parseIntArg(msg.Args[1])
	_, c.buf = parseIntArg(msg.Args[2])
	return nil
}

func (c *Connection) listen() error {
	log.Println("listening")
	rd := bufio.NewReader(c.tcp)
	for {
		pk, err := rd.ReadString('\n')
		if err != nil {
			return err
		}

		fields := strings.Split(strings.TrimSpace(pk), " ")
		msg := &Message{
			Name: fields[0],
			Args: fields[1:],
		}

		log.Printf("rcv: %s\n", msg)
		if msg.Name == "EVENT" {
			ch, del := c.ensureAsyncChan(msg.Args[1])
			defer del()
			ch <- msg
		} else {
			c.rcv <- msg
		}
	}
}

func (c *Connection) ensureAsyncChan(id string) (chan *Message, func()) {
	ch := c.async[id]
	if ch == nil {
		ch = make(chan *Message, 1)
		c.async[id] = ch
		return ch, func() {
			delete(c.async, id)
		}
	}
	return ch, func() {}
}
