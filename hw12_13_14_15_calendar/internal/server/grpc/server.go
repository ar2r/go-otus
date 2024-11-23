package grpcserver

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	pb "github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server/grpc/protobuf"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	conf              Config
	grpcServerService pb.EventServiceServer
	grpcServer        *grpc.Server
	logg              *slog.Logger
}

type Application interface { // TODO
}

func NewServer(logg *slog.Logger, conf Config, serviceServer pb.EventServiceServer) *Server {
	return &Server{
		conf:              conf,
		grpcServerService: serviceServer,
		logg:              logg,
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

	s.grpcServer = grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			CustomLoggingInterceptor(s.logg),
		),
	)

	pb.RegisterEventServiceServer(s.grpcServer, s.grpcServerService)
	reflection.Register(s.grpcServer) // Включил для отладки

	s.logg.Info(fmt.Sprintf("grpc server started on %s", lsn.Addr().String()))
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
	return nil
}
