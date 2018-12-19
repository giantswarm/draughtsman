package guest

import (
	"errors"
	"testing"
)

func Test_IsGuestAPINotAvailable(t *testing.T) {
	testCases := []struct {
		description   string
		errorMessage  string
		expectedMatch bool
	}{
		{
			description:   "case 1: dns not ready",
			errorMessage:  "dial tcp: lookup api.5xchu.aws.gigantic.io on 10.96.0.10:53: no such host",
			expectedMatch: true,
		},
		{
			description:   "case 2: dns not ready incorrect port",
			errorMessage:  "dial tcp: lookup api.5xchu.aws.gigantic.io on 10.96.0.10:443: no such host",
			expectedMatch: false,
		},
		{
			description:   "case 3: ingress not ready get request",
			errorMessage:  "Get https://api.5xchu.aws.gigantic.io: x509: certificate is valid for ingress.local, not api.5xchu.aws.gigantic.io:",
			expectedMatch: true,
		},
		{
			description:   "case 4: API not ready get EOF request",
			errorMessage:  "Get https://api.5xchu.aws.gigantic.io/api/v1/nodes: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 5: temporary issues with the master node serving the guest cluster API",
			errorMessage:  "Get https://api.8dnxs.g8s.gorgoth.gridscale.kvm.gigantic.io/api/v1/nodes: unexpected EOF",
			expectedMatch: true,
		},
		{
			description:   "case 6: temporary issues with the master node serving the guest cluster API",
			errorMessage:  "Get https://api.uth29.g8s.geckon.gridscale.kvm.gigantic.io/api/v1/nodes?timeout=30s: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 7: ingress not ready post request",
			errorMessage:  "Post https://api.5xchu.aws.gigantic.io: x509: certificate is valid for ingress.local, not api.5xchu.aws.gigantic.io:",
			expectedMatch: true,
		},
		{
			description:   "case 8: ingress not ready post different domain",
			errorMessage:  "Post https://api.5xchu.aws.gigantic.io: x509: certificate is valid for localhost, not api.5xchu.aws.gigantic.io:",
			expectedMatch: false,
		},
		{
			description:   "case 9: timeout getting namespace",
			errorMessage:  "Get https://api.3jwh2.k8s.aws.gigantic.io/api/v1/namespaces/giantswarm?timeout=30s: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 10: timeout getting service account",
			errorMessage:  "Post https://api.3jwh2.k8s.aws.gigantic.io/api/v1/namespaces/giantswarm/serviceaccounts?timeout=30s: EOF",
			expectedMatch: true,
		},
		{
			description:   "case 11: GET timeout establishing TLS handshake",
			errorMessage:  "Get https://api.08vka.k8s.gorgoth.gridscale.kvm.gigantic.io/api/v1/nodes?timeout=30s: net/http: TLS handshake timeout",
			expectedMatch: true,
		},
		{
			description:   "case 12: server is misbehaving due to TCP lookup",
			errorMessage:  "Get https://api.ci-wip-70f9b-5e958.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s: dial tcp: lookup api.ci-wip-70f9b-5e958.k8s.godsmack.westeurope.azure.gigantic.io on 10.96.0.10:53: server misbehaving",
			expectedMatch: true,
		},
		{
			description:   "case 13: request canceled while waiting for connection",
			errorMessage:  "Get https://api.ci-wip-2317d-c1c86.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s: net/http: request canceled while waiting for connection (Client.Timeout exceeded while awaiting headers)",
			expectedMatch: true,
		},
		{
			description:   "case 14: request canceled due to client timeout exceeded",
			errorMessage:  "Get https://api.06bhh.k8s.geckon.gridscale.kvm.gigantic.io/api/v1/nodes?timeout=30s: net/http: request canceled (Client.Timeout exceeded while awaiting headers)",
			expectedMatch: true,
		},
		{
			description:   "case 15: certificate signed by unknown authority",
			errorMessage:  "Get https://api.ci-cur-42bc2-cba40.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s: x509: certificate signed by unknown authority (possibly because of \"crypto/rsa: verification error\" while trying to verify candidate authority certificate \"ci-cur-42bc2-cba40.k8s.godsmack.westeurope.azure.gigantic.io\")",
			expectedMatch: true,
		},
		{
			description:   "case 16: Patch timeout establishing TLS handshake",
			errorMessage:  "Patch https://api.xca65.k8s.geckon.gridscale.kvm.gigantic.io/api/v1/nodes/worker-sruw7-689bd75b49-8gbtl?timeout=30s: net/http: TLS handshake timeout",
			expectedMatch: true,
		},
		{
			description:   "case 17: Get i/o timeout establishing TCP connection",
			errorMessage:  "Get https://api.wgrt8.k8s.godsmack.westeurope.azure.gigantic.io/api/v1/nodes?timeout=30s: dial tcp 40.113.146.2:443: i/o timeout",
			expectedMatch: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			err := errors.New(tc.errorMessage)
			result := IsAPINotAvailable(err)

			if result != tc.expectedMatch {
				t.Fatalf("expected %t, got %t", tc.expectedMatch, result)
			}
		})
	}
}
