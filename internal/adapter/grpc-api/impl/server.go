package adapter_grpc_api

import (
	context "context"
	"net"
	"sync"

	"github.com/goregion/hexago/internal/adapter/grpc-api/gen"
	"github.com/goregion/hexago/internal/entity"
	"github.com/goregion/hexago/pkg/tools"
	"github.com/pkg/errors"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

type Server struct {
	gen.UnimplementedOHLCServiceServer
	addr         string
	ohlcSessions sync.Map // symbol, grpc.ServerStreamingServer[OHLC]
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

func (s *Server) PublishOHLC(ctx context.Context, ohlc *entity.OHLC) error {
	s.ohlcSessions.Range(func(key, value any) bool {
		symbol := key.(string)
		stream := value.(grpc.ServerStreamingServer[gen.OHLC])
		if symbol == ohlc.Symbol {
			if err := stream.Send(mustMarshalOHLC(ohlc)); err != nil {
				return true // Just skip this stream, defer will clean it up eventually
			}
		}
		return true
	})
	return nil
}

func (s *Server) SubscribeToOHLCStream(request *gen.SubscribeToOHLCStreamRequest, stream grpc.ServerStreamingServer[gen.OHLC]) error {
	s.ohlcSessions.Store(request.Symbol, stream)
	defer s.ohlcSessions.Delete(request.Symbol)

	<-stream.Context().Done()
	return status.Errorf(codes.Canceled, "client canceled, abandoning")
}

func (s *Server) RunBlocked(ctx context.Context) error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return errors.Wrap(err, "failed to listen")
	}
	var grpcServer = grpc.NewServer()
	gen.RegisterOHLCServiceServer(grpcServer, s)

	err = tools.RunAsyncBlocked(ctx,
		func(context.Context) error {
			return grpcServer.Serve(listener)
		},
	)
	grpcServer.GracefulStop()
	if err != nil {
		return errors.Wrap(err, "grpc server exited with error")
	}
	return nil
}
