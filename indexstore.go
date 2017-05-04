package simian

type IndexStore interface {
	GetNode(h *IndexNodeHandle) (*IndexNode, error)
	GetRoot() *IndexNode
	SaveNode(n *IndexNode, f Fingerprint) error
}
