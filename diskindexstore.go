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

func (s *DiskIndexStore) fingerprintForChild(childPath string) (Fingerprint, error) {
	childFingerprint := Fingerprint{}
	childFingerprintFile := filepath.Join(childPath, nodeFingerprintFile)

	file, err := os.Open(childFingerprintFile)
	if err != nil {
		return childFingerprint, err
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		return childFingerprint, err
	}

	fingerprintBytes := make([]byte, fileInfo.Size(), fileInfo.Size())
	_, err = file.Read(fingerprintBytes)
	if err != nil {
		return childFingerprint, err
	}

	childFingerprint.UnmarshalBytes(fingerprintBytes)

	return childFingerprint, nil
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
				child, err := s.loadChild(n, info.Name())
				if err != nil {
					return err
				}

				n.registerChildByHandle(child)
			}
		}
	}

	return nil
}

func (s *DiskIndexStore) loadChild(n *IndexNode, childDirName string) (*IndexNodeHandle, error) {
	childPath := filepath.Join(n.path, childDirName)
	childFingerprint, err := s.fingerprintForChild(childPath)
	if err != nil {
		return nil, err
	}

	return &IndexNodeHandle{Path: childPath, Fingerprint: childFingerprint}, nil
}
