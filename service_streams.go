package telemetry

import (
	"sync"

	"github.com/diogo464/telemetry/internal/stream"
)

type serviceStreamId uint32

type serviceStream struct {
	stream   *stream.Stream
	streamId serviceStreamId
}

type serviceStreams struct {
	mu             sync.Mutex
	streams        map[serviceStreamId]*serviceStream
	nextID         serviceStreamId
	defaultOptions []stream.Option
}

func newServiceStreams(defaultOptions ...stream.Option) *serviceStreams {
	return &serviceStreams{
		streams:        make(map[serviceStreamId]*serviceStream),
		defaultOptions: defaultOptions,
	}
}

func (s *serviceStreams) create(options ...stream.Option) *serviceStream {
	s.mu.Lock()
	defer s.mu.Unlock()

	options = append(options, s.defaultOptions...)

	stream := stream.New(options...)

	id := s.nextID
	s.nextID++

	s.streams[id] = &serviceStream{
		stream:   stream,
		streamId: id,
	}

	return s.streams[id]
}

func (s *serviceStreams) has(id serviceStreamId) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.streams[id]
	return ok
}

func (s *serviceStreams) get(id serviceStreamId) *serviceStream {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.streams[id]
}
