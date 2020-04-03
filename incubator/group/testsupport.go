package group

import (
	"time"

	"github.com/cosmos/cosmos-sdk/baseapp"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/rootmulti"
	"github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	subspace "github.com/cosmos/cosmos-sdk/x/params/types"
	gogo "github.com/gogo/protobuf/types"
	fuzz "github.com/google/gofuzz"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/math"
	dbm "github.com/tendermint/tm-db"
)

func NewContext(keys ...sdk.StoreKey) sdk.Context {
	db := dbm.NewMemDB()
	cms := rootmulti.NewStore(db)
	for _, v := range keys {
		storeType := sdk.StoreTypeIAVL
		if _, ok := v.(*sdk.TransientStoreKey); ok {
			storeType = sdk.StoreTypeTransient
		}
		cms.MountStoreWithDB(v, storeType, db)
	}
	cms.SetPruning(types.PruneSyncable)
	if err := cms.LoadLatestVersion(); err != nil {
		panic(err)
	}
	return sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())
}

func createGroupKeeper() (Keeper, sdk.Context) {
	cdc := codec.NewHybridCodec(codec.New())
	pKey, pTKey := sdk.NewKVStoreKey(params.StoreKey), sdk.NewTransientStoreKey(params.TStoreKey)
	paramSpace := subspace.NewSubspace(cdc, pKey, pTKey, DefaultParamspace)

	groupKey := sdk.NewKVStoreKey(StoreKeyName)
	k := NewGroupKeeper(groupKey, paramSpace, baseapp.NewRouter(), &MockProposalI{})
	ctx := NewContext(pKey, pTKey, groupKey)
	k.setParams(ctx, DefaultParams())
	return k, ctx
}

func FuzzPositiveDec(m *sdk.Dec, c fuzz.Continue) {
	*m = sdk.NewDec(c.Rand.Int63())
}
func FuzzComment(m *string, c fuzz.Continue) {
	randString := c.RandString()
	*m = randString[0:math.MinInt(len(randString), defaultMaxCommentLength)]
}
func FuzzAddr(m *sdk.AccAddress, c fuzz.Continue) {
	*m = make([]byte, 20)
	c.Read(*m)
}
func FuzzPositiveDuration(m *gogo.Duration, c fuzz.Continue) {
	v := gogo.DurationProto(time.Duration(c.Int63()))
	*m = *v
}
func FuzzMember(m *Member, c fuzz.Continue) {
	FuzzAddr(&m.Address, c)
	FuzzPositiveDec(&m.Power, c)
	FuzzComment(&m.Comment, c)
}
func FuzzGroupMember(m *GroupMember, c fuzz.Continue) {
	m.Group = GroupID(c.RandUint64())
	FuzzAddr(&m.Member, c)
	FuzzPositiveDec(&m.Weight, c)
	FuzzComment(&m.Comment, c)
}
func FuzzChoice(m *Choice, c fuzz.Continue) {
	*m = Choice(c.Intn(len(Choice_name)-2) + 1)
}

type MockProposalI struct {
}

func (m MockProposalI) Marshal() ([]byte, error) {
	panic("implement me")
}

func (m MockProposalI) Unmarshal([]byte) error {
	panic("implement me")
}

func (m MockProposalI) GetBase() ProposalBase {
	panic("implement me")
}

func (m MockProposalI) SetBase(ProposalBase) {
	panic("implement me")
}

func (m MockProposalI) GetMsgs() []sdk.Msg {
	panic("implement me")
}

func (m MockProposalI) SetMsgs([]sdk.Msg) error {
	panic("implement me")
}

func (m MockProposalI) ValidateBasic() error {
	panic("implement me")
}
