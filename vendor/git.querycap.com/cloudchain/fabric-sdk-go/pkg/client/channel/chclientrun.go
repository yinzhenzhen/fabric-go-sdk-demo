/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package channel

import (
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/client/common/discovery/greylist"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/context"
	"git.querycap.com/cloudchain/fabric-sdk-go/pkg/common/providers/fab"
)

func newClient(channelContext context.Channel, membership fab.ChannelMembership, eventService fab.EventService, greylistProvider *greylist.Filter) Client {
	channelClient := Client{
		membership:   membership,
		eventService: eventService,
		greylist:     greylistProvider,
		context:      channelContext,
		metrics:      channelContext.GetMetrics(),
	}
	return channelClient
}
