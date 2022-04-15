package gm

import (
	"crypto/elliptic"
	"errors"
	"fmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/internal/github.com/hyperledger/fabric/bccsp"
	"github.com/tjfoc/gmsm/sm2"
	"github.com/tjfoc/gmsm/sm3"
	"github.com/tjfoc/gmsm/x509"
)

// 实现bccsp Key接口
type gmsm2PrivateKey struct {
	privKey *sm2.PrivateKey
}

func (k *gmsm2PrivateKey) Bytes() ([]byte, error) {
	return nil, errors.New("Not supported.")
}

func (k *gmsm2PrivateKey) SKI() []byte {
	if k.privKey == nil {
		return nil
	}

	//Marshall the public key
	raw := elliptic.Marshal(k.privKey.Curve, k.privKey.PublicKey.X, k.privKey.PublicKey.Y)

	// Hash it
	hash := sm3.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

func (k *gmsm2PrivateKey) Symmetric() bool {
	return false
}

func (k *gmsm2PrivateKey) Private() bool {
	return true
}

func (k *gmsm2PrivateKey) PublicKey() (bccsp.Key, error) {
	return &gmsm2PublicKey{&k.privKey.PublicKey}, nil
}

type gmsm2PublicKey struct {
	pubKey *sm2.PublicKey
}

func (k *gmsm2PublicKey) Bytes() (raw []byte, err error) {
	raw, err = x509.MarshalSm2PublicKey(k.pubKey)
	if err != nil {
		return nil, fmt.Errorf("Failed marshalling key [%s]", err)
	}
	return
}

func (k *gmsm2PublicKey) SKI() []byte {
	if k.pubKey == nil {
		return nil
	}

	//Marshall the public key
	raw := elliptic.Marshal(k.pubKey.Curve, k.pubKey.X, k.pubKey.Y)

	// Hash it
	hash := sm3.New()
	hash.Write(raw)
	return hash.Sum(nil)
}

func (k *gmsm2PublicKey) Symmetric() bool {
	return false
}

func (k *gmsm2PublicKey) Private() bool {
	return false
}

func (k *gmsm2PublicKey) PublicKey() (bccsp.Key, error) {
	return k, nil
}
