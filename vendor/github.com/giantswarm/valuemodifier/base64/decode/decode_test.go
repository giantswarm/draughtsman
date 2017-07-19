package decode

import (
	"testing"
)

func Test_Base64_Decode_Service_Modify(t *testing.T) {
	config := DefaultConfig()
	newService, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expected := []byte("hello world")
	value := []byte("aGVsbG8gd29ybGQ=") // "hello world"
	modified, err := newService.Modify(value)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if string(modified) != string(expected) {
		t.Fatal("expected", expected, "got", modified)
	}
}
