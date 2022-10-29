package monitor

import (
	"context"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
)

var (
	_ (monitorCommand) = (*monitorCommandDiscover)(nil)
	_ (monitorCommand) = (*monitorCommandDiscoverWithAddr)(nil)
	_ (monitorCommand) = (*monitorCommandPeerFailed)(nil)
)

type monitorCommand interface {
	execute(*Monitor2)
}

type Monitor2 struct {
	// Safe for use outside task
	command_sender chan<- monitorCommand
	host           host.Host
	opts           *options
	exporter       Exporter

	// Unsafe for use outside task
	command_receiver <-chan monitorCommand
	peers            map[peer.ID]*peerTask
}

func Start(ctx context.Context, o ...Option) (*Monitor2, error) {
	opts := defaults()
	if err := apply(opts, o...); err != nil {
		return nil, err
	}

	if opts.Host == nil {
		h, err := createDefaultHost(ctx)
		if err != nil {
			return nil, err
		}
		opts.Host = h
	}

	if opts.Exporter == nil {
		opts.Exporter = NewNoOpExporter()
	}

	command_channel := make(chan monitorCommand)
	m := &Monitor2{
		command_sender: command_channel,
		host:           opts.Host,
		opts:           opts,
		exporter:       opts.Exporter,

		command_receiver: command_channel,
		peers:            map[peer.ID]*peerTask{},
	}

	go m.run(ctx)
	return m, nil
}

func (m *Monitor2) Discover(ctx context.Context, pid peer.ID) {
	m.sendCommand(newMonitorCommandDiscover(pid))
}

func (m *Monitor2) DiscoverWithAddr(ctx context.Context, paddr peer.AddrInfo) {
	m.sendCommand(newMonitorCommandDiscoverWithAddr(paddr))
}

func (m *Monitor2) run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case cmd := <-m.command_receiver:
			cmd.execute(m)
		}
	}
}

func (m *Monitor2) sendCommand(cmd monitorCommand) {
	m.command_sender <- cmd
}

func (m *Monitor2) discover(pid peer.ID) {
	if pt, ok := m.peers[pid]; ok {
		pt.sendCommand(newPeerCommandResetErrors())
		return
	}

	m.peers[pid] = newPeerTask(pid, m.host, m.opts, m.exporter, m)
}

type monitorCommandDiscover struct {
	pid peer.ID
}

func newMonitorCommandDiscover(pid peer.ID) *monitorCommandDiscover {
	return &monitorCommandDiscover{
		pid: pid,
	}
}

// execute implements monitorCommand
func (c *monitorCommandDiscover) execute(m *Monitor2) {
	m.discover(c.pid)
}

type monitorCommandDiscoverWithAddr struct {
	paddr peer.AddrInfo
}

func newMonitorCommandDiscoverWithAddr(paddr peer.AddrInfo) *monitorCommandDiscoverWithAddr {
	return &monitorCommandDiscoverWithAddr{
		paddr: paddr,
	}
}

// execute implements monitorCommand
func (c *monitorCommandDiscoverWithAddr) execute(m *Monitor2) {
	m.host.Peerstore().AddAddr(c.paddr.ID, c.paddr.Addrs[0], peerstore.PermanentAddrTTL)
	m.discover(c.paddr.ID)
}

type monitorCommandPeerFailed struct {
	pid peer.ID
}

func newMonitorCommandPeerFailed(pid peer.ID) *monitorCommandPeerFailed {
	return &monitorCommandPeerFailed{pid}
}

// execute implements monitorCommand
func (c *monitorCommandPeerFailed) execute(m *Monitor2) {
	delete(m.peers, c.pid)
}
