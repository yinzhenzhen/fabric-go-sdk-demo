package cli

import (
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/channel"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/resmgmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/fab"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/fab/ccpackager/gopackager"
	"git.querycap.com/cloudchain/fabric-sdk-go/third_party/github.com/hyperledger/fabric/common/policydsl"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"strings"
)

// 安装链码
func (c *Client) InstallChaincode(ccName, ccPath, ccVersion string, peers []string) error {
	// pack the chaincode
	ccPkg, err := gopackager.NewCCPackage(ccPath, c.CCGoPath)
	if err != nil {
		return errors.WithMessage(err, "pack chaincode error")
	}

	// new request of installing chaincode
	req := resmgmt.InstallCCRequest{
		Name:    ccName,
		Path:    ccPath,
		Version: ccVersion,
		Package: ccPkg,
	}

	targetPeer := resmgmt.WithTargetEndpoints(peers...)
	resps, err := c.rc.InstallCC(req, targetPeer)
	if err != nil {
		return errors.WithMessage(err, "installCC error")
	}

	// check other errors
	var errs []error
	for _, resp := range resps {
		log.Printf("Install response status: %v", resp.Status)
		if resp.Status != http.StatusOK {
			errs = append(errs, errors.New(resp.Info))
		}
		if resp.Info == "already installed" {
			log.Printf("Chaincode %s already installed on peer: %s.\n",
				ccName+"-"+ccVersion, resp.Target)
			return nil
		}
	}

	if len(errs) > 0 {
		log.Printf("InstallCC errors: %v", errs)
		return errors.WithMessage(errs[0], "installCC first error")
	}
	return nil
}

// 实例化链码
func (c *Client) InstantiateChaincode(ccName, ccPath, ccVersion string, peer, channelID string, args []string) error {
	// endorser policy
	org1OrOrg2 := "OR('Org1MSP.member','Org2MSP.member')"
	ccPolicy, err := c.genPolicy(org1OrOrg2)
	if err != nil {
		return errors.WithMessage(err, "gen policy from string error")
	}

	// new request
	argsBytes := packArgs(args)
	req := resmgmt.InstantiateCCRequest{
		Name:    ccName,
		Path:    ccPath,
		Version: ccVersion,
		Args:    argsBytes,
		Policy:  ccPolicy,
	}

	// send request and handle response
	reqPeers := resmgmt.WithTargetEndpoints(peer)
	resp, err := c.rc.InstantiateCC(channelID, req, reqPeers)
	if err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return nil
		}
		return errors.WithMessage(err, "instantiate chaincode error")
	}

	log.Printf("Instantitate chaincode tx: %s", resp.TransactionID)
	return nil
}

// 升级链码
func (c *Client) UpgradeChaincode(ccName, ccPath, ccVersion string, peer, channelID string) (fab.TransactionID, error) {
	// endorser policy
	org1AndOrg2 := "OR('Org1MSP.member','Org2MSP.member')"
	ccPolicy, err := c.genPolicy(org1AndOrg2)
	if err != nil {
		return "", errors.WithMessage(err, "gen policy from string error")
	}

	// new request
	args := packArgs([]string{"init"})
	req := resmgmt.UpgradeCCRequest{
		Name:    ccName,
		Path:    ccPath,
		Version: ccVersion,
		Args:    args,
		Policy:  ccPolicy,
	}

	// send request and handle response
	reqPeers := resmgmt.WithTargetEndpoints(peer)
	resp, err := c.rc.UpgradeCC(channelID, req, reqPeers)
	if err != nil {
		return "", errors.WithMessage(err, "upgrade chaincode error")
	}

	log.Printf("Upgrade chaincode tx: %s", resp.TransactionID)
	return resp.TransactionID, nil
}

// 调用链码
func (c *Client) InvokeChaincode(ccName, funcName string, args []string, peers []string) error {
	argsByte := packArgs(args)
	req := channel.Request{
		ChaincodeID: ccName,
		Fcn:         funcName,
		Args:        argsByte,
	}

	// send request and handle response
	// peers is needed
	reqPeers := channel.WithTargetEndpoints(peers...)
	resp, err := c.cc.Execute(req, reqPeers)
	if err != nil {
		return errors.WithMessage(err, "invoke chaincode error")
	}
	log.Printf("Invoke chaincode response:\n"+
		"id: %v\nvalidate: %v\nchaincode status: %v\n\n",
		resp.TransactionID,
		resp.TxValidationCode,
		resp.ChaincodeStatus)

	return nil
}

// 查询链码
func (c *Client) QueryChaincode(ccName, funcName string, args []string, peer string) error {
	// new channel request for query
	argsByte := packArgs(args)
	req := channel.Request{
		ChaincodeID: ccName,
		Fcn:         funcName,
		Args:        argsByte,
	}

	// send request and handle response
	reqPeers := channel.WithTargetEndpoints(peer)
	resp, err := c.cc.Query(req, reqPeers)
	if err != nil {
		return errors.WithMessage(err, "query chaincode error")
	}

	log.Printf("Query chaincode tx response:\ntx: %s\nresult: %v\n\n",
		resp.TransactionID,
		string(resp.Payload))
	return nil
}

func (c *Client) genPolicy(policy string) (*common.SignaturePolicyEnvelope, error) {
	//if policy == "ANY" {
	//	return policydsl.SignedByAnyMember([]string{c.OrgName}), nil
	//}
	return policydsl.FromString(policy)
}

func (c *Client) Close() {
	c.SDK.Close()
}

func packArgs(args []string) [][]byte {
	r := [][]byte{}
	for _, s := range args {
		temp := []byte(s)
		r = append(r, temp)
	}
	return r
}
