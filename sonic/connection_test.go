package sonic

import "testing"

func TestConnection(t *testing.T) {
	c, err := Connect("localhost:1491", "SecretPassword")
	if err != nil {
		t.Fatal(err)
	}

	i, err := c.Ingest()
	if err != nil {
		t.Fatal(err)
	}

	c, err = Connect("localhost:1491", "SecretPassword")
	if err != nil {
		t.Fatal(err)
	}

	s, err := c.Search()
	if err != nil {
		t.Fatal(err)
	}

	err = i.Push("test", "bucket_a", "key", "hello world", "")
	if err != nil {
		t.Fatal(err)
	}

	res, err := s.Query("bucket_a", "test", "hello", QueryOptions{})
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatalf("Expected results: instead got %s", res)
	}

	t.Log(res)
}
