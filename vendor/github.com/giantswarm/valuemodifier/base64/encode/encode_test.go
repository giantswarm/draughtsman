package encode

import (
	"testing"
)

func Test_Base64_Encode_Service_Modify(t *testing.T) {
	config := DefaultConfig()
	newService, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expected := []byte("aGVsbG8gd29ybGQ=") // "hello world"
	value := []byte("hello world")
	modified, err := newService.Modify(value)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if string(modified) != string(expected) {
		t.Fatal("expected", string(expected), "got", string(modified))
	}
}
