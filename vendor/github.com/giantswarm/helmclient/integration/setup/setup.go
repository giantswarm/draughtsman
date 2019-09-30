// +build k8srequired

package setup

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/giantswarm/helmclient/integration/env"
	"github.com/giantswarm/k8sclient"
	"github.com/giantswarm/microerror"
)

func Setup(m *testing.M, config Config) {
	ctx := context.Background()

	exitCode, err := setup(ctx, m, config)
	if err != nil {
		config.Logger.LogCtx(ctx, "level", "error", "message", "", "stack", fmt.Sprintf("%#v", err))
		os.Exit(1)
	}

	os.Exit(exitCode)
}

func setup(ctx context.Context, m *testing.M, config Config) (int, error) {
	var err error
	teardown := !env.CircleCI() && !env.KeepResources()

	var k8sSetup *k8sclient.Setup
	{
		c := k8sclient.SetupConfig{
			Clients: config.CPK8sClients,
			Logger:  config.Logger,
		}

		k8sSetup, err = k8sclient.NewSetup(c)
		if err != nil {
			return 1, microerror.Mask(err)
		}
	}

	{
		err = k8sSetup.EnsureNamespaceCreated(ctx, tillerNamespace)
		if err != nil {
			return 1, microerror.Mask(err)
		}
		if teardown {
			defer func() {
				err := k8sSetup.EnsureNamespaceDeleted(ctx, tillerNamespace)
				if err != nil {
					config.Logger.LogCtx(ctx, "level", "error", "message", fmt.Sprintf("failed to delete namespace %#q", tillerNamespace), "stack", fmt.Sprintf("%#v", err))
				}
			}()
		}
	}

	return m.Run(), nil
}
