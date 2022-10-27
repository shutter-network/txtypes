// Copyright 2021 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/pkg/errors"
)

type TransactionData struct {
	Type hexutil.Uint64 `json:"type"`

	// LegacyTx
	From     *common.Address `json:"from"`
	Gas      *hexutil.Uint64 `json:"gas"`
	GasPrice *hexutil.Big    `json:"gasPrice"`
	Hash     common.Hash     `json:"hash"`
	Input    *hexutil.Bytes  `json:"input"`
	Nonce    *hexutil.Uint64 `json:"nonce"`
	To       *common.Address `json:"to"`
	Value    *hexutil.Big    `json:"value"`

	// AccessListTx
	AccessList *AccessList  `json:"accessList,omitempty"`
	ChainID    *hexutil.Big `json:"chainId,omitempty"`

	// DynamicFeeTx
	MaxFeePerGas         *hexutil.Big `json:"maxFeePerGas,omitempty"`
	MaxPriorityFeePerGas *hexutil.Big `json:"maxPriorityFeePerGas,omitempty"`

	// ShutterTx
	EncryptedPayload *hexutil.Bytes `json:"encryptedPayload,omitempty"`

	// BatchTx
	DecryptionKey *hexutil.Bytes  `json:"decryptionKey,omitempty"`
	Timestamp     *hexutil.Big    `json:"timestamp,omitempty"`
	Transactions  []hexutil.Bytes `json:"transactions,omitempty"`

	// ShutterTx and BatchTx
	BatchIndex    *hexutil.Uint64 `json:"batchIndex,omitempty"`
	L1BlockNumber *hexutil.Uint64 `json:"l1BlockNumber,omitempty"`

	// Optional information for included transactions
	BlockHash        *common.Hash    `json:"blockHash,omitempty"`
	BlockNumber      *hexutil.Big    `json:"blockNumber,omitempty"`
	TransactionIndex *hexutil.Uint64 `json:"transactionIndex,omitempty"`

	V *hexutil.Big `json:"v"`
	R *hexutil.Big `json:"r"`
	S *hexutil.Big `json:"s"`
}

func (td *TransactionData) ValidateRequiredFields(names ...string) (bool, error) {
	rv := reflect.ValueOf(td)
	rv = rv.Elem()

	inValidFieldNames := make([]string, 0)

	for _, fieldName := range names {
		f := rv.FieldByName(fieldName)
		if (f == reflect.Value{}) {
			// field was not found, programming error!
			return false, errors.Errorf("field '%s' is not defined in type", fieldName)
		}
		if !f.IsValid() || f.IsNil() {
			inValidFieldNames = append(inValidFieldNames, fieldName)
		}
	}
	if len(inValidFieldNames) > 0 {
		errStr := "required fields are nil: "
		for _, n := range inValidFieldNames {
			errStr = errStr + fmt.Sprintf(" %s,", n)
		}
		return false, errors.New(errStr)
	}
	return true, nil
}

func (t *Transaction) MarshalJSON() ([]byte, error) {
	// MarshalJSON marshals as JSON with a hash.
	enc := t.TransactionData()
	return json.Marshal(enc)
}

func (t *Transaction) TransactionData() (enc *TransactionData) {
	// These are set for all tx types.
	enc = new(TransactionData)
	enc.Hash = t.Hash()
	enc.Type = hexutil.Uint64(t.Type())

	// Other fields are set conditionally depending on tx type.
	switch tx := t.inner.(type) {
	case *LegacyTx:
		enc.Nonce = (*hexutil.Uint64)(&tx.Nonce)
		enc.Gas = (*hexutil.Uint64)(&tx.Gas)
		enc.GasPrice = (*hexutil.Big)(tx.GasPrice)
		enc.Value = (*hexutil.Big)(tx.Value)
		enc.Input = (*hexutil.Bytes)(&tx.Data)
		enc.To = t.To()
		enc.V = (*hexutil.Big)(tx.V)
		enc.R = (*hexutil.Big)(tx.R)
		enc.S = (*hexutil.Big)(tx.S)
	case *AccessListTx:
		enc.ChainID = (*hexutil.Big)(tx.ChainID)
		enc.AccessList = &tx.AccessList
		enc.Nonce = (*hexutil.Uint64)(&tx.Nonce)
		enc.Gas = (*hexutil.Uint64)(&tx.Gas)
		enc.GasPrice = (*hexutil.Big)(tx.GasPrice)
		enc.Value = (*hexutil.Big)(tx.Value)
		enc.Input = (*hexutil.Bytes)(&tx.Data)
		enc.To = t.To()
		enc.V = (*hexutil.Big)(tx.V)
		enc.R = (*hexutil.Big)(tx.R)
		enc.S = (*hexutil.Big)(tx.S)
	case *DynamicFeeTx:
		enc.ChainID = (*hexutil.Big)(tx.ChainID)
		enc.AccessList = &tx.AccessList
		enc.Nonce = (*hexutil.Uint64)(&tx.Nonce)
		enc.Gas = (*hexutil.Uint64)(&tx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(tx.GasFeeCap)
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(tx.GasTipCap)
		enc.Value = (*hexutil.Big)(tx.Value)
		enc.Input = (*hexutil.Bytes)(&tx.Data)
		enc.To = t.To()
		enc.V = (*hexutil.Big)(tx.V)
		enc.R = (*hexutil.Big)(tx.R)
		enc.S = (*hexutil.Big)(tx.S)
	case *ShutterTx:
		enc.ChainID = (*hexutil.Big)(tx.ChainID)
		enc.Nonce = (*hexutil.Uint64)(&tx.Nonce)
		enc.Gas = (*hexutil.Uint64)(&tx.Gas)
		enc.MaxFeePerGas = (*hexutil.Big)(tx.GasFeeCap)
		enc.MaxPriorityFeePerGas = (*hexutil.Big)(tx.GasTipCap)
		enc.EncryptedPayload = (*hexutil.Bytes)(&tx.EncryptedPayload)
		enc.L1BlockNumber = (*hexutil.Uint64)(&tx.L1BlockNumber)
		enc.BatchIndex = (*hexutil.Uint64)(&tx.BatchIndex)
		enc.To = t.To()
		enc.V = (*hexutil.Big)(tx.V)
		enc.R = (*hexutil.Big)(tx.R)
		enc.S = (*hexutil.Big)(tx.S)
	case *BatchTx:
		enc.ChainID = (*hexutil.Big)(tx.ChainID)
		if tx.Transactions != nil {
			enc.Transactions = make([]hexutil.Bytes, len(tx.Transactions))
			for k, v := range tx.Transactions {
				enc.Transactions[k] = hexutil.Bytes(v)
			}
		}
		enc.Timestamp = (*hexutil.Big)(tx.Timestamp)
		enc.DecryptionKey = (*hexutil.Bytes)(&tx.DecryptionKey)
		enc.L1BlockNumber = (*hexutil.Uint64)(&tx.L1BlockNumber)
		enc.BatchIndex = (*hexutil.Uint64)(&tx.BatchIndex)
		enc.To = t.To()
		enc.V = (*hexutil.Big)(tx.V)
		enc.R = (*hexutil.Big)(tx.R)
		enc.S = (*hexutil.Big)(tx.S)
	}
	return enc
}

// UnmarshalJSON unmarshals from JSON.
func (t *Transaction) UnmarshalJSON(input []byte) error {
	var dec TransactionData
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	return t.FromTransactionData(&dec)
}

func (t *Transaction) FromTransactionData(dec *TransactionData) error {
	// Decode / verify fields according to transaction type.
	var inner TxInner
	switch dec.Type {
	case LegacyTxType:
		var itx LegacyTx
		inner = &itx
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, true); err != nil {
				return err
			}
		}

	case AccessListTxType:
		var itx AccessListTx
		inner = &itx
		// Access list is optional for now.
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.GasPrice == nil {
			return errors.New("missing required field 'gasPrice' in transaction")
		}
		itx.GasPrice = (*big.Int)(dec.GasPrice)
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' in transaction")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return err
			}
		}

	case DynamicFeeTxType:
		var itx DynamicFeeTx
		inner = &itx
		// Access list is optional for now.
		if dec.AccessList != nil {
			itx.AccessList = *dec.AccessList
		}
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.To != nil {
			itx.To = dec.To
		}
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.MaxPriorityFeePerGas == nil {
			return errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = (*big.Int)(dec.MaxPriorityFeePerGas)
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = (*big.Int)(dec.MaxFeePerGas)
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.Value == nil {
			return errors.New("missing required field 'value' in transaction")
		}
		itx.Value = (*big.Int)(dec.Value)
		if dec.Input == nil {
			return errors.New("missing required field 'input' in transaction")
		}
		itx.Data = *dec.Input
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return err
			}
		}
	case ShutterTxType:
		var itx ShutterTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)
		if dec.Nonce == nil {
			return errors.New("missing required field 'nonce' in transaction")
		}
		itx.Nonce = uint64(*dec.Nonce)
		if dec.MaxPriorityFeePerGas == nil {
			return errors.New("missing required field 'maxPriorityFeePerGas' for txdata")
		}
		itx.GasTipCap = (*big.Int)(dec.MaxPriorityFeePerGas)
		if dec.MaxFeePerGas == nil {
			return errors.New("missing required field 'maxFeePerGas' for txdata")
		}
		itx.GasFeeCap = (*big.Int)(dec.MaxFeePerGas)
		if dec.Gas == nil {
			return errors.New("missing required field 'gas' for txdata")
		}
		itx.Gas = uint64(*dec.Gas)
		if dec.L1BlockNumber == nil {
			return errors.New("missing required field 'l1BlockNumber' in transaction")
		}
		itx.L1BlockNumber = uint64(*dec.L1BlockNumber)

		if dec.EncryptedPayload == nil {
			return errors.New("missing required field 'encryptedPayload' in transaction")
		}

		itx.EncryptedPayload = *dec.EncryptedPayload

		if dec.BatchIndex == nil {
			return errors.New("missing required field 'batchIndex' in transaction")
		}
		itx.BatchIndex = uint64(*dec.BatchIndex)
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}

		hasTo := bool(dec.To != nil)
		hasValue := bool(dec.Value != nil)
		hasInput := bool(dec.Input != nil)
		if hasTo || hasValue || hasInput {
			itx.Payload = &ShutterPayload{
				To: dec.To,
			}
			if hasInput {
				// optional
				itx.Payload.Data = *dec.Input
			}
			if !hasValue {
				// this is only required when there are other payload values set
				return errors.New("missing required nested field 'value' in transaction payload")
			}
			itx.Payload.Value = dec.Value.ToInt()
		}

		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return err
			}
		}
	case BatchTxType:
		var itx BatchTx
		inner = &itx
		if dec.ChainID == nil {
			return errors.New("missing required field 'chainId' in transaction")
		}
		itx.ChainID = (*big.Int)(dec.ChainID)

		if dec.Timestamp == nil {
			return errors.New("missing required field 'timestamp' in transaction")
		}
		itx.Timestamp = (*big.Int)(dec.Timestamp)

		if dec.Transactions == nil {
			return errors.New("missing required field 'transactions' in transaction")
		}
		itx.Transactions = make([][]byte, len(dec.Transactions))
		for i, txx := range dec.Transactions {
			itx.Transactions[i] = []byte(txx)
		}

		if dec.L1BlockNumber == nil {
			return errors.New("missing required field 'l1BlockNumber' in transaction")
		}
		itx.L1BlockNumber = uint64(*dec.L1BlockNumber)

		if dec.BatchIndex == nil {
			return errors.New("missing required field 'batchIndex' in transaction")
		}
		itx.BatchIndex = uint64(*dec.BatchIndex)
		if dec.V == nil {
			return errors.New("missing required field 'v' in transaction")
		}
		itx.V = (*big.Int)(dec.V)
		if dec.R == nil {
			return errors.New("missing required field 'r' in transaction")
		}
		itx.R = (*big.Int)(dec.R)
		if dec.S == nil {
			return errors.New("missing required field 's' in transaction")
		}
		itx.S = (*big.Int)(dec.S)
		withSignature := itx.V.Sign() != 0 || itx.R.Sign() != 0 || itx.S.Sign() != 0
		if withSignature {
			if err := sanityCheckSignature(itx.V, itx.R, itx.S, false); err != nil {
				return err
			}
		}

	default:
		return ErrTxTypeNotSupported
	}

	// Now set the inner transaction.
	t.setDecoded(inner, 0)

	// TODO: check hash here?
	return nil
}
