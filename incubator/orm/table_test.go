package orm

import (
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/modules/incubator/orm/testdata"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreate(t *testing.T) {
	specs := map[string]struct {
		srcObj Persistent
		expErr *errors.Error
	}{
		"happy path": {
			srcObj: &testdata.GroupMember{
				Group:  sdk.AccAddress(EncodeSequence(1)),
				Member: sdk.AccAddress([]byte("member-address")),
				Weight: 10,
			},
		},
		"wrong type": {
			srcObj: &testdata.GroupMetadata{},
			expErr: ErrType,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			storeKey := sdk.NewKVStoreKey("test")
			const anyPrefix = 0x10
			tableBuilder := NewTableBuilder(anyPrefix, storeKey, &testdata.GroupMember{}, Max255DynamicLengthIndexKeyCodec{})
			myTable := tableBuilder.Build()

			ctx := NewMockContext()
			err := myTable.Create(ctx, []byte("any-id"), spec.srcObj)

			require.True(t, spec.expErr.Is(err), err)
			shouldExists := spec.expErr == nil
			assert.Equal(t, shouldExists, myTable.Has(ctx, []byte("any-id")), fmt.Sprintf("expected %v", shouldExists))
		})
	}

}
func TestUpdate(t *testing.T) {
	specs := map[string]struct {
		src    Persistent
		expErr *errors.Error
	}{
		"happy path": {
			src: &testdata.GroupMember{
				Group:  sdk.AccAddress(EncodeSequence(1)),
				Member: sdk.AccAddress([]byte("member-address")),
				Weight: 9999,
			},
		},
		"wrong type": {
			src:    &testdata.GroupMetadata{},
			expErr: ErrType,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			storeKey := sdk.NewKVStoreKey("test")
			const anyPrefix = 0x10
			tableBuilder := NewTableBuilder(anyPrefix, storeKey, &testdata.GroupMember{}, Max255DynamicLengthIndexKeyCodec{})
			myTable := tableBuilder.Build()

			initValue := testdata.GroupMember{
				Group:  sdk.AccAddress(EncodeSequence(1)),
				Member: sdk.AccAddress([]byte("member-address")),
				Weight: 10,
			}
			ctx := NewMockContext()
			err := myTable.Create(ctx, []byte("any-id"), &initValue)
			require.NoError(t, err)

			// when
			err = myTable.Save(ctx, []byte("any-id"), spec.src)
			require.True(t, spec.expErr.Is(err), "got ", err)

			// then
			var loaded testdata.GroupMember
			require.NoError(t, myTable.GetOne(ctx, []byte("any-id"), &loaded))
			if spec.expErr == nil {
				assert.Equal(t, spec.src, &loaded)
			} else {
				assert.Equal(t, initValue, loaded)
			}
		})
	}

}
