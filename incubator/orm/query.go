package orm

import (
	"fmt"
	"reflect"
)

const (
	// KeyQueryMod means to query for exact match (key)
	KeyQueryMod = ""
	// PrefixQueryMod means to query for anything with this prefix
	PrefixQueryMod = "prefix"
	// RangeQueryMod means to expect complex range query
	// TODO: implement
	RangeQueryMod = "range"
)

const maxQueryResult = 50


type queryTableAdapter struct {
	table Table
}

func QueryTableAdapter(s TableExportable) queryTableAdapter {
	return queryTableAdapter{table: s.Table()}
}

func (t queryTableAdapter) Has(ctx HasKVStore, key []byte) bool {
	return t.table.Has(ctx, key)
}

func (t queryTableAdapter) Get(ctx HasKVStore, searchKey []byte) (Iterator, error) {
	return NewSingleValueIterator(ctx, TableRowGetter(t.table), searchKey), nil
}

func (t queryTableAdapter) PrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	return t.table.PrefixScan(ctx, start, end)
}

func (t queryTableAdapter) ReversePrefixScan(ctx HasKVStore, start []byte, end []byte) (Iterator, error) {
	return t.table.ReversePrefixScan(ctx, start, end)
}

func (t queryTableAdapter) GetModelType() reflect.Type {
	return t.table.model
}

type queryIndexAdapter struct {
	Index
	model reflect.Type
}

func QueryIndexAdapter(i Index, model reflect.Type) queryIndexAdapter {
	return queryIndexAdapter{Index: i, model: model}
}
func (t queryIndexAdapter) GetModelType() reflect.Type {
	return t.model
}

// ModelTypeExportable this is an extension point for custom index implementations
type ModelTypeExportable interface{
	Type() reflect.Type
}

func GetModelType(s interface{}) (reflect.Type, error) {
	switch obj := s.(type) {
	case TableExportable:
		return obj.Table().model, nil
	case MultiKeyIndex:
		return obj.rowGetter.Type(), nil
	case UniqueIndex:
		return obj.rowGetter.Type(), nil
	case ModelTypeExportable:
		return obj.Type(), nil
	}
	return nil, fmt.Errorf("unsupported type %T", s)
}
