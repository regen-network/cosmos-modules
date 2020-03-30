package group

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
	proto "github.com/gogo/protobuf/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestVoteNaturalKey(t *testing.T) {
	v := Vote{
		Proposal: 1,
		Voter:    []byte{0xff, 0xfe},
	}
	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 1, 0xff, 0xfe}, v.NaturalKey())
}

func TestGroupMetadataValidation(t *testing.T) {
	specs := map[string]struct {
		src    GroupMetadata
		expErr bool
	}{
		"all good": {
			src: GroupMetadata{
				Group:       1,
				Admin:       []byte("valid--admin-address"),
				Comment:     "any",
				Version:     1,
				TotalWeight: types.ZeroDec(),
			},
		},
		"invalid group": {
			src: GroupMetadata{
				Admin:       []byte("valid--admin-address"),
				Comment:     "any",
				Version:     1,
				TotalWeight: types.ZeroDec(),
			},
			expErr: true,
		},
		"invalid admin": {
			src: GroupMetadata{
				Group:       1,
				Admin:       []byte("invalid"),
				Comment:     "any",
				Version:     1,
				TotalWeight: types.ZeroDec(),
			},
			expErr: true,
		},
		"invalid version": {
			src: GroupMetadata{
				Group:       1,
				Admin:       []byte("valid--admin-address"),
				Comment:     "any",
				TotalWeight: types.ZeroDec(),
			},
			expErr: true,
		},
		"unset total weight": {
			src: GroupMetadata{
				Group:   1,
				Admin:   []byte("valid--admin-address"),
				Comment: "any",
				Version: 1,
			},
			expErr: true,
		},
		"negative total weight": {
			src: GroupMetadata{
				Group:       1,
				Admin:       []byte("valid--admin-address"),
				Comment:     "any",
				Version:     1,
				TotalWeight: types.NewDec(-1),
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

func TestGroupMemberValidation(t *testing.T) {
	specs := map[string]struct {
		src    GroupMember
		expErr bool
	}{
		"all good": {
			src: GroupMember{
				Group:   1,
				Member:  []byte("valid-member-address"),
				Weight:  types.OneDec(),
				Comment: "any",
			},
		},
		"invalid group": {
			src: GroupMember{
				Group:   0,
				Member:  []byte("valid-member-address"),
				Weight:  types.OneDec(),
				Comment: "any",
			},
			expErr: true,
		},
		"invalid address": {
			src: GroupMember{
				Group:   1,
				Member:  []byte("invalid-member-address"),
				Weight:  types.OneDec(),
				Comment: "any",
			},
			expErr: true,
		},
		"empy address": {
			src: GroupMember{
				Group:   1,
				Weight:  types.OneDec(),
				Comment: "any",
			},
			expErr: true,
		},
		"invalid weight": {
			src: GroupMember{
				Group:   1,
				Member:  []byte("valid-member-address"),
				Weight:  types.ZeroDec(),
				Comment: "any",
			},
			expErr: true,
		},
		"nil weight": {
			src: GroupMember{
				Group:   1,
				Member:  []byte("valid-member-address"),
				Comment: "any",
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

func TestGroupAccountMetadataBase(t *testing.T) {
	specs := map[string]struct {
		src    GroupAccountMetadataBase
		expErr bool
	}{
		"all good": {
			src: GroupAccountMetadataBase{
				Group:        1,
				GroupAccount: []byte("valid--group-address"),
				Admin:        []byte("valid--admin-address"),
				Comment:      "any",
				Version:      1,
			},
		},
		"invalid group": {
			src: GroupAccountMetadataBase{
				Group:        0,
				GroupAccount: []byte("valid--group-address"),
				Admin:        []byte("valid--admin-address"),
				Comment:      "any",
				Version:      1,
			},
			expErr: true,
		},
		"invalid group account address": {
			src: GroupAccountMetadataBase{
				Group:        1,
				GroupAccount: []byte("any-invalid-group-address"),
				Admin:        []byte("valid--admin-address"),
				Comment:      "any",
				Version:      1,
			},
			expErr: true,
		},
		"empty group account address": {
			src: GroupAccountMetadataBase{
				Group:   1,
				Admin:   []byte("valid--admin-address"),
				Comment: "any",
				Version: 1,
			},
			expErr: true,
		},
		"invalid admin account address": {
			src: GroupAccountMetadataBase{
				Group:        1,
				GroupAccount: []byte("valid--group-address"),
				Admin:        []byte("any-invalid-admin-address"),
				Comment:      "any",
				Version:      1,
			},
			expErr: true,
		},
		"empty admin account address": {
			src: GroupAccountMetadataBase{
				Group:        1,
				GroupAccount: []byte("valid--group-address"),
				Comment:      "any",
				Version:      1,
			},
			expErr: true,
		},
		"empty version number": {
			src: GroupAccountMetadataBase{
				Group:        1,
				GroupAccount: []byte("valid--group-address"),
				Admin:        []byte("valid--admin-address"),
				Comment:      "any",
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

func TestStdGroupAccountMetadata(t *testing.T) {
	specs := map[string]struct {
		src    StdGroupAccountMetadata
		expErr bool
	}{
		"all good": {
			src: StdGroupAccountMetadata{
				Base: GroupAccountMetadataBase{
					Group:        1,
					GroupAccount: []byte("valid--group-address"),
					Admin:        []byte("valid--admin-address"),
					Comment:      "any",
					Version:      1,
				},
				DecisionPolicy: StdDecisionPolicy{Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
					Threshold: types.ZeroDec(),
					Timout:    proto.Duration{Seconds: 1},
				}}},
			},
		},
		"invalid base": {
			src: StdGroupAccountMetadata{
				Base: GroupAccountMetadataBase{},
				DecisionPolicy: StdDecisionPolicy{Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
					Threshold: types.ZeroDec(),
					Timout:    proto.Duration{Seconds: 1},
				}}},
			},
			expErr: true,
		},
		"missing base": {
			src: StdGroupAccountMetadata{
				DecisionPolicy: StdDecisionPolicy{Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{
					Threshold: types.ZeroDec(),
					Timout:    proto.Duration{Seconds: 1},
				}}},
			},
			expErr: true,
		},
		"invalid decision policy": {
			src: StdGroupAccountMetadata{
				Base: GroupAccountMetadataBase{
					Group:        1,
					GroupAccount: []byte("valid--group-address"),
					Admin:        []byte("valid--admin-address"),
					Comment:      "any",
					Version:      1,
				},
				DecisionPolicy: StdDecisionPolicy{Sum: &StdDecisionPolicy_Threshold{&ThresholdDecisionPolicy{}}},
			},
			expErr: true,
		},
		"concrete decision policy not set": {
			src: StdGroupAccountMetadata{
				Base: GroupAccountMetadataBase{
					Group:        1,
					GroupAccount: []byte("valid--group-address"),
					Admin:        []byte("valid--admin-address"),
					Comment:      "any",
					Version:      1,
				},
				DecisionPolicy: StdDecisionPolicy{Sum: &StdDecisionPolicy_Threshold{}},
			},
			expErr: true,
		},
		"missing decision policy": {
			src: StdGroupAccountMetadata{
				Base: GroupAccountMetadataBase{
					Group:        1,
					GroupAccount: []byte("valid--group-address"),
					Admin:        []byte("valid--admin-address"),
					Comment:      "any",
					Version:      1,
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
