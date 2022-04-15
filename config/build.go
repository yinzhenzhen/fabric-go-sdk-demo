package config

import (
	"fmt"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/core"
	"github.com/yinzhenzhen/fabric-go-sdk-demo/config/backend"
)

// 配置信息列表
type Config struct {
	SDKConfigPath string
	Peers         []BaseNodeConfig
	Orders        []BaseNodeConfig
	OrgName       string
	OrgMspId      string
	OrgCryptoPath string
	UserName      string
	ChannelName   string
	PeerName      string
}

//基础节点信息列表
type BaseNodeConfig struct {
	NodeName    string
	NodeUrl     string
	EventUrl    string
	NodeTlsPath string
}

func BuildSdkConfig(nodeSvrConf *NodeSvrConfig, channelName, peerName string) core.ConfigProvider {
	if nodeSvrConf == nil {
		return nil
	}
	conf := transferConfig(nodeSvrConf)

	// 添加Order列表
	var orders []backend.OrdererConfig
	for _, item := range conf.Orders {
		orders = append(orders, backend.OrdererConfig{
			OrdererName:    item.NodeName,
			OrdererUrl:     item.NodeUrl,
			TlsCACertsPath: item.NodeTlsPath})
	}

	// 节点名称列表
	var peerNames []string
	// 添加Peer列表
	var peers []backend.PeerConfig
	for _, item := range conf.Peers {
		peer := backend.PeerConfig{
			PeerName:       item.NodeName,
			PeerUrl:        item.NodeUrl,
			PeerEventUrl:   item.EventUrl,
			TlsCACertsPath: item.NodeTlsPath}
		peers = append(peers, peer)
		peerNames = append(peerNames, item.NodeName)
	}

	// 添加组织
	org := backend.OrganizationConfig{
		OrgName:    conf.OrgName,
		MspId:      conf.OrgMspId,
		CryptoPath: conf.OrgCryptoPath,
		Peers:      peerNames}

	// 通道
	channel := backend.ChannelConfig{channelName, peerName}

	// 匹配
	var entity backend.EntityMatchersConfig
	for _, item := range conf.Peers {
		peer := backend.EntityPeer{PeerName: item.NodeName, PeerUrl: item.NodeUrl, PeerEventUrl: item.EventUrl}
		entity.EntityPeers = append(entity.EntityPeers, peer)
	}
	for _, item := range conf.Orders {
		order := backend.EntityOrder{OrderName: item.NodeName, OrderUrl: item.NodeUrl}
		entity.EntityOrders = append(entity.EntityOrders, order)
	}

	// 添加配置
	var s []SetOption
	s = append(s, SetOrderer(&orders))
	s = append(s, SetPeer(&peers))
	s = append(s, SetOrg(&org))
	s = append(s, SetChannel(&channel))
	s = append(s, SetEntityMatchers(&entity))
	return FromFile(conf.SDKConfigPath, s)
}

func transferConfig(nodeSvrConf *NodeSvrConfig) *Config {
	// 检验记账节点信息并进行添加
	var peerlist []BaseNodeConfig
	for _, item := range nodeSvrConf.PeerOrg.Peers {
		peerlist = append(peerlist, BaseNodeConfig{
			NodeName:    item.Name,
			NodeUrl:     fmt.Sprintf("grpcs://%s:%s", item.IP, item.NodePort),
			EventUrl:    fmt.Sprintf("grpcs://%s:%s", item.IP, item.EventPort),
			NodeTlsPath: nodeSvrConf.CertPath + item.TlsCACerts,
		})
	}

	// 检验排序节点信息并进行添加
	var ordererlist []BaseNodeConfig
	for _, item := range nodeSvrConf.OrdererOrg.Orderers {
		ordererlist = append(ordererlist, BaseNodeConfig{
			NodeName:    item.Name,
			NodeUrl:     fmt.Sprintf("grpcs://%s:%s", item.IP, item.NodePort),
			NodeTlsPath: nodeSvrConf.CertPath + item.TlsCACerts,
		})
	}

	conf := Config{
		SDKConfigPath: nodeSvrConf.SDKConfigPath,
		Peers:         peerlist,
		Orders:        ordererlist,
		OrgName:       nodeSvrConf.PeerOrg.OrgName,
		OrgMspId:      nodeSvrConf.PeerOrg.MspID,
		OrgCryptoPath: nodeSvrConf.PeerOrg.KeyStorePath,
		UserName:      nodeSvrConf.PeerOrg.SysUserName,
	}
	return &conf
}
