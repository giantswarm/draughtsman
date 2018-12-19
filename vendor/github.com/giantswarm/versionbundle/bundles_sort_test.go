package versionbundle

import (
	"reflect"
	"sort"
	"testing"
)

func TestBundlesSortByVersion(t *testing.T) {
	testCases := []struct {
		name          string
		bundles       []Bundle
		expectedOrder []Bundle
	}{
		{
			name: "case 0: sort 1.0.0, 2.0.0, 3.0.0",
			bundles: []Bundle{
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
			expectedOrder: []Bundle{
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
			bundles: []Bundle{
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
			expectedOrder: []Bundle{
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
			sort.Sort(SortBundlesByVersion(tc.bundles))
			if !reflect.DeepEqual(tc.bundles, tc.expectedOrder) {
				expectedOrderMsg := "["
				for _, b := range tc.expectedOrder {
					if len(expectedOrderMsg) > 1 {
						expectedOrderMsg += ", "
					}

					expectedOrderMsg += b.Version
				}

				expectedOrderMsg += "]"

				gotOrderMsg := "["
				for _, b := range tc.bundles {
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
