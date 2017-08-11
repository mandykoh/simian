package simian

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
)

var errResultLimitReached = errors.New("result limit reached")

type IndexNode struct {
	childFingerprints         []Fingerprint
	childFingerprintsByString map[string]*Fingerprint
	entries                   []*IndexEntry
}

func (node *IndexNode) Add(entry *IndexEntry, nodeFingerprint Fingerprint, childFingerprintSize int, index *Index) (*IndexNode, error) {

	fmt.Printf("Node Add %d\n", childFingerprintSize)

	childFingerprint := entry.FingerprintForSize(childFingerprintSize)

	if len(node.childFingerprints) == 0 {

		// We can go deeper and this new entry is sufficiently different to
		// the rest, so split this leaf node by turning entries into children.
		fmt.Printf("Max Diff: %f\n", node.maxChildDifferenceTo(entry.MaxFingerprint))
		if childFingerprintSize < index.maxFingerprintSize && node.maxChildDifferenceTo(entry.MaxFingerprint) > index.maxEntryDifference {
			fmt.Printf("Pushing entries to children\n")
			node.pushEntriesToChildren(nodeFingerprint, childFingerprintSize, index.Store)
			fmt.Printf("Done pushing entries to children\n")

		} else {
			fmt.Printf("Adding entry %s\n", nodeFingerprint.String())
			err := index.Store.AddEntry(entry, node, nodeFingerprint)
			if err != nil {
				return nil, err
			}
			fmt.Printf("Added entry\n")
			return node, nil
		}
	}

	child, err := index.Store.GetOrCreateChild(childFingerprint, node, nodeFingerprint)
	if err != nil {
		return nil, err
	}

	return child.Add(entry, childFingerprint, childFingerprintSize+1, index)
}

func (node *IndexNode) FindNearest(entry *IndexEntry, childFingerprintSize int, index *Index, maxResults int, maxDifference float64) ([]*IndexEntry, error) {
	results := make([]*IndexEntry, 0, maxResults)

	err := node.gatherNearest(entry, childFingerprintSize, index, maxDifference, &results)
	if err != nil && err != errResultLimitReached {
		return nil, err
	}

	return results, nil
}

func (node *IndexNode) MarshalJSON() ([]byte, error) {
	return json.Marshal(&indexNodeJSON{
		ChildFingerprints: node.childFingerprints,
		Entries:           node.entries,
	})
}

func (node *IndexNode) UnmarshalJSON(b []byte) error {
	var value indexNodeJSON
	err := json.Unmarshal(b, &value)
	if err != nil {
		return err
	}

	node.childFingerprints = value.ChildFingerprints

	node.childFingerprintsByString = make(map[string]*Fingerprint)
	for i := 0; i < len(node.childFingerprints); i++ {
		f := &node.childFingerprints[i]
		node.childFingerprintsByString[f.String()] = f
	}

	node.entries = value.Entries

	return nil
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

	fmt.Printf("%d gatherNearest %d\n", childFingerprintSize, len(node.entries))

	// Check for an exact matching child
	childFingerprint := entry.FingerprintForSize(childFingerprintSize)
	exactChildFingerprint, exactChildFingerprintExists := node.childFingerprintsByString[childFingerprint.String()]

	var exactChildFingerprintString string
	var exactChild *IndexNode
	if exactChildFingerprintExists {
		exactChildFingerprintString = exactChildFingerprint.String()

		var err error
		exactChild, err = index.Store.GetChild(childFingerprint, node)
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

	childFingerprints := make([]Fingerprint, len(node.childFingerprints))
	copy(childFingerprints, node.childFingerprints)

	// Need more results - find and sort all children by nearness
	sort.Sort(nodesByDifferenceToFingerprintWith(childFingerprints, childFingerprint))

	// fmt.Printf("Sorting %d children...\n", len(children))
	// for i, child := range children {
	// 	diff := child.fingerprint.Difference(entryFingerprint)
	// 	fmt.Printf("%d sorted child %d of %f (%d %d)\n", childFingerprintSize+1, i, diff, len(child.fingerprint.samples), len(entryFingerprint.samples))
	// }

	// Recursively gather from nearest children
	for i, cf := range childFingerprints {
		fmt.Printf("Visiting child %d\n", i)
		if exactChildFingerprintExists && cf.String() == exactChildFingerprintString {
			continue
		}

		childNode, err := index.Store.GetChild(cf, node)
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

func (node *IndexNode) pushEntriesToChildren(nodeFingerprint Fingerprint, childFingerprintSize int, store IndexStore) error {
	node.withEachEntry(func(entry *IndexEntry) error {
		childFingerprint := entry.FingerprintForSize(childFingerprintSize)
		child, err := store.GetOrCreateChild(childFingerprint, node, nodeFingerprint)
		if err != nil {
			return err
		}
		return store.AddEntry(entry, child, childFingerprint)
	})

	return store.RemoveEntries(node, nodeFingerprint)
}

func (node *IndexNode) registerChild(childFingerprint Fingerprint) {
	node.childFingerprints = append(node.childFingerprints, childFingerprint)
	node.childFingerprintsByString[childFingerprint.String()] = &node.childFingerprints[len(node.childFingerprints)-1]
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

type indexNodeJSON struct {
	ChildFingerprints []Fingerprint `json:"childFingerprints"`
	Entries           []*IndexEntry `json:"entries"`
}

type nodesByDifferenceToFingerprint struct {
	nodeFingerprints []Fingerprint
	differences      []float64
}

func (sorter *nodesByDifferenceToFingerprint) Len() int {
	return len(sorter.nodeFingerprints)
}

func (sorter *nodesByDifferenceToFingerprint) Less(i, j int) bool {
	return sorter.differences[i] < sorter.differences[j]
}

func (sorter *nodesByDifferenceToFingerprint) Swap(i, j int) {
	tmp := sorter.nodeFingerprints[i]
	sorter.nodeFingerprints[i] = sorter.nodeFingerprints[j]
	sorter.nodeFingerprints[j] = tmp

	tmpDiff := sorter.differences[i]
	sorter.differences[i] = sorter.differences[j]
	sorter.differences[j] = tmpDiff
}

func nodesByDifferenceToFingerprintWith(nodeFingerprints []Fingerprint, f Fingerprint) *nodesByDifferenceToFingerprint {
	differences := make([]float64, len(nodeFingerprints), len(nodeFingerprints))
	for i, nf := range nodeFingerprints {
		differences[i] = nf.Difference(f)
	}

	return &nodesByDifferenceToFingerprint{nodeFingerprints: nodeFingerprints, differences: differences}
}
