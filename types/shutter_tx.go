package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type ShutterTx struct {
	ChainID   *big.Int
	Nonce     uint64
	GasTipCap *big.Int
	GasFeeCap *big.Int
	Gas       uint64

	EncryptedPayload []byte
	BatchIndex       uint64
	L1BlockNumber    uint64

	// Optional, only set when decrypted
	// This is ignored in rlp encoding
	// and thus hashing
	Payload *ShutterPayload `rlp:"-"`

	// Signature values
	V *big.Int `json:"v" gencodec:"required"`
	R *big.Int `json:"r" gencodec:"required"`
	S *big.Int `json:"s" gencodec:"required"`
}

// copy creates a deep copy of the transaction data and initializes all fields.
func (tx *ShutterTx) copy() TxInner {
	cpy := &ShutterTx{
		Nonce: tx.Nonce,
		Gas:   tx.Gas,

		// These are copied below.
		EncryptedPayload: []byte{},
		ChainID:          new(big.Int),
		GasTipCap:        new(big.Int),
		GasFeeCap:        new(big.Int),
		BatchIndex:       tx.BatchIndex,
		L1BlockNumber:    tx.L1BlockNumber,
		V:                new(big.Int),
		R:                new(big.Int),
		S:                new(big.Int),
	}
	if tx.ChainID != nil {
		cpy.ChainID.Set(tx.ChainID)
	}
	if tx.GasTipCap != nil {
		cpy.GasTipCap.Set(tx.GasTipCap)
	}
	if tx.GasFeeCap != nil {
		cpy.GasFeeCap.Set(tx.GasFeeCap)
	}
	if tx.EncryptedPayload != nil {
		cpy.EncryptedPayload = common.CopyBytes(tx.EncryptedPayload)
	}
	if tx.Payload != nil {
		cpy.Payload = tx.Payload.Copy()
	}
	if tx.V != nil {
		cpy.V.Set(tx.V)
	}
	if tx.R != nil {
		cpy.R.Set(tx.R)
	}
	if tx.S != nil {
		cpy.S.Set(tx.S)
	}
	return cpy
}

// accessors for innerTx.
func (tx *ShutterTx) txType() byte           { return ShutterTxType }
func (tx *ShutterTx) chainID() *big.Int      { return tx.ChainID }
func (tx *ShutterTx) protected() bool        { return true }
func (tx *ShutterTx) accessList() AccessList { return nil }
func (tx *ShutterTx) data() []byte {
	if tx.Payload != nil {
		return tx.Payload.Data
	}
	return []byte{}
}
func (tx *ShutterTx) gas() uint64         { return tx.Gas }
func (tx *ShutterTx) gasFeeCap() *big.Int { return tx.GasFeeCap }
func (tx *ShutterTx) gasTipCap() *big.Int { return tx.GasTipCap }
func (tx *ShutterTx) gasPrice() *big.Int  { return tx.GasFeeCap }
func (tx *ShutterTx) value() *big.Int {
	if tx.Payload != nil {
		return tx.Payload.Value
	}
	return big.NewInt(0)
}

func (tx *ShutterTx) nonce() uint64 { return tx.Nonce }
func (tx *ShutterTx) to() *common.Address {
	if tx.Payload == nil {
		return nil
	}
	return tx.Payload.To
}
func (tx *ShutterTx) encryptedPayload() []byte { return tx.EncryptedPayload }
func (tx *ShutterTx) decryptionKey() []byte    { return nil }
func (tx *ShutterTx) batchIndex() uint64       { return tx.BatchIndex }
func (tx *ShutterTx) l1BlockNumber() uint64    { return tx.L1BlockNumber }
func (tx *ShutterTx) timestamp() *big.Int      { return nil }
func (tx *ShutterTx) transactions() [][]byte   { return nil }

func (tx *ShutterTx) rawSignatureValues() (v, r, s *big.Int) {
	return tx.V, tx.R, tx.S
}

func (tx *ShutterTx) setSignatureValues(chainID, v, r, s *big.Int) {
	tx.ChainID, tx.V, tx.R, tx.S = chainID, v, r, s
}
