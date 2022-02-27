package utils

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

// CopyDirToTar adds all files from the source directory to the open archive.
// Assumes that the archive writer lifetime is managed externally.
func CopyDirToTar(tw *tar.Writer, srcDir string) error {

	// Based on the following:
	// https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc0
	return filepath.Walk(srcDir, func(file string, fi os.FileInfo, err error) error {

		if err != nil {
			return errors.Wrapf(err, "error walking dir %s", srcDir)
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return errors.Wrapf(err, "failed to create header for ", fi.Name())
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, srcDir, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return errors.Wrapf(err, "failed to write header for %s", header.Name)
		}

		f, err := os.Open(file)
		if err != nil {
			return errors.Wrapf(err, "failed to open file %s", file)
		}

		if _, err := io.Copy(tw, f); err != nil {
			return errors.Wrapf(err, "failed to copy file %s to tar", file)
		}

		f.Close()

		return nil
	})
}

// Untar unpacks to a given directory given the path to the tarball.
// Based on https://medium.com/@skdomino/taring-untaring-files-in-go-6b07cf56bc07
func Untar(archive, outDir string) error {
	file, err := os.Open(archive)
	if err != nil {
		return errors.Wrapf(err, "couldn't open archive %s", archive)
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return errors.Wrapf(err, "failed to open gzip reader for %s", file.Name())
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		case err == io.EOF:
			return nil

		case err != nil:
			return errors.Wrap(err, "error walking archive")

		case header == nil:
			continue
		}

		target := filepath.Join(outDir, header.Name)

		switch header.Typeflag {

		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return errors.Wrapf(err, "error creating directory %s", target)
				}
			}

		case tar.TypeReg:
			outpath := path.Join(outDir, header.Name)
			if err := os.MkdirAll(path.Dir(outpath), 0755); err != nil {
				return errors.Wrapf(err, "failed to create directory %s", outpath)
			}

			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return errors.Wrapf(err, "error opening file %s", target)
			}

			if _, err := io.Copy(f, tr); err != nil {
				return errors.Wrapf(err, "error copying file %s", target)
			}

			f.Close()
		}
	}
}
