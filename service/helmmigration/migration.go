package helmmigration

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"time"
)

type Config struct {
	KubernetesClient kubernetes.Interface
	Logger           micrologger.Logger

	Repository     string
	HelmBinaryPath string
	ProjectList    []string
}

type HelmMigration struct {
	kubernetesClient kubernetes.Interface
	logger           micrologger.Logger

	repository  string
	helmBinary  string
	projectList []string
}

func New(c Config) (HelmMigration, error) {
	h := HelmMigration{
		kubernetesClient: c.KubernetesClient,
		logger:           c.Logger,

		repository:  c.Repository,
		helmBinary:  c.HelmBinaryPath,
		projectList: c.ProjectList,
	}

	return h, nil
}

func (h *HelmMigration) Migrate(ctx context.Context) error {
	h.logger.Debugf(ctx, "Checking remaining helm2 releases from draughtsman")
	projectList, err := h.listRemainingHelmRelease()
	if err != nil {
		return microerror.Mask(err)
	}

	if len(projectList) == 0 {
		h.logger.Debugf(ctx, "no helm2 releases from draughtsman, quitting now...")
		return nil
	}

	h.logger.Debugf(ctx, "migrating total %s helm2 releases from draughtsman", len(projectList))

	err = h.installHelm2to3Migration(projectList)
	if err != nil {
		return microerror.Mask(err)
	}

	b := backoff.NewExponential(2*time.Minute, 15*time.Second)
	n := backoff.NewNotifier(h.logger, context.Background())

	o := func() error {
		remainingHelmRelease, err := h.listRemainingHelmRelease()
		if err != nil {
			return microerror.Mask(err)
		}

		if len(remainingHelmRelease) == 0 {
			return nil
		} else {
			return microerror.Maskf(executionFailedError, "still remaining %d helm2 releases", len(remainingHelmRelease))
		}
	}

	err = backoff.RetryNotify(o, b, n)
	if err != nil {
		return microerror.Mask(err)
	}

	h.logger.Debugf(ctx, "migrated total %s helm2 releases from draughtsman", len(projectList))

	h.logger.Debugf(ctx, "deleting migration resources")

	err = h.deleteHelm2to3Migration()
	if err != nil {
		return microerror.Mask(err)
	}

	h.logger.Debugf(ctx, "deleted migration resources")

	return nil
}

func (h *HelmMigration) listRemainingHelmRelease() ([]string, error) {
	var remaining []string
	for _, project := range h.projectList {
		lo := metav1.ListOptions{
			LabelSelector: fmt.Sprintf("OWNER=TILLER,NAME=%s", project),
		}

		projectConfigMaps, err := h.kubernetesClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).List(lo)
		if err != nil {
			return nil, microerror.Mask(err)
		}

		if len(projectConfigMaps.Items) > 0 {
			remaining = append(remaining, project)
		}
	}

	return remaining, nil
}

func (h *HelmMigration) installHelm2to3Migration(projectList []string) error {
	_, err := h.runHelmCommand("install", "install",
		"draughtsman-helm-migration",
		"default-catalog/helm-2to3-migration",
		"--namespace",
		"giantswarm",
		"--set",
		fmt.Sprintf("\"releases=%s\"", strings.Join(projectList, ",")),
		"--set",
		fmt.Sprintf("image.registry=%s", h.repository),
	)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (h *HelmMigration) deleteHelm2to3Migration() error {
	_, err := h.runHelmCommand("delete", "delete",
		"draughtsman-helm-migration",
		"--namespace",
		"giantswarm",
	)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

// runHelmCommand runs the given Helm command.
func (h *HelmMigration) runHelmCommand(name string, args ...string) (string, error) {
	h.logger.Log("debug", "running helm command", "name", name)

	cmd := exec.Command(h.helmBinary, args...)

	var stdOutBuf, stdErrBuf bytes.Buffer
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf

	h.logger.Log(
		"debug", "ran helm command", "name", name,
		"stdout", stdOutBuf.String(), "stderr", stdErrBuf.String(),
	)
	err := cmd.Run()
	if err != nil {
		return "", microerror.Maskf(executionFailedError, "error output: %s", stdErrBuf.String())
	}

	if strings.Contains(stdOutBuf.String(), "Error") {
		return "", microerror.Maskf(executionFailedError, stdOutBuf.String())
	}

	return stdOutBuf.String(), nil
}
