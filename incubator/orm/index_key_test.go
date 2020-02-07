package orm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeIndexKey(t *testing.T) {
	specs := map[string]struct {
		srcKey   []byte
		srcRowID RowID
		expKey   []byte
	}{
		"example 1": {
			srcKey:   []byte{0x0, 0x1, 0x2},
			srcRowID: []byte{0x3, 0x4},
			expKey:   []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x2},
		},
		"example 2": {
			srcKey:   []byte{0x0, 0x1},
			srcRowID: []byte{0x2, 0x3, 0x4},
			expKey:   []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x3},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			got := makeIndexPrefixScanKey(spec.srcKey, spec.srcRowID)
			assert.Equal(t, spec.expKey, got)
		})
	}
}
func TestDecodeIndexKey(t *testing.T) {
	specs := map[string]struct {
		srcKey   []byte
		expRowID RowID
	}{
		"example 1": {
			srcKey:   []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x2},
			expRowID: []byte{0x3, 0x4},
		},
		"example 2": {
			srcKey:   []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x3},
			expRowID: []byte{0x2, 0x3, 0x4},
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			gotRow := stripRowIDFromIndexPrefixScanKey(spec.srcKey)
			assert.Equal(t, spec.expRowID, gotRow)
		})
	}
}
