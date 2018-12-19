package versionbundle

import (
	"reflect"
	"testing"
	"time"
)

func Test_Release_Changelogs(t *testing.T) {
	testCases := []struct {
		Bundles            []Bundle
		ExpectedChangelogs []Changelog
		ErrorMatcher       func(err error) bool
	}{
		// Test 0 ensures creating a release with a nil slice of bundles throws
		// an error when creating a new release type.
		{
			Bundles:            nil,
			ExpectedChangelogs: nil,
			ErrorMatcher:       IsInvalidConfig,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Bundles:            []Bundle{},
			ExpectedChangelogs: nil,
			ErrorMatcher:       IsInvalidConfig,
		},

		// Test 2 ensures computing the release changelogs when having a list
		// of one bundle given works as expected.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "kubernetes",
							Description: "description",
							Kind:        "fixed",
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
					Version: "0.0.1",
				},
			},
			ExpectedChangelogs: []Changelog{
				{
					Component:   "kubernetes",
					Description: "description",
					Kind:        "fixed",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 3 is the same as 2 but with a different changelogs.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "kubernetes",
							Description: "description",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "kube-dns",
							Version: "1.17.0",
						},
						{
							Name:    "calico",
							Version: "3.1.0",
						},
					},
					Name:    "kubernetes-operator",
					Version: "11.4.1",
				},
			},
			ExpectedChangelogs: []Changelog{
				{
					Component:   "kubernetes",
					Description: "description",
					Kind:        "changed",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures computing the release changelogs when having a list of
		// two bundles given works as expected.
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
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.1",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
			},
			ExpectedChangelogs: []Changelog{
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
				{
					Component:   "etcd",
					Description: "Etcd version updated.",
					Kind:        "changed",
				},
				{
					Component:   "kubernetes",
					Description: "Kubernetes version updated.",
					Kind:        "changed",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 5 is like 4 but with version bundles being flipped.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.1",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.2.0",
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
					Version: "0.1.0",
				},
			},
			ExpectedChangelogs: []Changelog{
				{
					Component:   "etcd",
					Description: "Etcd version updated.",
					Kind:        "changed",
				},
				{
					Component:   "kubernetes",
					Description: "Kubernetes version updated.",
					Kind:        "changed",
				},
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
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		config := ReleaseConfig{
			Bundles: tc.Bundles,
		}

		r, err := NewRelease(config)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}

		c := r.Changelogs()
		if !reflect.DeepEqual(c, tc.ExpectedChangelogs) {
			t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedChangelogs, c)
		}
	}
}

func Test_Release_Components(t *testing.T) {
	testCases := []struct {
		Bundles            []Bundle
		ExpectedComponents []Component
		ErrorMatcher       func(err error) bool
	}{
		// Test 0 ensures creating a release with a nil slice of bundles throws
		// an error when creating a new release type.
		{
			Bundles:            nil,
			ExpectedComponents: nil,
			ErrorMatcher:       IsInvalidConfig,
		},

		// Test 1 is the same as 0 but with an empty list of bundles.
		{
			Bundles:            []Bundle{},
			ExpectedComponents: nil,
			ErrorMatcher:       IsInvalidConfig,
		},

		// Test 2 ensures computing the release components when having a list
		// of one bundle given works as expected.
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
					Version: "0.0.1",
				},
			},
			ExpectedComponents: []Component{
				{
					Name:    "calico",
					Version: "1.1.0",
				},
				{
					Name:    "kube-dns",
					Version: "1.0.0",
				},
				{
					Name:    "kubernetes-operator",
					Version: "0.0.1",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 3 is the same as 2 but with a different components.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{},
					Components: []Component{
						{
							Name:    "kube-dns",
							Version: "1.17.0",
						},
						{
							Name:    "calico",
							Version: "3.1.0",
						},
					},
					Name:    "kubernetes-operator",
					Version: "11.4.1",
				},
			},
			ExpectedComponents: []Component{
				{
					Name:    "calico",
					Version: "3.1.0",
				},
				{
					Name:    "kube-dns",
					Version: "1.17.0",
				},
				{
					Name:    "kubernetes-operator",
					Version: "11.4.1",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 4 ensures computing the release components when having a list of
		// two bundles given works as expected.
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
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.1",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
			},
			ExpectedComponents: []Component{
				{
					Name:    "calico",
					Version: "1.1.0",
				},
				{
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
				{
					Name:    "etcd",
					Version: "3.2.0",
				},
				{
					Name:    "kube-dns",
					Version: "1.0.0",
				},
				{
					Name:    "kubernetes",
					Version: "1.7.1",
				},
				{
					Name:    "kubernetes-operator",
					Version: "0.1.0",
				},
			},
			ErrorMatcher: nil,
		},

		// Test 5 is like 4 but with version bundles being flipped.
		{
			Bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "etcd",
							Description: "Etcd version updated.",
							Kind:        "changed",
						},
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.1",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.2.0",
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
					Version: "0.1.0",
				},
			},
			ExpectedComponents: []Component{
				{
					Name:    "calico",
					Version: "1.1.0",
				},
				{
					Name:    "cloud-config-operator",
					Version: "0.2.0",
				},
				{
					Name:    "etcd",
					Version: "3.2.0",
				},
				{
					Name:    "kube-dns",
					Version: "1.0.0",
				},
				{
					Name:    "kubernetes",
					Version: "1.7.1",
				},
				{
					Name:    "kubernetes-operator",
					Version: "0.1.0",
				},
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		config := ReleaseConfig{
			Bundles: tc.Bundles,
		}

		r, err := NewRelease(config)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		}

		c := r.Components()
		if !reflect.DeepEqual(c, tc.ExpectedComponents) {
			t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedComponents, c)
		}
	}
}

func Test_Releases_GetNewestRelease(t *testing.T) {
	testCases := []struct {
		Releases        []Release
		ExpectedRelease Release
		ErrorMatcher    func(err error) bool
	}{
		// Test 0 ensures that a nil list throws an execution failed error.
		{
			Releases:        nil,
			ExpectedRelease: Release{},
			ErrorMatcher:    IsExecutionFailed,
		},

		// Test 1 ensures that the newest release can be found.
		{
			Releases: []Release{
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
			},
			ExpectedRelease: Release{
				bundles:    []Bundle{},
				changelogs: []Changelog{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
				version:   "0.1.0",
			},
			ErrorMatcher: nil,
		},

		// Test 2 is the same as 1 but with different releases.
		{
			Releases: []Release{
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
					version:   "0.2.0",
				},
			},
			ExpectedRelease: Release{
				bundles:    []Bundle{},
				changelogs: []Changelog{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
				version:   "0.2.0",
			},
			ErrorMatcher: nil,
		},

		// Test 3 is the same as 1 but with different releases.
		{
			Releases: []Release{
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
					version:   "0.2.0",
				},
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
			},
			ExpectedRelease: Release{
				bundles:    []Bundle{},
				changelogs: []Changelog{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
				version:   "0.2.0",
			},
			ErrorMatcher: nil,
		},

		// Test 4 is the same as 1 but with different releases.
		{
			Releases: []Release{
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
					version:   "0.2.0",
				},
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 40, 0, time.UTC),
					version:   "2.3.12",
				},
			},
			ExpectedRelease: Release{
				bundles:    []Bundle{},
				changelogs: []Changelog{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 40, 0, time.UTC),
				version:   "2.3.12",
			},
			ErrorMatcher: nil,
		},

		// Test 5 is the same as 1 but with different releases.
		{
			Releases: []Release{
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 20, 0, time.UTC),
					version:   "0.2.0",
				},
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 40, 0, time.UTC),
					version:   "2.3.12",
				},
				{
					bundles:    []Bundle{},
					changelogs: []Changelog{},
					components: []Component{
						{
							Name:    "calico",
							Version: "1.1.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.5",
						},
					},
					timestamp: time.Date(1970, time.January, 1, 0, 0, 10, 0, time.UTC),
					version:   "0.1.0",
				},
			},
			ExpectedRelease: Release{
				bundles:    []Bundle{},
				changelogs: []Changelog{},
				components: []Component{
					{
						Name:    "calico",
						Version: "1.1.0",
					},
					{
						Name:    "kubernetes",
						Version: "1.7.5",
					},
				},
				timestamp: time.Date(1970, time.January, 1, 0, 0, 40, 0, time.UTC),
				version:   "2.3.12",
			},
			ErrorMatcher: nil,
		},
	}

	for i, tc := range testCases {
		result, err := GetNewestRelease(tc.Releases)
		if tc.ErrorMatcher != nil {
			if !tc.ErrorMatcher(err) {
				t.Fatalf("test %d expected %#v got %#v", i, true, false)
			}
		} else if err != nil {
			t.Fatalf("test %d expected %#v got %#v", i, nil, err)
		} else {
			if !reflect.DeepEqual(result, tc.ExpectedRelease) {
				t.Fatalf("test %d expected %#v got %#v", i, tc.ExpectedRelease, result)
			}
		}
	}
}

func Test_Release_removeChangelogEntry(t *testing.T) {
	testCases := []struct {
		name              string
		release           Release
		changelogToRemove Changelog
		expectedRelease   Release
	}{
		{
			name:    "case 0: remove from empty list",
			release: Release{},
			changelogToRemove: Changelog{
				Component:   "foo-operator",
				Description: "changed bar",
				Kind:        KindChanged,
			},
			expectedRelease: Release{},
		},
		{
			name: "case 1: remove from list of one",
			release: Release{
				changelogs: []Changelog{
					{
						Component:   "foo-operator",
						Description: "changed bar",
						Kind:        KindChanged,
					},
				},
			},
			changelogToRemove: Changelog{
				Component:   "foo-operator",
				Description: "changed bar",
				Kind:        KindChanged,
			},
			expectedRelease: Release{
				changelogs: []Changelog{},
			},
		},
		{
			name: "case 2: remove from list of many",
			release: Release{
				changelogs: []Changelog{
					{
						Component:   "foo-operator",
						Description: "changed bar",
						Kind:        KindChanged,
					},
					{
						Component:   "foo-operator",
						Description: "added quux",
						Kind:        KindAdded,
					},
					{
						Component:   "bar-operator",
						Description: "fixed bug",
						Kind:        KindFixed,
					},
				},
			},
			changelogToRemove: Changelog{
				Component:   "foo-operator",
				Description: "changed bar",
				Kind:        KindChanged,
			},
			expectedRelease: Release{
				changelogs: []Changelog{
					{
						Component:   "foo-operator",
						Description: "added quux",
						Kind:        KindAdded,
					},
					{
						Component:   "bar-operator",
						Description: "fixed bug",
						Kind:        KindFixed,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tc.release.removeChangelogEntry(tc.changelogToRemove)
			if !reflect.DeepEqual(tc.release, tc.expectedRelease) {
				t.Fatalf("got %#v, expected %#v", tc.release, tc.expectedRelease)
			}
		})
	}
}
