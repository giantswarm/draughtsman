package versionbundle

import (
	"reflect"
	"sort"
	"testing"
)

func Test_Bundles_Contains(t *testing.T) {
	testCases := []struct {
		Bundles        []Bundle
		Bundle         Bundle
		ExpectedResult bool
	}{
		// Test 0 ensures that a nil list and an empty bundle result in false.
		{
			Bundles:        nil,
			Bundle:         Bundle{},
			ExpectedResult: false,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Bundles:        []Bundle{},
			Bundle:         Bundle{},
			ExpectedResult: false,
		},

		// Test 2 ensures a list containing one version bundle and a matching
		// version bundle results in true.
		{
			Bundles: []Bundle{
				{
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
			},
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
			ExpectedResult: true,
		},

		// Test 3 ensures a list containing two version bundle and a matching
		// version bundle results in true.
		{
			Bundles: []Bundle{
				{
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
				{
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
							Version: "1.2.0",
						},
						{
							Name:    "kube-dns",
							Version: "1.0.0",
						},
					},
					Name:    "kubernetes-operator",
					Version: "0.2.0",
				},
			},
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
						Version: "1.2.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			ExpectedResult: true,
		},

		// Test 4 ensures a list containing one version bundle and a version bundle
		// that does not match results in false.
		{
			Bundles: []Bundle{
				{
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
			},
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
						Version: "1.2.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.2.0",
			},
			ExpectedResult: false,
		},

		// Test 5 ensures a list containing two version bundle and a version bundle
		// that does not match results in false.
		{
			Bundles: []Bundle{
				{
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
				{
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
							Version: "1.2.0",
						},
						{
							Name:    "kube-dns",
							Version: "1.0.0",
						},
					},
					Name:    "kubernetes-operator",
					Version: "0.2.0",
				},
			},
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
						Version: "1.3.0",
					},
					{
						Name:    "kube-dns",
						Version: "1.0.0",
					},
				},
				Name:    "kubernetes-operator",
				Version: "0.3.0",
			},
			ExpectedResult: false,
		},
	}

	for i, tc := range testCases {
		result := Bundles(tc.Bundles).Contain(tc.Bundle)
		if result != tc.ExpectedResult {
			t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedResult, result)
		}
	}
}

func Test_Bundles_Copy(t *testing.T) {
	bundles := []Bundle{
		{
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
		{
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
			Version: "0.0.9",
		},
	}

	b1 := CopyBundles(bundles)
	b2 := CopyBundles(bundles)

	sort.Sort(SortBundlesByVersion(b2))

	if reflect.DeepEqual(b1, b2) {
		t.Fatalf("expected %#v got %#v", b1, b2)
	}
}

func Test_Bundles_GetBundleByName(t *testing.T) {
	testCases := []struct {
		Bundles        []Bundle
		Name           string
		ExpectedBundle Bundle
		ErrorMatcher   func(err error) bool
	}{
		// Test 0 ensures that a nil list and an empty name throws an execution
		// failed error.
		{
			Bundles:        nil,
			Name:           "",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 1 ensures that a nil list and a non-empty name throws an execution
		// failed error.
		{
			Bundles:        nil,
			Name:           "kubernetes-operator",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 2 ensures that a non-empty list and an empty name throws an execution
		// failed error.
		{
			Bundles: []Bundle{
				{
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
			},
			Name:           "",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 3 ensures that a non-empty list and an non-empty name throws a
		// not found errorn case the given name does not exist in the given list.
		{
			Bundles: []Bundle{
				{
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
			},
			Name:           "cert-operator",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsBundleNotFound,
		},

		// Test 4 is the same as 3 but with different version bundles.
		{
			Bundles: []Bundle{
				{
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
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.1.0",
				},
			},
			Name:           "cert-operator",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsBundleNotFound,
		},

		// Test 5 ensures that a bundle can be found.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.1.0",
				},
			},
			Name: "cloud-config-operator",
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Name:    "cloud-config-operator",
				Version: "0.1.0",
			},
			ErrorMatcher: nil,
		},

		// Test 6 is the same as 5 but with different bundles.
		{
			Bundles: []Bundle{
				{
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
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.1.0",
				},
			},
			Name: "cloud-config-operator",
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Name:    "cloud-config-operator",
				Version: "0.1.0",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		result, err := GetBundleByName(tc.Bundles, tc.Name)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if !reflect.DeepEqual(result, tc.ExpectedBundle) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedBundle, result)
			}
		}
	}
}

func Test_Bundles_GetBundleByNameForProvider(t *testing.T) {
	testCases := []struct {
		Bundles        []Bundle
		Name           string
		Provider       string
		ExpectedBundle Bundle
		ErrorMatcher   func(err error) bool
	}{
		// Test 0 ensures that a nil list and an empty name throws an execution
		// failed error.
		{
			Bundles:        nil,
			Name:           "",
			Provider:       "aws",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 1 ensures that a nil list and a non-empty name throws an execution
		// failed error.
		{
			Bundles:        nil,
			Name:           "kubernetes-operator",
			Provider:       "aws",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsBundleNotFound,
		},

		// Test 2 ensures that a non-empty list and an empty name throws an execution
		// failed error.
		{
			Bundles: []Bundle{
				{
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
			},
			Name:           "",
			Provider:       "",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 3 ensures that a non-empty list and an non-empty name throws a
		// not found errorn case the given name does not exist in the given list.
		{
			Bundles: []Bundle{
				{
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
					Name:     "kubernetes-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
			},
			Name:           "cert-operator",
			Provider:       "aws",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsBundleNotFound,
		},

		// Test 4 is the same as 3 but with different version bundles.
		{
			Bundles: []Bundle{
				{
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
					Name:     "kubernetes-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:     "cloud-config-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
			},
			Name:           "cert-operator",
			Provider:       "aws",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsBundleNotFound,
		},

		// Test 5 ensures that a bundle can be found.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:     "cloud-config-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
			},
			Name:     "cloud-config-operator",
			Provider: "aws",
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Name:     "cloud-config-operator",
				Provider: "aws",
				Version:  "0.1.0",
			},
			ErrorMatcher: nil,
		},

		// Test 6 is the same as 5 but with different bundles.
		{
			Bundles: []Bundle{
				{
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
					Name:     "kubernetes-operator",
					Provider: "azure",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:     "cloud-config-operator",
					Provider: "azure",
					Version:  "0.1.0",
				},
			},
			Name:     "cloud-config-operator",
			Provider: "azure",
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Name:     "cloud-config-operator",
				Provider: "azure",
				Version:  "0.1.0",
			},
			ErrorMatcher: nil,
		},

		// Test 7 is the same as 5 but with different bundles.
		{
			Bundles: []Bundle{
				{
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
					Name:     "kubernetes-operator",
					Provider: "azure",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:     "cloud-config-operator",
					Provider: "azure",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "aws",
					Name:       "cluster-operator",
					Version:    "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "azure",
					Name:       "cluster-operator",
					Version:    "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "kvm",
					Name:       "cluster-operator",
					Version:    "0.1.0",
				},
			},
			Name:     "cluster-operator",
			Provider: "azure",
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{},
				Provider:   "azure",
				Name:       "cluster-operator",
				Version:    "0.1.0",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		result, err := GetBundleByNameForProvider(tc.Bundles, tc.Name, tc.Provider)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if !reflect.DeepEqual(result, tc.ExpectedBundle) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedBundle, result)
			}
		}
	}
}

func Test_Bundles_GetNewestBundle(t *testing.T) {
	testCases := []struct {
		Bundles        []Bundle
		ExpectedBundle Bundle
		ErrorMatcher   func(err error) bool
	}{
		// Test 0 ensures that a nil list throws an execution failed error.
		{
			Bundles:        nil,
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsExecutionFailed,
		},

		// Test 1 ensures that the newest bundle can be found.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.1.0",
				},
			},
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Name:    "cloud-config-operator",
				Version: "0.1.0",
			},
			ErrorMatcher: nil,
		},

		// Test 2 is the same as 1 but with different bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
			},
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Name:    "cloud-config-operator",
				Version: "0.2.0",
			},
			ErrorMatcher: nil,
		},

		// Test 3 is the same as 1 but with different bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.1.0",
				},
			},
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Name:    "cloud-config-operator",
				Version: "0.2.0",
			},
			ErrorMatcher: nil,
		},

		// Test 4 is the same as 1 but with different bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "2.3.12",
				},
			},
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Name:    "cloud-config-operator",
				Version: "2.3.12",
			},
			ErrorMatcher: nil,
		},

		// Test 5 is the same as 1 but with different bundles.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "2.3.12",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.1.0",
				},
			},
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				Name:    "cloud-config-operator",
				Version: "2.3.12",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		result, err := GetNewestBundle(tc.Bundles)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if !reflect.DeepEqual(result, tc.ExpectedBundle) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedBundle, result)
			}
		}
	}
}

func Test_Bundles_GetNewestBundleForProvider(t *testing.T) {
	testCases := []struct {
		Bundles        []Bundle
		Provider       string
		ExpectedBundle Bundle
		ErrorMatcher   func(err error) bool
	}{
		// Test 0 verifies that newest bundle can be found for provider.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "aws",
					Name:       "cluster-operator",
					Version:    "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "azure",
					Name:       "cluster-operator",
					Version:    "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "kvm",
					Name:       "cluster-operator",
					Version:    "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "aws",
					Name:       "cluster-operator",
					Version:    "0.2.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "azure",
					Name:       "cluster-operator",
					Version:    "0.4.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "kvm",
					Name:       "cluster-operator",
					Version:    "0.3.0",
				},
			},
			Provider: "kvm",
			ExpectedBundle: Bundle{
				Changelogs: []Changelog{},
				Components: []Component{},
				Provider:   "kvm",
				Name:       "cluster-operator",
				Version:    "0.3.0",
			},
			ErrorMatcher: nil,
		},

		// Test 1 verifies that bundleNotFoundError is returned for missing
		// provider.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "aws",
					Name:       "cluster-operator",
					Version:    "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "azure",
					Name:       "cluster-operator",
					Version:    "0.1.0",
				},
				{
					Changelogs: []Changelog{},
					Components: []Component{},
					Provider:   "kvm",
					Name:       "cluster-operator",
					Version:    "0.1.0",
				},
			},
			Provider:       "bluemix",
			ExpectedBundle: Bundle{},
			ErrorMatcher:   IsBundleNotFound,
		},
	}

	for i, tc := range testCases {
		result, err := GetNewestBundleForProvider(tc.Bundles, tc.Provider)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if !reflect.DeepEqual(result, tc.ExpectedBundle) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedBundle, result)
			}
		}
	}
}

func Test_Bundles_Validate(t *testing.T) {
	testCases := []struct {
		Bundles      []Bundle
		ErrorMatcher func(err error) bool
	}{
		// Test 0 ensures that a nil list is not valid.
		{
			Bundles:      nil,
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Bundles:      []Bundle{},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 2 ensures validation of a list of version bundles where any version
		// bundle has no changelogs throws an error.
		{
			Bundles: []Bundle{
				{
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
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 3 is the same as 2 but with multiple version bundles.
		{
			Bundles: []Bundle{
				{
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
				{
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
					Version: "0.2.0",
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 4 ensures validation of a list of version bundles where a
		// version bundle has no components does not throw an error.
		{
			Bundles: []Bundle{
				{
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
			},
			ErrorMatcher: nil,
		},

		// Test 5 is the same as 4 but with multiple version bundles.
		{
			Bundles: []Bundle{
				{
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
				{
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
					Version: "0.2.0",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 6 ensures validation of a list of version bundles having the
		// different name and version not throws an error.
		{
			Bundles: []Bundle{
				{
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
				{
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
					Version: "0.2.0",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 7 ensures validation of a list of version bundles having duplicated
		// version bundles throws an error.
		{
			Bundles: []Bundle{
				{
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
				{
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
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 8 ensures validation of a list of version bundles having the same
		// version throws an error.
		{
			Bundles: []Bundle{
				{
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
				{
					Changelogs: []Changelog{
						{
							Component:   "kube-dns",
							Description: "Kube-DNS version updated.",
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
							Version: "1.1.0",
						},
					},
					Name:    "kubernetes-operator",
					Version: "0.1.0",
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},

		// Test 9 verifies that version increment is validated per provider.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Initial version for AWS.",
							Kind:        "added",
						},
					},
					Components: []Component{
						{
							Name:    "aws-operator",
							Version: "1.0.0",
						},
					},
					Provider: "aws",
					Name:     "cluster-operator",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Initial version for Azure.",
							Kind:        "added",
						},
					},
					Components: []Component{
						{
							Name:    "azure-operator",
							Version: "1.0.0",
						},
					},
					Provider: "azure",
					Name:     "cluster-operator",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Initial version for KVM.",
							Kind:        "added",
						},
					},
					Components: []Component{
						{
							Name:    "kvm-operator",
							Version: "1.0.0",
						},
					},
					Provider: "kvm",
					Name:     "cluster-operator",
					Version:  "0.1.0",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 10 like test 9 but verifies invalidBundlesError when there are
		// two Bundles for AWS provider and second of those has newer timestamp
		// and lower version number.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Initial version for AWS.",
							Kind:        "added",
						},
					},
					Components: []Component{
						{
							Name:    "aws-operator",
							Version: "1.0.0",
						},
					},
					Provider: "aws",
					Name:     "cluster-operator",
					Version:  "1.0.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Updated version for AWS.",
							Kind:        "updated",
						},
					},
					Components: []Component{
						{
							Name:    "aws-operator",
							Version: "1.0.0",
						},
					},
					Provider: "aws",
					Name:     "cluster-operator",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Initial version for Azure.",
							Kind:        "added",
						},
					},
					Components: []Component{
						{
							Name:    "azure-operator",
							Version: "1.0.0",
						},
					},
					Provider: "azure",
					Name:     "cluster-operator",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "Cluster Operator",
							Description: "Initial version for KVM.",
							Kind:        "added",
						},
					},
					Components: []Component{
						{
							Name:    "kvm-operator",
							Version: "1.0.0",
						},
					},
					Provider: "kvm",
					Name:     "cluster-operator",
					Version:  "0.1.0",
				},
			},
			ErrorMatcher: IsInvalidBundlesError,
		},
	}

	for i, tc := range testCases {
		err := Bundles(tc.Bundles).Validate()
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}
	}
}
