package simian

type IndexStore interface {
	GetNode(h *IndexNodeHandle) (*IndexNode, error)
	SaveNode(n *IndexNode, f Fingerprint) error
}
