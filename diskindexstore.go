package simian

import (
	"os"
	"path/filepath"
)

const nodeFingerprintFile = "fingerprint"

type DiskIndexStore struct {
	RootPath string
}

func (s *DiskIndexStore) GetNode(path string) (*IndexNode, error) {
	node := &IndexNode{
		path: path,
		childrenByFingerprint: make(map[string]*IndexNodeHandle),
	}

	err := s.loadAllChildren(node)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (s *DiskIndexStore) GetRoot() (*IndexNode, error) {
	return s.GetNode(s.RootPath)
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

func (s *DiskIndexStore) loadAllChildren(n *IndexNode) error {
	dir, err := os.Open(n.path)
	if err != nil {
		return err
	}
	defer dir.Close()

	for fileInfos, err := dir.Readdir(1); err == nil && len(fileInfos) > 0; fileInfos, err = dir.Readdir(1) {
		for _, info := range fileInfos {
			if info.IsDir() && info.Name() != nodeEntriesDir {
				child, err := n.loadChild(info.Name())
				if err != nil {
					return err
				}

				n.registerChildByHandle(child)
			}
		}
	}

	return nil
}
