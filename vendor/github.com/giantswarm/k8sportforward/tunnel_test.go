package k8sportforward

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/fortytw2/leaktest"
	"k8s.io/apimachinery/pkg/util/httpstream"
	"k8s.io/apimachinery/pkg/util/httpstream/spdy"
)

type fakeDialer struct {
	dialed             bool
	conn               httpstream.Connection
	err                error
	negotiatedProtocol string
}

func (d *fakeDialer) Dial(protocols ...string) (httpstream.Connection, string, error) {
	d.dialed = true
	return d.conn, d.negotiatedProtocol, d.err
}

func TestGoroutineLeak(t *testing.T) {
	defer leaktest.Check(t)()

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
	}))
	defer ts.Close()
	parsedURL, err := url.Parse(ts.URL)
	if err != nil {
		t.Fatalf("could not parse test server URL %v", err)
	}

	conn, err := net.Dial("tcp", parsedURL.Host)
	if err != nil {
		t.Fatalf("could not create connection to server %v", err)
	}

	streamConn, err := spdy.NewClientConnection(conn)
	if err != nil {
		t.Fatalf("could not create connection to server %v", err)
	}
	dialer := &fakeDialer{
		conn: streamConn,
	}
	config := tunnelConfig{
		Dialer: dialer,

		RemotePort: 80,
	}
	_, err = newTunnel(config)
	if err != nil {
		t.Fatalf("could not create tunnel %v", err)
	}
}
