package orm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/url"
	"reflect"
	"strconv"
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
	if h.Has(path) {
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

func (h *queryHandler) Has(path string) bool {
	_, ok := h.Routes[path]
	return ok
}

func (h *queryHandler) Handle(ctx sdk.Context, _ []string, req abci.RequestQuery) ([]byte, error) {
	q, err := splitPath(req.Path)
	if err != nil {
		return nil, err
	}
	queryIndex, ok := h.Routes[q.path]
	if !ok {
		return nil, errors.Wrapf(errors.ErrUnknownRequest, "unknown query path: %s", q.path)
	}
	it, err := DoQuery(ctx, queryIndex, q.mod, q.cursor, req.Data)
	if err != nil {
		return nil, err
	}
	result, err := buildResult(queryIndex, it, q.limit)
	if err != nil {
		return nil, err
	}
	return json.Marshal(&result)
}

type queryArgs struct {
	path   string
	mod    string
	cursor Cursor
	limit  uint8
}

// splitPath splits out the real path along with the query
// modifier (everything after the ?)
func splitPath(path string) (queryArgs, error) {
	q := queryArgs{mod: KeyQueryMod, limit: maxQueryResult}

	chunks := strings.SplitN(path, "?", 2)
	if len(chunks) != 2 {
		q.path = path
		return q, nil
	}
	q.path = chunks[0]
	query, err := url.ParseQuery(chunks[1])
	if err != nil {
		return q, err
	}
	q.cursor, err = ParseCursor(query.Get(CursorQueryArg))
	if err != nil {
		return q, errors.Wrap(err, "cursor")
	}
	for _, m := range []string{PrefixQueryMod, RangeQueryMod} {
		if _, ok := query[m]; ok {
			q.mod = m
			break
		}
	}

	if rawLimit := query.Get("limit"); rawLimit != "" {
		l, err := strconv.ParseUint(rawLimit, 10, 8)
		if err != nil {
			return q, errors.Wrap(err, "limit")
		}
		if new := uint8(l); new > 0 && new < q.limit {
			q.limit = new
		}
	}
	return q, nil
}

type RangeQueryParam struct {
	start, end []byte
}

func DoQuery(ctx HasKVStore, index Index, mod string, cursor Cursor, data []byte) (Iterator, error) {
	switch mod {
	case KeyQueryMod:
		return index.Get(ctx, data)
	case PrefixQueryMod:
		if len(data) == 0 {
			return nil, errors.Wrap(ErrArgument, "prefix")
		}
		start, end := prefixRange(data)
		if !cursor.InRange(start, end) {
			return nil, errors.Wrap(ErrArgument, "cursor not in range")
		}
		return index.PrefixScan(ctx, cursor.Move(start), end)
	case RangeQueryMod:
		var param RangeQueryParam
		if len(data) != 0 {
			if err := json.Unmarshal(data, &param); err != nil {
				return nil, errors.Wrap(err, "unmarshal param")
			}
		}
		if !cursor.InRange(param.start, param.end) {
			return nil, errors.Wrap(ErrArgument, "cursor not in range")
		}
		return index.PrefixScan(ctx, cursor.Move(param.start), param.end)
	default:
		return nil, errors.Wrapf(ErrArgument, "unknown mod: %s", mod)
	}
}

type QueryResult struct {
	Data    []Model `json:"data"`
	HasMore bool    `json:"has_more"`
	Cursor  Cursor  `json:"cursor,omitempty"`
}

func buildResult(queryIndex Queryable, it Iterator, limit uint8) (QueryResult, error) {
	limitIt := LimitIterator(it, int(limit))
	defer limitIt.Close()

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
		}
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
	cursor := it.NextPosition()
	hasMore := cursor != nil
	return QueryResult{Data: data, HasMore: hasMore, Cursor: cursor}, nil
}
