package tarsplitter

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

func createDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(path, 0750); err != nil {
			return err
		}
	}
	return nil
}

func SplitTar(destDir, prefix string, inTar io.Reader, useGzip bool, maxSplitSize int64) error {
	err := createDir(destDir)
	if err != nil {
		return fmt.Errorf("failed to create destination dir -- %v", err)
	}

	splitNum := 0
	splitSize := int64(0)
	baseName := prefix
	if ext := filepath.Ext(prefix); ext != "" {
		baseName = prefix[:len(prefix)-len(ext)]
	}
	splitPath := fmt.Sprintf("%s/%s_%06d.tar", destDir, baseName, splitNum)
	splitFile, err := os.Create(splitPath)
	if err != nil {
		return err
	}
	var gzw io.WriteCloser
	if useGzip {
		gzw = gzip.NewWriter(splitFile)
	} else {
		gzw = splitFile
	}
	currSplitWrter := tar.NewWriter(gzw)

	var tarReader *tar.Reader
	if useGzip {
		gzr, err := gzip.NewReader(inTar)
		if err != nil {
			return err
		}
		defer gzr.Close()
		tarReader = tar.NewReader(gzr)
	} else {
		tarReader = tar.NewReader(inTar)
	}
	if err != nil {
		return fmt.Errorf("failed to read tar file -- %v", err)
	}

	copyBuf := make([]byte, 500000000)

	readErr := func() error {
		for {
			header, err := tarReader.Next()
			switch {
			case err == io.EOF:
				return nil
			case err != nil:
				return err
			case header == nil:
				continue
			}

			fileSize := header.FileInfo().Size()
			if fileSize+splitSize > maxSplitSize {
				currSplitWrter.Close()
				if useGzip {
					gzw.Close()
				}
				splitFile.Close()
				splitNum += 1
				splitSize = int64(0)
				splitPath = fmt.Sprintf("%s/%s_%06d.tar", destDir, baseName, splitNum)
				splitFile, err = os.Create(splitPath)
				if err != nil {
					return err
				}
				if useGzip {
					gzw = gzip.NewWriter(splitFile)
				} else {
					gzw = splitFile
				}
				currSplitWrter = tar.NewWriter(gzw)
			}
			if err := currSplitWrter.WriteHeader(header); err != nil {
				return err
			}
			if _, err := io.CopyBuffer(currSplitWrter, tarReader, copyBuf); err != nil {
				return err
			}
			splitSize += fileSize
		}
	}()

	currSplitWrter.Close()
	gzw.Close()
	splitFile.Close()
	return readErr
}

func IsGzip(fn string) (bool, error) {
	f, err := os.Open(fn)
	if err != nil {
		return false, err
	}
	defer f.Close()
	magic := make([]byte, 2)
	_, err = f.Read(magic)
	if err != nil {
		return false, err
	}
	if magic[0] == 0x1f && magic[1] == 0x8b {
		return true, nil
	}
	return false, nil
}
