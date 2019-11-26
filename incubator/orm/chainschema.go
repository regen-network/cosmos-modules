package orm

import "reflect"

type ChainSchemaKeeper struct {
	moduleTable                  AutoUInt64Table
	subStoreTable                AutoUInt64Table
	subStoreByModuleAndNameIndex Index // this should be UInt64_String_Index or something
	typeTable                    AutoUInt64Table
	typeByNameIndex              Index
}

type ModuleInfo struct {
	Name string `json:"name"`
}

type SubStoreLayout struct {
	Module    uint64    `json:"module"`
	Prefix    byte      `json:"prefix"`
	Name      string    `json:"name"`
	KeyLayout KeyLayout `json:"key_layout"`
	ValueType uint64    `json:"value_type"`
}

type KeyLayout struct {
	KeyTypes []KeyType
}

type KeyType uint

const (
	UInt64KT KeyType = 0
	BytesKT  KeyType = 1
	StringKT KeyType = 2
)

type TypeInfo struct {
	Kind   reflect.Kind `json:"kind"`
	Fields []StructField
}

type StructField struct {
}
