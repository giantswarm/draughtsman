// +build k8srequired

package charttarball

import (
	"os"
	"path/filepath"

	"github.com/giantswarm/microerror"
	"github.com/mholt/archiver"
)

const (
	fixturesDir = "/e2e/fixtures"
)

// Create creates a tar.gz archive for a directory name relative to
// /e2e/fixtures (the directory passed with --test-dir flag is mounted under
// /e2e path in the test container). It returns created archive path.
//
// NOTE: The created archive should be deleted at the end of the test
// preferably with `defer os.Remove(path)`.
func Create(chartDirName string) (string, error) {
	chartDirPath := filepath.Join(fixturesDir, chartDirName)
	tarballPath := filepath.Join(fixturesDir, "tmp", chartDirName+".tar.gz")

	{
		info, err := os.Stat(chartDirPath)
		if os.IsNotExist(err) {
			return "", microerror.Maskf(executionFailedError, "directory %#q does not exist", chartDirPath)
		} else if err != nil {
			return "", microerror.Mask(err)
		}

		if !info.IsDir() {
			return "", microerror.Maskf(executionFailedError, "file %#q is not a directory", chartDirPath)
		}
	}

	{
		_, err := os.Stat(tarballPath)
		if os.IsNotExist(err) {
			// Fall trough.
		} else if err != nil {
			return "", microerror.Mask(err)
		} else {
			return "", microerror.Maskf(executionFailedError, "file %#q already exists", tarballPath)
		}
	}

	{
		dir := filepath.Dir(tarballPath)
		if dir != "." && dir != "/" {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				return "", microerror.Mask(err)
			}
		}
	}

	{
		err := archiver.TarGz.Make(tarballPath, []string{chartDirPath})
		if err != nil {
			if err != nil {
				return "", microerror.Mask(err)
			}

		}
	}

	return tarballPath, nil
}
