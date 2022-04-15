package main

import (
	"fmt"
	"github.com/yinzhenzhen/fabric-go-sdk-demo/cli"
	"github.com/yinzhenzhen/fabric-go-sdk-demo/config"
)

var (
	//org1Peers = []string{"peer0.org1.example.com"}
	//org2Peers = []string{"peer0.org2.example.com"}
	cloudchainOrg1Peers = []string{"peer0.org1.cloudchain-dev-gm.rktl.xyz", "peer1.org1.cloudchain-dev-gm.rktl.xyz", "peer2.org1.cloudchain-dev-gm.rktl.xyz"}
	//cloudchainOrg2Peers = []string{"peer0.org2.cloudchain-dev-gm.rktl.xyz","peer1.org2.cloudchain-dev-gm.rktl.xyz","peer2.org2.cloudchain-dev-gm.rktl.xyz"}
)

func main() {

	channelId := "cloudchain"
	channelConfigPath := "/root/go/src/git.querycap.com/cloudchain/fabric-depoly-resource/channel-artifacts/cloudchain.tx"

	node1Config, err := config.LoadNodeSvrConfig()
	//node2Config, err := config.LoadNodeSvrConfig()
	if err != nil {
		panic(err)
	}

	org1Client := cli.NewClient(node1Config, channelId, cloudchainOrg1Peers[0], channelConfigPath)
	//org2Client := cli.NewClient(node2Config, channelId, org2Peers[0], channelConfigPath)

	defer org1Client.Close()
	//defer org2Client.Close()

	//err = org1Client.CreateChannel(channelId, channelConfigPath)

	//err = org1Client.JoinChannel(channelId, org1Peers)
	//err = org2Client.JoinChannel(channelId, org2Peers)

	//ccName := "mycc"
	//ccPath := "github.com/hyperledger/fabric/examples/chaincode/go/mycc"
	//ccVersion := "1.0"
	//initArgs := []string{"init","a","100","b","200"}

	var ccName = "universal-deposit"
	var ccPath = "git.querycap.com/cloudchain/cc-market/universal-deposit"
	var ccVersion = "1.0.0"
	initArgs := []string{"init"}

	err = org1Client.InstallChaincode(ccName, ccPath, ccVersion, cloudchainOrg1Peers)
	if err != nil {
		fmt.Println("InstallChaincode panic")
		panic(err)
	}
	//err = org2Client.InstallChaincode(ccName, ccPath, ccVersion, org2Peers)

	err = org1Client.InstantiateChaincode(ccName, ccPath, ccVersion, cloudchainOrg1Peers[0], channelId, initArgs)
	if err != nil {
		fmt.Println("InstantiateChaincode panic")
		panic(err)
	}
	//InstallAndUpgradeCC("1.2", org1Client, org2Client)

	// 调用链码
	//funcName := "invoke"
	//args := []string{"b", "a", "20"}
	//err = org1Client.InvokeChaincode(ccName, funcName, args, org1Peers)
	//err = org2Client.InvokeChaincode(ccName, funcName, args, org2Peers)

	// 查询链码
	funcName := "query"
	args := []string{"b"}
	err = org1Client.QueryChaincode(ccName, funcName, args, cloudchainOrg1Peers[0])
	//err = org2Client.QueryChaincode(ccName, funcName, args, org2Peers[0])

	if err != nil {
		panic(err)
	}
}
