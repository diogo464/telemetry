package main

import (
	"fmt"

	"github.com/diogo464/telemetry"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/urfave/cli/v2"
)

func clientFromContext(c *cli.Context) (*telemetry.Client, error) {
	switch c.String(FLAG_CONN_TYPE.Name) {
	case "libp2p":

		info, err := peer.AddrInfoFromString(c.String(FLAG_HOST.Name))
		if err != nil {
			return nil, err
		}
		h, err := libp2p.New(libp2p.NoListenAddrs)
		if err != nil {
			return nil, err
		}
		fmt.Println(info)
		h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)
		return telemetry.NewClient(c.Context, telemetry.WithClientLibp2pDial(h, info.ID))
	case "tcp":
		return telemetry.NewClient(c.Context, telemetry.WithClientGrpcDial(c.String(FLAG_HOST.Name)))
	default:
		return nil, fmt.Errorf("unknown connection type: %s", c.String(FLAG_CONN_TYPE.Name))
	}
}
