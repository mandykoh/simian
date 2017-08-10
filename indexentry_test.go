package simian

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestIndexEntry(t *testing.T) {

	t.Run("JSON serialisation", func(t *testing.T) {

		t.Run("should roundtrip all fields", func(t *testing.T) {

			entry := &IndexEntry{
				MaxFingerprint: Fingerprint{samples: []uint8{0xF0, 0xF0, 0xF0, 0xF0}},
				Attributes:     make(map[string]interface{}),
			}
			entry.Attributes["some key"] = "some value"
			entry.Attributes["some other key"] = "some other value"

			jsonBytes, err := json.Marshal(entry)
			if err != nil {
				t.Fatalf("Error marshalling JSON: %v", err)
			}

			var result *IndexEntry
			err = json.Unmarshal(jsonBytes, &result)
			if err != nil {
				t.Fatalf("Error unmarshalling JSON: %v", err)
			}

			if distance := result.MaxFingerprint.Distance(entry.MaxFingerprint); distance != 0 {
				t.Errorf("Expected no difference in fingerprints but got %d", distance)
			}
			if !reflect.DeepEqual(entry.Attributes, result.Attributes) {
				t.Errorf("Expected attributes to match but got %v", result.Attributes)
			}
		})
	})
}
