package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/modules/incubator/orm"

	"time"
)

// TODO: what makes sense here?
const MaxCommentSize = 256

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
	if len(m.Comment) > MaxCommentSize {
		return sdkerrors.Wrap(ErrMaxLimit, "comment size")
	}
	return nil
}
