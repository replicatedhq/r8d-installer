package utils

import (
	"os"
	"path"

	"github.com/pkg/errors"
)

func MoveFile(src, dstDir string) error {
	err := os.MkdirAll(dstDir, 0755)
	if err != nil {
		return errors.Wrapf(err, "failed to create directory %s", dstDir)
	}

	filename := path.Base(src)
	err = os.Rename(src, path.Join(dstDir, filename))
	if err != nil {
		return errors.Wrapf(err, "failed to move file %s to %s", src, dstDir)
	}

	return nil
}
