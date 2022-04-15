package gm

import (
	"errors"
	"git.querycap.com/cloudchain/fabric-sdk-go/internal/github.com/hyperledger/fabric/bccsp"
	"github.com/tjfoc/gmsm/sm3"
)

// 实现bccsp Key接口
type gmsm4PrivateKey struct {
	privKey    []byte
	exportable bool
}

// Bytes converts this key to its byte representation,
// if this operation is allowed.
func (k *gmsm4PrivateKey) Bytes() (raw []byte, err error) {
	if k.exportable {
		return k.privKey, nil
	}

	return nil, errors.New("Not supported")
}

// SKI returns the subject key identifier of this key.
func (k *gmsm4PrivateKey) SKI() (ski []byte) {
	hash := sm3.New()
	hash.Write([]byte{0x01})
	hash.Write(k.privKey)
	return hash.Sum(nil)
}

// Symmetric returns true if this key is a symmetric key,
// false if this key is asymmetric
func (k *gmsm4PrivateKey) Symmetric() bool {
	return true
}

// Private returns true if this key is a private key,
// false otherwise.
func (k *gmsm4PrivateKey) Private() bool {
	return true
}

// PublicKey returns the corresponding public key part of an asymmetric public/private key pair.
// This method returns an error in symmetric key schemes.
func (k *gmsm4PrivateKey) PublicKey() (bccsp.Key, error) {
	return nil, errors.New("Cannot call this method on a symmetric key")
}
