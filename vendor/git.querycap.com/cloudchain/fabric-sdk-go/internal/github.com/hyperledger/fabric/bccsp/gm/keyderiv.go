package gm

import (
	"errors"
	"git.querycap.com/cloudchain/fabric-sdk-go/internal/github.com/hyperledger/fabric/bccsp"
)

//定义国密 Key的驱动 ，实现 KeyDeriver 接口
type smPublicKeyKeyDeriver struct{}

func (kd *smPublicKeyKeyDeriver) KeyDeriv(k bccsp.Key, opts bccsp.KeyDerivOpts) (dk bccsp.Key, err error) {

	return nil, errors.New("Not implemented")

}
