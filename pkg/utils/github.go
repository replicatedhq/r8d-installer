package utils

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/google/go-github/v42/github"
	"github.com/pkg/errors"
) // with go modules enabled (GO111MODULE=on or outside GOPATH)

var client *github.Client

func init() {
	client = github.NewClient(nil)
}

// DownloadAssetFromGithubRelease returns a handle to the asset file for the given owner/repo and release tag.
// It's assume the caller owns the temporary file handle returned.
// TODO (dans): Add checksum validation
func DownloadAssetFromGithubRelease(owner, repo, tag, assetName string) (string, error) {

	release, _, err := client.Repositories.GetReleaseByTag(context.Background(), owner, repo, tag)
	if err != nil {
		return "", errors.Wrap(err, "failed to get release")
	} else if release == nil {
		return "", errors.Errorf("release %s/%s:%s not found", owner, repo, tag)
	}

	var id int64
	for asset := range release.Assets {
		if *release.Assets[asset].Name == assetName {
			id = *release.Assets[asset].ID
			break
		}
	}

	if id == 0 {
		return "", errors.Errorf("asset %s not found in release %s/%s:%s", assetName, owner, repo, tag)
	}

	file, err := os.Create(path.Join(os.TempDir(), assetName))
	if err != nil {
		return "", errors.Wrap(err, "failed to open temp file for download")
	}
	defer file.Close()

	readCloser, _, err := client.Repositories.DownloadReleaseAsset(context.Background(), owner, repo, id, http.DefaultClient)
	if err != nil {
		return "", errors.Wrap(err, "failed to get release asset")
	}

	_, err = io.Copy(file, readCloser)
	if err != nil {
		return "", errors.Wrap(err, "failed copy asset to temp file")
	}

	return file.Name(), nil
}
