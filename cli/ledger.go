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

	// ???????????????
	blockInfo := &BlockInfo{}
	// ?????????
	blockInfo.BlockNum = block.Header.Number
	// ?????????????????????
	blockInfo.PrevBlockHash = hex.EncodeToString(block.Header.PreviousHash)
	// ????????????
	blockInfo.BlockHash = hex.EncodeToString(utils.BlockHeaderHash(block.Header))
	// ????????????
	blockInfo.DataHash = hex.EncodeToString(block.Header.DataHash)
	// ????????????
	blockInfo.BlockSize = uint64(len(block.String()))

	var trans []*Transaction
	for _, envBytes := range block.Data.Data {
		var env *common.Envelope
		var err error

		// ??????????????????Envelope
		if env, err = utils.GetEnvelopeFromBlock(envBytes); err != nil {
			return nil
		}
		tran, err := parseTransaction(env)
		if err != nil {
			log.Debug("parseTransaction err %v", err)
		}
		trans = append(trans, tran)
	}

	// ????????????
	blockInfo.BlockTime = trans[len(trans)-1].TxTime

	// ?????????
	blockInfo.TxCount = int32(len(trans))

	blockInfo.Trans = trans
	return blockInfo
}

func parseTransaction(env *common.Envelope) (*Transaction, error) {

	tran := &Transaction{}
	var err error
	var payload *common.Payload

	// ???Envelope??????payload
	if payload, err = utils.UnmarshalPayload(env.Payload); err != nil {
		return nil, err
	}

	var chdr *common.ChannelHeader
	// ???payload???????????????????????????
	if chdr, err = utils.UnmarshalChannelHeader(payload.Header.ChannelHeader); err != nil {
		return nil, err
	}

	// ??????id
	tran.TxId = chdr.TxId
	// ????????????
	tran.ValidationCode = ""
	// ????????????
	tran.Type = strconv.Itoa(int(chdr.Type))
	// ????????????
	//tran.TransSize = uint64(len(string(envBytes)))
	// ????????????
	tran.TxTime = time.Unix(chdr.Timestamp.Seconds, int64(chdr.Timestamp.Nanos)).Format(time.UnixDate)

	//var cche *peer.ChaincodeHeaderExtension
	//if cche, err = gmutils.GetChaincodeHeaderExtension(payload.Header); err != nil {
	//	return nil, err
	//}
	//// ????????????
	//tran.ChainCodeName = cche.ChaincodeId.Name + ":" + cche.ChaincodeId.Version

	if common.HeaderType(chdr.Type) == common.HeaderType_ENDORSER_TRANSACTION {

		var tx *peer.Transaction
		// ???payload???????????????data
		if tx, err = utils.UnmarshalTransaction(payload.Data); err != nil {
			return nil, err
		}

		if len(tx.Actions) > 0 {

			action := tx.Actions[0]
			// ????????????????????????????????????
			AShdr, err := utils.UnmarshalSignatureHeader(action.Header)

			if err != nil {
				return nil, err
			}

			//????????????????????????
			var subject string //mspid
			if _, subject, err = decodeSerializedIdentity(AShdr.Creator); err != nil {
				return nil, err
			}
			// ?????????mspid
			tran.CreaterMSPId = subject

			//var capayload *peer.ChaincodeActionPayload
			var ca *peer.ChaincodeAction
			var cap *peer.ChaincodeActionPayload
			// ????????????????????????payload
			if cap, ca, err = utils.GetPayloads(action); err != nil {
				return nil, err
			}

			// ??????????????????
			//tran.Input = string(cap.ChaincodeProposalPayload)
			tran.Input = base64.StdEncoding.EncodeToString(cap.ChaincodeProposalPayload)
			// ????????????
			tran.ChainCodeName = ca.ChaincodeId.Name
			// ??????id
			tran.ChainCodeId = ca.ChaincodeId.Name + ":" + ca.ChaincodeId.Version
			// ??????
			tran.Status = ca.Response.Status
			// ??????????????????
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
