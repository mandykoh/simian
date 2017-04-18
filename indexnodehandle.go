package simian

type indexNodeHandle struct {
	path        string
	fingerprint Fingerprint
}

func (h *indexNodeHandle) Node() *IndexNode {
	return &IndexNode{path: h.path}
}
