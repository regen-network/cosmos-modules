package group

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateGroup(t *testing.T) {
	k := NewGroupKeeper(sdk.NewKVStoreKey(StoreKeyName))
	members := []Member{{
		Address: sdk.AccAddress([]byte("member-address")),
		Power:   sdk.NewDec(1),
		Comment: "first",
	}}
	ctx := orm.NewMockContext()
	id, err := k.CreateGroup(ctx, []byte("admin-address"), members, "test")
	require.NoError(t, err)
	assert.Equal(t, GroupID(1), id)
}
