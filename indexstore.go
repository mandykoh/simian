package simian

type IndexStore interface {
	AddEntry(entry *IndexEntry, node *IndexNode, nodeFingerprint Fingerprint) error
	Close() error
	GetChild(f Fingerprint, parent *IndexNode) (*IndexNode, error)
	GetOrCreateChild(f Fingerprint, parent *IndexNode, parentFingerprint Fingerprint) (*IndexNode, error)
	GetRoot() (*IndexNode, error)
	RemoveEntries(node *IndexNode, nodeFingerprint Fingerprint) error
}
