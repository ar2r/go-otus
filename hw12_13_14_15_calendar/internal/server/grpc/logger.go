package grpcserver

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"time"

	"github.com/ar2r/go-otus/hw12_13_14_15_calendar/internal/server"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

func CustomLoggingInterceptor(logg *slog.Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()
		p, ok := peer.FromContext(ctx)
		clientIP := "unknown"
		if ok {
			clientIP, _, _ = net.SplitHostPort(p.Addr.String())
			if ip, err := server.NormalizeIPv4(clientIP); err == nil {
				clientIP = ip
			}
		}

		md, ok := metadata.FromIncomingContext(ctx)
		userAgent := "unknown_user_agent"
		if ok {
			if ua, exists := md["user-agent"]; exists && len(ua) > 0 {
				userAgent = ua[0]
			}
		}

		resp, err := handler(ctx, req)
		statusCode := grpcCodeToHTTPStatus(status.Code(err))
		duration := time.Since(start)

		logg.Info(fmt.Sprintf("%s [%s] %s %s %s %d %d \"%s\"",
			clientIP,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			"GRPC",
			info.FullMethod,
			"HTTP/2.0",
			statusCode,
			duration.Milliseconds(),
			userAgent,
		))

		return resp, err
	}
}

func grpcCodeToHTTPStatus(code codes.Code) int {
	switch code {
	case codes.OK:
		return 200
	case codes.Canceled:
		return 499
	case codes.Unknown:
		return 500
	case codes.InvalidArgument:
		return 400
	case codes.DeadlineExceeded:
		return 504
	case codes.NotFound:
		return 404
	case codes.AlreadyExists:
		return 409
	case codes.PermissionDenied:
		return 403
	case codes.ResourceExhausted:
		return 429
	case codes.FailedPrecondition:
		return 400
	case codes.Aborted:
		return 409
	case codes.OutOfRange:
		return 400
	case codes.Unimplemented:
		return 501
	case codes.Internal:
		return 500
	case codes.Unavailable:
		return 503
	case codes.DataLoss:
		return 500
	case codes.Unauthenticated:
		return 401
	default:
		return 500
	}
}
