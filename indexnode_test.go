package simian

import (
	"encoding/json"
	"testing"
)

func TestIndexNode(t *testing.T) {

	t.Run("JSON serialisation", func(t *testing.T) {

		t.Run("should roundtrip all fields", func(t *testing.T) {
			n := &IndexNode{
				path: "some-path",
				childFingerprintsByString: make(map[string]*Fingerprint),
			}

			n.registerChild(Fingerprint{samples: []uint8{0x10, 0x20, 0x30, 0x40}})
			n.registerChild(Fingerprint{samples: []uint8{0x50, 0x60, 0x70, 0x80}})

			entry1 := &IndexEntry{
				MaxFingerprint: Fingerprint{samples: []uint8{1, 2, 3, 4, 5, 6, 7, 8, 9}},
				Attributes:     make(map[string]interface{}),
			}
			n.registerEntry(entry1)

			entry2 := &IndexEntry{
				MaxFingerprint: Fingerprint{samples: []uint8{10, 11, 12, 13, 14, 15, 16, 17, 18}},
				Attributes:     make(map[string]interface{}),
			}
			n.registerEntry(entry2)

			jsonBytes, err := json.Marshal(n)
			if err != nil {
				t.Fatalf("Error marshalling JSON: %v", err)
			}

			var result *IndexNode
			err = json.Unmarshal(jsonBytes, &result)
			if err != nil {
				t.Fatalf("Error unmarshalling JSON: %v", err)
			}

			if result.path != n.path {
				t.Errorf("Expected path '%s' but got '%s'", n.path, result.path)
			}

			if actual, expected := len(result.childFingerprints), len(n.childFingerprints); actual != expected {
				t.Fatalf("Expected %d child fingerprints but got %d", expected, actual)
			}
			for i := 0; i < len(result.childFingerprints); i++ {
				actual := result.childFingerprints[i].String()
				expected := n.childFingerprints[i].String()

				if actual != expected {
					t.Errorf("Expected fingerprint '%s' but got '%s'", expected, actual)
				}
			}

			if actual, expected := len(result.childFingerprintsByString), len(n.childFingerprintsByString); actual != expected {
				t.Fatalf("Expected %d child fingerprints mapped by string but got %d", expected, actual)
			}
			for k, v := range n.childFingerprintsByString {
				actual := result.childFingerprintsByString[k].String()
				expected := v.String()

				if actual != expected {
					t.Errorf("Expected fingerprint '%s' but got '%s'", expected, actual)
				}
			}

			if actual, expected := len(result.entries), len(n.entries); actual != expected {
				t.Fatalf("Expected %d entries but got %d", expected, actual)
			}
			for i := 0; i < len(result.entries); i++ {
				actualBytes, err := json.Marshal(result.entries[i])
				if err != nil {
					t.Fatalf("Error marshalling entry: %v", err)
				}
				actual := string(actualBytes)

				expectedBytes, err := json.Marshal(n.entries[i])
				if err != nil {
					t.Fatalf("Error marshalling entry: %v", err)
				}
				expected := string(expectedBytes)

				if actual != expected {
					t.Errorf("Expected entry '%s' but got '%s'", expected, actual)
				}
			}
		})
	})
}
