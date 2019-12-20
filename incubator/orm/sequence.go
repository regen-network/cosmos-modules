package orm

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ Sequence = &sequence{}

// sequenceStorageKey is a fix key to read/ write data on the storage layer
var sequenceStorageKey = []byte{0x1}

// sequence is a persistent unique key generator based on a counter.
type sequence struct {
	storeKey sdk.StoreKey
	prefix   []byte
}

func NewSequence(storeKey sdk.StoreKey, prefix []byte) *sequence {
	return &sequence{
		prefix:   prefix,
		storeKey: storeKey,
	}
}

// NextVal increments the counter by one and returns the value.
func (s sequence) NextVal(ctx HasKVStore) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), s.prefix)
	// TODO: store does not return an error. inconsistent method signature above
	v := store.Get(sequenceStorageKey)
	seq := DecodeSequence(v)
	seq += 1
	store.Set(sequenceStorageKey, EncodeSequence(seq))
	return seq, nil
}

// CurVal returns the last value used. 0 if none.
func (s sequence) CurVal(ctx HasKVStore) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), s.prefix)
	// TODO: store does not return an error. inconsistent method signature above
	v := store.Get(sequenceStorageKey)
	return DecodeSequence(v), nil
}

// DecodeSequence converts the binary representation into an Uint64 value.
func DecodeSequence(bz []byte) uint64 {
	if bz == nil {
		return 0
	}
	val := binary.BigEndian.Uint64(bz)
	return val
}

// EncodeSequence converts the sequence value into the binary representation.
func EncodeSequence(val uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, val)
	return bz
}
