package sonic

import (
	"fmt"
	"reflect"
	"strings"
)

// Message represents a message to/from Sonic
type Message struct {
	Name    string
	Args    []string
	Text    string
	Options []MessageOption
}

func (m Message) String() string {
	str := strings.Builder{}
	str.WriteString(m.Name)

	for _, arg := range m.Args {
		if arg != "" {
			str.WriteString(" " + arg)
		}
	}

	if len(m.Text) > 0 {
		str.WriteString(" " + fmt.Sprintf("%q", m.Text))
	}

	for _, opt := range m.Options {
		s := opt.String()
		if s != "" {
			str.WriteString(" " + s)
		}
	}

	return str.String()
}

// Split splits this message into chunks where text <= size
func (m Message) Split(size int) []Message {
	l := len(m.Text)
	if l <= size || size < 1 {
		return []Message{m}
	}

	parts := []string{m.Text[:l/2], m.Text[l/2:]}
	msgs := []Message{}
	for _, part := range parts {
		newMessage := Message{
			Name:    m.Name,
			Args:    m.Args,
			Text:    part,
			Options: m.Options,
		}
		msgs = append(msgs, newMessage.Split(size)...)
	}

	return msgs
}

// MessageOption represents options in the format name(value)
type MessageOption struct {
	Name  string
	Value interface{}
}

func (o MessageOption) String() string {
	v := reflect.ValueOf(o.Value)
	if v.Interface() == reflect.Zero(v.Type()).Interface() {
		return ""
	}

	return fmt.Sprintf("%s(%v)", o.Name, o.Value)
}
