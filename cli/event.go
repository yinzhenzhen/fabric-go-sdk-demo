package cli

import (
	"fmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/channel"
	"log"
	"time"
)

func (c *Client) InvokeChaincodeForEvent(ccName, funcName string, args []string, peers []string) error {

	eventID := "eventInvoke"

	// Add data that will be visible in the proposal, like a description of the invoke request
	transientDataMap := make(map[string][]byte)
	transientDataMap["result"] = []byte("Transient data in hello invoke")

	reg, notifier, err := c.ec.RegisterChaincodeEvent(ccName, eventID)
	if err != nil {
		return err
	}
	defer c.ec.Unregister(reg)

	req := channel.Request{
		ChaincodeID:  ccName,
		Fcn:          funcName,
		Args:         packArgs(args),
		TransientMap: transientDataMap,
	}

	// Create a request (proposal) and send it
	resp, err := c.cc.Execute(req)
	if err != nil {
		return fmt.Errorf("failed to move funds: %v", err)
	}

	// Wait for the result of the submission
	select {
	case ccEvent := <-notifier:
		fmt.Printf("Received CC event: %v\n", ccEvent)
	case <-time.After(time.Second * 20):
		return fmt.Errorf("did NOT receive CC event for eventId(%s)", eventID)
	}

	log.Printf("Invoke chaincode response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)

	return nil
}
