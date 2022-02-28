// +build deps

package openebs

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/pkg/errors"
	"github.com/replicatedhq/r8d-installer/pkg/utils"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"sigs.k8s.io/yaml"
)

// GetName returns "openebs"
func (o *OpenEBS) GetName() string {
	return "openebs"
}

// GetVersion returns the OpenEBS version
func (o *OpenEBS) GetVersion() string {
	return o.Version
}

// GetManifest merges all the YAML files for the OpenEBS "lite" configuration into a single YAML.
// Based on the Local PV Hostpath Install: https://openebs.io/docs/user-guides/localpv-hostpath#install
func (o *OpenEBS) GetManifests() (string, error) {
	content, err := getOpenEBSManifests(o.GetVersion())
	if err != nil {
		return "", errors.Wrap(err, "failed to get openebs manifests")
	}

	path := path.Join(os.TempDir(), "openebs.yaml")
	file, err := os.Create(path)
	if err != nil {
		return "", errors.Wrapf(err, "failed to create openebs manifests file %s", path)
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return "", errors.Wrapf(err, "failed to write openebs manifests to file %s", path)
	}

	return path, nil
}

// GetManifest returns a file path to the compressed airgap images for openebs.
// It's assume that the caller owns the file after calling.
// It pulls the images from the manifest files and statically adds the openebs linux-tools
// image. Since the manifests are not tagged, an error is thrown if we don't get back
// a tag we are expecting.
func (o *OpenEBS) GetImageArchive() (string, error) {

	content, err := getOpenEBSManifests(o.GetVersion())
	if err != nil {
		return "", errors.Wrap(err, "failed to get openebs manifests")
	}

	images := utils.GetImages(content)

	// Hack: this is the only thing to do to make sure the manifest didn't
	// change on the `gh-pages` branch.
	if !imageExists(images, fmt.Sprintf("openebs/provisioner-localpv:%s", o.GetVersion())) {
		return "", errors.New(fmt.Sprintf("Manifest tags do not match the desired version %s."+
			"Likely a new manifest was published for a new version of openebs."+
			"Try updating.", o.GetVersion()))
	}

	images = append(images, fmt.Sprintf("openebs/linux-utils:%s", o.GetVersion()))

	archivePath, err := utils.CreateArchive(o.GetName(), images)
	if err != nil {
		return "", errors.Wrap(err, "failed to create image archive")
	}

	// Rename the generic archive
	destPath := path.Join(os.TempDir(), "openebs-image-archive.tar.zst")
	if err = os.Rename(archivePath, destPath); err != nil {
		return "", errors.Wrap(err, "failed to rename openebs image archive")
	}

	return destPath, nil
}

// GetBinaries returns nil as OpenEBS does not have any associated executables1.
func (o *OpenEBS) GetBinaries() ([]string, error) {
	return nil, nil
}

func getOpenEBSManifests(version string) (string, error) {
	var content strings.Builder

	storageClassContent, err := utils.GetSourceFileFromGithubRelease("openebs", "charts", "gh-pages", "openebs-lite-sc.yaml")
	if err != nil {
		return "", errors.Wrap(err, "failed to get openebs storage class manifest from github")
	}

	storageClassContent, err = addDefaultStorageClass(storageClassContent)
	content.WriteString(storageClassContent)

	content.WriteString("---\n")

	operatorContent, err := utils.GetSourceFileFromGithubRelease("openebs", "charts", "gh-pages", "openebs-operator-lite.yaml")
	if err != nil {
		return "", errors.Wrap(err, "failed to get openebs operator manifest from github")
	}
	content.WriteString(operatorContent)

	// TODO (dans): since we can't grab the manifests by tag, verify the manifest have the expected images tags and bail if not

	return content.String(), nil
}

func addDefaultStorageClass(content string) (string, error) {

	var manifest unstructured.Unstructured
	if err := yaml.Unmarshal([]byte(content), &manifest); err != nil {
		return "", errors.Wrap(err, "failed to unmarshal openebs storage class manifest")
	}

	annotations := manifest.GetAnnotations()
	annotations["storageclass.kubernetes.io/is-default-class"] = "true"
	manifest.SetAnnotations(annotations)

	contentBuffer, err := yaml.Marshal(manifest.Object)
	if err != nil {
		return "", errors.Wrap(err, "failed to marshal storage class manifest")
	}

	return string(contentBuffer), nil
}

// imageExists checks if the image exists in the list of images
func imageExists(images []string, image string) bool {
	for _, i := range images {
		if i == image {
			return true
		}
	}
	return false
}
