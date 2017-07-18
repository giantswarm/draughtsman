package encrypt

import (
	"strings"
	"testing"
)

func Test_GPG_Encrypt_Service_Modify(t *testing.T) {
	config := DefaultConfig()
	config.Pass = "foo"
	newService, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	// Note that we cannot predict the outcome of the encryption due to the fact
	// how GPG works. Below we can only assume to have a proper GPG message by
	// comparing its length and verify the prefix and suffix of the GPG message is
	// as expected.
	expected := []byte(`-----BEGIN PGP SIGNATURE-----

xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
xxxxxxxxxxxxxxxxxxxxxxxxxxxx
xxxxx
-----END PGP SIGNATURE-----`) // "hello world"
	value := []byte("hello world")
	modified, err := newService.Modify(value)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	le := len(string(expected))
	lm := len(string(modified))
	if lm != le {
		t.Fatal("expected", le, "got", lm)
	}
	if !strings.HasPrefix(string(modified), "-----BEGIN PGP SIGNATURE-----\n\n") {
		t.Fatal("expected", true, "got", false)
	}
	if !strings.HasSuffix(string(modified), "\n-----END PGP SIGNATURE-----") {
		t.Fatal("expected", true, "got", false)
	}
}
