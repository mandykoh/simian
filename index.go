package simian

import (
	"fmt"
	"image"
	"math"
	"os"
	"sort"
)

const rootFingerprintSize = 1

type Index struct {
	Store              IndexStore
	maxFingerprintSize int
	maxEntryDifference float64
}

func (i *Index) Add(image image.Image, metadata interface{}) (key string, err error) {
	entry, err := NewIndexEntry(image, i.maxFingerprintSize)
	if err != nil {
		return "", nil
	}

	root, err := i.Store.GetRoot()
	if err != nil {
		return "", err
	}

	var rootFingerprint Fingerprint

	_, err = root.Add(entry, rootFingerprint, rootFingerprintSize+1, i)
	if err != nil {
		return "", err
	}

	fmt.Printf("Root node has %d children and %d entries\n", len(root.childFingerprints), len(root.entries))

	return "", nil
}

func (i *Index) Close() error {
	return i.Store.Close()
}

func (i *Index) FindNearest(image image.Image, maxResults int, maxDifference float64) ([]*IndexEntry, error) {
	entry, err := NewIndexEntry(image, i.maxFingerprintSize)
	if err != nil {
		return nil, nil
	}

	root, err := i.Store.GetRoot()
	if err != nil {
		return nil, err
	}

	results, err := root.FindNearest(entry, rootFingerprintSize+1, i, maxResults, math.Max(maxDifference, i.maxEntryDifference))
	if err != nil {
		return nil, err
	}
	sort.Sort(entriesByDifferenceToEntryWith(results, entry))

	return results, err
}

func NewIndex(path string, maxFingerprintSize int, maxEntryDifference float64) (*Index, error) {
	err := os.MkdirAll(path, 0700)
	if err != nil {
		return nil, err
	}

	indexStore, err := NewDiskIndexStore(path)
	if err != nil {
		return nil, err
	}

	return &Index{
		Store:              indexStore,
		maxFingerprintSize: maxFingerprintSize,
		maxEntryDifference: maxEntryDifference,
	}, err
}

type entriesByDifferenceToEntry struct {
	entries     []*IndexEntry
	differences []float64
}

func (sorter *entriesByDifferenceToEntry) Len() int {
	return len(sorter.entries)
}

func (sorter *entriesByDifferenceToEntry) Less(i, j int) bool {
	return sorter.differences[i] < sorter.differences[j]
}

func (sorter *entriesByDifferenceToEntry) Swap(i, j int) {
	tmpEntry := sorter.entries[i]
	sorter.entries[i] = sorter.entries[j]
	sorter.entries[j] = tmpEntry

	tmpDiff := sorter.differences[i]
	sorter.differences[i] = sorter.differences[j]
	sorter.differences[j] = tmpDiff
}

func entriesByDifferenceToEntryWith(entries []*IndexEntry, target *IndexEntry) *entriesByDifferenceToEntry {
	differences := make([]float64, len(entries), len(entries))
	for i, entry := range entries {
		differences[i] = entry.MaxFingerprint.Difference(target.MaxFingerprint)
	}

	return &entriesByDifferenceToEntry{entries: entries, differences: differences}
}
