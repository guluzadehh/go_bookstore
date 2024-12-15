package usergrpc

import (
	"context"
	"errors"
	"log/slog"

	"github.com/go-playground/validator/v10"
	user_grpc "github.com/guluzadehh/go_bookstore/protos/gen/go/user"
	"github.com/guluzadehh/go_bookstore/services/user/internal/config"
	"github.com/guluzadehh/go_bookstore/services/user/internal/domain/models"
	grpchandler "github.com/guluzadehh/go_bookstore/services/user/internal/grpc"
	"github.com/guluzadehh/go_bookstore/services/user/internal/lib/sl"
	"github.com/guluzadehh/go_bookstore/services/user/internal/lib/utilsgrpc"
	"github.com/guluzadehh/go_bookstore/services/user/internal/lib/validators"
	service "github.com/guluzadehh/go_bookstore/services/user/internal/services"
	typesgrpc "github.com/guluzadehh/go_bookstore/services/user/internal/types/grpc"
	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

type UserService interface {
	CreateUser(ctx context.Context, email string, password string) (*models.User, error)
}

type UserHandler struct {
	*grpchandler.Handler
	cfg  *config.Config
	srvc UserService
	user_grpc.UnimplementedUserServiceServer
}

func New(log *slog.Logger, cfg *config.Config, userService UserService) *UserHandler {
	return &UserHandler{
		Handler: grpchandler.NewHandler(log),
		cfg:     cfg,
		srvc:    userService,
	}
}

func Register(grpcServer *grpc.Server, userHandler *UserHandler) {
	user_grpc.RegisterUserServiceServer(grpcServer, userHandler)
}

type CreateUserRequest struct {
	Email        string `json:"email" validate:"required,email"`
	Password     string `json:"password" validate:"required,min=5,passwordpattern"`
	ConfPassword string `json:"conf_password" validate:"required,eqfield=Password"`
}

func (h *UserHandler) CreateUser(
	ctx context.Context,
	req *user_grpc.CreateUserRequest,
) (*user_grpc.CreateUserResponse, error) {
	const op = "grpc.auth.CreateUser"

	log := h.Log.With(slog.String("op", op))

	v := validator.New()
	v.RegisterValidation("passwordpattern", validators.PasswordPatternValidator)

	if err := v.StructCtx(ctx, &CreateUserRequest{
		Email:        req.GetEmail(),
		Password:     req.GetPassword(),
		ConfPassword: req.GetConfirmPassword(),
	}); err != nil {
		validateErr, ok := err.(validator.ValidationErrors)
		if !ok {
			log.Error("error happened while validating request", sl.Err(err))
			return nil, utilsgrpc.UnexpectedError()
		}

		log.Info("invalid request")

		return nil, utilsgrpc.WithDetails(
			codes.Internal,
			"validation error",
			&errdetails.BadRequest{
				FieldViolations: utilsgrpc.ValidationError(validateErr),
			})
	}

	user, err := h.srvc.CreateUser(ctx, req.Email, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrEmailExists) {
			return nil, utilsgrpc.WithDetails(
				codes.InvalidArgument,
				"user exists",
				&errdetails.BadRequest_FieldViolation{
					Field:       "email",
					Description: "email is already being used",
				})
		}

		return nil, utilsgrpc.UnexpectedError()
	}

	return &user_grpc.CreateUserResponse{
		User: typesgrpc.NewUser(user),
	}, nil
}
