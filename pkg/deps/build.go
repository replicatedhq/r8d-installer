package deps

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/replicatedhq/r8d-installer/pkg/component"
	"github.com/replicatedhq/r8d-installer/pkg/component/kots"
	"github.com/replicatedhq/r8d-installer/pkg/utils"
)

// You might see an error here in your editor. Make sure your editor has the `deps` build tag enabled.
// https://www.ryanchapin.com/configuring-vscode-to-use-build-tags-in-golang-to-separate-integration-and-unit-test-code/
func convertManifestToComponents(manifest Manifest) ([]component.Builder, error) {
	components := []component.Builder{}

	if manifest.RKE2 == "" {
		return nil, errors.New("no RKE2 version specified")
	}
	// components = append(components, rke2.New(manifest.RKE2))

	if manifest.KOTS == "" {
		return nil, errors.New("no KOTS version specified")
	}
	components = append(components, kots.New(manifest.KOTS))

	// TODO (dans): other components here

	return components, nil
}

func processBinaries(logger *log.Logger, buildable component.Builder) error {
	logger.Printf("└── processing binaries for %s", buildable.GetName())
	paths, err := buildable.GetBinaries()
	if err != nil {
		return errors.Wrap(err, "failed to get binaries")
	}

	dst := fmt.Sprintf(assetPath, buildable.GetName(), "binaries")

	for _, src := range paths {
		err = utils.MoveFile(src, dst)
		if err != nil {
			return errors.Wrapf(err, "failed to move %s to %s", src, dst)
		}
	}

	return nil
}

func processImagesArchive(logger *log.Logger, buildable component.Builder) error {
	logger.Printf("└── processing image archive for %s", buildable.GetName())
	src, err := buildable.GetImageArchive()
	if err != nil {
		return errors.Wrap(err, "failed to get image archive")
	}

	if src == "" {
		return nil
	}

	dst := fmt.Sprintf(assetPath, buildable.GetName(), "images")

	err = utils.MoveFile(src, dst)
	if err != nil {
		return errors.Wrapf(err, "failed to move %s to %s", src, dst)
	}

	return nil
}

func processManifests(logger *log.Logger, buildable component.Builder) error {
	logger.Printf("└── processing manifests for %s", buildable.GetName())
	yaml, err := buildable.GetManifests()
	if err != nil {
		return errors.Wrap(err, "failed to get manifests")
	}

	if yaml == "" {
		return nil
	}

	dir := fmt.Sprintf(assetPath, buildable.GetName(), "manifests")
	err = os.MkdirAll(dir, 0755)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dir)
	}

	dst := path.Join(dir, buildable.GetName()+".yaml")

	err = os.WriteFile(dst, []byte(yaml), 0644)
	if err != nil {
		return errors.Wrapf(err, "failed to write manifests file %s", dst)
	}

	return nil
}
