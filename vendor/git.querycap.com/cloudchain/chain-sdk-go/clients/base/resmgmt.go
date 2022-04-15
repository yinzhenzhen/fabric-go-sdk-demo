package base

import (
	"io"

	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/peer"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/resmgmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/fab"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/msp"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/fab/resource"
)

func (c *Client) ResmgmtQueryChannels(
	options ...resmgmt.RequestOption,
) (*peer.ChannelQueryResponse, error) {
	return c.resmgmtClient.QueryChannels(options...)
}

func (c *Client) ResmgmtSaveChannel(
	request resmgmt.SaveChannelRequest,
	options ...resmgmt.RequestOption,
) (resmgmt.SaveChannelResponse, error) {
	return c.resmgmtClient.SaveChannel(request, options...)
}

func (c *Client) ResmgmtJoinChannel(
	channelID string,
	options ...resmgmt.RequestOption,
) error {
	return c.resmgmtClient.JoinChannel(channelID, options...)
}

func (c *Client) ResmgmtQueryInstalledChaincodes(
	options ...resmgmt.RequestOption,
) (*peer.ChaincodeQueryResponse, error) {
	return c.resmgmtClient.QueryInstalledChaincodes(options...)
}

func (c *Client) ResmgmtInstallCC(
	request resmgmt.InstallCCRequest,
	options ...resmgmt.RequestOption,
) ([]resmgmt.InstallCCResponse, error) {
	return c.resmgmtClient.InstallCC(request, options...)
}

func (c *Client) ResmgmtQueryInstantiatedChaincodes(
	channelID string,
	options ...resmgmt.RequestOption,
) (*peer.ChaincodeQueryResponse, error) {
	return c.resmgmtClient.QueryInstantiatedChaincodes(channelID, options...)
}

func (c *Client) ResmgmtInstantiateCC(
	channelID string,
	request resmgmt.InstantiateCCRequest,
	options ...resmgmt.RequestOption,
) (resmgmt.InstantiateCCResponse, error) {
	return c.resmgmtClient.InstantiateCC(channelID, request, options...)
}

func (c *Client) ResmgmtUpgradeCC(
	channelID string,
	request resmgmt.UpgradeCCRequest,
	options ...resmgmt.RequestOption,
) (resmgmt.UpgradeCCResponse, error) {
	return c.resmgmtClient.UpgradeCC(channelID, request, options...)
}

func (c *Client) ResmgmtQueryConfigBlockFromOrderer(
	channelID string,
	options ...resmgmt.RequestOption,
) (*common.Block, error) {
	return c.resmgmtClient.QueryConfigBlockFromOrderer(channelID, options...)
}

func (c *Client) ResmgmtQueryConfigFromOrderer(
	channelID string,
	options ...resmgmt.RequestOption,
) (fab.ChannelCfg, error) {
	return c.resmgmtClient.QueryConfigFromOrderer(channelID, options...)
}

func (c *Client) ResmgmtCreateConfigSignatureFromReader(
	signer msp.SigningIdentity,
	channelConfig io.Reader,
) (*common.ConfigSignature, error) {
	return c.resmgmtClient.CreateConfigSignatureFromReader(signer, channelConfig)
}

func (c *Client) ResmgmtCreateConfigSignatureDataFromReader(
	signer msp.SigningIdentity,
	channelConfig io.Reader,
) (resource.ConfigSignatureData, error) {
	return c.resmgmtClient.CreateConfigSignatureDataFromReader(signer, channelConfig)
}

func (c *Client) ResmgmtQueryCollectionsConfig(
	channelID string,
	chaincodeName string,
	options ...resmgmt.RequestOption,
) (*peer.CollectionConfigPackage, error) {
	return c.resmgmtClient.QueryCollectionsConfig(channelID, chaincodeName, options...)
}
