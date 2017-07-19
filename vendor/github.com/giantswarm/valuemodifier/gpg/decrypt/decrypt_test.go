package decrypt

import (
	"testing"
)

func Test_GPG_Decrypt_Service_Modify(t *testing.T) {
	config := DefaultConfig()
	config.Pass = "foo"
	newService, err := New(config)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	expected := []byte("hello world")
	value := []byte(`-----BEGIN PGP SIGNATURE-----

wx4EBwMIbxESvPWYOgBgEmcsCe70T3fWYMXUAO/SBZHS4AHk6qQ8xikHnoCiBasb
8HKnT+FFOOCN4P3hujrg8+I+yH/84GfjXeQwavZMLtTgAeGQ8eCB4GXgguSqXR+p
l5T01NrWN5NQZM+H4ngibenh8GwA
=C3bS
-----END PGP SIGNATURE-----`) // "hello world"
	modified, err := newService.Modify(value)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}

	if string(modified) != string(expected) {
		t.Fatal("expected", expected, "got", modified)
	}
}
