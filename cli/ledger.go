package cli

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"git.querycap.com/cloudchain/chain-sdk-go/clients/base"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/fab"
	"github.com/cloudflare/cfssl/log"
	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric-protos-go/msp"
	"github.com/hyperledger/fabric-protos-go/peer"
	"github.com/tjfoc/gmsm/x509"
	"github.com/yinzhenzhen/fabric-go-sdk-demo/protos/utils"
	"strconv"
	"time"
)

func (c *Client) QueryChannelConfig() {
	resp, err := c.lc.QueryConfig()
	if err != nil {
		fmt.Printf("Failed to queryConfig: %s", err)
		return
	}
	fmt.Println("ChannelID: ", resp.ID())
	fmt.Println("Channel Orderers: ", resp.Orderers())
	fmt.Println("Channel Versions: ", resp.Versions())
}

func (c *Client) QueryChannelInfo() {
	resp, err := c.lc.QueryInfo()
	if err != nil {
		fmt.Printf("Failed to queryInfo: %s", err)
		return
	}
	fmt.Println("BlockChainInfo:", resp.BCI)
	fmt.Println("Endorser:", resp.Endorser)
	fmt.Println("Status:", resp.Status)
}

func (c *Client) QueryBlock(number uint64) *BlockInfo {
	block, err := c.lc.QueryBlock(number)
	if err != nil {
		fmt.Printf("Failed to QueryBlock: %s", err)
		return nil
	}
	info := parseBlock(block)
	return info
}

func (c *Client) QueryBlock2(number uint64) *base.Block {
	block, err := c.lc.QueryBlock(number)
	if err != nil {
		fmt.Printf("Failed to QueryBlock: %s", err)
		return nil
	}
	info, err := base.DecodeBlock(block)
	if err != nil {
		fmt.Printf("Failed to ParseBlock: %s", err)
		return nil
	}
	return info
}

func (c *Client) QueryBlockByHash(hash []byte) *BlockInfo {
	block, err := c.lc.QueryBlockByHash(hash)
	if err != nil {
		fmt.Printf("Failed to QueryBlockByHash: %s", err)
		return nil
	}
	info := parseBlock(block)
	fmt.Println("blockinfo:", info)
	return info
}

func (c *Client) QueryTransaction(hash string) *Transaction {
	pt, err := c.lc.QueryTransaction(fab.TransactionID(hash))
	if err != nil {
		fmt.Printf("Failed to QueryBlockByHash: %s", err)
		return nil
	}

	info, err := parseTransaction(pt.TransactionEnvelope)
	if err != nil {
		log.Errorf("parse transaction err: %v", err)
		return nil
	}
	fmt.Println("blockinfo:", info)
	return info
}

func parseBlock(block *common.Block) *BlockInfo {

	// 转化为区块
	blockInfo := &BlockInfo{}
	// 区块号
	blockInfo.BlockNum = block.Header.Number
	// 上一个区块哈希
	blockInfo.PrevBlockHash = hex.EncodeToString(block.Header.PreviousHash)
	// 区块哈希
	blockInfo.BlockHash = hex.EncodeToString(utils.BlockHeaderHash(block.Header))
	// 数据哈希
	blockInfo.DataHash = hex.EncodeToString(block.Header.DataHash)
	// 区块大小
	blockInfo.BlockSize = uint64(len(block.String()))

	var trans []*Transaction
	for _, envBytes := range block.Data.Data {
		var env *common.Envelope
		var err error

		// 从区块中提取Envelope
		if env, err = utils.GetEnvelopeFromBlock(envBytes); err != nil {
			return nil
		}
		tran, err := parseTransaction(env)
		if err != nil {
			log.Debug("parseTransaction err %v", err)
		}
		trans = append(trans, tran)
	}

	// 落块时间
	blockInfo.BlockTime = trans[len(trans)-1].TxTime

	// 交易数
	blockInfo.TxCount = int32(len(trans))

	blockInfo.Trans = trans
	return blockInfo
}

func parseTransaction(env *common.Envelope) (*Transaction, error) {

	tran := &Transaction{}
	var err error
	var payload *common.Payload

	// 从Envelope提取payload
	if payload, err = utils.UnmarshalPayload(env.Payload); err != nil {
		return nil, err
	}

	var chdr *common.ChannelHeader
	// 从payload的头部提取通道头部
	if chdr, err = utils.UnmarshalChannelHeader(payload.Header.ChannelHeader); err != nil {
		return nil, err
	}

	// 交易id
	tran.TxId = chdr.TxId
	// 交易编码
	tran.ValidationCode = ""
	// 交易类型
	tran.Type = strconv.Itoa(int(chdr.Type))
	// 交易大小
	//tran.TransSize = uint64(len(string(envBytes)))
	// 交易时间
	tran.TxTime = time.Unix(chdr.Timestamp.Seconds, int64(chdr.Timestamp.Nanos)).Format(time.UnixDate)

	//var cche *peer.ChaincodeHeaderExtension
	//if cche, err = gmutils.GetChaincodeHeaderExtension(payload.Header); err != nil {
	//	return nil, err
	//}
	//// 链码名称
	//tran.ChainCodeName = cche.ChaincodeId.Name + ":" + cche.ChaincodeId.Version

	if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {

		var tx *peer.Transaction
		// 从payload的头部提取data
		if tx, err = utils.UnmarshalTransaction(payload.Data); err != nil {
			return nil, err
		}

		if len(tx.Actions) > 0 {

			action := tx.Actions[0]
			// 从第一个交易提取签名头部
			AShdr, err := utils.UnmarshalSignatureHeader(action.Header)

			if err != nil {
				return nil, err
			}

			//获取交易的提交者
			var subject string //mspid
			if _, subject, err = decodeSerializedIdentity(AShdr.Creator); err != nil {
				return nil, err
			}
			// 创建者mspid
			tran.CreaterMSPId = subject

			//var capayload *peer.ChaincodeActionPayload
			var ca *peer.ChaincodeAction
			var cap *peer.ChaincodeActionPayload
			// 从第一个交易提取payload
			if cap, ca, err = utils.GetPayloads(action); err != nil {
				return nil, err
			}

			// 交易输入信息
			//tran.Input = string(cap.ChaincodeProposalPayload)
			tran.Input = base64.StdEncoding.EncodeToString(cap.ChaincodeProposalPayload)
			// 链码名称
			tran.ChainCodeName = ca.ChaincodeId.Name
			// 链码id
			tran.ChainCodeId = ca.ChaincodeId.Name + ":" + ca.ChaincodeId.Version
			// 状态
			tran.Status = ca.Response.Status
			// 交易输出信息
			//tran.Output = string(ca.Response.Payload)
			tran.Output = base64.StdEncoding.EncodeToString(ca.Response.Payload)
		}

	}

	return tran, nil

}

func decodeSerializedIdentity(creator []byte) (string, string, error) {

	si := &msp.SerializedIdentity{}
	err := proto.Unmarshal(creator, si)
	if err != nil {
		return "", "", err
	}

	mspId := si.Mspid
	dcert, _ := pem.Decode(si.IdBytes)
	x509Cert, err := x509.ParseCertificate(dcert.Bytes)
	if err != nil {
		return "", "", err
	}

	subject := x509Cert.Subject.CommonName
	return mspId, subject, nil
}
