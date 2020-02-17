package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"time"
)

type GroupID uint64

type ProposalID uint64

type DecisionPolicy interface {
	Allow(tally Tally, totalPower sdk.Dec, votingDuration time.Duration) bool
}

