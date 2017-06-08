package simian

type IndexStore interface {
	GetChild(f Fingerprint, parent *IndexNode) (*IndexNode, error)
	GetOrCreateChild(f Fingerprint, parent *IndexNode) (*IndexNode, error)
	GetRoot() (*IndexNode, error)
}
