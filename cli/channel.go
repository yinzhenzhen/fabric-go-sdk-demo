package cli

import (
	"fmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/resmgmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/pkg/errors"
	"log"
	"os"
)

func (c *Client) CreateChannel(channelId, channelConfigPath string) error {

	// channel.tx
	f, err := os.Open(channelConfigPath)

	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("failed to open channel config: %s\n", err))
	}
	defer f.Close()

	req := resmgmt.SaveChannelRequest{
		ChannelID:         channelId,
		ChannelConfigPath: channelConfigPath,
	}

	//reqOrders := resmgmt.WithOrdererEndpoint(orderer)
	resp, err := c.rc.SaveChannel(req, resmgmt.WithRetry(retry.DefaultResMgmtOpts))

	if err != nil {
		return errors.WithMessage(err, "createChannel error")
	}

	if resp.TransactionID == "" {
		return errors.WithMessage(err, "Failed to save channel")
	}
	log.Printf("create channel tx response:\ntx: %s\n",
		resp.TransactionID)

	return nil
}

func (c *Client) JoinChannel(channelId string, peers []string) error {

	reqPeers := resmgmt.WithTargetEndpoints(peers...)

	err := c.rc.JoinChannel(channelId, reqPeers)

	if err != nil {
		return errors.WithMessage(err, "Failed to join channel")
	}

	log.Printf("join channel success")

	return nil
}
