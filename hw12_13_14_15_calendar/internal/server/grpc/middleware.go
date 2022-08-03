package grpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

var ErrPeerFromContext = status.Error(codes.Internal, "getting peer fail")

func LoggingInterceptor(logger Logger) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (response interface{}, err error) {
		p, ok := peer.FromContext(ctx)
		if !ok {
			logger.Error(ErrPeerFromContext)

			return response, ErrPeerFromContext
		}

		start := time.Now()
		response, err = handler(ctx, req)

		logger.LogGRPCRequest(
			status.Code(err),
			info.FullMethod,
			p.Addr.String(),
			time.Since(start),
		)

		return response, err
	}
}