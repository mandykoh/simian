package simian

type IndexStore interface {
	GetNode(path string) (*IndexNode, error)
	GetRoot() (*IndexNode, error)
	SaveNode(n *IndexNode, f Fingerprint) error
}
