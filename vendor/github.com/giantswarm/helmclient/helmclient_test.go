package helmclient

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/giantswarm/micrologger/microloggertest"
	corev1 "k8s.io/api/core/v1"
	helmclient "k8s.io/helm/pkg/helm"
	helmchart "k8s.io/helm/pkg/proto/hapi/chart"
	helmrelease "k8s.io/helm/pkg/proto/hapi/release"
)

func Test_DeleteRelease(t *testing.T) {
	testCases := []struct {
		description  string
		namespace    string
		releaseName  string
		releases     []*helmrelease.Release
		errorMatcher func(error) bool
	}{
		{
			description:  "case 0: try to delete non-existent release",
			releaseName:  "chart-operator",
			releases:     []*helmrelease.Release{},
			errorMatcher: IsReleaseNotFound,
		},
		{
			description: "case 1: delete basic release",
			releaseName: "chart-operator",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "chart-operator",
					Namespace: "default",
				}),
			},
		},
		{
			description: "case 2: try to delete release with wrong name",
			releaseName: "chart-operator",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "node-exporter",
					Namespace: "default",
				}),
			},
			errorMatcher: IsReleaseNotFound,
		},
	}

	for _, tc := range testCases {
		ctx := context.Background()

		t.Run(tc.description, func(t *testing.T) {
			helm := Client{
				helmClient: &helmclient.FakeClient{
					Rels: tc.releases,
				},
				logger: microloggertest.New(),
			}
			err := helm.DeleteRelease(ctx, tc.releaseName)

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

func Test_GetReleaseContent(t *testing.T) {
	testCases := []struct {
		description     string
		releaseName     string
		releases        []*helmrelease.Release
		expectedContent *ReleaseContent
		errorMatcher    func(error) bool
	}{
		{
			description: "case 0: basic match with deployed status",
			releaseName: "chart-operator",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "chart-operator",
					Namespace: "default",
				}),
			},
			expectedContent: &ReleaseContent{
				Name:   "chart-operator",
				Status: "DEPLOYED",
				Values: map[string]interface{}{
					// Note: Values cannot be configured via the Helm mock client.
					"name": "value",
				},
			},
			errorMatcher: nil,
		},
		{
			description: "case 1: basic match with failed status",
			releaseName: "chart-operator",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:       "chart-operator",
					Namespace:  "default",
					StatusCode: helmrelease.Status_FAILED,
				}),
			},
			expectedContent: &ReleaseContent{
				Name:   "chart-operator",
				Status: "FAILED",
				Values: map[string]interface{}{
					"name": "value",
				},
			},
			errorMatcher: nil,
		},
		{
			description: "case 2: chart not found",
			releaseName: "missing",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name: "chart-operator",
				}),
			},
			expectedContent: nil,
			errorMatcher:    IsReleaseNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			helm := Client{
				helmClient: &helmclient.FakeClient{
					Rels: tc.releases,
				},
				logger: microloggertest.New(),
			}
			result, err := helm.GetReleaseContent(ctx, tc.releaseName)

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

			if !reflect.DeepEqual(result, tc.expectedContent) {
				t.Fatalf("Release == %q, want %q", result, tc.expectedContent)
			}
		})
	}
}

func Test_GetReleaseHistory(t *testing.T) {
	testCases := []struct {
		description     string
		releaseName     string
		releases        []*helmrelease.Release
		expectedHistory *ReleaseHistory
		errorMatcher    func(error) bool
	}{
		{
			description: "case 0: basic match with version",
			releaseName: "chart-operator",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "chart-operator",
					Namespace: "default",
					Chart: &helmchart.Chart{
						Metadata: &helmchart.Metadata{
							AppVersion: "1.0.0",
							Version:    "0.1.0",
						},
					},
				}),
			},
			expectedHistory: &ReleaseHistory{
				AppVersion:  "1.0.0",
				Description: "Release mock",
				Name:        "chart-operator",
				// LastDeployed is hardcoded in the fake Helm Client.
				LastDeployed: time.Unix(242085845, 0).UTC(),
				Version:      "0.1.0",
			},
			errorMatcher: nil,
		},
		{
			description: "case 1: different version",
			releaseName: "chart-operator",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "chart-operator",
					Namespace: "default",
					Chart: &helmchart.Chart{
						Metadata: &helmchart.Metadata{
							AppVersion: "2.0.0",
							Version:    "1.0.0-rc1",
						},
					},
				}),
			},
			expectedHistory: &ReleaseHistory{
				AppVersion:  "2.0.0",
				Description: "Release mock",
				Name:        "chart-operator",
				// LastDeployed is hardcoded in the fake Helm Client.
				LastDeployed: time.Unix(242085845, 0).UTC(),
				Version:      "1.0.0-rc1",
			},
			errorMatcher: nil,
		},
		{
			description: "case 2: not found",
			releaseName: "missing",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "chart-operator",
					Namespace: "default",
					Chart: &helmchart.Chart{
						Metadata: &helmchart.Metadata{
							Version: "1.0.0-rc1",
						},
					},
				}),
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "chart-operator",
					Namespace: "default",
					Chart: &helmchart.Chart{
						Metadata: &helmchart.Metadata{
							Version: "1.0.0-rc1",
						},
					},
				}),
			},
			expectedHistory: nil,
			errorMatcher:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			helm := Client{
				helmClient: &helmclient.FakeClient{
					Rels: tc.releases,
				},
				logger: microloggertest.New(),
			}
			result, err := helm.GetReleaseHistory(ctx, tc.releaseName)

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

			if !reflect.DeepEqual(result, tc.expectedHistory) {
				t.Fatalf("Release == %q, want %q", result, tc.expectedHistory)
			}
		})
	}
}

func Test_Client_InstallFromTarball(t *testing.T) {
	testCases := []struct {
		description   string
		namespace     string
		releases      []*helmrelease.Release
		expectedError bool
	}{
		{
			description: "basic install, empty releases",
			namespace:   "default",
			releases:    []*helmrelease.Release{},
		},
		{
			description: "other release in same namespace",
			namespace:   "default",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "my-chart",
					Namespace: "default",
				}),
			},
		},
		{
			description: "same release in same namespace",
			namespace:   "default",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "test-chart",
					Namespace: "default",
				}),
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			helm := Client{
				helmClient: &helmclient.FakeClient{
					Rels: tc.releases,
				},
				logger: microloggertest.New(),
			}
			// helm fake client does not actually use the tarball.
			err := helm.InstallReleaseFromTarball(ctx, "/path", tc.namespace, helmclient.ReleaseName("test-chart"))

			switch {
			case err == nil && !tc.expectedError:
				// correct; carry on
			case err != nil && !tc.expectedError:
				t.Fatalf("error == %#v, want nil", err)
			case err == nil && tc.expectedError:
				t.Fatalf("error == nil, want non-nil")
			case !tc.expectedError:
				t.Fatalf("error == %#v, want matching", err)
			}
		})
	}
}

func Test_Client_ListReleaseContents(t *testing.T) {
	testCases := []struct {
		description      string
		releases         []*helmrelease.Release
		expectedContents []*ReleaseContent
		errorMatcher     func(error) bool
	}{
		{
			description:      "case 0: no releases",
			releases:         []*helmrelease.Release{},
			expectedContents: []*ReleaseContent{},
			errorMatcher:     nil,
		},
		{
			description: "case 1: one release",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "foobar",
					Namespace: "default",
				}),
			},
			expectedContents: []*ReleaseContent{
				{
					Name:   "foobar",
					Status: "DEPLOYED",
					Values: map[string]interface{}{
						"name": "value",
					},
				},
			},
			errorMatcher: nil,
		},
		{
			description: "case 2: two releases, in two namespaces",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "foobar",
					Namespace: "default",
				}),
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "jabberwocky",
					Namespace: "not-default",
				}),
			},
			expectedContents: []*ReleaseContent{
				{
					Name:   "foobar",
					Status: "DEPLOYED",
					Values: map[string]interface{}{
						"name": "value",
					},
				},
				{
					Name:   "jabberwocky",
					Status: "DEPLOYED",
					Values: map[string]interface{}{
						"name": "value",
					},
				},
			},
			errorMatcher: nil,
		},
		{
			description: "case 3: two releases, one successful, one failed",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:       "foobar",
					Namespace:  "default",
					StatusCode: helmrelease.Status_DEPLOYED,
				}),
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:       "jabberwocky",
					Namespace:  "default",
					StatusCode: helmrelease.Status_FAILED,
				}),
			},
			expectedContents: []*ReleaseContent{
				{
					Name:   "foobar",
					Status: "DEPLOYED",
					Values: map[string]interface{}{
						"name": "value",
					},
				},
				{
					Name:   "jabberwocky",
					Status: "FAILED",
					Values: map[string]interface{}{
						"name": "value",
					},
				},
			},
			errorMatcher: nil,
		},
		{
			description: "case 4: two releases of the same chart with different versions",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:       "foobar",
					Namespace:  "default",
					StatusCode: helmrelease.Status_FAILED,
					Version:    1,
				}),
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:       "foobar",
					Namespace:  "default",
					StatusCode: helmrelease.Status_DEPLOYED,
					Version:    2,
				}),
			},
			expectedContents: []*ReleaseContent{
				{
					Name:   "foobar",
					Status: "DEPLOYED",
					Values: map[string]interface{}{
						"name": "value",
					},
				},
			},
			errorMatcher: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			helm := Client{
				helmClient: &helmclient.FakeClient{
					Rels: tc.releases,
				},
				logger: microloggertest.New(),
			}
			result, err := helm.ListReleaseContents(context.Background())

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

			if !reflect.DeepEqual(result, tc.expectedContents) {
				t.Fatalf("Releases == %#v, want %#v", result, tc.expectedContents)
			}
		})
	}
}

func Test_UpdateReleaseFromTarball(t *testing.T) {
	testCases := []struct {
		description  string
		namespace    string
		releaseName  string
		releases     []*helmrelease.Release
		errorMatcher func(error) bool
	}{
		{
			description:  "try to update non-existent release",
			releaseName:  "chart-operator",
			releases:     []*helmrelease.Release{},
			errorMatcher: IsReleaseNotFound,
		},
		{
			description: "update basic release",
			releaseName: "chart-operator",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name: "chart-operator",
				}),
			},
		},
		{
			description: "try to update release with wrong name",
			releaseName: "chart-operator",
			releases: []*helmrelease.Release{
				helmclient.ReleaseMock(&helmclient.MockReleaseOptions{
					Name:      "node-exporter",
					Namespace: "default",
				}),
			},
			errorMatcher: IsReleaseNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			ctx := context.Background()

			helm := Client{
				helmClient: &helmclient.FakeClient{
					Rels: tc.releases,
				},
				logger: microloggertest.New(),
			}
			// helm fake client does not actually use the tarball.
			err := helm.UpdateReleaseFromTarball(ctx, tc.releaseName, "/path")

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

func Test_isTillerInvalidVersion(t *testing.T) {
	testCases := []struct {
		name         string
		tillerPod    *corev1.Pod
		errorMatcher func(error) bool
	}{
		{
			name: "case 0: tiller pod is up to date",
			tillerPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: fmt.Sprintf("%s/%s", defaultTillerImageRegistry, defaultTillerImageName),
						},
					},
				},
			},
		},
		{
			name: "case 1: tiller pod is newer",
			tillerPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "quay.io/giantswarm/tiller:v9.8.7",
						},
					},
				},
			},
			errorMatcher: IsTillerInvalidVersion,
		},
		{
			name: "case 2: tiller pod is outdated",
			tillerPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "quay.io/giantswarm/tiller:v2.8.1",
						},
					},
				},
			},
			errorMatcher: IsTillerInvalidVersion,
		},
		{
			name: "case 3: tiller image is an outdated release candidate",
			tillerPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "quay.io/giantswarm/tiller:v2.8.0-rc.1",
						},
					},
				},
			},
			errorMatcher: IsTillerInvalidVersion,
		},
		{
			name: "case 4: tiller image has no version tag so we upgrade",
			tillerPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "quay.io/giantswarm/tiller",
						},
					},
				},
			},
			errorMatcher: IsTillerInvalidVersion,
		},
		{
			name: "case 5: tiller image uses latest tag so we upgrade",
			tillerPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "quay.io/giantswarm/tiller:latest",
						},
					},
				},
			},
			errorMatcher: IsTillerInvalidVersion,
		},
		{
			name: "case 6: tiller image tag format is invalid",
			tillerPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "quay.io/giantswarm/tiller:2.x.1",
						},
					},
				},
			},
			errorMatcher: IsExecutionFailed,
		},
		{
			name: "case 7: tiller image tag format is invalid",
			tillerPod: &corev1.Pod{
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "quay.io/giantswarm/tiller:4.3.2.1",
						},
					},
				},
			},
			errorMatcher: IsExecutionFailed,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := validateTillerVersion(tc.tillerPod, fmt.Sprintf("%s/%s", defaultTillerImageRegistry, defaultTillerImageName))

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
