package group

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	auth "github.com/cosmos/cosmos-sdk/x/auth/types"
	proto "github.com/gogo/protobuf/types"
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

func TestMsgCreateGroupAccountStd(t *testing.T) {
	_, _, myAddr := auth.KeyTestPubAddr()

	specs := map[string]struct {
		src    MsgCreateGroupAccountStd
		expErr bool
	}{
		"all good with minimum fields set": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Admin: myAddr, Group: 1},
				DecisionPolicy: StdDecisionPolicy{
					Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
						Threshold: sdk.OneDec(),
						Timout:    proto.Duration{Seconds: 1},
					}},
				},
			},
		},
		"zero threshold allowed": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Admin: myAddr, Group: 1},
				DecisionPolicy: StdDecisionPolicy{
					Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
						Threshold: sdk.ZeroDec(),
						Timout:    proto.Duration{Seconds: 1},
					}},
				},
			},
		},
		"admin required": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Group: 1},
				DecisionPolicy: StdDecisionPolicy{
					Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
						Threshold: sdk.ZeroDec(),
						Timout:    proto.Duration{Seconds: 1},
					}},
				},
			},
			expErr: true,
		},
		"valid admin required": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Admin: []byte("invalid-address"), Group: 1},
				DecisionPolicy: StdDecisionPolicy{
					Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
						Threshold: sdk.ZeroDec(),
						Timout:    proto.Duration{Seconds: 1},
					}},
				},
			},
			expErr: true,
		},
		"group required": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Admin: myAddr},
				DecisionPolicy: StdDecisionPolicy{
					Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
						Threshold: sdk.ZeroDec(),
						Timout:    proto.Duration{Seconds: 1},
					}},
				},
			},
			expErr: true,
		},
		"decision policy required": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Admin: myAddr, Group: 1},
			},
			expErr: true,
		},
		"decision policy without timout": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Admin: myAddr, Group: 1},
				DecisionPolicy: StdDecisionPolicy{
					Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
						Threshold: sdk.ZeroDec(),
					}},
				},
			},
			expErr: true,
		},
		"decision policy with invalid timout": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Admin: myAddr, Group: 1},
				DecisionPolicy: StdDecisionPolicy{
					Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
						Threshold: sdk.ZeroDec(),
						Timout:    proto.Duration{Seconds: -1},
					}},
				},
			},
			expErr: true,
		},
		"decision policy without threshold": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Admin: myAddr, Group: 1},
				DecisionPolicy: StdDecisionPolicy{
					Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
						Timout: proto.Duration{Seconds: 1},
					}},
				},
			},
			expErr: true,
		},
		"decision policy with negative threshold": {
			src: MsgCreateGroupAccountStd{
				Base: MsgCreateGroupAccountBase{Admin: myAddr, Group: 1},
				DecisionPolicy: StdDecisionPolicy{
					Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
						Threshold: sdk.NewDec(-1),
						Timout:    proto.Duration{Seconds: 1},
					}},
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
