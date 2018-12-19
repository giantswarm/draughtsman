package versionbundle

import (
	"reflect"
	"sort"
	"testing"
)

func TestReleasesSortByVersion(t *testing.T) {
	testCases := []struct {
		name          string
		releases      []Release
		expectedOrder []Release
	}{
		{
			name: "case 0: sort 1.0.0, 2.0.0, 3.0.0",
			releases: []Release{
				{
					version: "2.0.0",
				},
				{
					version: "1.0.0",
				},
				{
					version: "3.0.0",
				},
			},
			expectedOrder: []Release{
				{
					version: "1.0.0",
				},
				{
					version: "2.0.0",
				},
				{
					version: "3.0.0",
				},
			},
		},
		{
			name: "case 1: sort 1.0.0, 1.10.1, 1.2.10",
			releases: []Release{
				{
					version: "1.10.1",
				},
				{
					version: "1.0.0",
				},
				{
					version: "1.2.10",
				},
			},
			expectedOrder: []Release{
				{
					version: "1.0.0",
				},
				{
					version: "1.2.10",
				},
				{
					version: "1.10.1",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sort.Sort(SortReleasesByVersion(tc.releases))
			if !reflect.DeepEqual(tc.releases, tc.expectedOrder) {
				expectedOrderMsg := "["
				for _, b := range tc.expectedOrder {
					if len(expectedOrderMsg) > 1 {
						expectedOrderMsg += ", "
					}

					expectedOrderMsg += b.version
				}

				expectedOrderMsg += "]"

				gotOrderMsg := "["
				for _, b := range tc.releases {
					if len(gotOrderMsg) > 1 {
						gotOrderMsg += ", "
					}
					gotOrderMsg += b.version
				}
				gotOrderMsg += "]"

				t.Fatalf("expected order: %s, got: %s", expectedOrderMsg, gotOrderMsg)
			}
		})
	}
}
