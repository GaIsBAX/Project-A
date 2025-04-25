package auth

import (
	"context"

	sso "github.com/GaIsBax/protos/gen/go/Project_A"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Auth - интерфейс для аутентификации и авторизации
type Auth interface {
	// Login - аутентификация
	Login(
		ctx context.Context,
		name string,
		password string,
		appID int,
	) (token string, err error)

	// RegisterNewUser - регистрация нового пользователя
	RegisterNewUser(
		ctx context.Context,
		name string,
		email string,
		password string,
	) (userID int64, err error)

	// IsAdmin - проверка, является ли пользователь администратором
	IsAdmin(ctx context.Context, userID int64) (bool, error)
}

type serverAPI struct {
	sso.UnimplementedAuthServer
	auth Auth
}

// Register регистрирует grpc хендлеры для аутентификации
func Register(gRPC *grpc.Server, auth Auth) {
	sso.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

const (
	emptyValue = 0
)

// Login - аутентификация пользователя
func (s *serverAPI) Login(ctx context.Context, req *sso.LoginRequest) (*sso.LoginResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.AppId == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "app_id is required")
	}

	token, err := s.auth.Login(ctx, req.GetName(), req.GetPassword(), int(req.AppId))
	if err != nil {
		return nil, status.Error(codes.Internal, "internal server")
	}

	return &sso.LoginResponse{
		Token: token,
	}, nil
}

// Register - регистрация нового пользователя
func (s *serverAPI) Register(ctx context.Context, req *sso.RegisterRequest) (*sso.RegisterResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "name is required")
	}

	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetName(), req.GetEmail(), req.Password)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to register user")
	}

	return &sso.RegisterResponse{UserId: userID}, nil
}

// IsAdmin - проверка, является ли пользователь администратором
func (s *serverAPI) IsAdmin(ctx context.Context, req *sso.IsAdminRequest) (*sso.IsAdminResponse, error) {
	if req.UserId == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to check admin status")
	}

	return &sso.IsAdminResponse{IsAdmin: isAdmin}, nil

}
