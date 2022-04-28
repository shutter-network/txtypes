package types

type TxDataExtension interface {
	encryptedPayload() []byte
	decryptionKey() []byte
	batchIndex() []byte
}

func (tx *Transaction) EncryptedPayload() []byte {
	return tx.inner.encryptedPayload()
}

func (tx *Transaction) DecryptionKey() []byte {
	return tx.inner.decryptionKey()
}

func (tx *Transaction) BatchIndex() []byte {
	return tx.inner.batchIndex()
}

func (*AccessListTx) encryptedPayload() []byte {
	return nil
}

func (*AccessListTx) decryptionKey() []byte {
	return nil
}

func (*AccessListTx) batchIndex() []byte {
	return nil
}

func (*DynamicFeeTx) encryptedPayload() []byte {
	return nil
}

func (*DynamicFeeTx) decryptionKey() []byte {
	return nil
}

func (*DynamicFeeTx) batchIndex() []byte {
	return nil
}

func (*LegacyTx) encryptedPayload() []byte {
	return nil
}

func (*LegacyTx) decryptionKey() []byte {
	return nil
}

func (*LegacyTx) batchIndex() []byte {
	return nil
}
