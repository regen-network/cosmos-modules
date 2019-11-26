package orm

type sequence struct {
	key StoreKey
}

func (s sequence) NextVal(ctx HasKVStore) (uint64, error) {
	panic("implement me")
}

func (s sequence) CurVal(ctx HasKVStore) (uint64, error) {
	panic("implement me")
}

