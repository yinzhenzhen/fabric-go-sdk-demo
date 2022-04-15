/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package defcore

import (
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/logging"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/core"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/fab"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/core/logging/api"

	cryptosuiteimpl "git.querycap.com/cloudchain/fabric-sdk-go/pkg/core/cryptosuite/bccsp/gm"
	signingMgr "git.querycap.com/cloudchain/fabric-sdk-go/pkg/fab/signingmgr"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/fabsdk/provider/fabpvdr"

	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/core/logging/modlog"
)

var logger = logging.NewLogger("fabsdk")

// ProviderFactory represents the default SDK provider factory.
type ProviderFactory struct {
}

// NewProviderFactory returns the default SDK provider factory.
func NewProviderFactory() *ProviderFactory {
	f := ProviderFactory{}
	return &f
}

// CreateCryptoSuiteProvider returns a new default implementation of BCCSP
func (f *ProviderFactory) CreateCryptoSuiteProvider(config core.CryptoSuiteConfig) (core.CryptoSuite, error) {
	if config.SecurityProvider() != "gm" {
		logger.Warnf("default provider factory doesn't support '%s' crypto provider", config.SecurityProvider())
	}
	cryptoSuiteProvider, err := cryptosuiteimpl.GetSuiteByConfig(config)
	return cryptoSuiteProvider, err
}

// CreateSigningManager returns a new default implementation of signing manager
func (f *ProviderFactory) CreateSigningManager(cryptoProvider core.CryptoSuite) (core.SigningManager, error) {
	return signingMgr.New(cryptoProvider)
}

// CreateInfraProvider returns a new default implementation of fabric primitives
func (f *ProviderFactory) CreateInfraProvider(config fab.EndpointConfig) (fab.InfraProvider, error) {
	return fabpvdr.New(config), nil
}

// NewLoggerProvider returns a new default implementation of a logger backend
// This function is separated from the factory to allow logger creation first.
func NewLoggerProvider() api.LoggerProvider {
	return modlog.LoggerProvider()
}
