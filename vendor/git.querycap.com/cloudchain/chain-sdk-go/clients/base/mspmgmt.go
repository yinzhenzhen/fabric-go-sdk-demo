package base

import (
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/msp"
)

func (c *Client) MspmgmtGetSigningIdentity(
	id string,
) (msp.SigningIdentity, error) {
	return c.mspmgmtClient.GetSigningIdentity(id)
}
