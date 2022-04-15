package base

import (
	"fmt"

	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/channel"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/event"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/ledger"
	mspmgmt "git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/msp"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/resmgmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/core"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/core/config"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/fabsdk"
	"github.com/sirupsen/logrus"
)

type Client struct {
	ConfigPath    string            `env:""`
	ConfigData    []byte            `env:"-"`
	Organization  string            `env:""`
	Username      string            `env:""`
	ChannelID     string            `env:""`
	fabricSDK     *fabsdk.FabricSDK `env:"-"`
	resmgmtClient *resmgmt.Client   `env:"-"`
	mspmgmtClient *mspmgmt.Client   `env:"-"`
	ledgerClient  *ledger.Client    `env:"-"`
	channelClient *channel.Client   `env:"-"`
	eventClient   *event.Client     `env:"-"`
}

func NewClient(opts ...Option) (*Client, error) {
	c := &Client{}
	for _, opt := range opts {
		opt(c)
	}
	c.SetDefaults()
	if err := c.setup(); err != nil {
		return nil, fmt.Errorf("failed to setup: %v", err)
	}
	return c, nil
}

func (c *Client) SetDefaults() {
}

func (c *Client) Init() {
	if len(c.ConfigPath) == 0 {
		return
	}
	if err := c.setup(); err != nil {
		logrus.Fatalf("failed to setup: %v", err)
	}
}

func (c *Client) setup() error {
	var configProvider core.ConfigProvider
	if len(c.ConfigData) != 0 {
		configProvider = config.FromRaw(c.ConfigData, "yaml")
	} else {
		configProvider = config.FromFile(c.ConfigPath)
	}
	fabricSDK, err := fabsdk.New(configProvider)
	if err != nil {
		return fmt.Errorf("failed to new fabricSDK: %v", err)
	}
	c.fabricSDK = fabricSDK
	clientProvider := c.fabricSDK.Context(
		fabsdk.WithOrg(c.Organization),
		fabsdk.WithUser(c.Username),
	)
	resmgmtClient, err := resmgmt.New(clientProvider)
	if err != nil {
		return fmt.Errorf("failed to new resmgmtClient: %v", err)
	}
	c.resmgmtClient = resmgmtClient
	mspmgmtClient, err := mspmgmt.New(clientProvider)
	if err != nil {
		return fmt.Errorf("failed to new mspmgmtClient: %v", err)
	}
	c.mspmgmtClient = mspmgmtClient

	if len(c.ChannelID) != 0 {
		channelProvider := c.fabricSDK.ChannelContext(
			c.ChannelID,
			fabsdk.WithOrg(c.Organization),
			fabsdk.WithUser(c.Username),
		)
		ledgerClient, err := ledger.New(channelProvider)
		if err != nil {
			return fmt.Errorf("failed to new ledgerClient: %v", err)
		}
		c.ledgerClient = ledgerClient
		channelClient, err := channel.New(channelProvider)
		if err != nil {
			return fmt.Errorf("failed to new channelClient: %v", err)
		}
		c.channelClient = channelClient
		eventClient, err := event.New(channelProvider)
		if err != nil {
			return fmt.Errorf("failed to new eventClient: %v", err)
		}
		c.eventClient = eventClient
	}

	return nil
}

func (c *Client) Close() {
	if c.fabricSDK != nil {
		c.fabricSDK.Close()
	}
}
