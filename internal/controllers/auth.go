package controllers

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	authv1 "weather_back_api_getway/pkg/auth_v1"
)

type Auth interface {
	Login(ctx context.Context,
		login string,
		password string,
	) (token string, err error)
	Register(ctx context.Context,
		login string,
		password string,
		email string,
	) (uuid string, err error)
}

type serverAPI struct {
	authv1.UnimplementedAuthV1Server
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	authv1.RegisterAuthV1Server(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(
	ctx context.Context,
	req *authv1.LoginRequest,
) (*authv1.LoginResponse, error) {
	if err := validateLogin(req); err != nil {
		return nil, err
	}

	token, err := s.auth.Login(ctx, req.GetLogin(), req.GetPassword())
	if err != nil {
		//TODO: обработка ошибок
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(
	ctx context.Context,
	req *authv1.RegisterRequest,
) (*authv1.RegisterResponse, error) {
	if err := validateRegister(req); err != nil {
		return nil, err
	}

	uuid, err := s.auth.Register(ctx, req.Login, req.Password, req.Email)
	if err != nil {
		//TODO: обработка ошибок
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &authv1.RegisterResponse{
		Uuid: uuid,
	}, nil
}

func validateLogin(req *authv1.LoginRequest) error {
	if req.GetLogin() == "" {
		return status.Error(codes.InvalidArgument, "login required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}

	return nil
}

func validateRegister(req *authv1.RegisterRequest) error {
	if req.GetEmail() == "" {
		return status.Error(codes.InvalidArgument, "email required")
	}

	if req.GetLogin() == "" {
		return status.Error(codes.InvalidArgument, "login required")
	}

	if req.GetPassword() == "" {
		return status.Error(codes.InvalidArgument, "password required")
	}

	return nil
}
