package group

import "github.com/cosmos/modules/incubator/orm"

// AccountCondition returns a condition to build a group account address.
func AccountCondition(id uint64) Condition {
	return NewCondition("group", "account", orm.EncodeSequence(id))
}
