package versionbundle

import (
	"testing"
)

func Test_Bundle_IsMajorUpgrade(t *testing.T) {
	testCases := []struct {
		Bundle         Bundle
		Other          Bundle
		ErrorMatcher   func(err error) bool
		ExpectedResult bool
	}{
		// Test 0 ensures that empty bundles throw an error.
		{
			Bundle:         Bundle{},
			Other:          Bundle{},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 1 is the same as 0 but with the lower bundle being empty.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			Other:          Bundle{},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 2 is the same as 0 but with the higher bundle being empty.
		{
			Bundle: Bundle{},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 3 ensures the same version results in false.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 4 ensures version bundles of different authorities cannot be
		// compared and result in an error.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "ingress-operator",
				Version: "0.1.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 5 ensures a smaller major version is not considered a major upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "1.1.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 6 is the same as 5 but with different versions.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "5.7.18",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "4.17.7",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 7 ensures a smaller minor version is not considered a major upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 8 ensures a smaller patch version is not considered a major upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 9 ensures a bigger major version is considered a major upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "1.0.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 10 is the same as 9 but with different versions.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "5.17.8",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "11.8.17",
			},
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 11 ensures a bigger minor version is not considered a major upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.3.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 12 ensures a bigger patch version is not considered a major upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.2",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},
	}

	for i, tc := range testCases {
		result, err := tc.Bundle.IsMajorUpgrade(tc.Other)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if result != tc.ExpectedResult {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedResult, result)
			}
		}
	}
}

func Test_Bundle_IsMinorUpgrade(t *testing.T) {
	testCases := []struct {
		Bundle         Bundle
		Other          Bundle
		ErrorMatcher   func(err error) bool
		ExpectedResult bool
	}{
		// Test 0 ensures that empty bundles throw an error.
		{
			Bundle:         Bundle{},
			Other:          Bundle{},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 1 is the same as 0 but with the lower bundle being empty.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			Other:          Bundle{},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 2 is the same as 0 but with the higher bundle being empty.
		{
			Bundle: Bundle{},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 3 ensures the same version results in false.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 4 ensures version bundles of different authorities cannot be
		// compared and result in an error.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "ingress-operator",
				Version: "0.1.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 5 ensures a smaller major version is not considered a minor upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "1.1.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 6 is the same as 5 but with different versions.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "5.7.18",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "4.17.7",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 7 ensures a smaller minor version is not considered a minor upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 8 ensures a smaller patch version is not considered a minor upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 9 ensures a bigger major version is not considered a minor upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "1.0.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 10 ensures a bigger minor version is considered a minor upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.3.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 11 is the same as 10 but with different versions.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "5.7.18",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "5.17.7",
			},
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 12 ensures a bigger patch version is not considered a minor upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.2",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},
	}

	for i, tc := range testCases {
		result, err := tc.Bundle.IsMinorUpgrade(tc.Other)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if result != tc.ExpectedResult {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedResult, result)
			}
		}
	}
}

func Test_Bundle_IsPatchUpgrade(t *testing.T) {
	testCases := []struct {
		Bundle         Bundle
		Other          Bundle
		ErrorMatcher   func(err error) bool
		ExpectedResult bool
	}{
		// Test 0 ensures that empty bundles throw an error.
		{
			Bundle:         Bundle{},
			Other:          Bundle{},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 1 is the same as 0 but with the lower bundle being empty.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			Other:          Bundle{},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 2 is the same as 0 but with the higher bundle being empty.
		{
			Bundle: Bundle{},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 3 ensures the same version results in false.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 4 version bundles of different authorities cannot be compared and
		// result in an error.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "ingress-operator",
				Version: "0.1.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			ErrorMatcher:   IsInvalidBundleError,
			ExpectedResult: false,
		},

		// Test 5 ensures a smaller major version is not considered a patch upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "1.1.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 6 ensures a smaller minor version is not considered a patch upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 7 ensures a smaller patch version is not considered a patch upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 8 ensures a bigger major version is not considered a patch upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "1.0.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 9 ensures a bigger minor version is not considered a patch upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.3.0",
			},
			ErrorMatcher:   nil,
			ExpectedResult: false,
		},

		// Test 10 ensures a bigger patch version is considered a patch upgrade.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.1",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.2",
			},
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},

		// Test 11 is the same as 10 but with a versions.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "2.33.5",
			},
			Other: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "2.33.8",
			},
			ErrorMatcher:   nil,
			ExpectedResult: true,
		},
	}

	for i, tc := range testCases {
		result, err := tc.Bundle.IsPatchUpgrade(tc.Other)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if result != tc.ExpectedResult {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedResult, result)
			}
		}
	}
}

func Test_Bundle_Validate(t *testing.T) {
	testCases := []struct {
		Bundle       Bundle
		ErrorMatcher func(err error) bool
	}{
		// Test 0 ensures that an empty version bundle is not valid.
		{
			Bundle:       Bundle{},
			ErrorMatcher: IsInvalidBundleError,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{},
				Name:       "",
				Version:    "",
			},
			ErrorMatcher: IsInvalidBundleError,
		},

		// Test 2 ensures a version bundle without changelogs throws an error.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher: IsInvalidBundleError,
		},

		// Test 3 ensures a version bundle without components does not throw an
		// error.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
					{
						Component:   "kubernetes",
						Description: "Kubernetes version requirements changed due to calico update.",
						Kind:        "changed",
					},
				},
				Components: []Component{},
				Name:       "kubernetes-operator",
				Version:    "0.1.0",
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures a version bundle without version throws an error.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
					{
						Component:   "kubernetes",
						Description: "Kubernetes version requirements changed due to calico update.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "",
			},
			ErrorMatcher: IsInvalidBundleError,
		},

		// Test 5 ensures an invalid version throws an error.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
					{
						Component:   "kubernetes",
						Description: "Kubernetes version requirements changed due to calico update.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "foo",
			},
			ErrorMatcher: IsInvalidBundleError,
		},

		// Test 6 is the same as 5 but with a different version.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
					{
						Component:   "kubernetes",
						Description: "Kubernetes version requirements changed due to calico update.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "1.2.3.4",
			},
			ErrorMatcher: IsInvalidBundleError,
		},

		// Test 7 ensures a valid version bundle does not throw an error.
		{
			Bundle: Bundle{
				Changelogs: []Changelog{
					{
						Component:   "calico",
						Description: "Calico version updated.",
						Kind:        "changed",
					},
					{
						Component:   "kubernetes",
						Description: "Kubernetes version requirements changed due to calico update.",
						Kind:        "changed",
					},
				},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.1.0",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		err := tc.Bundle.Validate()
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}
	}
}
