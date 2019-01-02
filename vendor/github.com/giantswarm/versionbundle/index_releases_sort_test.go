package versionbundle

import (
	"reflect"
	"sort"
	"testing"
)

func TestIndexReleasesSortByVersion(t *testing.T) {
	testCases := []struct {
		name          string
		releases      []IndexRelease
		expectedOrder []IndexRelease
	}{
		{
			name: "case 0: sort 1.0.0, 2.0.0, 3.0.0",
			releases: []IndexRelease{
				{
					Version: "2.0.0",
				},
				{
					Version: "1.0.0",
				},
				{
					Version: "3.0.0",
				},
			},
			expectedOrder: []IndexRelease{
				{
					Version: "1.0.0",
				},
				{
					Version: "2.0.0",
				},
				{
					Version: "3.0.0",
				},
			},
		},
		{
			name: "case 1: sort 1.0.0, 1.10.1, 1.2.10",
			releases: []IndexRelease{
				{
					Version: "1.10.1",
				},
				{
					Version: "1.0.0",
				},
				{
					Version: "1.2.10",
				},
			},
			expectedOrder: []IndexRelease{
				{
					Version: "1.0.0",
				},
				{
					Version: "1.2.10",
				},
				{
					Version: "1.10.1",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sort.Sort(SortIndexReleasesByVersion(tc.releases))
			if !reflect.DeepEqual(tc.releases, tc.expectedOrder) {
				expectedOrderMsg := "["
				for _, b := range tc.expectedOrder {
					if len(expectedOrderMsg) > 1 {
						expectedOrderMsg += ", "
					}

					expectedOrderMsg += b.Version
				}

				expectedOrderMsg += "]"

				gotOrderMsg := "["
				for _, b := range tc.releases {
					if len(gotOrderMsg) > 1 {
						gotOrderMsg += ", "
					}
					gotOrderMsg += b.Version
				}
				gotOrderMsg += "]"

				t.Fatalf("expected order: %s, got: %s", expectedOrderMsg, gotOrderMsg)
			}
		})
	}
}
