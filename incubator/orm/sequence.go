package orm

import (
	"encoding/binary"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var _ Sequence = &sequence{}

type sequence struct {
	storeKey sdk.StoreKey
	prefix   []byte
	seqKey   []byte
}

func NewSequence(storeKey sdk.StoreKey, prefix []byte) *sequence {
	return &sequence{
		prefix: prefix,
		storeKey: storeKey,
		seqKey:[]byte("seq"), // todo: should seq key also be a short one?
	}
}

// NextVal increments the counter and returns the value.
func (s sequence) NextVal(ctx HasKVStore) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), s.prefix)
	// TODO: store does not return an error. inconsistent method signature above
	v := store.Get(s.seqKey)
	seq := decodeSequence(v)
	seq += 1
	store.Set(s.seqKey, encodeSequence(seq))
	return seq, nil
}

// CurVal returns the last value used. 0 if none.
func (s sequence) CurVal(ctx HasKVStore) (uint64, error) {
	store := prefix.NewStore(ctx.KVStore(s.storeKey), s.prefix)
	// TODO: store does not return an error. inconsistent method signature above
	v := store.Get(s.seqKey)
	return decodeSequence(v), nil
}

func decodeSequence(bz []byte) uint64 {
	if bz == nil {
		return 0
	}
	val := binary.BigEndian.Uint64(bz)
	return val
}

func encodeSequence(val uint64) []byte {
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, val)
	return bz
}
