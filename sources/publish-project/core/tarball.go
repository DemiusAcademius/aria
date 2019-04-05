package core

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CreateTarball scan source directory and create .tar.gz
func CreateTarball(source string, dockerfile []byte) ([]byte, error) {
	var fileinfo, err = os.Stat(source)
	if os.IsNotExist(err) {
		return nil, fmt.Errorf("source `%s` does not exists", source)
	}

	tarBuffer := new(bytes.Buffer)
	gzWriter := gzip.NewWriter(tarBuffer)
	tarWriter := tar.NewWriter(gzWriter)

	// add Dockerfile
	header := &tar.Header{
		Name:    "Dockerfile",
		Size:    int64(len(dockerfile)),
		ModTime: time.Now(),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("could not write into tarball header for Dockerfile: %v", err)
	}
	if _, err := tarWriter.Write(dockerfile); err != nil {
		return nil, fmt.Errorf("could not write into tarball for Dockerfile: %v", err)
	}
	// end add Dockerfile

	var ret error

	if fileinfo.IsDir() {
		// add to tar.gz directory
		sourceLen := len(source)
		ret = filepath.Walk(source, func(path string, f os.FileInfo, err error) error {
			if err != nil {
				return fmt.Errorf("error walk dir%s: %v", path, err)
			}

			if f.IsDir() {
				return nil
			}

			filename := path[sourceLen+1:]

			if err = addFileToTarball(filename, path, f, tarWriter); err != nil {
				return err
			}

			return nil
		})
	} else {
		ret = addFileToTarball(filepath.Base(source), source, fileinfo, tarWriter)
	}

	tarWriter.Close()
	gzWriter.Close()

	if ret == nil {
		return tarBuffer.Bytes(), nil
	}

	return nil, ret
}

func addFileToTarball(filename, path string, f os.FileInfo, tarWriter *tar.Writer) error {
	fp, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("can not open source file %s: %v", path, err)
	}
	defer fp.Close()

	dat, err := ioutil.ReadAll(fp)
	if err != nil {
		return fmt.Errorf("can not read file %s: %v", path, err)
	}
	datalen := int64(len(dat))
	if datalen != f.Size() {
		return fmt.Errorf("size of readed file: %d <> fs file size: %d", datalen, f.Size())
	}

	header := &tar.Header{
		Name:    strings.Replace(filename, "\\", "/", -1),
		Size:    datalen,
		ModTime: f.ModTime(),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return fmt.Errorf("could not write into tarball header for file `%s`: %v", filename, err)
	}
	if _, err := tarWriter.Write(dat); err != nil {
		return fmt.Errorf("could not write into tarball for file `%s`: %v", filename, err)
	}

	return nil
}