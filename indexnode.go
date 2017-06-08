package simian

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

var errResultLimitReached = errors.New("result limit reached")

type IndexNode struct {
	path                  string
	children              []*IndexNodeHandle
	childrenByFingerprint map[string]*IndexNodeHandle
	entries               []*IndexEntry
}

func (node *IndexNode) Add(entry *IndexEntry, childFingerprintSize int, index *Index) (*IndexNode, error) {

	fmt.Printf("Node[%s] Add %d\n", node.path, childFingerprintSize)

	entryFingerprint := entry.FingerprintForSize(childFingerprintSize)

	if len(node.children) == 0 {

		// We can go deeper and this new entry is sufficiently different to
		// the rest, so split this leaf node by turning entries into children.
		fmt.Printf("Max Diff: %f\n", node.maxChildDifferenceTo(entry.MaxFingerprint))
		if childFingerprintSize < index.maxFingerprintSize && node.maxChildDifferenceTo(entry.MaxFingerprint) > index.maxEntryDifference {
			fmt.Printf("Pushing entries to children\n")
			node.pushEntriesToChildren(childFingerprintSize, index.Store)

		} else {
			fmt.Printf("Adding entry %s\n", node.path)
			err := index.Store.AddEntry(entry, node)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Added entry\n")
			return node, nil
		}
	}

	child, err := index.Store.GetOrCreateChild(entryFingerprint, node)
	if err != nil {
		return nil, err
	}

	return child.Add(entry, childFingerprintSize+1, index)
}

func (node *IndexNode) FindNearest(entry *IndexEntry, childFingerprintSize int, index *Index, maxResults int, maxDifference float64) ([]*IndexEntry, error) {
	results := make([]*IndexEntry, 0, maxResults)

	err := node.gatherNearest(entry, childFingerprintSize, index, maxDifference, &results)
	if err != errResultLimitReached {
		return nil, err
	}

	return results, nil
}

func (node *IndexNode) addSimilarEntriesTo(entries *[]*IndexEntry, fingerprint Fingerprint, maxDifference float64) error {
	fmt.Printf("addSimilarEntriesTo\n")

	return node.withEachEntry(func(entry *IndexEntry) error {
		if len(*entries) >= cap(*entries) {
			fmt.Printf("Max results hit\n")
			return errResultLimitReached
		}

		diff := entry.MaxFingerprint.Difference(fingerprint)
		if diff <= maxDifference {
			fmt.Printf("Found %d of difference %f\n", len(*entries), diff)
			*entries = append(*entries, entry)
		} else {
			fmt.Printf("Max difference hit at %f\n", diff)
			return errResultLimitReached
		}

		return nil
	})
}

func (node *IndexNode) gatherNearest(entry *IndexEntry, childFingerprintSize int, index *Index, maxDifference float64, results *[]*IndexEntry) error {

	fmt.Printf("%d gatherNearest\n", childFingerprintSize)

	// Check for an exact matching child
	entryFingerprint := entry.FingerprintForSize(childFingerprintSize)
	exactChildHandle := node.childrenByFingerprint[entryFingerprint.String()]

	var exactChild *IndexNode
	if exactChildHandle != nil {
		var err error
		exactChild, err = index.Store.GetChild(entryFingerprint, node)
		if err != nil {
			return err
		}
	}

	// One exists - recursively search it
	if exactChild != nil {
		err := exactChild.gatherNearest(entry, childFingerprintSize+1, index, maxDifference, results)
		if err != nil {
			return err
		}

		err = exactChild.addSimilarEntriesTo(results, entry.MaxFingerprint, maxDifference)
		if err != nil {
			return err
		}
	}

	children := make([]*IndexNodeHandle, len(node.children))
	copy(children, node.children)

	// Need more results - find and sort all children by nearness
	sort.Sort(nodesByDifferenceToFingerprintWith(children, entryFingerprint))

	// fmt.Printf("Sorting %d children...\n", len(children))
	// for i, child := range children {
	// 	diff := child.fingerprint.Difference(entryFingerprint)
	// 	fmt.Printf("%d sorted child %d of %f (%d %d)\n", childFingerprintSize+1, i, diff, len(child.fingerprint.samples), len(entryFingerprint.samples))
	// }

	// Recursively gather from nearest children
	for i, child := range children {
		fmt.Printf("Visiting child %d\n", i)
		if exactChild != nil && child.Path == exactChild.path {
			continue
		}

		childNode, err := index.Store.GetChild(child.Fingerprint, node)
		if err != nil {
			return err
		}

		err = childNode.gatherNearest(entry, childFingerprintSize+1, index, maxDifference, results)
		if err != nil {
			return err
		}

		err = childNode.addSimilarEntriesTo(results, entry.MaxFingerprint, maxDifference)
		if err != nil {
			return err
		}
	}

	return nil
}

func (node *IndexNode) maxChildDifferenceTo(f Fingerprint) float64 {
	maxDifference := 0.0

	node.withEachEntry(func(entry *IndexEntry) error {
		diff := entry.MaxFingerprint.Difference(f)
		maxDifference = math.Max(diff, maxDifference)
		return nil
	})

	return maxDifference
}

func (node *IndexNode) pushEntriesToChildren(childFingerprintSize int, store IndexStore) error {
	node.withEachEntry(func(entry *IndexEntry) error {
		entryFingerprint := entry.FingerprintForSize(childFingerprintSize)
		child, err := store.GetOrCreateChild(entryFingerprint, node)
		if err != nil {
			return err
		}
		store.AddEntry(entry, child)
		return nil
	})

	return store.RemoveEntries(node)
}

func (node *IndexNode) registerChild(child *IndexNode, f Fingerprint) {
	childHandle := &IndexNodeHandle{Path: child.path, Fingerprint: f}
	node.registerChildByHandle(childHandle)
}

func (node *IndexNode) registerChildByHandle(childHandle *IndexNodeHandle) {
	node.children = append(node.children, childHandle)
	node.childrenByFingerprint[childHandle.Fingerprint.String()] = childHandle
}

func (node *IndexNode) registerEntry(entry *IndexEntry) {
	node.entries = append(node.entries, entry)
}

func (node *IndexNode) removeEntries() {
	node.entries = nil
}

func (node *IndexNode) withEachEntry(action func(*IndexEntry) error) error {
	for _, entry := range node.entries {
		err := action(entry)
		if err != nil {
			return err
		}
	}

	return nil
}

type nodesByDifferenceToFingerprint struct {
	nodes       []*IndexNodeHandle
	differences []float64
}

func (sorter *nodesByDifferenceToFingerprint) Len() int {
	return len(sorter.nodes)
}

func (sorter *nodesByDifferenceToFingerprint) Less(i, j int) bool {
	return sorter.differences[i] < sorter.differences[j]
}

func (sorter *nodesByDifferenceToFingerprint) Swap(i, j int) {
	tmpNode := sorter.nodes[i]
	sorter.nodes[i] = sorter.nodes[j]
	sorter.nodes[j] = tmpNode

	tmpDiff := sorter.differences[i]
	sorter.differences[i] = sorter.differences[j]
	sorter.differences[j] = tmpDiff
}

func nodesByDifferenceToFingerprintWith(nodes []*IndexNodeHandle, f Fingerprint) *nodesByDifferenceToFingerprint {
	differences := make([]float64, len(nodes), len(nodes))
	for i, n := range nodes {
		differences[i] = n.Fingerprint.Difference(f)
	}

	return &nodesByDifferenceToFingerprint{nodes: nodes, differences: differences}
}
