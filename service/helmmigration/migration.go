package helmmigration

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	"github.com/giantswarm/micrologger"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	chartName   = "default-catalog/helm-2to3-migration"
	releaseName = "draughtsman-helm-migration"
)

type Config struct {
	KubernetesClient kubernetes.Interface
	Logger           micrologger.Logger

	HelmBinaryPath string
	ProjectList    []string
	Repository     string
}

type HelmMigration struct {
	kubernetesClient kubernetes.Interface
	logger           micrologger.Logger

	helmBinary  string
	repository  string
	projectList []string
}

func New(c Config) (*HelmMigration, error) {
	if c.KubernetesClient == nil {
		return nil, microerror.Maskf(invalidConfigError, "kubernetesClient must not be empty")
	}
	if c.Logger == nil {
		return nil, microerror.Maskf(invalidConfigError, "logger must not be empty")
	}

	if c.HelmBinaryPath == "" {
		return nil, microerror.Maskf(invalidConfigError, "helmBinaryPath must not be empty")
	}
	if c.Repository == "" {
		return nil, microerror.Maskf(invalidConfigError, "repository must not be empty")
	}
	if c.ProjectList == nil {
		return nil, microerror.Maskf(invalidConfigError, "projectList must not be empty")
	}

	h := HelmMigration{
		kubernetesClient: c.KubernetesClient,
		logger:           c.Logger,

		repository:  c.Repository,
		helmBinary:  c.HelmBinaryPath,
		projectList: c.ProjectList,
	}

	return &h, nil
}

func (h *HelmMigration) Migrate(ctx context.Context) error {
	h.logger.Debugf(ctx, "checking remaining helm2 releases from draughtsman")
	projectList, err := h.listRemainingHelmRelease()
	if err != nil {
		return microerror.Mask(err)
	}

	if len(projectList) == 0 {
		h.logger.Debugf(ctx, "no helm2 releases from draughtsman")

		err = h.ensureTillerDeleted(ctx)
		if err != nil {
			return microerror.Mask(err)
		}

		err = h.ensureBackupConfigMapsDeleted(ctx)
		if err != nil {
			return microerror.Mask(err)
		}

		return nil
	}

	h.logger.Debugf(ctx, "migrating total %d helm2 releases from draughtsman", len(projectList))

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

func (h *HelmMigration) deleteHelm2to3Migration() error {
	_, err := h.runHelmCommand("delete", "delete",
		releaseName,
		"--namespace",
		"giantswarm",
	)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}

func (h *HelmMigration) ensureTillerDeleted(ctx context.Context) error {
	serviceAccounts := []string{"tiller", "tiller-giantswarm"}
	for _, account := range serviceAccounts {
		h.logger.Debugf(ctx, "deleting service account %#q in namespace %#q", account, metav1.NamespaceSystem)
		err := h.kubernetesClient.CoreV1().ServiceAccounts(metav1.NamespaceSystem).Delete(account, &metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// no-op
		} else if err != nil {
			return microerror.Mask(err)
		}
		h.logger.Debugf(ctx, "deleted service account %#q in namespace %#q", account, metav1.NamespaceSystem)
	}

	clusterRoleBindings := []string{"tiller", "tiller-giantswarm-psp", "tiller-kube-system"}
	for _, c := range clusterRoleBindings {
		h.logger.Debugf(ctx, "deleting cluster role binding %#q", c)
		err := h.kubernetesClient.RbacV1().ClusterRoleBindings().Delete(c, &metav1.DeleteOptions{})
		if apierrors.IsNotFound(err) {
			// no-op
		} else if err != nil {
			return microerror.Mask(err)
		}
		h.logger.Debugf(ctx, "deleted cluster role binding %#q", c)
	}

	pspName := "tiller-giantswarm-psp"
	h.logger.Debugf(ctx, "deleting pod security policy %#q", pspName)
	err := h.kubernetesClient.PolicyV1beta1().PodSecurityPolicies().Delete(pspName, &metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		// no-op
	} else if err != nil {
		return microerror.Mask(err)
	}

	h.logger.Debugf(ctx, "deleted pod security policy %#q", pspName)

	networkPolicyName := "tiller-giantswarm"
	h.logger.Debugf(ctx, "deleting network policy %#q in namespace %#q", networkPolicyName, metav1.NamespaceSystem)
	err = h.kubernetesClient.NetworkingV1().NetworkPolicies(metav1.NamespaceSystem).Delete(networkPolicyName, &metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		// no-op
	} else if err != nil {
		return microerror.Mask(err)
	}

	h.logger.Debugf(ctx, "deleted network policy %#q in namespace %#q", networkPolicyName, metav1.NamespaceSystem)

	h.logger.Debugf(ctx, "deleting deployment `tiller-deploy` in namespace %#q", metav1.NamespaceSystem)
	err = h.kubernetesClient.AppsV1().Deployments(metav1.NamespaceSystem).Delete("tiller-deploy", &metav1.DeleteOptions{})
	if apierrors.IsNotFound(err) {
		// no-op
	} else if err != nil {
		return microerror.Mask(err)
	}
	h.logger.Debugf(ctx, "deleted deployment `tiller-deploy` in namespace %#q", metav1.NamespaceSystem)

	return nil
}

func (h *HelmMigration) ensureBackupConfigMapsDeleted(ctx context.Context) error {
	lo := metav1.ListOptions{
		LabelSelector: "OWNER=HELM-BACKUP-V2",
	}

	cms, err := h.kubernetesClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).List(lo)
	if err != nil {
		return microerror.Mask(err)
	}

	h.logger.Debugf(ctx, "deleting %d helm v2 backup configmaps", len(cms.Items))

	err = h.kubernetesClient.CoreV1().ConfigMaps(metav1.NamespaceSystem).DeleteCollection(&metav1.DeleteOptions{}, lo)
	if err != nil {
		return microerror.Mask(err)
	}

	h.logger.Debugf(ctx, "deleted %d helm v2 backup configmaps", len(cms.Items))

	return nil
}

func (h *HelmMigration) installHelm2to3Migration(projectList []string) error {
	_, err := h.runHelmCommand("install", "install",
		releaseName,
		chartName,
		"--namespace",
		"giantswarm",
		"--set",
		fmt.Sprintf("releases={%s}", strings.Join(projectList, ",")),
		"--set",
		fmt.Sprintf("image.registry=%s", h.repository),
	)
	if err != nil {
		return microerror.Mask(err)
	}

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

// runHelmCommand runs the given Helm command.
func (h *HelmMigration) runHelmCommand(name string, args ...string) (string, error) {
	h.logger.Log("debug", "running helm command", "name", name)

	cmd := exec.Command(h.helmBinary, args...)

	var stdOutBuf, stdErrBuf bytes.Buffer
	cmd.Stdout = &stdOutBuf
	cmd.Stderr = &stdErrBuf

	h.logger.Log(
		"debug", "ran helm command", "name", name,
		"command", cmd.Args,
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
