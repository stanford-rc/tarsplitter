package tarsplitter

import (
	"compress/gzip"
	"errors"
	"os"
)

type WriteCloser struct {
	fh      *os.File
	gz      *gzip.Writer
	written int64
}

func NewWriteCloser(fh *os.File, useGzip bool) *WriteCloser {
	p := &WriteCloser{
		fh: fh,
	}

	if useGzip {
		p.gz = gzip.NewWriter(p.fh)
	}

	return p
}

func (p *WriteCloser) Write(buf []byte) (n int, err error) {
	if p.gz != nil {
		n, err = p.gz.Write(buf)
	} else {
		n, err = p.fh.Write(buf)
	}

	p.written += int64(n)

	return n, err
}

func (p *WriteCloser) Written() int64 {
	return p.written
}

func (p *WriteCloser) Close() error {
	var errs []error

	if p.gz != nil {
		errs = append(errs, p.gz.Close())
	}

	errs = append(errs, p.fh.Close())

	return errors.Join(errs...)
}

func (p *WriteCloser) Delete() error {
	return os.Remove(p.fh.Name())
}
