package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/modules/incubator/orm"
	"gopkg.in/yaml.v2"

	"time"
)

type GroupID uint64

func (g GroupID) Byte() []byte {
	return orm.EncodeSequence(uint64(g))
}

type ProposalID uint64

func (p ProposalID) Byte() []byte {
	return orm.EncodeSequence(uint64(p))
}

type DecisionPolicy interface {
	Allow(tally Tally, totalPower sdk.Dec, votingDuration time.Duration) bool
}

func (g GroupMember) NaturalKey() []byte {
	result := make([]byte, 0, 8+len(g.Member))
	copy(result[0:8], orm.EncodeSequence(uint64(g.Group)))
	result = append(result, g.Member...)
	return result
}

func (g GroupAccountMetadataBase) NaturalKey() []byte {
	return g.GroupAccount
}

func (g StdGroupAccountMetadata) NaturalKey() []byte {
	return g.Base.NaturalKey()
}

func (v Vote) NaturalKey() []byte {
	result := make([]byte, 0, 8+len(v.Voter))
	copy(result[0:8], orm.EncodeSequence(uint64(v.Proposal)))
	result = append(result, v.Voter...)
	return result
}

func (m Member) ValidateBasic() error {
	if m.Address.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "address")
	}
	return nil
}

const defaultMaxCommentLength = 255

// Parameter keys
var (
	ParamMaxCommentLength = []byte("MaxCommentLength")
)

// DefaultParams returns the default parameters for the group module.
func DefaultParams() Params {
	return Params{
		MaxCommentLength: defaultMaxCommentLength,
	}
}

func (p Params) String() string {
	out, _ := yaml.Marshal(p)
	return string(out)
}

// ParamSetPairs returns the parameter set pairs.
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(ParamMaxCommentLength, &p.MaxCommentLength, noopValidator()),
	}
}
func (p Params) Validate() error {
	return nil
}

func noopValidator() subspace.ValueValidatorFn {
	return func(value interface{}) error { return nil }
}
