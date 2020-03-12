package orm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	abci "github.com/tendermint/tendermint/abci/types"
)

type Queryable interface {
	Index
	GetModelType() reflect.Type
}

type queryHandler struct {
	Routes map[string]Queryable
}

func NewQueryHandler() *queryHandler {
	return &queryHandler{Routes: make(map[string]Queryable)}
}

func (h *queryHandler) AddTableRoute(path string, table Table) {
	if h.Has(path){
		panic(fmt.Sprintf("path %q exists already", path))
	}
	h.Routes[path] = QueryTableAdapter(table)
}

func (h *queryHandler) AddIndexRoute(path string, index Index) {
	if h.Has(path) {
		panic(fmt.Sprintf("path %q exists already", path))
	}
	tp, err := GetModelType(index)
	if err != nil {
		panic(err)
	}
	h.Routes[path] = QueryIndexAdapter(index, tp)
}

func (h *queryHandler) Has(path string) bool{
	_, ok := h.Routes[path]
	return ok
}

func (h *queryHandler) Handle(ctx sdk.Context, _ []string, req abci.RequestQuery) ([]byte, error) {
	path, mod := splitPath(req.Path)
	queryIndex, ok := h.Routes[path]
	if !ok {
		return nil, errors.Wrapf(errors.ErrUnknownRequest, "unknown query path: %s", path)
	}
	it, err := DoQuery(ctx, queryIndex, mod, req.Data)
	if err != nil {
		return nil, err
	}
	result, err := buildResult(queryIndex, it)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&result)
}

// splitPath splits out the real path along with the query
// modifier (everything after the ?)
func splitPath(path string) (string, string) {
	var mod string
	chunks := strings.SplitN(path, "?", 2)
	if len(chunks) == 2 {
		path = chunks[0]
		mod = chunks[1]
	}
	return path, mod
}

type RangeQueryParam struct {
	start, end []byte
}

func DoQuery(ctx HasKVStore, index Index, mod string, data []byte) (Iterator, error) {
	switch mod {
	case KeyQueryMod:
		return index.Get(ctx, data)
	case PrefixQueryMod:
		if len(data) == 0 {
			return nil, errors.Wrap(ErrArgument, "prefix")
		}
		start, end := prefixRange(data)
		return index.PrefixScan(ctx, start, end)
	case RangeQueryMod:
		var param RangeQueryParam
		if len(data) != 0 {
			if err := json.Unmarshal(data, &param); err != nil {
				return nil, errors.Wrap(err, "unmarshal param")
			}
		}
		return index.PrefixScan(ctx, param.start, param.end)
	default:
		return nil, errors.Wrapf(ErrArgument, "unknown mod: %s", mod)
	}
}


type QueryResult struct {
	Data    []Model `json:"data"`
	HasMore bool        `json:"has_more"`
}

func buildResult(queryIndex Queryable, it Iterator) (QueryResult, error) {
	limitIt := LimitIterator(it, maxQueryResult)
	defer it.Close()

	mrsh := jsonpb.Marshaler{}
	var data []Model
OUTER:
	for {
		d := reflect.New(queryIndex.GetModelType()).Interface().(Persistent)
		id, err := limitIt.LoadNext(d)
		switch {
		case ErrIteratorDone.Is(err):
			break OUTER
		case err != nil:
			return QueryResult{}, err
		default:
			obj, ok := d.(proto.Message)
			if !ok {
				return QueryResult{}, errors.Wrap(ErrArgument, "not proto message")
			}
			var buf bytes.Buffer
			if err := mrsh.Marshal(&buf, obj); err != nil {
				return QueryResult{}, errors.Wrap(err, "marshal group entity from proto")
			}
			data = append(data, Model{Key: id, Value: buf.Bytes()})
		}
	}
	return QueryResult{Data: data, HasMore: limitIt.Remaining() == 0}, nil
}