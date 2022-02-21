package deps

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/replicatedhq/r8d-installer/pkg/component"
	"github.com/replicatedhq/r8d-installer/pkg/component/rke2"
	"github.com/replicatedhq/r8d-installer/pkg/utils"
)

const assetPath = "./pkg/component/%s/assets/%s"

type Manifest struct {
	RKE2 string
}

func Build(manifest Manifest) error {
	log.Default()

	components, err := convertManifestToComponents(manifest)
	if err != nil {
		return errors.Wrap(err, "failed to convert manifest to components")
	}

	for _, buildable := range components {

		err = processBinaries(buildable)
		if err != nil {
			return errors.Wrapf(err, "failed to process binaries for %s", buildable.GetName())
		}

		_, err = buildable.GetManifest()
		if err != nil {
			return errors.Wrapf(err, "failed to get manifest for %s", buildable.GetName())
		}
		// copy to assets

		_, err = buildable.GetImageArchive()
		if err != nil {
			return errors.Wrapf(err, "failed to get manifest for %s", buildable.GetName())
		}
		// copy to assets

	}

	return nil

	// create the Components
	// iterate through Components and build them
	// 		get the binaries
	// 		generate the manifests (using binaries as necessary)
	// 		generate the images bundles. Download, archive and compress

}

func Update(old Manifest) error {
	return nil
}

// You might see an error here in your editor. Make sure your editor has the `deps` build tag enabled.
// https://www.ryanchapin.com/configuring-vscode-to-use-build-tags-in-golang-to-separate-integration-and-unit-test-code/
func convertManifestToComponents(manifest Manifest) ([]component.Builder, error) {
	components := []component.Builder{}

	if manifest.RKE2 == "" {
		return nil, errors.New("no RKE2 version specified")
	}
	components = append(components, rke2.New(manifest.RKE2))

	// TODO (dans): other components here

	return components, nil
}

func processBinaries(buildable component.Builder) error {

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
