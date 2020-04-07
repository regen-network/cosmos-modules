package testdata

import (
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func TxDecoder(cdc codec.Marshaler) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, error) {
		var tx Transaction

		if len(txBytes) == 0 {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, "tx bytes are empty")
		}

		err := cdc.UnmarshalBinaryBare(txBytes, &tx)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrTxDecode, err.Error())
		}

		return &tx, nil
	}
}

type TXFactory struct{}

func (x TXFactory) NewTx() tx.ClientTx {
	return &ClientTXAdaptor{Transaction{}}
}

// ClientTXAdaptor should not exist in the ideal world.
type ClientTXAdaptor struct {
	Transaction
}

func (c ClientTXAdaptor) GetSignatures() []sdk.Signature {
	r := make([]sdk.Signature, len(c.Transaction.Signatures))
	for i := range c.Transaction.Signatures {
		r[i] = c.Transaction.Signatures[i]
	}
	return r
}

func (c ClientTXAdaptor) GetFee() sdk.Fee {
	return c.Transaction.Fee
}
