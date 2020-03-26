package group

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/modules/incubator/orm"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
	"google.golang.org/grpc"
)

type QuerierServer struct {
	sd *grpc.ServiceDesc
	ss interface{}
}

func (q *QuerierServer) RegisterService(sd *grpc.ServiceDesc, ss interface{}) {
	q.sd = sd
	q.ss = ss
}

func (q *QuerierServer) Querier() sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		path0 := path[0]
		for _, md := range q.sd.Methods {
			if md.MethodName == path0 {
				res, err := md.Handler(q.ss, WithSDKContext(ctx.Context(), ctx), func(i interface{}) error {
					protoMsg := i.(proto.Message)
					return proto.Unmarshal(req.Data, protoMsg)
				}, nil)
				if err != nil {
					return nil, err
				}
				protoMsg := res.(proto.Message)
				return proto.Marshal(protoMsg)
			}
		}
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown query path: %s", path[0])
	}
}

var _ Server = &QuerierServer{}

type Querier struct {
	Keeper
}

type sdkContextKeyType string

const sdkContextKey sdkContextKeyType = "sdk-context"

func SDKContext(ctx context.Context) sdk.Context {
	return ctx.Value(sdkContextKey).(sdk.Context)
}

func WithSDKContext(ctx context.Context, sdkCtx sdk.Context) context.Context {
	return context.WithValue(ctx, sdkContextKey, sdkCtx)
}

func (q Querier) GetGroupMetadata(ctx context.Context, req *GroupMetadataRequest) (res *GroupMetadata, err error) {
	err = q.groupTable.GetOne(SDKContext(ctx), orm.EncodeSequence(req.Group), res)
	return res, err
}

func (q Querier) GetGroupAccountMetadata(context.Context, *GroupAccountMetadataRequest) (*AnyGroupAccountMetadata, error) {
	panic("implement me")
}

func (q Querier) GetGroupsByMember(context.Context, *ByAddressRequest) (*GroupsList, error) {
	panic("implement me")
}

func (q Querier) GetGroupsByOwner(context.Context, *ByAddressRequest) (*GroupsList, error) {
	panic("implement me")
}

func (q Querier) GetGroupAccountsByGroup(context.Context, *ByGroupRequest) (*GroupAccountsList, error) {
	panic("implement me")
}

func (q Querier) GetGroupAccountsByOwner(context.Context, *ByAddressRequest) (*GroupAccountsList, error) {
	panic("implement me")
}

func (q Querier) GetProposal(context.Context, *ProposalRequest) (*AnyProposal, error) {
	panic("implement me")
}

func (q Querier) GetProposalsByGroupAccount(context.Context, *ByGroupAccountRequest) (*ProposalList, error) {
	panic("implement me")
}

func (q Querier) GetVotes(context.Context, *VotesRequest) (*VoteList, error) {
	panic("implement me")
}

var _ QueryServer = Querier{}
