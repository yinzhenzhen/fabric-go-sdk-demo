package base

import (
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/channel"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/channel/invoke"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/fab"
)

func (c *Client) ChannelQuery(
	request channel.Request,
	options ...channel.RequestOption,
) (channel.Response, error) {
	return c.channelClient.Query(request, options...)
}

func (c *Client) ChannelExecute(
	request channel.Request,
	options ...channel.RequestOption,
) (channel.Response, error) {
	return c.channelClient.Execute(request, options...)
}

func (c *Client) ChannelInvokeHandler(
	handler invoke.Handler,
	request channel.Request,
	options ...channel.RequestOption,
) (channel.Response, error) {
	return c.channelClient.InvokeHandler(handler, request, options...)
}

func (c *Client) ChannelRegisterChaincodeEvent(
	chaincodeID string,
	eventFilter string,
) (fab.Registration, <-chan *fab.CCEvent, error) {
	return c.channelClient.RegisterChaincodeEvent(chaincodeID, eventFilter)
}

func (c *Client) ChannelUnregisterChaincodeEvent(
	registration fab.Registration,
) {
	c.channelClient.UnregisterChaincodeEvent(registration)
}
