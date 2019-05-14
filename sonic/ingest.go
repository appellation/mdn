package sonic

import (
	"strconv"
)

// Ingest represents a Sonic connection in ingest mode
type Ingest struct {
	Conn *Connection
}

// Push pushes data to Sonic
func (i *Ingest) Push(collection, bucket, object, text, locale string) (err error) {
	// fmt.Println(collection, bucket, object, locale)
	m := &Message{
		Name: "PUSH",
		Args: []string{collection, bucket, object},
		Text: text,
		Options: []MessageOption{
			{"LANG", locale},
		},
	}

	m, err = i.Conn.Send(m)
	if err != nil {
		return
	}

	if m.Name != "OK" {
		err = ErrUnexpectedResponse
	}
	return
}

// Count counts indexed search data; bucket and object are optional
func (i *Ingest) Count(collection, bucket, object string) (c int, err error) {
	// object cannot be specified if bucket is empty
	if bucket == "" && object != "" {
		err = ErrInvalidOptions
		return
	}

	m, err := i.Conn.Send(&Message{
		Name: "COUNT",
		Args: []string{collection, bucket, object},
	})

	if err != nil {
		return
	}

	if m.Name != "RESULT" {
		err = ErrUnexpectedResponse
		return
	}

	res, err := strconv.ParseInt(m.Args[0], 10, 0)
	c = int(res)
	return
}

// FlushCollection flushes the collection
func (i *Ingest) FlushCollection(collection string) (c int, err error) {
	m, err := i.Conn.Send(&Message{
		Name: "FLUSHC",
		Args: []string{collection},
	})
	if err != nil {
		return
	}

	if m.Name != "RESULT" {
		err = ErrUnexpectedResponse
		return
	}

	res, err := strconv.ParseInt(m.Args[0], 10, 0)
	c = int(res)
	return
}
