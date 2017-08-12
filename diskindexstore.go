package simian

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/mandykoh/keva"
)

const nodeFingerprintFile = "fingerprint"
const nodeEntriesDir = "entries"
const thumbnailsDir = "thumbnails"

type DiskIndexStore struct {
	rootPath string
	nodes    *keva.Store
}

func (s *DiskIndexStore) AddEntry(entry *IndexEntry, node *IndexNode, nodeFingerprint Fingerprint) error {
	err := entry.saveThumbnail(s.pathForThumbnail(entry))
	if err != nil {
		return err
	}

	node.registerEntry(entry)

	fmt.Printf("AddEntry - Saving [%s] %d %d\n", nodeFingerprint.String(), len(node.childFingerprints), len(node.entries))
	return s.nodes.Put(nodeFingerprint.String(), node)
}

func (s *DiskIndexStore) Close() error {
	return s.nodes.Close()
}

func (s *DiskIndexStore) GetChild(f Fingerprint, parent *IndexNode) (*IndexNode, error) {
	var node IndexNode

	err := s.nodes.Get(f.String(), &node)
	if err == keva.ErrValueNotFound {
		return nil, nil

	} else if err == nil {
		err = s.loadThumbnails(&node)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, err
	}

	return &node, nil
}

func (s *DiskIndexStore) GetOrCreateChild(f Fingerprint, parent *IndexNode, parentFingerprint Fingerprint) (*IndexNode, error) {
	fmt.Printf("GetOrCreateChild() %s\n", f.String())

	nodeKey := f.String()

	var node IndexNode
	err := s.nodes.Get(nodeKey, &node)

	if err == keva.ErrValueNotFound {
		fmt.Printf("Creating child\n")

		node = IndexNode{
			childFingerprintsByString: make(map[string]*Fingerprint),
		}

		fmt.Printf("GetOrCreateChild - Saving [%s] %d %d\n", nodeKey, len(node.childFingerprints), len(node.entries))
		err = s.nodes.Put(nodeKey, &node)
		if err != nil {
			return nil, err
		}

		parent.registerChild(f)
		fmt.Printf("GetOrCreateChild - Parent - Saving [%s] %d %d\n", parentFingerprint.String(), len(parent.childFingerprints), len(parent.entries))
		err = s.nodes.Put(parentFingerprint.String(), parent)
		if err != nil {
			return nil, err
		}

	} else if err == nil {
		err = s.loadThumbnails(&node)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, err
	}

	return &node, nil
}

func (s *DiskIndexStore) GetRoot() (*IndexNode, error) {
	var rootKey = Fingerprint{}.String()

	var root IndexNode
	err := s.nodes.Get(rootKey, &root)

	if err == keva.ErrValueNotFound {
		fmt.Printf("Root node not found - creating it\n")
		root = IndexNode{
			childFingerprintsByString: make(map[string]*Fingerprint),
		}

	} else if err == nil {
		fmt.Printf("Found root node with %d children and %d entries\n", len(root.childFingerprints), len(root.entries))

		err = s.loadThumbnails(&root)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, err
	}

	return &root, nil
}

func (s *DiskIndexStore) RemoveEntries(node *IndexNode, nodeFingerprint Fingerprint) error {
	node.removeEntries()
	fmt.Printf("RemoveEntries - Saving [%s] %d %d\n", nodeFingerprint.String(), len(node.childFingerprints), len(node.entries))
	return s.nodes.Put(nodeFingerprint.String(), node)
}

func (s *DiskIndexStore) loadThumbnails(n *IndexNode) error {
	return n.withEachEntry(func(entry *IndexEntry) error {
		return entry.loadThumbnail(s.pathForThumbnail(entry))
	})
}

func (s *DiskIndexStore) pathForThumbnail(entry *IndexEntry) string {
	thumbnailHash := sha256.Sum256(entry.MaxFingerprint.Bytes())
	thumbnailHex := hex.EncodeToString(thumbnailHash[:])
	return path.Join(s.rootPath, thumbnailsDir, thumbnailHex[0:2], thumbnailHex[2:4], thumbnailHex[4:])
}

func NewDiskIndexStore(rootPath string) (*DiskIndexStore, error) {
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
