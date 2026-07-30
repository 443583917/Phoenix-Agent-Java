package client

import (
	"context"

	pb "github.com/phoenix-agent-go/rpc/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Clients struct {
	Agent     pb.AgentServiceClient
	Privilege pb.PrivilegeServiceClient
	Data      pb.DataServiceClient
	conn      *grpc.ClientConn
}

func NewClients(target string) (*Clients, error) {
	conn, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &Clients{
		Agent:     pb.NewAgentServiceClient(conn),
		Privilege: pb.NewPrivilegeServiceClient(conn),
		Data:      pb.NewDataServiceClient(conn),
		conn:      conn,
	}, nil
}

func (c *Clients) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

func (c *Clients) ListAgents(ctx context.Context) ([]*pb.AgentInfo, error) {
	resp, err := c.Agent.ListAgents(ctx, &pb.ListAgentsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Agents, nil
}
