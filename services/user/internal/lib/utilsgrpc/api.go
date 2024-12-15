package utilsgrpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/runtime/protoiface"
)

func UnexpectedError() error {
	return status.Error(codes.Internal, "unexpected error")
}

func WithDetails(code codes.Code, msg string, details ...protoiface.MessageV1) error {
	st := status.New(code, msg)
	ds, err := st.WithDetails(details...)
	if err != nil {
		return st.Err()
	}

	return ds.Err()
}
