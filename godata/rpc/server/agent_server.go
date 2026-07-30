package server

import (
	"context"

	"github.com/phoenix-agent-go/internal/service"
	pb "github.com/phoenix-agent-go/rpc/proto"
)

type AgentServer struct {
	pb.UnimplementedAgentServiceServer
	dataSvc *service.DataService
}

func NewAgentServer(dataSvc *service.DataService) *AgentServer {
	return &AgentServer{dataSvc: dataSvc}
}

func (s *AgentServer) ListAgents(ctx context.Context, req *pb.ListAgentsRequest) (*pb.ListAgentsResponse, error) {
	agents, err := s.dataSvc.ListAgent(ctx)
	if err != nil {
		return nil, err
	}

	var infos []*pb.AgentInfo
	for _, a := range agents {
		infos = append(infos, &pb.AgentInfo{
			Id:          a.ID,
			Sn:          a.Sn,
			Name:        a.Name,
			Type:        a.Type,
			Status:      a.Status,
			Description: a.Description,
		})
	}

	return &pb.ListAgentsResponse{Agents: infos}, nil
}

func (s *AgentServer) StreamCall(req *pb.AgentCallRequest, stream pb.AgentService_StreamCallServer) error {
	return stream.Send(&pb.AgentCallResponse{
		Content:   "gRPC agent call not yet implemented",
		EventType: "message",
		Done:      true,
	})
}
