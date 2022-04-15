package factory

import (
	"errors"
	"fmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/internal/github.com/hyperledger/fabric/bccsp"
	"git.querycap.com/cloudchain/fabric-sdk-go/internal/github.com/hyperledger/fabric/bccsp/gm"
)

const (
	// SoftwareBasedFactoryName is the name of the factory of the software-based BCCSP implementation
	GMSoftwareBasedFactoryName = "GM"
)

type GMFactory struct{}

func (f *GMFactory) Name() string {
	return GMSoftwareBasedFactoryName
}

func (f *GMFactory) Get(config *SwOpts) (bccsp.BCCSP, error) {

	// Validate arguments
	if config == nil {
		return nil, errors.New("Invalid channel-artifacts. It must not be nil.")
	}

	gmOpts := config

	var ks bccsp.KeyStore
	if gmOpts.DummyKeystore != nil {
		ks = gm.NewDummyKeyStore()
	} else if gmOpts.FileKeystore != nil {
		fks, err := gm.NewFileBasedKeyStore(nil, gmOpts.FileKeystore.KeyStorePath, false)
		if err != nil {
			return nil, fmt.Errorf("Failed to initialize gm software key store: %s", err)
		}
		ks = fks
	} else {
		// Default to DummyKeystore
		ks = gm.NewDummyKeyStore()
	}

	return gm.New(gmOpts.SecLevel, "GMSM3", ks)

}

// SwOpts contains options for the SWFactory
type SwOpts struct {
	// Default algorithms when not specified (Deprecated?)
	SecLevel   int    `mapstructure:"security" json:"security" yaml:"Security"`
	HashFamily string `mapstructure:"hash" json:"hash" yaml:"Hash"`

	// Keystore Options
	Ephemeral     bool               `mapstructure:"tempkeys,omitempty" json:"tempkeys,omitempty"`
	FileKeystore  *FileKeystoreOpts  `mapstructure:"filekeystore,omitempty" json:"filekeystore,omitempty" yaml:"FileKeyStore"`
	DummyKeystore *DummyKeystoreOpts `mapstructure:"dummykeystore,omitempty" json:"dummykeystore,omitempty"`
	InmemKeystore *InmemKeystoreOpts `mapstructure:"inmemkeystore,omitempty" json:"inmemkeystore,omitempty"`
}

// Pluggable Keystores, could add JKS, P12, etc..
type FileKeystoreOpts struct {
	KeyStorePath string `mapstructure:"keystore" yaml:"KeyStore"`
}

type DummyKeystoreOpts struct{}

// InmemKeystoreOpts - empty, as there is no config for the in-memory keystore
type InmemKeystoreOpts struct{}
