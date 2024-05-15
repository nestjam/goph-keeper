package model

type MasterKey struct {
	key []byte
}

func NewMasterKey(key []byte) *MasterKey {
	return &MasterKey{
		key: key,
	}
}
