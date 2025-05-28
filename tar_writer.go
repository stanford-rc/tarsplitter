package tarsplitter

import (
	"archive/tar"
	"errors"
)

type TarWriter struct {
	wc *WriteCloser
	tw *tar.Writer
}

func NewTarWriter(wc *WriteCloser) *TarWriter {
	return &TarWriter{
		wc: wc,
		tw: tar.NewWriter(wc),
	}
}

func (p *TarWriter) WriteHeader(hdr *tar.Header) error {
	return p.tw.WriteHeader(hdr)
}

func (p *TarWriter) Write(buf []byte) (int, error) {
	return p.tw.Write(buf)
}

func (p *TarWriter) Close() error {
	err := errors.Join(p.tw.Close(), p.wc.Close())
	if err != nil {
		err = errors.Join(err, p.wc.Delete())
	}
	return err
}

func (p *TarWriter) Written() int64 {
	return p.wc.Written()
}
