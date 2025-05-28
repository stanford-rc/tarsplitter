package tarsplitter

import (
	"archive/tar"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
)

// size of a tar header per https://en.wikipedia.org/wiki/Tar_(computing)
const tarHdrSize = 512

// TarSplitter reads a tar input and writes to one or more split tar files
// based on the maxBytes limit
type TarSplitter struct {
	maxBytes int64
	useGzip  bool

	baseDir  string
	baseName string

	splitNum int

	mu *sync.Mutex
}

// NewTarSplitter initializes a TarSplitter that will write its files to
// {baseDir}/{baseName}_{splitNum}.{tar | tar.gz}, limiting each file to
// maxBytes bytes in size.
func NewTarSplitter(baseDir, baseName string, maxBytes int64, useGzip bool) (*TarSplitter, error) {
	if err := createDir(baseDir); err != nil {
		return nil, err
	}

	return &TarSplitter{
		maxBytes: maxBytes,
		useGzip:  useGzip,

		baseDir:  baseDir,
		baseName: baseName,

		splitNum: 0,

		mu: &sync.Mutex{},
	}, nil
}

// Split the input tar file
func (p *TarSplitter) Split(r io.Reader) (err error) {
	// buf will be used to copy from source to dest
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)

	// tr handles reading tar headers and records
	var tr *tar.Reader
	if p.useGzip {
		gzr, err := gzip.NewReader(r)
		if err != nil {
			return err
		}

		defer gzr.Close()

		tr = tar.NewReader(gzr)
	} else {
		tr = tar.NewReader(r)
	}

	// tw handles writing the tar records to disk, we will initialize it
	// once we've read one tar header from tr
	var tw *TarWriter

	// process loop read the input tar one record at a time
	for {
		// read the next tar record header
		hdr, err := tr.Next()
		if err != nil {
			// on EOF just break out of the processing loop
			if errors.Is(err, io.EOF) {
				break
			}

			// non-EOF error encountered, close tw if it is
			// initialized and return any errors
			if tw != nil {
				return errors.Join(err, tw.Close())
			} else {
				return err
			}
		}

		// we have a record available, if out has not been created yet
		// open the writer and wrap it with tw
		if tw == nil {
			tw, err = p.nextTarWriter()
			if err != nil {
				return err
			}
		}

		// determine the number of bytes to read and create an
		// io.LimitedReader to read that size
		tarRecSize := hdr.FileInfo().Size()

		// if this record would put us over the maxBytes limit, close
		// down this writer and re-open with the next sub-tar writer
		if tw.Written()+(tarHdrSize+tarRecSize) > p.maxBytes {
			if err = tw.Close(); err != nil {
				return err
			}

			if tw, err = p.nextTarWriter(); err != nil {
				return err
			}
		}

		// write the tar header and record to tw
		err = tw.WriteHeader(hdr)
		if err != nil {
			return errors.Join(err, tw.Close())
		}

		_, err = io.CopyBuffer(tw, io.LimitReader(tr, tarRecSize), buf)
		if err != nil {
			return errors.Join(err, tw.Close())
		}
	}

	// close final tar writer
	if tw != nil {
		return tw.Close()
	}

	return nil
}

func (p *TarSplitter) nextTarWriter() (*TarWriter, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	defer func() {
		p.splitNum += 1
	}()

	ext := ""
	if p.useGzip {
		ext = ".gz"
	}

	outputPath := fmt.Sprintf("%s_%06d.tar%s",
		filepath.Join(p.baseDir, p.baseName), p.splitNum, ext)

	fh, err := os.Create(outputPath)
	if err != nil {
		return nil, err
	}

	wc := NewWriteCloser(fh, p.useGzip)

	return NewTarWriter(wc), nil
}
