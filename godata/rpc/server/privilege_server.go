package server

import (
	"context"
	"strconv"

	"github.com/phoenix-agent-go/internal/model"
	"github.com/phoenix-agent-go/internal/service"
	pb "github.com/phoenix-agent-go/rpc/proto"
)

type PrivilegeServer struct {
	pb.UnimplementedPrivilegeServiceServer
	svc *service.PrivilegeService
}

func NewPrivilegeServer(svc *service.PrivilegeService) *PrivilegeServer {
	return &PrivilegeServer{svc: svc}
}

func (s *PrivilegeServer) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	dto := model.LoginInfoDTO{
		Username: req.Username,
		Password: req.Password,
		Type:     req.LoginType,
	}

	userInfo, err := s.svc.Login(ctx, dto, req.Ip)
	if err != nil {
		return nil, err
	}

	return &pb.LoginResponse{
		UserId:   userInfo.UserID,
		Username: userInfo.Username,
		Roles:    userInfo.Roles,
	}, nil
}

func (s *PrivilegeServer) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	user, err := s.svc.GetUserByID(ctx, req.UserId)
	if err != nil {
		return nil, err
	}

	return &pb.GetUserResponse{
		UserId:   user.ID,
		Username: user.Username,
		RealName: user.RealName,
		Status:   strconv.Itoa(user.Status),
	}, nil
}
