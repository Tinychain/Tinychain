package account

import (
	"errors"
	"github.com/tinychain/tinychain/common"
	"github.com/tinychain/tinychain/crypto"
)

var (
	log              = common.GetLogger("account")
	ErrCreateKeyPair = errors.New("failed to create key pair")
	ErrCreateAcc     = errors.New("failed to create account")
	ErrUnlockAcc     = errors.New("failed to unlock account in wallet")
	ErrNotFoundAcc   = errors.New("account is not in wallet. plz unlock it first")
	ErrNotUnlock     = errors.New("account has not been unlocked")
	ErrSignTx        = errors.New("failed to sign transaction")
)

type Account struct {
	Address common.Address
	key     *Key
}

// New account by private key
func NewAccountWithKey(privKey crypto.PrivateKey) (*Account, error) {
	key := &Key{privKey, privKey.GetPublic()}
	addr, err := common.GenAddrByPrivkey(privKey)
	if err != nil {
		return nil, err
	}
	return &Account{addr, key}, nil
}

func NewAccount() (*Account, error) {
	key, err := NewKeyPairs()
	if err != nil {
		log.Errorf(ErrCreateKeyPair.Error())
		return nil, ErrCreateKeyPair
	}
	addr, err := common.GenAddrByPubkey(key.pubKey)
	if err != nil {
		return nil, err
	}
	return &Account{addr, key}, nil
}

func (ac *Account) PrivKey() crypto.PrivateKey {
	return ac.key.privKey
}

func (ac *Account) PubKey() crypto.PublicKey {
	return ac.key.pubKey
}
