package group

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/types"
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
