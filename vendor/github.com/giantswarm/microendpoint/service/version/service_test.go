package version

import (
	"context"
	"reflect"
	"runtime"
	"testing"

	"github.com/giantswarm/versionbundle"
)

func Test_Get(t *testing.T) {
	testCases := []struct {
		description                       string
		gitCommit                         string
		name                              string
		source                            string
		version                           string
		versionBundles                    []versionbundle.Bundle
		errorExpected                     bool
		errorExpectedDuringInitialization bool
		result                            Response
	}{
		// Case 0. A valid configuration.
		{
			description:                       "test desc",
			gitCommit:                         "b6bf741b5c34be4fff51d944f973318d8b078284",
			name:                              "api",
			source:                            "microkit",
			version:                           "1.0.0",
			versionBundles:                    nil,
			errorExpected:                     false,
			errorExpectedDuringInitialization: false,
			result: Response{
				Description:    "test desc",
				GitCommit:      "b6bf741b5c34be4fff51d944f973318d8b078284",
				GoVersion:      runtime.Version(),
				Name:           "api",
				OSArch:         runtime.GOOS + "/" + runtime.GOARCH,
				Source:         "microkit",
				Version:        "1.0.0",
				VersionBundles: nil,
			},
		},

		// Case 1. Same as 1 but with version bundles.
		{
			description: "test desc",
			gitCommit:   "b6bf741b5c34be4fff51d944f973318d8b078284",
			name:        "api",
			source:      "microkit",
			version:     "1.0.0",
			versionBundles: []versionbundle.Bundle{
				{
					Changelogs: []versionbundle.Changelog{
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
					Components: []versionbundle.Component{
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
					Changelogs: []versionbundle.Changelog{
						{
							Component:   "kubernetes",
							Description: "Kubernetes version updated.",
							Kind:        "changed",
						},
					},
					Components: []versionbundle.Component{
						{
							Name:    "etcd",
							Version: "3.2.0",
						},
						{
							Name:    "kubernetes",
							Version: "1.7.2",
						},
					},
					Name:    "cloud-config-operator",
					Version: "0.3.0",
				},
			},
			errorExpected:                     false,
			errorExpectedDuringInitialization: false,
			result: Response{
				Description: "test desc",
				GitCommit:   "b6bf741b5c34be4fff51d944f973318d8b078284",
				GoVersion:   runtime.Version(),
				Name:        "api",
				OSArch:      runtime.GOOS + "/" + runtime.GOARCH,
				Source:      "microkit",
				Version:     "1.0.0",
				VersionBundles: []versionbundle.Bundle{
					{
						Changelogs: []versionbundle.Changelog{
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
						Components: []versionbundle.Component{
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
						Changelogs: []versionbundle.Changelog{
							{
								Component:   "kubernetes",
								Description: "Kubernetes version updated.",
								Kind:        "changed",
							},
						},
						Components: []versionbundle.Component{
							{
								Name:    "etcd",
								Version: "3.2.0",
							},
							{
								Name:    "kubernetes",
								Version: "1.7.2",
							},
						},
						Name:    "cloud-config-operator",
						Version: "0.3.0",
					},
				},
			},
		},

		// Case 2. Ensure version bundle validation during service initialization.
		//
		// NOTE that changelogs and components are required.
		{
			description: "test desc",
			gitCommit:   "b6bf741b5c34be4fff51d944f973318d8b078284",
			name:        "api",
			source:      "microkit",
			versionBundles: []versionbundle.Bundle{
				{
					Changelogs: []versionbundle.Changelog{},
					Components: []versionbundle.Component{},
					Name:       "cloud-config-operator",
					Version:    "0.2.0",
				},
			},
			errorExpected:                     false,
			errorExpectedDuringInitialization: true,
			result:                            Response{},
		},

		// Case 3. Missing git commit.
		{
			description:                       "test desc",
			gitCommit:                         "",
			name:                              "microendpoint",
			source:                            "microkit",
			errorExpected:                     true,
			errorExpectedDuringInitialization: false,
			result:                            Response{},
		},
	}

	for i, tc := range testCases {
		config := Config{
			Description:    tc.description,
			GitCommit:      tc.gitCommit,
			Name:           tc.name,
			Source:         tc.source,
			Version:        tc.version,
			VersionBundles: tc.versionBundles,
		}

		service, err := New(config)
		if tc.errorExpectedDuringInitialization {
			if err == nil {
				t.Fatal("case", i, "expected", "error", "got", nil)
			}
		} else {
			if !tc.errorExpected {
				response, err := service.Get(context.TODO(), Request{})
				if !tc.errorExpected && err != nil {
					t.Fatal("case", i, "expected", nil, "got", err)
				}

				if !reflect.DeepEqual(*response, tc.result) {
					t.Fatal("case", i, "expected", tc.result, "got", response)
				}
			}
		}
	}
}
