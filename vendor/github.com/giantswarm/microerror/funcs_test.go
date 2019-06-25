package microerror

import (
	"errors"
	"regexp"
	"strconv"
	"testing"
)

func Test_Stack(t *testing.T) {
	testCases := []struct {
		name                string
		inputErr            error
		expectedStackRegexp string
	}{
		{
			name:                "case 0: annotated microerror error",
			inputErr:            Maskf(&Error{Kind: "testKind"}, "annotation"),
			expectedStackRegexp: `^\[\{[a-zA-Z_/-]*/src/github.com/giantswarm/microerror/funcs_test.go:\d+: annotation\} \{test kind\}\]$`,
		},
		{
			name:                "case 1: non annotated microerror error",
			inputErr:            &Error{Kind: "testKind"},
			expectedStackRegexp: "^test kind$",
		},
		{
			name:                "case 2: external error",
			inputErr:            errors.New("external error"),
			expectedStackRegexp: "^external error$",
		},
	}

	for i, tc := range testCases {
		re, err := regexp.Compile(tc.expectedStackRegexp)
		if err != nil {
			t.Fatalf("err = %q, want nil", Stack(err))
		}

		t.Run(strconv.Itoa(i), func(t *testing.T) {
			stack := Stack(tc.inputErr)
			if !re.MatchString(stack) {
				t.Fatalf("stack = %q, want matching regexp %#q", stack, tc.expectedStackRegexp)
			}
		})
	}
}
