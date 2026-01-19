package sync

import (
	"bytes"
	"encoding/json"
	"sort"
)

// OrderedMap marshals map keys in a stable order.
type OrderedMap[T any] map[string]T

// MarshalJSON encodes a map with sorted keys for deterministic output.
func (m OrderedMap[T]) MarshalJSON() ([]byte, error) {
	if m == nil {
		return []byte("null"), nil
	}
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var buffer bytes.Buffer
	buffer.WriteString("{")
	for i, key := range keys {
		if i > 0 {
			buffer.WriteString(",")
		}
		keyJSON, err := json.Marshal(key)
		if err != nil {
			return nil, err
		}
		valueJSON, err := json.Marshal(m[key])
		if err != nil {
			return nil, err
		}
		buffer.Write(keyJSON)
		buffer.WriteString(":")
		buffer.Write(valueJSON)
	}
	buffer.WriteString("}")
	return buffer.Bytes(), nil
}
