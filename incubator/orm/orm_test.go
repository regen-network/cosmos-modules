package orm

import (
	"reflect"
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTypeSafeRowGetter(t *testing.T) {
	storeKey := sdk.NewKVStoreKey("test")
	cdc := codec.New()
	ctx := NewMockContext()
	const prefixKey = 0x2
	store := prefix.NewStore(ctx.KVStore(storeKey), []byte{prefixKey})
	store.Set(EncodeSequence(1), cdc.MustMarshalBinaryBare("foo"))

	specs := map[string]struct {
		srcRowID     uint64
		srcModelType reflect.Type
		expObj       interface{}
		expErr       *errors.Error
	}{
		"happy path": {
			srcRowID:     1,
			srcModelType: reflect.TypeOf(""),
			expObj:       "foo",
		},
		"unknown rowID should return ErrNotFound": {
			srcRowID:     999,
			srcModelType: reflect.TypeOf(""),
			expErr:       ErrNotFound,
		},
		"wrong type should cause ErrType": {
			srcRowID:     1,
			srcModelType: reflect.TypeOf(1),
			expErr:       ErrType,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			getter := NewTypeSafeRowGetter(storeKey, prefixKey, cdc, spec.srcModelType)
			var loadedObj string
			key, err := getter(ctx, spec.srcRowID, &loadedObj)
			if spec.expErr != nil {
				require.True(t, spec.expErr.Is(err))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, EncodeSequence(spec.srcRowID), key)
			assert.Equal(t, spec.expObj, loadedObj)
		})
	}
}
