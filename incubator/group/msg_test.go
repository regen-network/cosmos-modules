package group

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMsgCreateGroupValidation(t *testing.T) {
	_, _, myAddr := auth.KeyTestPubAddr()
	_, _, myOtherAddr := auth.KeyTestPubAddr()

	specs := map[string]struct {
		src    MsgCreateGroup
		expErr bool
	}{
		"all good with minimum fields set": {
			src: MsgCreateGroup{Admin: myAddr},
		},
		"all good with a member": {
			src: MsgCreateGroup{
				Admin: myAddr,
				Members: []Member{
					{Address: myAddr, Power: types.NewDec(1)},
				},
			},
		},
		"all good with multiple members": {
			src: MsgCreateGroup{
				Admin: myAddr,
				Members: []Member{
					{Address: myAddr, Power: types.NewDec(1)},
					{Address: myOtherAddr, Power: types.NewDec(2)},
				},
			},
		},
		"admin required": {
			src:    MsgCreateGroup{},
			expErr: true,
		},
		"valid admin required": {
			src: MsgCreateGroup{
				Admin: []byte("invalid-address"),
			},
			expErr: true,
		},
		"duplicate member addresses not allowed": {
			src: MsgCreateGroup{
				Admin: myAddr,
				Members: []Member{
					{Address: myAddr, Power: types.NewDec(1)},
					{Address: myAddr, Power: types.NewDec(2)},
				},
			},
			expErr: true,
		},
		"negative member's power not allowed": {
			src: MsgCreateGroup{
				Admin: myAddr,
				Members: []Member{
					{Address: myAddr, Power: types.NewDec(-1)},
				},
			},
			expErr: true,
		},
		"empty member's power not allowed": {
			src: MsgCreateGroup{
				Admin:   myAddr,
				Members: []Member{{Address: myAddr}},
			},
			expErr: true,
		},
		"zero member's power not allowed": {
			src: MsgCreateGroup{
				Admin:   myAddr,
				Members: []Member{{Address: myAddr, Power: sdk.ZeroDec()}},
			},
			expErr: true,
		},
		"member address required": {
			src: MsgCreateGroup{
				Admin: myAddr,
				Members: []Member{
					{Power: types.NewDec(1)},
				},
			},
			expErr: true,
		},
		"valid member address required": {
			src: MsgCreateGroup{
				Admin: myAddr,
				Members: []Member{
					{Address: []byte("invalid-address"), Power: types.NewDec(1)},
				},
			},
			expErr: true,
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			err := spec.src.ValidateBasic()
			if spec.expErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
func TestMsgCreateGroupSigner(t *testing.T) {
	_, _, myAddr := auth.KeyTestPubAddr()
	assert.Equal(t, []sdk.AccAddress{myAddr}, MsgCreateGroup{Admin: myAddr}.GetSigners())
}
