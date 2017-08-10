package simian

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/mandykoh/keva"
)

const nodeFingerprintFile = "fingerprint"
const nodeEntriesDir = "entries"
const thumbnailsDir = "thumbnails"

type DiskIndexStore struct {
	rootPath string
	nodes    *keva.Store
}

func (s *DiskIndexStore) AddEntry(entry *IndexEntry, node *IndexNode) error {
	entriesDir := filepath.Join(node.path, nodeEntriesDir)
	os.Mkdir(entriesDir, os.ModePerm)

	err := entry.save(entriesDir, s.pathForThumbnail(entry))
	if err != nil {
		return err
	}

	node.registerEntry(entry)
	return nil
}

func (s *DiskIndexStore) GetChild(f Fingerprint, parent *IndexNode) (*IndexNode, error) {
	path := s.childPathForFingerprint(f, parent.path)
	return s.getNodeByPath(path)
}

func (s *DiskIndexStore) GetOrCreateChild(f Fingerprint, parent *IndexNode) (*IndexNode, error) {
	fmt.Printf("GetOrCreateChild()\n")
	childPath := s.childPathForFingerprint(f, parent.path)

	node, err := s.getNodeByPath(childPath)
	if err != nil {
		return nil, err
	}

	if node == nil {
		fmt.Printf("Creating child\n")

		node = &IndexNode{
			path: childPath,
			childFingerprintsByString: make(map[string]*Fingerprint),
		}

		err := s.saveNode(node, f)
		if err != nil {
			return nil, err
		}

		parent.registerChild(f)
	}

	return node, nil
}

func (s *DiskIndexStore) GetRoot() (*IndexNode, error) {
	return s.getNodeByPath(path.Join(s.rootPath, "legacy"))
}

func (s *DiskIndexStore) RemoveEntries(node *IndexNode) error {
	entriesDir := filepath.Join(node.path, nodeEntriesDir)
	err := os.RemoveAll(entriesDir)
	if err != nil {
		return err
	}

	node.removeEntries()
	return nil
}

func (s *DiskIndexStore) childPathForFingerprint(f Fingerprint, parentPath string) string {
	fingerprintHash := sha256.Sum256(f.Bytes())
	childDirName := hex.EncodeToString(fingerprintHash[:8])
	return filepath.Join(parentPath, childDirName)
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

func (s *DiskIndexStore) getNodeByPath(path string) (*IndexNode, error) {

	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}

	node := &IndexNode{
		path: path,
		childFingerprintsByString: make(map[string]*Fingerprint),
	}

	err = s.loadAllChildren(node)
	if err != nil {
		return nil, err
	}

	err = s.loadAllEntries(node)
	if err != nil {
		return nil, err
	}

	return node, nil
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

				n.registerChild(child)
			}
		}
	}

	return nil
}

func (s *DiskIndexStore) loadAllEntries(n *IndexNode) error {
	entriesDir := filepath.Join(n.path, nodeEntriesDir)

	dir, err := os.Open(entriesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer dir.Close()

	for fileInfos, err := dir.Readdir(1); err == nil && len(fileInfos) > 0; fileInfos, err = dir.Readdir(1) {
		for _, fileInfo := range fileInfos {
			if strings.HasSuffix(fileInfo.Name(), ".entry") {
				entry, err := NewIndexEntryFromFile(filepath.Join(entriesDir, fileInfo.Name()))
				if err != nil {
					return err
				}
				err = entry.loadThumbnail(s.pathForThumbnail(entry))
				if err != nil {
					return err
				}

				n.registerEntry(entry)
			}
		}
	}

	return nil
}

func (s *DiskIndexStore) loadChild(n *IndexNode, childDirName string) (Fingerprint, error) {
	childPath := filepath.Join(n.path, childDirName)
	return s.fingerprintForChild(childPath)
}

func (s *DiskIndexStore) pathForThumbnail(entry *IndexEntry) string {
	thumbnailHash := sha256.Sum256(entry.MaxFingerprint.Bytes())
	thumbnailHex := hex.EncodeToString(thumbnailHash[:])
	return path.Join(s.rootPath, thumbnailsDir, thumbnailHex[0:2], thumbnailHex[2:4], thumbnailHex[4:])
}

func (s *DiskIndexStore) saveNode(n *IndexNode, f Fingerprint) error {
	fmt.Printf("Saving node %s\n", n.path)

	os.Mkdir(n.path, os.FileMode(0700))

	// Save the actual (non-truncated) fingerprint
	fingerprintFile := filepath.Join(n.path, nodeFingerprintFile)
	file, err := os.Create(fingerprintFile)
	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(f.Bytes())
	return err
}

func NewDiskIndexStore(rootPath string) (*DiskIndexStore, error) {
	legacyNodesDir := path.Join(rootPath, "legacy")
	err := os.MkdirAll(legacyNodesDir, os.FileMode(0700))
	if err != nil {
		return nil, err
	}

	thumbnailsDir := path.Join(rootPath, thumbnailsDir)
	os.MkdirAll(thumbnailsDir, os.FileMode(0700))

	nodeStore, err := keva.NewStore(path.Join(rootPath, "nodes"))
	if err != nil {
		return nil, err
	}

	return &DiskIndexStore{
		rootPath: rootPath,
		nodes:    nodeStore,
	}, nil
}
