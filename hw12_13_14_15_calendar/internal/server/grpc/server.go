package grpcserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/config"
	pb "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc/protobuf"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	conf              config.GrpcServerConf
	grpcServerService pb.EventServiceServer
	grpcServer        *grpc.Server
}

type Logger interface {
	InfoRaw(msg string)
	Info(msg string)
	Error(msg string)
}

type Application interface { // TODO
}

func NewServer(conf config.GrpcServerConf, serviceServer pb.EventServiceServer) *Server {
	return &Server{
		conf:              conf,
		grpcServerService: serviceServer,
	}
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (s *Server) Run() error {
	listenAddr := fmt.Sprintf("%s:%d", s.conf.Host, s.conf.Port)
	lsn, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return err

	}
	opts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall, logging.FinishCall),
	}
	logg := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{}))
	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logg), opts...),
		),
	)

	pb.RegisterEventServiceServer(s.grpcServer, s.grpcServerService)
	reflection.Register(s.grpcServer) // Включил для отладки

	logg.Info(fmt.Sprintf("grpc server started on %s", lsn.Addr().String()))
	if err := s.grpcServer.Serve(lsn); err != nil {
		log.Fatal(err)
	}
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	_, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}
	// logg.Info("HTTP server stopped")
	return nil
}
