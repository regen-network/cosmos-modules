package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/orm"

	"time"
)

type GroupID uint64

type ProposalID uint64

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
