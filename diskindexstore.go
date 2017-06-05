package simian

import (
	"os"
	"path/filepath"
)

const nodeFingerprintFile = "fingerprint"

type DiskIndexStore struct {
	RootPath string
}

func (s *DiskIndexStore) GetNode(path string, f Fingerprint) (*IndexNode, error) {
	return &IndexNode{path: path}, nil
}

func (s *DiskIndexStore) GetRoot() *IndexNode {
	return &IndexNode{path: s.RootPath}
}

func (s *DiskIndexStore) SaveNode(n *IndexNode, f Fingerprint) (err error) {

	os.Mkdir(n.path, os.FileMode(0522))

	// Save the actual (non-truncated) fingerprint
	fingerprintFile := filepath.Join(n.path, nodeFingerprintFile)
	file, err := os.Create(fingerprintFile)
	if err != nil {
		return
	}

	defer file.Close()

	_, err = file.Write(f.Bytes())
	return
}
