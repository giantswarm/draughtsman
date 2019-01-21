package versionbundle

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/giantswarm/micrologger/microloggertest"
)

func urlMustParse(v string) *URL {
	u, err := url.Parse(v)
	if err != nil {
		panic(err)
	}

	return &URL{
		URL: u,
	}
}

func Test_buildReleases(t *testing.T) {
	testCases := []struct {
		name             string
		indexReleases    []IndexRelease
		bundles          []Bundle
		expectedReleases []Release
		errorMatcher     func(error) bool
	}{
		{
			name: "case 0: build one release",
			indexReleases: []IndexRelease{
				{
					Active: true,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://app-controller:8000/"),
							Name:     "app-controller",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://app-controller:8000/"),
							Name:     "appcatalog-controller",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "1.0.0",
				},
			},
			bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "cert-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
					},
					Name:    "cert-operator",
					Version: "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "cluster-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
					},
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "cluster-operator",
							Description: "Change feature x.",
							Kind:        KindChanged,
						},
					},
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.2.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "cluster-operator",
							Description: "Change AWS feature k.",
							Kind:        KindChanged,
						},
					},
					Name:     "cluster-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kvm-operator",
							Description: "Change feature y.",
							Kind:        KindChanged,
						},
					},
					Name:    "kvm-operator",
					Version: "1.2.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kvm-operator",
							Description: "Change feature z.",
							Kind:        KindChanged,
						},
					},
					Name:    "kvm-operator",
					Version: "1.4.2",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kvm-operator",
							Description: "New component q.",
							Kind:        KindAdded,
						},
					},
					Name:    "kvm-operator",
					Version: "2.2.1",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "app-controller",
							Description: "Added initial version.",
							Kind:        KindAdded,
						},
					},
					Name:    "app-controller",
					Version: "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "appcatalog-controller",
							Description: "Added initial version.",
							Kind:        KindAdded,
						},
					},
					Name:    "appcatalog-controller",
					Version: "0.1.0",
				},
			},
			expectedReleases: []Release{
				{
					active: true,
					bundles: []Bundle{
						{
							Changelogs: []Changelog{
								{
									Component:   "app-controller",
									Description: "Added initial version.",
									Kind:        KindAdded,
								},
							},
							Name:    "app-controller",
							Version: "0.1.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "appcatalog-controller",
									Description: "Added initial version.",
									Kind:        KindAdded,
								},
							},
							Name:    "appcatalog-controller",
							Version: "0.1.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "cert-operator",
									Description: "First release version.",
									Kind:        KindAdded,
								},
							},
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "cluster-operator",
									Description: "First release version.",
									Kind:        KindAdded,
								},
							},
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "kvm-operator",
									Description: "New component q.",
									Kind:        KindAdded,
								},
							},
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					changelogs: []Changelog{
						{
							Component:   "app-controller",
							Description: "Added initial version.",
							Kind:        KindAdded,
						},
						{
							Component:   "appcatalog-controller",
							Description: "Added initial version.",
							Kind:        KindAdded,
						},
						{
							Component:   "cert-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
						{
							Component:   "cluster-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
						{
							Component:   "kvm-operator",
							Description: "New component q.",
							Kind:        KindAdded,
						},
					},
					components: []Component{
						{
							Name:    "app-controller",
							Version: "0.1.0",
						},
						{
							Name:    "appcatalog-controller",
							Version: "0.1.0",
						},
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:    "cluster-operator",
							Version: "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					timestamp: time.Date(2018, time.April, 16, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: build two releases",
			indexReleases: []IndexRelease{
				{
					Active: true,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "1.0.0",
				},
				{
					Active: true,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 22, 12, 00, 0, 0, time.UTC),
					Version: "1.1.0",
				},
			},
			bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "cert-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
					},
					Name:    "cert-operator",
					Version: "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "cluster-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
					},
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "cluster-operator",
							Description: "Change feature x.",
							Kind:        KindChanged,
						},
					},
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.2.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "cluster-operator",
							Description: "Change AWS feature k.",
							Kind:        KindChanged,
						},
					},
					Name:     "cluster-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kvm-operator",
							Description: "Change feature y.",
							Kind:        KindChanged,
						},
					},
					Name:    "kvm-operator",
					Version: "1.2.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kvm-operator",
							Description: "Change feature z.",
							Kind:        KindChanged,
						},
					},
					Name:    "kvm-operator",
					Version: "1.4.2",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kvm-operator",
							Description: "New component q.",
							Kind:        KindAdded,
						},
					},
					Name:    "kvm-operator",
					Version: "2.2.1",
				},
			},
			expectedReleases: []Release{
				{
					active: true,
					bundles: []Bundle{
						{
							Changelogs: []Changelog{
								{
									Component:   "cert-operator",
									Description: "First release version.",
									Kind:        KindAdded,
								},
							},
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "cluster-operator",
									Description: "First release version.",
									Kind:        KindAdded,
								},
							},
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "kvm-operator",
									Description: "New component q.",
									Kind:        KindAdded,
								},
							},
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					changelogs: []Changelog{
						{
							Component:   "cert-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
						{
							Component:   "cluster-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
						{
							Component:   "kvm-operator",
							Description: "New component q.",
							Kind:        KindAdded,
						},
					},
					components: []Component{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:    "cluster-operator",
							Version: "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					timestamp: time.Date(2018, time.April, 16, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
				{
					active: true,
					bundles: []Bundle{
						{
							Changelogs: []Changelog{
								{
									Component:   "cert-operator",
									Description: "First release version.",
									Kind:        KindAdded,
								},
							},
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "cluster-operator",
									Description: "Change feature x.",
									Kind:        KindChanged,
								},
							},
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "kvm-operator",
									Description: "New component q.",
									Kind:        KindAdded,
								},
							},
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					changelogs: []Changelog{
						{
							Component:   "cert-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
						{
							Component:   "cluster-operator",
							Description: "Change feature x.",
							Kind:        KindChanged,
						},

						{
							Component:   "kvm-operator",
							Description: "New component q.",
							Kind:        KindAdded,
						},
					},
					components: []Component{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:    "cluster-operator",
							Version: "0.2.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					timestamp: time.Date(2018, time.April, 22, 12, 0, 0, 0, time.UTC),
					version:   "1.1.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: try to build two release but miss one bundle for second one",
			indexReleases: []IndexRelease{
				{
					Active: true,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "1.0.0",
				},
				{
					Active: true,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.4.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 22, 12, 00, 0, 0, time.UTC),
					Version: "1.1.0",
				},
			},
			bundles: []Bundle{
				{
					Changelogs: []Changelog{
						{
							Component:   "cert-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
					},
					Name:    "cert-operator",
					Version: "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "cluster-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
					},
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "cluster-operator",
							Description: "Change feature x.",
							Kind:        KindChanged,
						},
					},
					Name:     "cluster-operator",
					Provider: "kvm",
					Version:  "0.2.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "cluster-operator",
							Description: "Change AWS feature k.",
							Kind:        KindChanged,
						},
					},
					Name:     "cluster-operator",
					Provider: "aws",
					Version:  "0.1.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kvm-operator",
							Description: "Change feature y.",
							Kind:        KindChanged,
						},
					},
					Name:    "kvm-operator",
					Version: "1.2.0",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kvm-operator",
							Description: "Change feature z.",
							Kind:        KindChanged,
						},
					},
					Name:    "kvm-operator",
					Version: "1.4.2",
				},
				{
					Changelogs: []Changelog{
						{
							Component:   "kvm-operator",
							Description: "New component q.",
							Kind:        KindAdded,
						},
					},
					Name:    "kvm-operator",
					Version: "2.2.1",
				},
			},
			expectedReleases: []Release{
				{
					active: true,
					bundles: []Bundle{
						{
							Changelogs: []Changelog{
								{
									Component:   "cert-operator",
									Description: "First release version.",
									Kind:        KindAdded,
								},
							},
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "cluster-operator",
									Description: "First release version.",
									Kind:        KindAdded,
								},
							},
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Changelogs: []Changelog{
								{
									Component:   "kvm-operator",
									Description: "New component q.",
									Kind:        KindAdded,
								},
							},
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					changelogs: []Changelog{
						{
							Component:   "cert-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
						{
							Component:   "cluster-operator",
							Description: "First release version.",
							Kind:        KindAdded,
						},
						{
							Component:   "kvm-operator",
							Description: "New component q.",
							Kind:        KindAdded,
						},
					},
					components: []Component{
						{
							Name:    "cert-operator",
							Version: "0.1.0",
						},
						{
							Name:    "cluster-operator",
							Version: "0.1.0",
						},
						{
							Name:    "kvm-operator",
							Version: "2.2.1",
						},
					},
					timestamp: time.Date(2018, time.April, 16, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
			},
			errorMatcher: nil,
		},
	}

	logger := microloggertest.New()

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			releases, err := buildReleases(logger, tc.indexReleases, tc.bundles)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}

			if !reflect.DeepEqual(releases, tc.expectedReleases) {
				t.Fatalf("releases don't match expectedReleases; got:\n%#v\n\nexpected:\n%#v\n\n", releases, tc.expectedReleases)
			}
		})
	}
}

func Test_deduplicateReleaseChangelog(t *testing.T) {
	testCases := []struct {
		name             string
		releases         []Release
		expectedReleases []Release
		errorMatcher     func(error) bool
	}{
		{
			name: "case 0: simple linear changelog history without duplicates",
			releases: []Release{
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 2, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature y",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 12, 12, 0, 0, 0, time.UTC),
					version:   "2.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature z",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 22, 12, 0, 0, 0, time.UTC),
					version:   "3.0.0",
				},
			},
			expectedReleases: []Release{
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 2, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature y",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 12, 12, 0, 0, 0, time.UTC),
					version:   "2.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature z",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 22, 12, 0, 0, 0, time.UTC),
					version:   "3.0.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: simple linear changelog history with one duplicate",
			releases: []Release{
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 2, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
						{
							Component:   "bar-operator",
							Description: "changed feature k",
							Kind:        KindChanged,
						},
					},
					timestamp: time.Date(2018, time.May, 18, 12, 0, 0, 0, time.UTC),
					version:   "1.0.1",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature z",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 22, 12, 0, 0, 0, time.UTC),
					version:   "3.0.0",
				},
			},
			expectedReleases: []Release{
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 2, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "bar-operator",
							Description: "changed feature k",
							Kind:        KindChanged,
						},
					},
					timestamp: time.Date(2018, time.May, 18, 12, 0, 0, 0, time.UTC),
					version:   "1.0.1",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature z",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 22, 12, 0, 0, 0, time.UTC),
					version:   "3.0.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: introduction of patch to bar-operator",
			releases: []Release{
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
						{
							Component:   "bar-operator",
							Description: "new feature y",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 2, 12, 0, 0, 0, time.UTC),
					version:   "1.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
						{
							Component:   "bar-operator",
							Description: "changed feature y",
							Kind:        KindChanged,
						},
					},
					timestamp: time.Date(2018, time.May, 28, 12, 0, 0, 0, time.UTC),
					version:   "1.0.1",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
						{
							Component:   "bar-operator",
							Description: "new feature y",
							Kind:        KindAdded,
						},
						{
							Component:   "baz-operator",
							Description: "new feature quux",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 22, 12, 0, 0, 0, time.UTC),
					version:   "2.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
						{
							Component:   "bar-operator",
							Description: "changed feature y",
							Kind:        KindChanged,
						},
						{
							Component:   "baz-operator",
							Description: "new feature quux",
							Kind:        KindAdded,
						},
					},
					timestamp: time.Date(2018, time.May, 28, 15, 0, 0, 0, time.UTC),
					version:   "2.0.1",
				},
			},
			expectedReleases: []Release{
				{
					changelogs: []Changelog{
						{
							Component:   "foo-operator",
							Description: "new feature x",
							Kind:        KindAdded,
						},
						{
							Component:   "bar-operator",
							Description: "new feature y",
							Kind:        KindAdded,
						},
					},
					version: "1.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "bar-operator",
							Description: "changed feature y",
							Kind:        KindChanged,
						},
					},
					version: "1.0.1",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "baz-operator",
							Description: "new feature quux",
							Kind:        KindAdded,
						},
					},
					version: "2.0.0",
				},
				{
					changelogs: []Changelog{
						{
							Component:   "bar-operator",
							Description: "changed feature y",
							Kind:        KindChanged,
						},
					},
					version: "2.0.1",
				},
			},
			errorMatcher: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			filteredReleases := deduplicateReleaseChangelog(tc.releases)

			gotChangelogs := make([]string, 0)
			for _, r := range filteredReleases {
				gotChangelogs = append(gotChangelogs, fmt.Sprintf("Version %s: [", r.Version()))
				for _, clog := range r.Changelogs() {
					gotChangelogs = append(gotChangelogs, clog.String())
				}
				gotChangelogs = append(gotChangelogs, "]")
			}

			expectedChangelogs := make([]string, 0)
			for _, r := range tc.expectedReleases {
				expectedChangelogs = append(expectedChangelogs, fmt.Sprintf("Version %s: [", r.Version()))
				for _, clog := range r.Changelogs() {
					expectedChangelogs = append(expectedChangelogs, clog.String())
				}
				expectedChangelogs = append(expectedChangelogs, "]")
			}

			got := "[" + strings.Join(gotChangelogs, ", ") + "]"
			expected := "[" + strings.Join(expectedChangelogs, ", ") + "]"

			if got != expected {
				t.Fatalf("\ngot:\t\t%s\nexpected:\t%s", got, expected)
			}
		})
	}
}

func Test_findPreviousRelease(t *testing.T) {
	testCases := []struct {
		name            string
		r0              Release
		releases        []Release
		expectedRelease Release
	}{
		{
			name: "case 0: return empty Release when releases is empty",
			r0: Release{
				timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
			},
			releases:        []Release{},
			expectedRelease: Release{},
		},
		{
			name: "case 1: return empty Release when releases contains only current release",
			r0: Release{
				timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
			},
			releases: []Release{
				{
					timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
				},
			},
			expectedRelease: Release{},
		},
		{
			name: "case 2: return correct release when releases contains two releases",
			r0: Release{
				timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
			},
			releases: []Release{
				{
					timestamp: time.Date(2018, time.May, 23, 12, 0, 0, 0, time.UTC),
				},
				{
					timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
				},
			},
			expectedRelease: Release{
				timestamp: time.Date(2018, time.May, 23, 12, 0, 0, 0, time.UTC),
			},
		},
		{
			name: "case 3: return correct release when releases contains two older releases",
			r0: Release{
				timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
			},
			releases: []Release{
				{
					timestamp: time.Date(2018, time.May, 18, 12, 0, 0, 0, time.UTC),
				},
				{
					timestamp: time.Date(2018, time.May, 23, 12, 0, 0, 0, time.UTC),
				},
				{
					timestamp: time.Date(2018, time.May, 25, 12, 0, 0, 0, time.UTC),
				},
				{
					timestamp: time.Date(2018, time.May, 26, 12, 0, 0, 0, time.UTC),
				},
			},
			expectedRelease: Release{
				timestamp: time.Date(2018, time.May, 23, 12, 0, 0, 0, time.UTC),
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			release := findPreviousRelease(tc.r0, tc.releases)

			if !reflect.DeepEqual(release, tc.expectedRelease) {
				t.Fatalf("\ngot:\n%v\n\nexpected:\n%v\n\n", release, tc.expectedRelease)
			}
		})
	}
}

func Test_validateReleaseAuthorityEndpoints(t *testing.T) {
	testCases := []struct {
		name         string
		releases     []IndexRelease
		errorMatcher func(error) bool
	}{
		{
			name: "case 0: success with only single release",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: success with multiple releases",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: failure with single release containing one with nil authorities",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active:      false,
					Authorities: nil,
					Date:        time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version:     "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 3: failure with single release containing one with empty authorities list",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active:      false,
					Authorities: make([]Authority, 0),
					Date:        time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version:     "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 4: failure with single release containing one authority without endpoint",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: nil,
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 5: failure with single release containing one authority without name",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 6: failure with single release containing one authority without version",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateReleaseAuthorities(tc.releases)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}
		})
	}
}

func Test_validateReleaseDates(t *testing.T) {
	testCases := []struct {
		name         string
		releases     []IndexRelease
		errorMatcher func(error) bool
	}{
		{
			name: "case 0: success with only single release",
			releases: []IndexRelease{
				{
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "1.0.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: success with multiple unique releases",
			releases: []IndexRelease{
				{
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "3.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 20, 13, 12, 00, 00, time.UTC),
					Version: "2.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 19, 13, 12, 00, 00, time.UTC),
					Version: "1.0.0",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: failure with one release that has empty date",
			releases: []IndexRelease{
				{
					Date:    time.Time{},
					Version: "1.0.0",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 3: failure with multiple releases including one with empty date",
			releases: []IndexRelease{
				{
					Date:    time.Time{},
					Version: "4.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "3.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 20, 13, 12, 00, 00, time.UTC),
					Version: "2.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 19, 13, 12, 00, 00, time.UTC),
					Version: "1.0.0",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 4: failure with duplicate release date in two releases",
			releases: []IndexRelease{
				{
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "4.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 20, 13, 12, 00, 00, time.UTC),
					Version: "3.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 20, 13, 12, 00, 00, time.UTC),
					Version: "2.0.0",
				},
				{
					Date:    time.Date(2018, time.May, 19, 13, 12, 00, 00, time.UTC),
					Version: "1.0.0",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateReleaseDates(tc.releases)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}
		})
	}
}

func Test_validateUniqueReleases(t *testing.T) {
	testCases := []struct {
		name         string
		releases     []IndexRelease
		errorMatcher func(error) bool
	}{
		{
			name: "case 0: success with only single release",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 1: success with multiple unique releases",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: nil,
		},
		{
			name: "case 2: failure with multiple releases including one duplicate version",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 3: failure with multiple releases including one with duplicate version contents",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1-duplicate",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.March, 27, 12, 00, 0, 0, time.UTC),
					Version: "2.4.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
		{
			name: "case 4: failure with multiple releases including multiple duplicate version",
			releases: []IndexRelease{
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.5",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.3.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.3.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.May, 21, 13, 12, 00, 00, time.UTC),
					Version: "2.6.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.2.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
				{
					Active: false,
					Authorities: []Authority{
						{
							Endpoint: urlMustParse("http://cert-operator:8000/"),
							Name:     "cert-operator",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://cluster-operator:8000/"),
							Name:     "cluster-operator",
							Provider: "kvm",
							Version:  "0.1.0",
						},
						{
							Endpoint: urlMustParse("http://kvm-operator:8000/"),
							Name:     "kvm-operator",
							Version:  "2.2.1",
						},
					},
					Date:    time.Date(2018, time.April, 16, 12, 00, 0, 0, time.UTC),
					Version: "2.5.1",
				},
			},
			errorMatcher: IsInvalidRelease,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateUniqueReleases(tc.releases)

			switch {
			case err == nil && tc.errorMatcher == nil:
				// correct; carry on
			case err != nil && tc.errorMatcher == nil:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.errorMatcher != nil:
				t.Fatalf("error == nil, want non-nil")
			case !tc.errorMatcher(err):
				t.Fatalf("error == %#v, want matching", err)
			}
		})
	}
}
