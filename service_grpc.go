package telemetry

import (
	"context"

	"github.com/diogo464/telemetry/internal/pb"
	"github.com/diogo464/telemetry/internal/stream"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	TAG_UPLOAD     = "upload"
	TAG_DOWNLOAD   = "download"
	TAG_STREAM     = "getdatapoints"
	TAG_GETRECORDS = "getrecords"
)

var (
	ErrBlocked              = status.Errorf(codes.Unavailable, "blocked")
	ErrStreamNotAvailable   = status.Errorf(codes.NotFound, "stream not available")
	ErrPropertyNotAvailable = status.Errorf(codes.NotFound, "property not available")
	ErrCaptureNotAvailable  = status.Errorf(codes.NotFound, "capture not available")
	ErrEventNotAvailable    = status.Errorf(codes.NotFound, "event not available")
)

func (s *Service) GetSession(ctx context.Context, req *pb.GetSessionRequest) (*pb.GetSessionResponse, error) {
	return &pb.GetSessionResponse{
		Uuid: s.session.String(),
	}, nil
}

func (s *Service) GetMetricDescriptors(req *pb.GetMetricDescriptorsRequest, srv pb.Telemetry_GetMetricDescriptorsServer) error {
	s.metrics.mu.Lock()
	pbdescs := make([]*pb.MetricDescriptor, 0, len(s.metrics.descriptors))
	pbdescs = append(pbdescs, s.metrics.descriptors...)
	s.metrics.mu.Unlock()

	for _, desc := range pbdescs {
		if err := srv.Send(desc); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetMetrics(req *pb.GetMetricsRequest, srv pb.Telemetry_GetMetricsServer) error {
	stream := s.metrics.stream
	since := req.GetSequenceNumberSince()
	return grpcSendStreamSegments(stream, int(since), srv)
}

func (s *Service) GetPropertyDescriptors(req *pb.GetPropertyDescriptorsRequest, srv pb.Telemetry_GetPropertyDescriptorsServer) error {
	properties := s.properties

	properties.mu.Lock()
	pbdescs := make([]*pb.PropertyDescriptor, 0, len(properties.properties))
	for idx, p := range properties.properties {
		pbdescs = append(pbdescs, &pb.PropertyDescriptor{
			Id:          uint32(idx),
			Scope:       p.pbproperty.GetScope(),
			Name:        p.pbproperty.GetName(),
			Description: p.pbproperty.GetDescription(),
		})
	}
	properties.mu.Unlock()

	for _, desc := range pbdescs {
		if err := srv.Send(desc); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetProperties(req *pb.GetPropertiesRequest, srv pb.Telemetry_GetPropertiesServer) error {
	properties := s.properties

	properties.mu.Lock()
	pbprops := make([]*pb.Property, 0, len(properties.properties))
	for _, p := range properties.properties {
		pbprops = append(pbprops, p.pbproperty)
	}
	properties.mu.Unlock()

	for _, v := range pbprops {
		if err := srv.Send(v); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetCaptureDescriptors(req *pb.GetCaptureDescriptorsRequest, srv pb.Telemetry_GetCaptureDescriptorsServer) error {
	captures := s.captures

	captures.mu.Lock()
	descriptors := make([]*pb.CaptureDescriptor, 0, len(captures.captures))
	for _, c := range captures.captures {
		descriptors = append(descriptors, c.pbdescriptor)
	}
	captures.mu.Unlock()

	for _, v := range descriptors {
		if err := srv.Send(v); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) GetCapture(req *pb.GetCaptureRequest, srv pb.Telemetry_GetCaptureServer) error {
	captures := s.captures

	captures.mu.Lock()
	capture := captures.captures[req.GetId()]
	captures.mu.Unlock()

	if capture == nil {
		return ErrCaptureNotAvailable
	}
	return grpcSendStreamSegments(capture.stream, int(req.GetSequenceNumberSince()), srv)
}

func (s *Service) GetEventDescriptors(req *pb.GetEventDescriptorsRequest, srv pb.Telemetry_GetEventDescriptorsServer) error {
	events := s.events
	events.mu.Lock()
	descriptors := make([]*pb.EventDescriptor, 0, len(events.events))
	for _, e := range events.events {
		descriptors = append(descriptors, e.descriptor)
	}
	events.mu.Unlock()

	for _, desc := range descriptors {
		if err := srv.Send(desc); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) GetEvent(req *pb.GetEventRequest, srv pb.Telemetry_GetEventServer) error {
	events := s.events

	events.mu.Lock()
	e := events.events[req.GetId()]
	if e == nil {
		return ErrEventNotAvailable
	}
	evstream := e.stream
	events.mu.Unlock()

	return grpcSendStreamSegments(evstream, int(req.GetSequenceNumberSince()), srv)
}

func grpcSendStreamSegments(stream *stream.Stream, since int, srv grpcSegmentSender) error {
	for {
		segments := stream.Segments(since, 128)
		if len(segments) == 0 {
			break
		}
		since += len(segments)
		for _, segment := range segments {
			err := srv.Send(&pb.StreamSegment{
				SequenceNumber: uint32(segment.SeqN),
				Data:           segment.Data,
			})
			if err != nil {
				return err
			}
		}
	}
	return nil
}
