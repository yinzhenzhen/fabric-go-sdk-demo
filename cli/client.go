package cli

import (
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/channel"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/event"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/ledger"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/resmgmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/fabsdk"
	"github.com/yinzhenzhen/fabric-go-sdk-demo/config"
	"log"
	"os"
)

type Client struct {
	SDK *fabsdk.FabricSDK
	rc  *resmgmt.Client
	cc  *channel.Client
	lc  *ledger.Client
	ec  *event.Client

	CCGoPath string

	ChannelConfigPath string
}

func NewClient(svrConfig *config.NodeSvrConfig, channelId, peer, channelConfigPath string) *Client {
	c := Client{
		CCGoPath:          os.Getenv("GOPATH"),
		ChannelConfigPath: channelConfigPath,
	}

	cp := config.BuildSdkConfig(svrConfig, channelId, peer)

	// create sdk
	sdk, err := fabsdk.New(cp)
	if err != nil {
		log.Panicf("failed to create fabric sdk: %s", err)
	}
	c.SDK = sdk
	log.Println("Initialized fabric sdk")

	// 构建资源管理客户端
	rcp := sdk.Context(fabsdk.WithUser(svrConfig.PeerOrg.SysUserName),
		fabsdk.WithOrg(svrConfig.PeerOrg.OrgName))
	rc, err := resmgmt.New(rcp)
	if err != nil {
		log.Panicf("failed to create resource client: %s", err)
	}
	c.rc = rc
	log.Println("Initialized resource client")

	// 构建通道客户端
	ccp := sdk.ChannelContext(channelId,
		fabsdk.WithUser(svrConfig.PeerOrg.SysUserName),
		fabsdk.WithOrg(svrConfig.PeerOrg.OrgName))
	cc, err := channel.New(ccp)
	if err != nil {
		log.Printf("failed to create channel client: %s", err)
		return &c
	}
	c.cc = cc
	log.Println("Initialized channel client")

	// 账本客户端
	lc, err := ledger.New(ccp)
	if err != nil {
		log.Printf("failed to create ledger client: %s", err)
		return &c
	}
	c.lc = lc
	log.Println("Initialized ledger client")

	ec, err := event.New(ccp)
	if err != nil {
		log.Printf("failed to create ledger client: %s", err)
		return &c
	}
	c.ec = ec
	log.Println("Initialized event client")

	return &c
}
