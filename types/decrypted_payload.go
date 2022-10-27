package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
)

type ShutterPayload struct {
	To    *common.Address `rlp:"nil"`
	Data  []byte
	Value *big.Int
}

func (p *ShutterPayload) Copy() *ShutterPayload {
	cpy := new(ShutterPayload)
	cpy.Data = common.CopyBytes(p.Data)
	if p.To != nil {
		addr := common.BytesToAddress(p.To.Bytes())
		cpy.To = &addr
	}
	if p.Value != nil {
		cpy.Value = new(big.Int).Set(p.Value)
	}
	return cpy
}

func (p *ShutterPayload) Encode() ([]byte, error) {
	return rlp.EncodeToBytes(*p)
}

func DecodeShutterPayload(b []byte) (*ShutterPayload, error) {
	p := &ShutterPayload{}
	err := rlp.DecodeBytes(b, p)
	return p, err
}
