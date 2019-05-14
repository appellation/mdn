package sonic

// Search represents a search connection
type Search struct {
	Conn      *Connection
	msgs      chan Message
	listening bool
}

// QueryOptions represents options to send with the query
type QueryOptions struct {
	Limit, Offset int
	Lang          string
}

// Query sends a search query
func (s *Search) Query(bucket, collection, text string, opts QueryOptions) (results []string, err error) {
	msg := &Message{
		Name: "QUERY",
		Args: []string{bucket, collection},
		Text: text,
		Options: []MessageOption{
			{"LIMIT", opts.Limit},
			{"OFFSET", opts.Offset},
			{"LANG", opts.Lang},
		},
	}
	msg, err = s.Conn.SendAsync(msg)
	if err != nil {
		return
	}

	if len(msg.Args) < 2 {
		results = []string{}
	} else {
		results = msg.Args[2:]
	}

	return
}
