package config

import (
	"fmt"
	"github.com/spf13/viper"
	"strings"
)

type NodeSvrConfig struct {
	PeerOrg       PeerOrg    `yaml:"PeerOrg"`
	OrdererOrg    OrdererOrg `yaml:"OrdererOrg"`
	CA            CAInfo     `yaml:"CA"`
	CertPath      string     `yaml:"CertPath"`
	SDKConfigPath string     `yaml:"SDKConfigPath"`
}

type CAInfo struct {
	Address      string `yaml:"Address"`
	OrgName      string `yaml:"OrgName"`
	EnrollId     string `yaml:"EnrollId"`
	EnrollSecret string `yaml:"EnrollSecret"`
}

type PeerOrg struct {
	MspID        string     `yaml:"MspID"`
	OrgName      string     `yaml:"OrgName"`
	KeyStorePath string     `yaml:"KeyStorePath"`
	SysUserName  string     `yaml:"SysUserName"`
	Event        bool       `yaml:"Event"`
	Peers        []PeerItem `yaml:"Peers"`
}

type OrdererOrg struct {
	MspID    string        `yaml:"MspID"`
	Orderers []OrdererItem `yaml:"Orderers"`
}

type PeerItem struct {
	Name           string `yaml:"Name"`
	IP             string `yaml:"IP"`
	NodePort       string `yaml:"NodePort"`
	EventPort      string `yaml:"EventPort"`
	MspCACerts     string `yaml:"MspCACerts"`
	TlsCACerts     string `yaml:"TlsCACerts"`
	FabricDataPath string `yaml:"FabricDataPath"`
}

type OrdererItem struct {
	Name       string `yaml:"Name"`
	IP         string `yaml:"IP"`
	NodePort   string `yaml:"NodePort"`
	TlsCACerts string `yaml:"TlsCACerts"`
}

func LoadNodeSvrConfig() (*NodeSvrConfig, error) {
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
	var conf NodeSvrConfig
	if err := EnhancedExactUnmarshal(v, &conf); err != nil {
		return nil, fmt.Errorf("反序列化配置信息失败: %s", err)
	}
	return &conf, nil
}
