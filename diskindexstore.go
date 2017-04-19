package simian

import (
	"os"
	"path/filepath"
)

const nodeFingerprintFile = "fingerprint"

type DiskIndexStore struct {
}

func (s *DiskIndexStore) GetNode(h *IndexNodeHandle) (*IndexNode, error) {
	return &IndexNode{path: h.Path}, nil
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
