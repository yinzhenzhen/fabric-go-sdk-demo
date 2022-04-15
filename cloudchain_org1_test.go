package main

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/yinzhenzhen/fabric-go-sdk-demo/cli"
	"github.com/yinzhenzhen/fabric-go-sdk-demo/config"
	"strings"
	"testing"
)

var channelId = "cloudchain"
var cloudchainOrg1Client *cli.Client

var channelConfigPath = "/root/go/src/git.querycap.com/cloudchain/fabric-depoly-resource/channel-artifacts/cloudchain.tx"
var ccName = "universal-deposit"
var ccPath = "git.querycap.com/cloudchain/cc-market/universal-deposit"

//var ccName = "carbonchain-business"
//var ccPath = "git.querycap.com/cloudchain/cc-market/carbonchain-business"
var ccVersion = "1.0.0"
var initArgs = []string{"init"}

func init() {
	cloudchainOrg1Config, err := LoadNodeSvrConfig_Clouchain_Org1()
	if err != nil {
		panic(err)
	}
	cloudchainOrg1Client = cli.NewClient(cloudchainOrg1Config, channelId, cloudchainOrg1Peers[0], channelConfigPath)
}

func LoadNodeSvrConfig_Clouchain_Org1() (*config.NodeSvrConfig, error) {
	v := viper.New()
	v.SetConfigName("cloudchain_org1_svrconfig")
	v.SetConfigType("yaml")
	v.AddConfigPath("./conf")
	v.AddConfigPath(".")
	v.SetEnvPrefix("FABRICSM")
	v.AutomaticEnv()
	replacer := strings.NewReplacer(".", "_")
	v.SetEnvKeyReplacer(replacer)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("读取配置失败: %s", err)
	}
	var conf config.NodeSvrConfig
	if err := config.EnhancedExactUnmarshal(v, &conf); err != nil {
		return nil, fmt.Errorf("反序列化配置信息失败: %s", err)
	}
	return &conf, nil
}

func Test_CreateChannel_Clouchain_Org1(t *testing.T) {
	err := cloudchainOrg1Client.CreateChannel(channelId, channelConfigPath)
	if err != nil {
		panic(err)
	}
}

func Test_JoinChannel_Clouchain_Org1(t *testing.T) {
	err := cloudchainOrg1Client.JoinChannel(channelId, cloudchainOrg1Peers)
	if err != nil {
		panic(err)
	}
}

func Test_InstallCC_Clouchain_Org1(t *testing.T) {
	err := cloudchainOrg1Client.InstallChaincode(ccName, ccPath, ccVersion, cloudchainOrg1Peers)
	if err != nil {
		panic(err)
	}
}

func Test_InstantiateCC_Clouchain_Org1(t *testing.T) {
	err := cloudchainOrg1Client.InstantiateChaincode(ccName, ccPath, ccVersion, cloudchainOrg1Peers[0], channelId, initArgs)
	if err != nil {
		panic(err)
	}
}

func Test_InvokeCC_Clouchain_Org1(t *testing.T) {
	funcName := "CreateApplication"
	args := []string{"111", "应用1", "张三1"}
	err := cloudchainOrg1Client.InvokeChaincode(ccName, funcName, args, cloudchainOrg1Peers)
	if err != nil {
		panic(err)
	}
}

func Test_InvokeCCForEvent_Clouchain_Org1(t *testing.T) {
	funcName := "invoke"
	args := []string{"b", "a", "20"}
	err := cloudchainOrg1Client.InvokeChaincodeForEvent(ccName, funcName, args, cloudchainOrg1Peers)
	if err != nil {
		panic(err)
	}
}

func Test_QueryCC_Clouchain_Org1(t *testing.T) {
	funcName := "GetApplicationInfo"
	args := []string{"111"}
	err := cloudchainOrg1Client.QueryChaincode(ccName, funcName, args, cloudchainOrg1Peers[0])
	if err != nil {
		panic(err)
	}
}

func Test_QueryChannelConfig_Clouchain_Org1(t *testing.T) {
	cloudchainOrg1Client.QueryChannelConfig()
	cloudchainOrg1Client.QueryChannelInfo()
}

func Test_QueryBlock_Clouchain_Org1(t *testing.T) {
	for i := 1; i <= 5; i++ {
		info := cloudchainOrg1Client.QueryBlock(uint64(i))
		if info != nil {
			fmt.Println(i)
			//fmt.Println("blockinfo:", info)
			fmt.Println("block PrevBlockHash: " + info.PrevBlockHash)
			fmt.Println("block BlockHash: " + info.BlockHash)
		}
	}
}

func Test_QueryBlock_Clouchain_Org1_2(t *testing.T) {
	for i := 1; i <= 5; i++ {
		info := cloudchainOrg1Client.QueryBlock2(uint64(i))
		if info != nil {
			fmt.Println(i)
			//fmt.Println("blockinfo:", info)
			fmt.Println("block PrevBlockHash: " + info.Header.PreviousHash)
			fmt.Println("block BlockHash: " + info.BlockHash)
		}
	}
}

func Test_QueryBlockByHash_Clouchain_Org1(t *testing.T) {
	info := cloudchainOrg1Client.QueryBlockByHash([]byte(""))
	if info != nil {
		fmt.Println(info)
	}
}

func Test_QueryTransaction_Clouchain_Org1(t *testing.T) {
	info := cloudchainOrg1Client.QueryTransaction("0b4e64e3bc96caf3533664990c8a39814ce447d1af9d6843760d6b371bc692c5")
	if info != nil {
		fmt.Println(info)
	}
}
