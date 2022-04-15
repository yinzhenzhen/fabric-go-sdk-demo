/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabsdk

import (
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/core/logging/api"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/core/logging/modlog"
	sdkApi "git.querycap.com/cloudchain/fabric-sdk-go/pkg/fabsdk/api"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/fabsdk/factory/defcore"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/fabsdk/factory/defmsp"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/fabsdk/factory/defsvc"
)

type defPkgSuite struct{}

func (ps *defPkgSuite) Core() (sdkApi.CoreProviderFactory, error) {
	return defcore.NewProviderFactory(), nil
}

func (ps *defPkgSuite) MSP() (sdkApi.MSPProviderFactory, error) {
	return defmsp.NewProviderFactory(), nil
}

func (ps *defPkgSuite) Service() (sdkApi.ServiceProviderFactory, error) {
	return defsvc.NewProviderFactory(), nil
}

func (ps *defPkgSuite) Logger() (api.LoggerProvider, error) {
	return modlog.LoggerProvider(), nil
}
