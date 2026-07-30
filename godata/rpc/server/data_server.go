package server

import (
	"context"

	"github.com/phoenix-agent-go/internal/service"
	pb "github.com/phoenix-agent-go/rpc/proto"
)

type DataServer struct {
	pb.UnimplementedDataServiceServer
	dataSvc *service.DataService
}

func NewDataServer(dataSvc *service.DataService) *DataServer {
	return &DataServer{dataSvc: dataSvc}
}

func (s *DataServer) GetDatasource(ctx context.Context, req *pb.GetDatasourceRequest) (*pb.GetDatasourceResponse, error) {
	ds, err := s.dataSvc.GetDatasourceByID(ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &pb.GetDatasourceResponse{
		Id:     ds.ID,
		Name:   ds.Name,
		Type:   ds.Type,
		Status: ds.Status,
	}, nil
}

func (s *DataServer) Chat(req *pb.ChatRequest, stream pb.DataService_ChatServer) error {
	return stream.Send(&pb.ChatResponse{
		Content:   "gRPC chat not yet implemented",
		EventType: "message",
		Done:      true,
	})
}
