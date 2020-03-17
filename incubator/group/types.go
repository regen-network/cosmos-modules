package group

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/params/subspace"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/gogo/protobuf/types"
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"

	"time"
)

type GroupID uint64

func (p GroupID) Uint64() uint64 {
	return uint64(p)
}

func (p GroupID) Empty() bool {
	return p == 0
}

func (g GroupID) Bytes() []byte {
	return orm.EncodeSequence(uint64(g))
}

type ProposalID uint64

func (p ProposalID) Bytes() []byte {
	return orm.EncodeSequence(uint64(p))
}

func (p ProposalID) Uint64() uint64 {
	return uint64(p)
}
func (p ProposalID) Empty() bool {
	return p == 0
}

type DecisionPolicyResult struct {
	Allow bool
	Final bool
}

// DecisionPolicy persistent ruleset to determine the result of election on a proposal.
type DecisionPolicy interface {
	orm.Persistent
	Allow(tally Tally, totalPower sdk.Dec, votingDuration time.Duration) (DecisionPolicyResult, error)
}

// Allow allow a proposal to pass when the threshold is exceeded by yes votes before the timeout.
func (p ThresholdDecisionPolicy) Allow(tally Tally, totalPower sdk.Dec, votingDuration time.Duration) (DecisionPolicyResult, error) {
	timeout, err := types.DurationFromProto(&p.Timout)
	if err != nil {
		return DecisionPolicyResult{}, err
	}
	if timeout < votingDuration {
		return DecisionPolicyResult{Allow: false, Final: true}, nil
	}
	if tally.YesCount.GT(p.Threshold) {
		return DecisionPolicyResult{Allow: true, Final: true}, nil
	}
	undecided := totalPower.Sub(tally.TotalCounts())
	if tally.YesCount.Add(undecided).LTE(p.Threshold) {
		return DecisionPolicyResult{Allow: false, Final: true}, nil
	}
	return DecisionPolicyResult{Allow: false, Final: false}, nil
}

func (p ThresholdDecisionPolicy) ValidateBasic() error {
	if p.Threshold.IsNil() {
		return errors.Wrap(ErrEmpty, "threshold")
	}
	if p.Threshold.LT(sdk.ZeroDec()) {
		return errors.Wrap(ErrInvalid, "threshold")
	}
	timeout, err := types.DurationFromProto(&p.Timout)
	if err != nil {
		return errors.Wrap(err, "timeout")
	}

	if timeout <= time.Nanosecond {
		return errors.Wrap(ErrInvalid, "timeout")
	}
	return nil
}

func (g GroupMember) NaturalKey() []byte {
	result := make([]byte, 8, 8+len(g.Member))
	copy(result[0:8], g.Group.Bytes())
	result = append(result, g.Member...)
	return result
}

func (g GroupAccountMetadataBase) NaturalKey() []byte {
	return g.GroupAccount
}

func (g GroupAccountMetadataBase) ValidateBasic() error {
	if g.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if err := sdk.VerifyAddressFormat(g.Admin); err != nil {
		return sdkerrors.Wrap(err, "admin")
	}
	if g.GroupAccount.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "group account")
	}
	if err := sdk.VerifyAddressFormat(g.GroupAccount); err != nil {
		return sdkerrors.Wrap(err, "group account")
	}

	if g.Group == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group")
	}
	if g.Version == 0 {
		return sdkerrors.Wrap(ErrEmpty, "version")
	}

	return nil
}
func (s StdGroupAccountMetadata) NaturalKey() []byte {
	return s.Base.NaturalKey()
}

func (s StdGroupAccountMetadata) ValidateBasic() error {
	if err := s.Base.ValidateBasic(); err != nil {
		return errors.Wrap(err, "base")
	}
	policy := s.DecisionPolicy.GetDecisionPolicy()
	if policy == nil {
		return errors.Wrap(ErrEmpty, "policy")
	}
	if err := policy.ValidateBasic(); err != nil {
		return errors.Wrap(err, "policy")
	}
	return nil
}

func (v Vote) NaturalKey() []byte {
	result := make([]byte, 8, 8+len(v.Voter))
	copy(result[0:8], v.Proposal.Bytes())
	result = append(result, v.Voter...)
	return result
}

func (g Vote) ValidateBasic() error {
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

func (m GroupMetadata) ValidateBasic() error {
	if m.Group.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "group")
	}
	if m.Admin.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "admin")
	}
	if err := sdk.VerifyAddressFormat(m.Admin); err != nil {
		return sdkerrors.Wrap(err, "admin")
	}
	if m.TotalWeight.IsNil() || m.TotalWeight.LT(sdk.ZeroDec()) {
		return sdkerrors.Wrap(ErrInvalid, "total weight")
	}
	if m.Version == 0 {
		return sdkerrors.Wrap(ErrEmpty, "version")
	}
	return nil
}

func (m GroupMember) ValidateBasic() error {
	if m.Group.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "group")
	}
	if m.Member.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "address")
	}
	if m.Weight.IsNil() || m.Weight.LTE(sdk.ZeroDec()) {
		return sdkerrors.Wrap(ErrInvalid, "member power")
	}
	if err := sdk.VerifyAddressFormat(m.Member); err != nil {
		return sdkerrors.Wrap(err, "address")
	}
	return nil
}

func (p ProposalBase) ValidateBasic() error {
	if p.GroupAccount.Empty() {
		return sdkerrors.Wrap(ErrEmpty, "group account")
	}
	if err := sdk.VerifyAddressFormat(p.GroupAccount); err != nil {
		return sdkerrors.Wrap(err, "group account")
	}
	if len(p.Proposers) == 0 {
		return sdkerrors.Wrap(ErrEmpty, "proposers")
	}
	index := make(map[string]struct{}, len(p.Proposers))
	for _, p := range p.Proposers {
		if err := sdk.VerifyAddressFormat(p); err != nil {
			return sdkerrors.Wrap(err, "proposer")
		}
		if _, exists := index[string(p)]; exists {
			return sdkerrors.Wrapf(ErrDuplicate, "proposer %q", p.String())
		}
		index[string(p)] = struct{}{}
	}
	if p.SubmittedAt.Seconds == 0 && p.SubmittedAt.Nanos == 0 {
		return sdkerrors.Wrap(ErrEmpty, "submitted at")
	}
	if p.GroupVersion == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group version")
	}
	if p.GroupAccountVersion == 0 {
		return sdkerrors.Wrap(ErrEmpty, "group account version")
	}
	if p.Status == ProposalBase_PROPOSAL_STATUS_INVALID {
		return sdkerrors.Wrap(ErrEmpty, "status")
	}
	if p.Result == ProposalBase_PROPOSAL_RESULT_INVALID {
		return sdkerrors.Wrap(ErrEmpty, "result")
	}
	if err := p.VoteState.ValidateBasic(); err != nil {
		return errors.Wrap(err, "vote state")
	}
	if p.Timeout.Seconds == 0 && p.Timeout.Nanos == 0 {
		return sdkerrors.Wrap(ErrEmpty, "timeout")
	}
	return nil
}

func (t *Tally) Sub(vote Vote, weight sdk.Dec) error {
	if weight.LTE(sdk.ZeroDec()) {
		return errors.Wrap(ErrInvalid, "weight must be greater than 0")
	}
	switch vote.Choice {
	case Choice_YES:
		t.YesCount = t.YesCount.Sub(weight)
	case Choice_NO:
		t.NoCount = t.NoCount.Sub(weight)
	case Choice_ABSTAIN:
		t.AbstainCount = t.AbstainCount.Sub(weight)
	case Choice_VETO:
		t.VetoCount = t.VetoCount.Sub(weight)
	default:
		return errors.Wrapf(ErrInvalid, "unknown choice %s", vote.Choice.String())
	}
	return nil
}

func (t *Tally) Add(vote Vote, weight sdk.Dec) error {
	if weight.LTE(sdk.ZeroDec()) {
		return errors.Wrap(ErrInvalid, "weight must be greater than 0")
	}
	switch vote.Choice {
	case Choice_YES:
		t.YesCount = t.YesCount.Add(weight)
	case Choice_NO:
		t.NoCount = t.NoCount.Add(weight)
	case Choice_ABSTAIN:
		t.AbstainCount = t.AbstainCount.Add(weight)
	case Choice_VETO:
		t.VetoCount = t.VetoCount.Add(weight)
	default:
		return errors.Wrapf(ErrInvalid, "unknown choice %s", vote.Choice.String())
	}
	return nil
}

// TotalCounts is the sum of all weights.
func (t Tally) TotalCounts() sdk.Dec {
	return t.YesCount.Add(t.NoCount).Add(t.AbstainCount).Add(t.VetoCount)
}

func (t Tally) ValidateBasic() error {
	switch {
	case t.YesCount.IsNil():
		return errors.Wrap(ErrInvalid, "yes count nil")
	case t.YesCount.LT(sdk.ZeroDec()):
		return errors.Wrap(ErrInvalid, "yes count negative")
	case t.NoCount.IsNil():
		return errors.Wrap(ErrInvalid, "no count nil")
	case t.NoCount.LT(sdk.ZeroDec()):
		return errors.Wrap(ErrInvalid, "no count negative")
	case t.AbstainCount.IsNil():
		return errors.Wrap(ErrInvalid, "abstain count nil")
	case t.AbstainCount.LT(sdk.ZeroDec()):
		return errors.Wrap(ErrInvalid, "abstain count negative")
	case t.VetoCount.IsNil():
		return errors.Wrap(ErrInvalid, "veto count nil")
	case t.VetoCount.LT(sdk.ZeroDec()):
		return errors.Wrap(ErrInvalid, "veto count negative")
	}
	return nil
}

func (g GenesisState) String() string {
	out, _ := yaml.Marshal(g)
	return string(out)
}
