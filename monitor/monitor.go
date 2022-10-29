package monitor

// import (
// 	"context"
// 	"time"

// 	"github.com/diogo464/telemetry"
// 	"github.com/diogo464/telemetry/internal/actionqueue"
// 	"github.com/diogo464/telemetry/monitor/pb"
// 	"github.com/libp2p/go-libp2p-core/host"
// 	"github.com/libp2p/go-libp2p-core/peer"
// 	"github.com/sirupsen/logrus"
// )

// const (
// 	ActionDiscover   actionKind = "discover"
// 	ActionTelemetry  actionKind = "telemetry"
// 	ActionBandwidth  actionKind = "bandwidth"
// 	ActionRemovePeer actionKind = "remove_peer"
// )

// type Cursors = map[string]int

// type actionKind string

// type action struct {
// 	kind actionKind
// 	pid  peer.ID
// }

// type taskResult struct {
// 	kind   actionKind
// 	pid    peer.ID
// 	result interface{}
// }

// type taskTelemetryResult struct {
// 	session telemetry.Session
// 	cursors Cursors
// 	err     error
// }

// type taskBandwidthResult struct {
// 	err error
// }

// type taskProviderRecordsResult struct {
// 	err error
// }

// type peerState struct {
// 	id            peer.ID
// 	failedAttemps int
// 	lastSession   telemetry.Session
// 	cursors       Cursors
// 	ctx           context.Context
// 	cancel        context.CancelFunc
// }

// type Monitor struct {
// 	pb.UnimplementedMonitorServer
// 	h        host.Host
// 	ctx      context.Context
// 	opts     *options
// 	exporter Exporter

// 	peers   map[peer.ID]*peerState
// 	actions *actionqueue.Queue[action]
// 	caction chan actionqueue.Action[action]
// 	cresult chan *taskResult
// }

// func NewMonitor(ctx context.Context, o ...Option) (*Monitor, error) {
// 	opts := defaults()
// 	if err := apply(opts, o...); err != nil {
// 		return nil, err
// 	}

// 	if opts.Host == nil {
// 		h, err := createDefaultHost(ctx)
// 		if err != nil {
// 			return nil, err
// 		}
// 		opts.Host = h
// 	}

// 	if opts.Exporter == nil {
// 		opts.Exporter = NewNoOpExporter()
// 	}

// 	logrus.Debug("options: ", *opts)

// 	return &Monitor{
// 		h:        opts.Host,
// 		ctx:      ctx,
// 		opts:     opts,
// 		exporter: opts.Exporter,
// 		peers:    make(map[peer.ID]*peerState),
// 		actions:  actionqueue.NewQueue[action](),
// 		caction:  make(chan actionqueue.Action[action]),
// 		cresult:  make(chan *taskResult),
// 	}, nil
// }

// func (s *Monitor) Close() {
// 	for _, state := range s.peers {
// 		state.cancel()
// 	}
// 	s.h.Close()
// }

// func (s *Monitor) Run(ctx context.Context) {
// LOOP:
// 	for {
// 		logrus.Debug("Monitor main loop")
// 		actionTimer := s.actions.TimerUntilAction()
// 		select {
// 		case <-actionTimer.C:
// 			action := s.actions.Pop()
// 			s.processAction(action)
// 		case tresult := <-s.cresult:
// 			s.processTaskResult(tresult)
// 		case action := <-s.caction:
// 			s.actions.Push(action)
// 		case <-ctx.Done():
// 			break LOOP
// 		}
// 	}
// }

// func (s *Monitor) processAction(action *action) {
// 	switch action.kind {
// 	case ActionDiscover:
// 		logrus.WithField("peer", action.pid).Debug("action discovery")
// 		s.onActionDiscover(action.pid)
// 	case ActionTelemetry:
// 		logrus.WithField("peer", action.pid).Debug("action telemetry")
// 		s.onActionTelemetry(action.pid)
// 	case ActionBandwidth:
// 		logrus.WithField("peer", action.pid).Debug("action bandwidth")
// 		s.onActionBandwidth(action.pid)
// 	case ActionRemovePeer:
// 		logrus.WithField("peer", action.pid).Debug("action remove peer")
// 		s.onActionRemove(action.pid)
// 	default:
// 		logrus.Error("unhandled action, kind = ", action.kind)
// 	}
// }

// func (s *Monitor) processTaskResult(tresult *taskResult) {
// 	switch tresult.kind {
// 	case ActionTelemetry:
// 		s.onTaskResultTelemetry(tresult.pid, tresult.result.(*taskTelemetryResult))
// 	case ActionBandwidth:
// 		s.onTaskResultBandwidth(tresult.pid, tresult.result.(*taskBandwidthResult))
// 	default:
// 		logrus.Error("unhandled task result, kind = ", tresult.kind)
// 	}
// }

// func (s *Monitor) onActionDiscover(p peer.ID) {
// 	if state, ok := s.peers[p]; ok {
// 		state.failedAttemps = 0
// 	} else {
// 		ctx, cancel := context.WithCancel(context.Background())
// 		s.peers[p] = &peerState{
// 			id:            p,
// 			failedAttemps: 0,
// 			lastSession:   telemetry.InvalidSession,
// 			cursors:       make(map[string]int),
// 			ctx:           ctx,
// 			cancel:        cancel,
// 		}

// 		if s.opts.CollectEnabled {
// 			s.actions.Push(actionqueue.Now(&action{
// 				kind: ActionTelemetry,
// 				pid:  p,
// 			}))
// 		}

// 		if s.opts.BandwidthEnabled {
// 			s.actions.Push(actionqueue.After(&action{
// 				kind: ActionBandwidth,
// 				pid:  p,
// 			}, time.Second*60))
// 		}
// 	}
// }

// func (s *Monitor) onActionTelemetry(p peer.ID) {
// 	if state, ok := s.peers[p]; ok {
// 		go func() {
// 			result := s.taskCollectTelemetry(state.ctx, state.id, state.lastSession, state.cursors)
// 			s.cresult <- &taskResult{
// 				kind:   ActionTelemetry,
// 				pid:    p,
// 				result: result,
// 			}
// 		}()
// 	}
// }

// func (s *Monitor) onActionBandwidth(p peer.ID) {
// 	if state, ok := s.peers[p]; ok {
// 		go func() {
// 			result := s.taskBandwidth(state.ctx, state.id)
// 			s.cresult <- &taskResult{
// 				kind:   ActionBandwidth,
// 				pid:    p,
// 				result: result,
// 			}
// 		}()
// 	}
// }

// func (s *Monitor) onActionRemove(p peer.ID) {
// 	if state, ok := s.peers[p]; ok {
// 		state.cancel()
// 		delete(s.peers, p)
// 	}
// }

// func (s *Monitor) taskCollectTelemetry(pctx context.Context, pid peer.ID, lastSession telemetry.Session, cursors Cursors) *taskTelemetryResult {
// 	ctx, cancel := context.WithTimeout(pctx, s.opts.CollectTimeout)
// 	defer cancel()

// 	logrus.WithField("peer", pid).Debug("creating client")
// 	client, err := telemetry.NewClient(s.h, pid)
// 	if err != nil {
// 		return &taskTelemetryResult{err: err}
// 	}
// 	defer client.Close()

// 	logrus.WithField("peer", pid).Debug("getting session info")
// 	session, err := client.GetSession(ctx)
// 	if err != nil {
// 		return &taskTelemetryResult{err: err}
// 	}

// 	if session != lastSession {
// 		cursors = make(map[string]int)
// 		lastSession = session
// 	}

// 	logrus.WithField("peer", pid).Debug("exporting datapoints from ", cursors)

// 	streamDescriptors, err := client.AvailableStreams(ctx)
// 	if err != nil {
// 		return &taskTelemetryResult{err: err}
// 	}

// 	for _, descriptor := range streamDescriptors {
// 		segments := 0
// 		schan := make(chan telemetry.StreamSegment)
// 		go client.StreamSegments(ctx, descriptor.Name, cursors[descriptor.Name], schan)

// 		for segment := range schan {
// 			s.exporter.ExportStream(pid, session, descriptor, segment)
// 			segments += 1
// 		}

// 		logrus.
// 			WithField("peer", pid).
// 			WithField("stream", descriptor.Name).
// 			WithField("segments", segments).
// 			Debug("exported stream")
// 	}

// 	return &taskTelemetryResult{
// 		session: session,
// 		cursors: cursors,
// 		err:     nil,
// 	}
// }

// func (s *Monitor) taskBandwidth(pctx context.Context, pid peer.ID) *taskBandwidthResult {
// 	ctx, cancel := context.WithTimeout(pctx, s.opts.BandwidthTimeout)
// 	defer cancel()

// 	client, err := telemetry.NewClient(s.h, pid)
// 	if err != nil {
// 		return &taskBandwidthResult{err: err}
// 	}
// 	defer client.Close()

// 	session, err := client.GetSession(ctx)
// 	if err != nil {
// 		return &taskBandwidthResult{err: err}
// 	}

// 	bandwidth, err := client.Bandwidth(ctx, telemetry.DEFAULT_BANDWIDTH_PAYLOAD_SIZE)
// 	if err != nil {
// 		return &taskBandwidthResult{err: err}
// 	}

// 	s.exporter.ExportBandwidth(pid, session, bandwidth)

// 	return &taskBandwidthResult{err: nil}
// }

// func (s *Monitor) onTaskResultTelemetry(pid peer.ID, tresult *taskTelemetryResult) {
// 	s.onTaskResultCommon(pid, ActionTelemetry, s.opts.CollectPeriod, tresult.err)
// 	if state, ok := s.peers[pid]; ok && tresult.err == nil {
// 		state.cursors = tresult.cursors
// 		state.lastSession = tresult.session
// 	}
// }

// func (s *Monitor) onTaskResultBandwidth(pid peer.ID, tresult *taskBandwidthResult) {
// 	s.onTaskResultCommon(pid, ActionBandwidth, s.opts.BandwidthPeriod, tresult.err)
// }

// func (s *Monitor) onTaskResultProviderRecords(pid peer.ID, tresult *taskProviderRecordsResult) {
// 	s.onTaskResultCommon(pid, ActionProviderRecords, s.opts.ProviderRecordsPeriod, tresult.err)
// }

// func (s *Monitor) onTaskResultCommon(pid peer.ID, kind actionKind, interval time.Duration, err error) {
// 	if state, ok := s.peers[pid]; ok {
// 		if err == nil {
// 			logrus.WithField("peer", pid).Debug("running action ", kind, " in ", interval)
// 			state.failedAttemps = 0
// 			s.actions.Push(actionqueue.After(&action{
// 				kind: kind,
// 				pid:  pid,
// 			}, interval))
// 		} else {
// 			state.failedAttemps += 1
// 			logrus.WithField("peer", pid).Debug("removing peer")
// 			if state.failedAttemps >= s.opts.MaxFailedAttemps {
// 				s.actions.Push(actionqueue.Now(&action{
// 					kind: ActionRemovePeer,
// 					pid:  pid,
// 				}))
// 			} else {
// 				logrus.WithField("peer", pid).Debug("retrying action ", kind, " in ", s.opts.RetryInterval, ". error: ", err)
// 				s.actions.Push(actionqueue.After(&action{
// 					kind: kind,
// 					pid:  pid,
// 				}, s.opts.RetryInterval))
// 			}
// 		}
// 	}
// }
