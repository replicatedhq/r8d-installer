package utils

import (
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// MoveFile moves a file from src to the destination directory, keeping the
// same file name.
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

// MergeYAML creates the destination file and merges any YAML in the
// source directory into a multi-document YAML file. The file is closed
// for writing on success. It supports nested directories.
func MergeYAML(srcDir string, dst string) error {
	mergedYAML, err := os.Create(dst)
	if err != nil {
		return errors.Wrapf(err, "failed to create file %s", dst)
	}
	defer mergedYAML.Close()

	err = filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to walk %s", srcDir)
		}

		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(info.Name(), ".yaml") {
			srcFile, err := os.Open(path)
			if err != nil {
				return errors.Wrapf(err, "failed to open file %s", path)
			}

			fmt.Fprintf(mergedYAML, "---\n")

			_, err = io.Copy(mergedYAML, srcFile)
			if err != nil {
				return errors.Wrapf(err, "failed to copy file %s to %s", path, dst)
			}
		}

		return nil
	})
	if err != nil {
		return errors.Wrapf(err, "failed to walk %s", srcDir)
	}
	return nil
}
