module github.com/hyperledger/fabric-go-sdk-demo

go 1.16

require (
	git.querycap.com/cloudchain/chain-sdk-go v1.3.0-gm.0.20210706022209-9b72f8e0959d
	git.querycap.com/cloudchain/fabric-sdk-go v1.0.0-tjgm-alpha
	github.com/cloudflare/cfssl v1.6.0
	github.com/golang/protobuf v1.5.2
	github.com/hyperledger/fabric-protos-go v0.0.0-20210528200356-82833ecdac31
	github.com/mitchellh/mapstructure v1.4.1
	github.com/onsi/gomega v1.13.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/spf13/viper v1.8.1
	github.com/tjfoc/gmsm v1.4.0
)

replace (
	github.com/go-kit/kit => github.com/go-kit/kit v0.8.0
	github.com/mitchellh/mapstructure => github.com/mitchellh/mapstructure v1.2.2
)
