package versionbundle

import "testing"

func Test_VerifyBundleIDMatchAuthorityBundleID(t *testing.T) {
	testCases := []struct {
		name          string
		bundle        Bundle
		authority     Authority
		expectedMatch bool
	}{
		{
			name: "case 0: success with all fields set",
			bundle: Bundle{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			expectedMatch: true,
		},
		{
			name: "case 1: success with Name and Version fields set",
			bundle: Bundle{
				Name:    "foo-operator",
				Version: "1.2.9",
			},
			authority: Authority{
				Name:    "foo-operator",
				Version: "1.2.9",
			},
			expectedMatch: true,
		},
		{
			name: "case 2: success with all fields set, space prefix in authority name",
			bundle: Bundle{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:     "  foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			expectedMatch: true,
		},
		{
			name: "case 3: success with all fields set, space prefix in bundle name",
			bundle: Bundle{
				Name:     "  foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:     "  foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			expectedMatch: true,
		},
		{
			name: "case 4: success with all fields set, space suffix in authority name",
			bundle: Bundle{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:     "foo-operator   ",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			expectedMatch: true,
		},
		{
			name: "case 5: success with all fields set, space suffix in bundle name",
			bundle: Bundle{
				Name:     "foo-operator   ",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:     "  foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			expectedMatch: true,
		},
		{
			name: "case 6: fail with different provider",
			bundle: Bundle{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:     "foo-operator",
				Provider: "kvm",
				Version:  "1.2.9",
			},
			expectedMatch: false,
		},
		{
			name: "case 7: fail with different version",
			bundle: Bundle{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.10",
			},
			expectedMatch: false,
		},
		{
			name: "case 8: fail with different name",
			bundle: Bundle{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:     "bar-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			expectedMatch: false,
		},
		{
			name: "case 9: fail with other one missing provider",
			bundle: Bundle{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:    "foo-operator",
				Version: "1.2.9",
			},
			expectedMatch: false,
		},
		{
			name: "case 10: fail with other one missing version",
			bundle: Bundle{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Name:     "foo-operator",
				Provider: "hal9000",
			},
			expectedMatch: false,
		},
		{
			name: "case 11: fail with other one missing name",
			bundle: Bundle{
				Name:     "foo-operator",
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			authority: Authority{
				Provider: "hal9000",
				Version:  "1.2.9",
			},
			expectedMatch: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			match := tc.bundle.ID() == tc.authority.BundleID()
			if match != tc.expectedMatch {
				t.Fatalf("expectedMatch: %v, got %v when bundle.ID() == %s and authority.BundleID() == %s", tc.expectedMatch, match, tc.bundle.ID(), tc.authority.BundleID())
			}
		})
	}
}
