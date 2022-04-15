package gm

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"github.com/tjfoc/gmsm/x509"
	"reflect"

	"git.querycap.com/cloudchain/fabric-sdk-go/internal/github.com/hyperledger/fabric/bccsp"
	"git.querycap.com/cloudchain/fabric-sdk-go/internal/github.com/hyperledger/fabric/bccsp/utils"

	"github.com/tjfoc/gmsm/sm2"
)

type x509PublicKeyImportOptsKeyImporter struct {
	bccsp *CSP
}

func (ki *x509PublicKeyImportOptsKeyImporter) KeyImport(raw interface{}, opts bccsp.KeyImportOpts) (bccsp.Key, error) {
	//x509Cert, ok := raw.(*x509.Certificate)
	sm2Cert, ok := raw.(*x509.Certificate)
	if !ok {
		return nil, errors.New("Invalid raw material. Expected *gm.Certificate.")
	}

	//pk := x509Cert.PublicKey
	pk := sm2Cert.PublicKey

	switch pk.(type) {
	case sm2.PublicKey:
		fmt.Printf("bccsp gm keyimport pk is sm2.PublicKey")
		sm2PublicKey, ok := pk.(sm2.PublicKey)
		if !ok {
			return nil, errors.New("Parse interface []  to sm2 pk error")
		}
		der, err := x509.MarshalSm2PublicKey(&sm2PublicKey)
		if err != nil {
			return nil, errors.New("MarshalSm2PublicKey error")
		}

		return ki.bccsp.KeyImporters[reflect.TypeOf(&bccsp.GMSM2PublicKeyImportOpts{})].KeyImport(
			der,
			&bccsp.GMSM2PublicKeyImportOpts{Temporary: opts.Ephemeral()})
	case *sm2.PublicKey:
		fmt.Printf("bccsp gm keyimport pk is *sm2.PublicKey")
		return ki.bccsp.KeyImporters[reflect.TypeOf(&bccsp.GMSM2PublicKeyImportOpts{})].KeyImport(
			pk,
			&bccsp.GMSM2PublicKeyImportOpts{Temporary: opts.Ephemeral()})
	// todo
	case *ecdsa.PublicKey:
		return ki.bccsp.KeyImporters[reflect.TypeOf(&bccsp.ECDSAGoPublicKeyImportOpts{})].KeyImport(
			pk,
			&bccsp.ECDSAGoPublicKeyImportOpts{Temporary: opts.Ephemeral()})
	default:
		return nil, errors.New("Certificate's public key type not recognized. Supported keys: [GMSM2]")
	}
}

//实现内部的 KeyImporter 接口
type gmsm4ImportKeyOptsKeyImporter struct{}

func (*gmsm4ImportKeyOptsKeyImporter) KeyImport(raw interface{}, opts bccsp.KeyImportOpts) (k bccsp.Key, err error) {
	sm4Raw, ok := raw.([]byte)
	if !ok {
		return nil, errors.New("Invalid raw material. Expected byte array.")
	}

	if sm4Raw == nil {
		return nil, errors.New("Invalid raw material. It must not be nil.")
	}

	return &gmsm4PrivateKey{utils.Clone(sm4Raw), false}, nil
}

type gmsm2PrivateKeyImportOptsKeyImporter struct{}

func (*gmsm2PrivateKeyImportOptsKeyImporter) KeyImport(raw interface{}, opts bccsp.KeyImportOpts) (k bccsp.Key, err error) {

	der, ok := raw.([]byte)
	if !ok {
		return nil, errors.New("[GMSM2PrivateKeyImportOpts] Invalid raw material. Expected byte array.")
	}

	if len(der) == 0 {
		return nil, errors.New("[GMSM2PrivateKeyImportOpts] Invalid raw. It must not be nil.")
	}

	// lowLevelKey, err := utils.DERToPrivateKey(der)
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed converting PKIX to GMSM2 public key [%s]", err)
	// }

	// gmsm2SK, ok := lowLevelKey.(*sm2.PrivateKey)
	// if !ok {
	// 	return nil, errors.New("Failed casting to GMSM2 private key. Invalid raw material.")
	// }

	//gmsm2SK, err :=  sm2.ParseSM2PrivateKey(der)
	gmsm2SK, err := x509.ParsePKCS8UnecryptedPrivateKey(der)

	if err != nil {
		return nil, fmt.Errorf("Failed converting to GMSM2 private key [%s]", err)
	}

	return &gmsm2PrivateKey{gmsm2SK}, nil
}

type gmsm2PublicKeyImportOptsKeyImporter struct{}

func (*gmsm2PublicKeyImportOptsKeyImporter) KeyImport(raw interface{}, opts bccsp.KeyImportOpts) (k bccsp.Key, err error) {
	der, ok := raw.([]byte)
	if !ok {
		return nil, errors.New("[GMSM2PublicKeyImportOpts] Invalid raw material. Expected byte array.")
	}

	if len(der) == 0 {
		return nil, errors.New("[GMSM2PublicKeyImportOpts] Invalid raw. It must not be nil.")
	}

	//lowLevelKey, err := utils.DERToPrivateKey(der)
	//if err != nil {
	//	return nil, fmt.Errorf("Failed converting PKIX to GMSM2 public key [%s]", err)
	//}
	//
	//gmsm2SK, ok := lowLevelKey.(*sm2.PrivateKey)
	//if !ok {
	//	return nil, errors.New("Failed casting to GMSM2 private key. Invalid raw material.")
	//}

	gmsm2SK, err := x509.ParseSm2PublicKey(der)
	if err != nil {
		return nil, fmt.Errorf("Failed converting to GMSM2 public key [%s]", err)
	}

	return &gmsm2PublicKey{gmsm2SK}, nil
}

type ecdsaGoPublicKeyImportOptsKeyImporter struct{}

func (*ecdsaGoPublicKeyImportOptsKeyImporter) KeyImport(raw interface{}, opts bccsp.KeyImportOpts) (k bccsp.Key, err error) {
	lowLevelKey, ok := raw.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("Invalid raw material. Expected *ecdsa.PublicKey.")
	}

	return &ecdsaPublicKey{lowLevelKey}, nil
}
