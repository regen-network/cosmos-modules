package orm

import (
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeSafeRowGetter(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	ctx := NewMockContext()
	const prefixKey = 0x2
	store := prefix.NewStore(ctx.KVStore(storeKey), []byte{prefixKey})
	md := GroupMetadata{Description: "foo"}
	bz, err := md.Marshal()
	require.NoError(t, err)
	store.Set(EncodeSequence(1), bz)

	specs := map[string]struct {
		srcRowID     uint64
		srcModelType reflect.Type
		expObj       interface{}
		expErr       *errors.Error
	}{
		"happy path": {
			srcRowID:     1,
			srcModelType: reflect.TypeOf(GroupMetadata{}),
			expObj:       GroupMetadata{Description: "foo"},
		},
		"unknown rowID should return ErrNotFound": {
			srcRowID:     999,
			srcModelType: reflect.TypeOf(GroupMetadata{}),
			expErr:       ErrNotFound,
		},
		"wrong type should cause ErrType": {
			srcRowID:     1,
			srcModelType: reflect.TypeOf(GroupMember{}),
			expErr:       ErrType,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			getter := NewTypeSafeRowGetter(storeKey, prefixKey, spec.srcModelType)
			var loadedObj GroupMetadata
			err := getter(ctx, EncodeSequence(spec.srcRowID), &loadedObj)
			if spec.expErr != nil {
				require.True(t, spec.expErr.Is(err), err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, spec.expObj, loadedObj)
		})
	}
}
