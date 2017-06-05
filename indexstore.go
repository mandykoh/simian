package simian

type IndexStore interface {
	GetNode(path string, f Fingerprint) (*IndexNode, error)
	GetRoot() *IndexNode
	SaveNode(n *IndexNode, f Fingerprint) error
}
