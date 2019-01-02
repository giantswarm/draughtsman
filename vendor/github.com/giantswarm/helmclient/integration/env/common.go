// +build k8srequired

package env

import (
	"fmt"
	"os"
)

const (
	EnvVarCircleCI      = "CIRCLECI"
	EnvVarCircleSHA     = "CIRCLE_SHA1"
	EnvVarKeepResources = "KEEP_RESOURCES"
	EnvVarTestDir       = "TEST_DIR"
)

var (
	circleCI      string
	circleSHA     string
	keepResources string
	testDir       string
)

func init() {
	circleCI = os.Getenv(EnvVarCircleCI)

	circleSHA = os.Getenv(EnvVarCircleSHA)
	if circleSHA == "" {
		panic(fmt.Sprintf("env var %#q must not be empty", EnvVarCircleSHA))
	}

	keepResources = os.Getenv(EnvVarKeepResources)

	testDir = os.Getenv(EnvVarTestDir)
	if testDir == "" {
		panic(fmt.Sprintf("env var %#q must not be empty", EnvVarTestDir))
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

func TestDir() string {
	return testDir
}
