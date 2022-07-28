package telemetry

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"

	"git.d464.sh/uni/telemetry/pb"
	"git.d464.sh/uni/telemetry/utils"
	"github.com/gogo/protobuf/types"
	"github.com/libp2p/go-libp2p-core/host"
	"github.com/libp2p/go-libp2p-core/peer"
	gostream "github.com/libp2p/go-libp2p-gostream"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var ErrInvalidResponse = fmt.Errorf("invalid response")
var ErrNotUsingLibp2p = fmt.Errorf("not using libp2p")

type Client struct {
	// Can be null if we are not connected using libp2p
	h host.Host
	p peer.ID

	c *grpc.ClientConn
}

func NewClient(h host.Host, p peer.ID) (*Client, error) {
	conn, err := grpc.Dial(
		"",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithContextDialer(func(ctx context.Context, s string) (net.Conn, error) {
			conn, err := gostream.Dial(ctx, h, p, ID_TELEMETRY)
			return conn, err
		}))

	if err != nil {
		return nil, err
	}

	return &Client{
		h: h,
		p: p,
		c: conn,
	}, nil
}

func NewClient2(target string) (*Client, error) {
	conn, err := grpc.Dial(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Client{
		h: nil,
		p: "",
		c: conn,
	}, nil
}

func (c *Client) Close() {
	c.c.Close()
}

func (c *Client) SessionInfo(ctx context.Context) (*SessionInfo, error) {
	client, err := c.newGrpcClient()
	if err != nil {
		return nil, err
	}

	response, err := client.GetSessionInfo(ctx, &types.Empty{})
	if err != nil {
		return nil, err
	}

	session, err := ParseSession(response.GetSession())
	if err != nil {
		return nil, err
	}

	return &SessionInfo{
		Session:  session,
		BootTime: utils.TimeFromPB(response.GetBootTime()),
	}, nil
}

func (c *Client) AvailableStreams(ctx context.Context) ([]StreamDescriptor, error) {
	client, err := c.newGrpcClient()
	if err != nil {
		return nil, err
	}

	stream, err := client.GetAvailableStreams(ctx, &pb.GetAvailableStreamsRequest{})
	if err != nil {
		return nil, err
	}

	streams := make([]StreamDescriptor, 0)
	for {
		avail, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		streams = append(streams, StreamDescriptor{
			Name:     avail.GetName(),
			Period:   utils.DurationFromPB(avail.GetPeriod()),
			Encoding: avail.GetEncoding(),
		})
	}

	return streams, nil
}

func (c *Client) Segments(ctx context.Context, stream string, since int) ([]StreamSegment, error) {
	client, err := c.newGrpcClient()
	if err != nil {
		return nil, err
	}

	srv, err := client.GetStream(ctx, &pb.GetStreamRequest{
		Stream: stream,
		Seqn:   uint32(since),
	})
	if err != nil {
		return nil, err
	}

	segments := make([]StreamSegment, 0)

	for {
		seg, err := srv.Recv()
		if err == io.EOF {
			break
		}

		segments = append(segments, StreamSegment{
			SeqN: int(seg.GetSeqn()),
			Data: seg.GetData(),
		})
	}

	return segments, nil
}

func (c *Client) StreamSegments(ctx context.Context, stream string, since int, ch chan<- StreamSegment) error {
	client, err := c.newGrpcClient()
	if err != nil {
		return err
	}

	srv, err := client.GetStream(ctx, &pb.GetStreamRequest{
		Stream:    stream,
		Seqn:      uint32(since),
		Keepalive: 1,
	})
	if err != nil {
		return err
	}

	for {
		seg, err := srv.Recv()
		if err != nil {
			if err == io.EOF || err == ctx.Err() {
				break
			}
			return err
		}
		ch <- StreamSegment{
			SeqN: int(seg.GetSeqn()),
			Data: seg.GetData(),
		}
	}

	return nil
}

func (c *Client) AvailableProperties(ctx context.Context) ([]PropertyDescriptor, error) {
	client, err := c.newGrpcClient()
	if err != nil {
		return nil, err
	}

	response, err := client.GetAvailableProperties(context.Background(), &pb.GetAvailablePropertiesRequest{})
	if err != nil {
		return nil, err
	}

	properties := make([]PropertyDescriptor, 0)
	for {
		prop, err := response.Recv()
		if err != nil {
			if err == io.EOF || err == ctx.Err() {
				break
			}
			return nil, err
		}
		properties = append(properties, PropertyDescriptor{
			Name:     prop.GetName(),
			Encoding: prop.GetEncoding(),
		})
	}

	return properties, nil
}

func (c *Client) Property(ctx context.Context, property string) ([]byte, error) {
	client, err := c.newGrpcClient()
	if err != nil {
		return nil, err
	}

	response, err := client.GetProperty(ctx, &pb.GetPropertyRequest{
		Property: property,
	})
	if err != nil {
		return nil, err
	}

	var writer bytes.Buffer
	for {
		data, err := response.Recv()
		if err != nil {
			if err == io.EOF || err == ctx.Err() {
				break
			}
			return nil, err
		}
		writer.Write(data.GetData())
	}
	return writer.Bytes(), nil
}

func (c *Client) SystemInfo(ctx context.Context) (*SystemInfo, error) {
	client, err := c.newGrpcClient()
	if err != nil {
		return nil, err
	}

	response, err := client.GetSystemInfo(ctx, &types.Empty{})
	if err != nil {
		return nil, err
	}

	return &SystemInfo{
		OS:     response.Os,
		Arch:   response.Arch,
		NumCPU: response.Numcpu,
	}, nil
}

func (c *Client) Download(ctx context.Context, payload uint32) (uint32, error) {
	if c.h == nil {
		return 0, ErrNotUsingLibp2p
	}

	stream, err := c.h.NewStream(ctx, c.p, ID_DOWNLOAD)
	if err != nil {
		return 0, err
	}

	if err := utils.WriteU32(stream, payload); err != nil {
		return 0, err
	}

	if _, err := io.Copy(stream, io.LimitReader(utils.NullReader{}, int64(payload))); err != nil {
		return 0, err
	}

	rate, err := utils.ReadU32(stream)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func (c *Client) Upload(ctx context.Context, payload uint32) (uint32, error) {
	if c.h == nil {
		return 0, ErrNotUsingLibp2p
	}

	stream, err := c.h.NewStream(ctx, c.p, ID_UPLOAD)
	if err != nil {
		return 0, err
	}

	if err := utils.WriteU32(stream, payload); err != nil {
		return 0, err
	}

	if _, err := io.Copy(io.Discard, io.LimitReader(stream, int64(payload))); err != nil {
		return 0, err
	}

	rate, err := utils.ReadU32(stream)
	if err != nil {
		return 0, err
	}

	return rate, nil
}

func (c *Client) Bandwidth(ctx context.Context, payload uint32) (Bandwidth, error) {
	download, err := c.Download(ctx, payload)
	if err != nil {
		return Bandwidth{}, err
	}
	upload, err := c.Upload(ctx, payload)
	if err != nil {
		return Bandwidth{}, err
	}
	return Bandwidth{
		UploadRate:   upload,
		DownloadRate: download,
	}, nil
}

func (c *Client) Debug(ctx context.Context) (*Debug, error) {
	client, err := c.newGrpcClient()
	if err != nil {
		return nil, err
	}

	pbdbg, err := client.GetDebug(ctx, &types.Empty{})
	if err != nil {
		return nil, err
	}

	streams := make([]DebugStream, 0, len(pbdbg.GetStreams()))
	for _, pbs := range pbdbg.GetStreams() {
		streams = append(streams, DebugStream{
			Name:      pbs.Name,
			UsedSize:  pbs.Used,
			TotalSize: pbs.Total,
		})
	}

	return &Debug{
		Streams: streams,
	}, nil
}

func (c *Client) newGrpcClient() (pb.TelemetryClient, error) {
	return pb.NewTelemetryClient(c.c), nil
}
