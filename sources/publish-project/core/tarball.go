package core

import (
	"strings"
	"time"
	"io/ioutil"
	"archive/tar"
	"compress/gzip"
	"bytes"
	"fmt"
	"path/filepath"
	"os"
)

// CreateTarball scan source directory and create .tar.gz
func CreateTarball(source string, dockerfile []byte) ([]byte, error) {
	tarBuffer := new(bytes.Buffer)
	gzWriter := gzip.NewWriter(tarBuffer)
	tarWriter := tar.NewWriter(gzWriter)

	// add Dockerfile
	header := &tar.Header{
		Name: "Dockerfile",
		Size: int64(len(dockerfile)),
		ModTime: time.Now(),
	}
	if err := tarWriter.WriteHeader(header); err != nil {
		return nil, fmt.Errorf("could not write into tarball header for Dockerfile: %v", err)
	}
	if _, err := tarWriter.Write(dockerfile); err != nil {
		return nil, fmt.Errorf("could not write into tarball for Dockerfile: %v", err)
	}
	// end add Dockerfile
	
	sourceLen := len(source)
	var ret = filepath.Walk(source, func(path string, f os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("error walk dir%s: %v", path, err)
		}

		if f.IsDir() {
			return nil
		}

		filename := path[sourceLen+1:]
		fp, err := os.Open(path)
		if err != nil {
			return fmt.Errorf("can not open source file %s: %v", path, err)	
		}
		defer fp.Close()

		dat, err := ioutil.ReadAll(fp)
		if err != nil {
			return fmt.Errorf("can not read file %s: %v", source, err)
		}
		datalen := int64(len(dat))
		if datalen != f.Size() {
			return fmt.Errorf("size of readed file: %d <> fs file size: %d", datalen, f.Size())
		}

		header := &tar.Header{
			Name: strings.Replace(filename,"\\","/",-1),
			Size: datalen,
			ModTime: f.ModTime(),
		}
		if err := tarWriter.WriteHeader(header); err != nil {
			return fmt.Errorf("could not write into tarball header for file `%s`: %v", filename, err)
		}
		if _, err := tarWriter.Write(dat); err != nil {
			return fmt.Errorf("could not write into tarball for file `%s`: %v", filename, err)
		}
	
		return nil
	})

	if ret == nil {
		tarWriter.Close()
		gzWriter.Close()
		return tarBuffer.Bytes(), nil
	}

	return nil, ret
}
