package base

import (
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/fab"
)

func (c *Client) EventRegisterChaincodeEvent(
	chaincodeID string,
	eventFilter string,
) (fab.Registration, <-chan *fab.CCEvent, error) {
	return c.eventClient.RegisterChaincodeEvent(chaincodeID, eventFilter)
}

func (c *Client) EventRegisterBlockEvent(
	blockFilters ...fab.BlockFilter,
) (fab.Registration, <-chan *fab.BlockEvent, error) {
	return c.eventClient.RegisterBlockEvent(blockFilters...)
}

func (c *Client) EventRegisterFilteredBlockEvent() (fab.Registration, <-chan *fab.FilteredBlockEvent, error) {
	return c.eventClient.RegisterFilteredBlockEvent()
}

func (c *Client) EventRegisterTxStatusEvent(
	txID string,
) (fab.Registration, <-chan *fab.TxStatusEvent, error) {
	return c.eventClient.RegisterTxStatusEvent(txID)
}

func (c *Client) EventUnregister(
	registration fab.Registration,
) {
	c.eventClient.Unregister(registration)
}
