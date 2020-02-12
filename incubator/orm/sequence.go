package orm

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// sequenceStorageKey is a fix key to read/ write data on the storage layer
var sequenceStorageKey = []byte{0x1}

// sequence is a persistent unique key generator based on a counter.
type Sequence struct {
	storeKey sdk.StoreKey
	prefix   byte
}

func NewSequence(storeKey sdk.StoreKey, prefix byte) *Sequence {
	return &Sequence{
		prefix:   prefix,
		storeKey: storeKey,
	}
}

// NextVal increments the counter by one and returns the value.
func (s Sequence) NextVal(ctx HasKVStore) uint64 {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), []byte{s.prefix})
	v := store.Get(sequenceStorageKey)
	seq := DecodeSequence(v)
	seq++
	store.Set(sequenceStorageKey, EncodeSequence(seq))
	return seq
}

// CurVal returns the last value used. 0Nex if none.
func (s Sequence) CurVal(ctx HasKVStore) uint64 {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), []byte{s.prefix})
	v := store.Get(sequenceStorageKey)
	return DecodeSequence(v)
}

// DecodeSequence converts the binary representation into an Uint64 value.
func DecodeSequence(bz []byte) uint64 {
	if bz == nil {
		return 0
	}
	val := binary.BigEndian.Uint64(bz)
	return val
}

// EncodedSeqLength number of bytes used for the binary representation of a sequence value.
const EncodedSeqLength = 8

// EncodeSequence converts the sequence value into the binary representation.
func EncodeSequence(val uint64) []byte {
	bz := make([]byte, EncodedSeqLength)
	binary.BigEndian.PutUint64(bz, val)
	return bz
}
