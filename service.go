package telemetry

import (
	"context"
	"fmt"
	"net"

	"github.com/diogo464/telemetry/internal/bpool"
	"github.com/diogo464/telemetry/internal/otlp_exporter"
	"github.com/diogo464/telemetry/internal/pb"
	"github.com/diogo464/telemetry/internal/stream"
	"github.com/libp2p/go-libp2p/core/host"
	sdk_metric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc"
)

type ServiceAccessType string

const (
	ServiceAccessPublic     ServiceAccessType = "public"
	ServiceAccessRestricted ServiceAccessType = "restricted"
	ServiceAccessDisabled   ServiceAccessType = "disabled"
)

type Service struct {
	pb.UnimplementedTelemetryServer
	// current session, randomly generated uuid
	session        Session
	opts           *serviceOptions
	host           host.Host
	meter_provider *serviceMeterProvider
	bufferPool     *bpool.Pool

	ctx        context.Context
	cancel     context.CancelFunc
	grpcServer *grpc.Server

	serviceAcl *serviceAccessControl
	streams    *serviceStreams
	metrics    *serviceMetrics
	properties *serviceProperties
	events     *serviceEvents

	downloadBlocker *requestBlocker
	uploadBlocker   *requestBlocker
}

func NewService(h host.Host, os ...ServiceOption) (*Service, error) {
	opts := serviceDefaults()
	err := serviceApply(opts, os...)
	if err != nil {
		return nil, err
	}

	bufferPool := bpool.New(bpool.WithAllocSize(64*1024), bpool.WithMaxSize(8*1024*1024))
	streams := newServiceStreams(
		stream.WithActiveBufferLifetime(opts.activeBufferDuration),
		stream.WithSegmentLifetime(opts.windowDuration),
		stream.WithPool(bufferPool),
	)

	ctx, cancel := context.WithCancel(context.Background())

	t := &Service{
		session:        RandomSession(),
		opts:           opts,
		host:           h,
		meter_provider: nil,
		bufferPool:     bufferPool,

		ctx:    ctx,
		cancel: cancel,

		serviceAcl: newServiceAccessControl(opts.serviceAccessType, opts.serviceAccessWhitelist),
		streams:    streams,
		metrics: newServiceMetrics(streams.create(&pb.StreamType{
			Type: &pb.StreamType_Metric{},
		}).stream),
		properties: newServiceProperties(),
		events:     newServiceEvents(streams),

		downloadBlocker: newRequestBlocker(),
		uploadBlocker:   newRequestBlocker(),
	}

	if opts.enableBandwidth {
		h.SetStreamHandler(ID_UPLOAD, t.uploadHandler)
		h.SetStreamHandler(ID_DOWNLOAD, t.downloadHandler)
	}

	var listener net.Listener
	if opts.listener == nil {
		listener, err = newServiceListener(h, ID_TELEMETRY, t.serviceAcl)
		if err != nil {
			return nil, err
		}
	} else {
		listener = opts.listener
	}

	grpc_server := grpc.NewServer()
	pb.RegisterTelemetryServer(grpc_server, t)
	t.grpcServer = grpc_server

	go func() {
		err = grpc_server.Serve(listener)
		if err != nil {
			fmt.Println(err)
		}
		log.Info("grpc server stopped")
	}()

	res := resource.Default()
	if opts.otelResource != nil {
		res = opts.otelResource
	}

	otlpExporter := otlp_exporter.New(t.metrics.stream)
	sdk_meter_provider := sdk_metric.NewMeterProvider(
		sdk_metric.WithResource(res),
		sdk_metric.WithReader(
			sdk_metric.NewPeriodicReader(
				otlpExporter,
				sdk_metric.WithInterval(opts.metricsPeriod),
			),
		),
		sdk_metric.WithView(opts.otelViews...),
	)
	t.meter_provider = newServiceMeterProvider(t, sdk_meter_provider)

	return t, nil
}

func (s *Service) Context() context.Context {
	return s.ctx
}

func (s *Service) Close() {
	s.grpcServer.GracefulStop()
	s.cancel()
}
