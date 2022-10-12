package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type DecryptedPayload struct {
	To    *common.Address `rlp:"nil"`
	Data  []byte
	Value *big.Int
}

func (p *DecryptedPayload) AsMessage(tx *Transaction, signer Signer) (Message, error) {
	sender, err := signer.Sender(tx)
	if err != nil {
		return Message{}, err
	}
	return NewMessage(
		sender,        // from
		p.To,          // to
		tx.Nonce(),    // nonce
		p.Value,       // amount
		tx.Gas(),      // gas limit
		big.NewInt(0), // gas price
		big.NewInt(0), // gas fee cap
		big.NewInt(0), // gas tip cap
		p.Data,        // data
		nil,           // access list
		false,         // is fake
	), nil
}

func (p *DecryptedPayload) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(*p)
}

func DecodeDecryptedPayload(b []byte) (*DecryptedPayload, error) {
	p := &DecryptedPayload{}
	err := rlp.DecodeBytes(b, p)
	return p, err
}
