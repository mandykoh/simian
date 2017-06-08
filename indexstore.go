package simian

type IndexStore interface {
	AddEntry(entry *IndexEntry, node *IndexNode) error
	GetChild(f Fingerprint, parent *IndexNode) (*IndexNode, error)
	GetOrCreateChild(f Fingerprint, parent *IndexNode) (*IndexNode, error)
	GetRoot() (*IndexNode, error)
	RemoveEntries(node *IndexNode) error
}
