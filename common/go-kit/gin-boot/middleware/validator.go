package middleware

import (
	"context"
	"github.com/yuansuan/ticp/common/go-kit/gin-boot/middleware/validatepb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validatorError interface {
	Field() string
	Reason() string
	Key() bool
	Cause() error
	ErrorName() string
}

type validator interface {
	Validate() error
}

func validate(req interface{}) error {
	switch v := req.(type) {
	case validator:
		if err := v.Validate(); err != nil {
			status := status.New(codes.InvalidArgument, err.Error())
			if verr, ok := err.(validatorError); ok {
				status, err = status.WithDetails(&validatepb.ValidateError{
					Field:  verr.Field(),
					Reason: verr.Reason(),
				})
			}

			return status.Err()
		}
	}
	return nil
}

func ValidatorUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if err := validate(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}
}
