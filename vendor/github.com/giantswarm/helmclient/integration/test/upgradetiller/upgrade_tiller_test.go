// +build k8srequired

package upgradetiller

import (
	"context"
	"testing"
	"time"

	"github.com/giantswarm/backoff"
	"github.com/giantswarm/microerror"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUpgradeTiller(t *testing.T) {
	ctx := context.Background()

	currentTillerImage := "quay.io/giantswarm/tiller:v2.14.3"

	labelSelector := "app=helm,name=tiller"
	tillerNamespace := "giantswarm"

	// Upgrade tiller to the current image.
	{
		err := config.HelmClient.EnsureTillerInstalled(ctx)
		if err != nil {
			t.Fatalf("could not install tiller %#v", err)
		}

		tillerImage, err := getTillerImage(ctx, tillerNamespace, labelSelector)
		if err != nil {
			t.Fatalf("could not get tiller image %#v", err)
		}
		if tillerImage != currentTillerImage {
			t.Fatalf("tiller image == %#q, want %#q", tillerImage, currentTillerImage)
		}
	}
}

func getTillerDeployment(ctx context.Context, namespace string, labelSelector string) (*appsv1.Deployment, error) {
	var d *appsv1.Deployment
	{
		o := func() error {
			lo := metav1.ListOptions{
				LabelSelector: labelSelector,
			}
			l, err := config.CPK8sClients.K8sClient().AppsV1().Deployments(namespace).List(lo)
			if err != nil {
				return microerror.Mask(err)
			}

			if len(l.Items) != 1 {
				return microerror.Maskf(executionFailedError, "cannot get deployment for %#q %#q found %d, want 1", namespace, labelSelector, len(l.Items))
			}

			d = &l.Items[0]
			if d.Status.AvailableReplicas != 1 && d.Status.ReadyReplicas != 1 {
				return microerror.Maskf(executionFailedError, "tiller deployment not ready %d available %d ready, want 1", d.Status.AvailableReplicas, d.Status.ReadyReplicas)
			}

			return nil
		}

		b := backoff.NewExponential(2*time.Minute, 5*time.Second)
		n := backoff.NewNotifier(config.Logger, ctx)

		err := backoff.RetryNotify(o, b, n)
		if err != nil {
			return nil, microerror.Mask(err)
		}
	}

	return d, nil
}

func getTillerImage(ctx context.Context, namespace, labelSelector string) (string, error) {
	d, err := getTillerDeployment(ctx, namespace, labelSelector)
	if err != nil {
		return "", microerror.Mask(err)
	}

	if len(d.Spec.Template.Spec.Containers) != 1 {
		return "", microerror.Maskf(executionFailedError, "Spec.Template.Spec.Containers == %d, want 1", len(d.Spec.Template.Spec.Containers))
	}

	tillerImage := d.Spec.Template.Spec.Containers[0].Image
	if tillerImage == "" {
		return "", microerror.Maskf(executionFailedError, "tiller image is empty")
	}

	return tillerImage, nil
}

func updateTillerImage(ctx context.Context, namespace, labelSelector, tillerImage string) error {
	d, err := getTillerDeployment(ctx, namespace, labelSelector)
	if err != nil {
		return microerror.Mask(err)
	}

	d.Spec.Template.Spec.Containers[0].Image = tillerImage
	_, err = config.CPK8sClients.K8sClient().AppsV1().Deployments(namespace).Update(d)
	if err != nil {
		return microerror.Mask(err)
	}

	return nil
}
