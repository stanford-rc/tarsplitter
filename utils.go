package tarsplitter

import (
	"os"
	"sync"
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, 4*1024*1024)
	},
}

func createDir(path string) error {
	if _, err := os.Stat(path); err != nil {
		if err := os.MkdirAll(path, 0750); err != nil {
			return err
		}
	}
	return nil
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
