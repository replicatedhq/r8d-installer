package deps

import (
	"log"

	"github.com/pkg/errors"
)

const assetPath = "./pkg/component/%s/assets/%s"

type Manifest struct {
	RKE2    string
	KOTS    string
	OpenEBS string
}

func Build(manifest Manifest) error {
	logger := log.Default()

	components, err := convertManifestToComponents(manifest)
	if err != nil {
		return errors.Wrap(err, "failed to convert manifest to components")
	}

	for _, buildable := range components {
		logger.Printf("building assets for %s", buildable.GetName())

		err = processBinaries(logger, buildable)
		if err != nil {
			return errors.Wrapf(err, "failed to process binaries for %s", buildable.GetName())
		}

		err = processImagesArchive(logger, buildable)
		if err != nil {
			return errors.Wrapf(err, "failed to get image archive for %s", buildable.GetName())
		}

		err = processManifests(logger, buildable)
		if err != nil {
			return errors.Wrapf(err, "failed to get manifests for %s", buildable.GetName())
		}
	}

	return nil
}

func Update(old Manifest) error {
	return nil
}
