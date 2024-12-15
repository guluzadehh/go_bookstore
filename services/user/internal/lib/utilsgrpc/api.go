package utilsgrpc

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
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

func ValidationError(errs validator.ValidationErrors) []*errdetails.BadRequest_FieldViolation {
	var errDetails []*errdetails.BadRequest_FieldViolation

	for _, err := range errs {
		field := err.Field()

		var msg string

		switch err.ActualTag() {
		case "required":
			msg = fmt.Sprintf("field %s is required.", field)
		case "min":
			msg = fmt.Sprintf("field %s min length must be %s.", field, err.Param())
		case "max":
			msg = fmt.Sprintf("field %s max length must be %s.", field, err.Param())
		case "alphanum":
			msg = fmt.Sprintf("field %s must contain both letters and numbers.", field)
		case "containsany":
			msg = fmt.Sprintf("field %s must contain on of the following characters: %s.", field, err.Param())
		case "eqfield":
			msg = fmt.Sprintf("field %s is not equal to %s field.", field, err.Param())
		case "passwordpattern":
			msg = "field password must contain at least one letter, one number, and one special character."
		default:
			msg = fmt.Sprintf("field %s is invalid", field)
		}

		errDetails = append(errDetails, &errdetails.BadRequest_FieldViolation{
			Field:       field,
			Description: msg,
		})
	}

	return errDetails
}
