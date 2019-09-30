// +build k8srequired

package env

import (
	"fmt"
	"os"
)

const (
	EnvVarCircleCI      = "CIRCLECI"
	EnvVarCircleSHA     = "CIRCLE_SHA1"
	EnvVarE2EKubeconfig = "E2E_KUBECONFIG"
	EnvVarKeepResources = "KEEP_RESOURCES"
)

var (
	circleCI      string
	circleSHA     string
	keepResources string
	kubeconfig    string
)

func init() {
	circleCI = os.Getenv(EnvVarCircleCI)
	circleSHA = os.Getenv(EnvVarCircleSHA)
	if circleSHA == "" {
		panic(fmt.Sprintf("env var %#q must not be empty", EnvVarCircleSHA))
	}

	keepResources = os.Getenv(EnvVarKeepResources)

	kubeconfig = os.Getenv(EnvVarE2EKubeconfig)
	if kubeconfig == "" {
		panic(fmt.Sprintf("env var %#q must not be empty", EnvVarE2EKubeconfig))
	}
}

func CircleCI() bool {
	return circleCI == "true"
}

func CircleSHA() string {
	return circleSHA
}

func KeepResources() bool {
	return keepResources == "true"
}

func KubeConfigPath() string {
	return kubeconfig
}
