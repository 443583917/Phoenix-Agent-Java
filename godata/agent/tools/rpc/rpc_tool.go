package rpc

import (
	"context"
	"encoding/json"
	"fmt"

	pb "github.com/phoenix-agent-go/rpc/proto"
	"github.com/phoenix-agent-go/rpc/client"
	"trpc.group/trpc-go/trpc-agent-go/tool"
)

type RPCTool struct {
	clients *client.Clients
}

func NewRPCTool(clients *client.Clients) *RPCTool {
	return &RPCTool{clients: clients}
}

func (t *RPCTool) Declaration() *tool.Declaration {
	return &tool.Declaration{
		Name:        "rpc_call",
		Description: "Call a remote service via gRPC. Supported methods: list_agents",
		InputSchema: &tool.Schema{
			Type:     "object",
			Required: []string{"method"},
			Properties: map[string]*tool.Schema{
				"method": {
					Type:        "string",
					Description: "RPC method to call (list_agents)",
					Enum:        []any{"list_agents"},
				},
			},
		},
	}
}

type rpcParams struct {
	Method string `json:"method"`
}

func (t *RPCTool) Call(ctx context.Context, args any) (any, error) {
	argsJSON, err := json.Marshal(args)
	if err != nil {
		return nil, fmt.Errorf("marshal args: %w", err)
	}

	var params rpcParams
	if err := json.Unmarshal(argsJSON, &params); err != nil {
		return nil, fmt.Errorf("unmarshal args: %w", err)
	}

	switch params.Method {
	case "list_agents":
		return t.listAgents(ctx)
	default:
		return nil, fmt.Errorf("unknown method: %s", params.Method)
	}
}

func (t *RPCTool) listAgents(ctx context.Context) (any, error) {
	if t.clients == nil {
		return nil, fmt.Errorf("RPC clients not initialized")
	}
	agents, err := t.clients.ListAgents(ctx)
	if err != nil {
		return nil, err
	}
	var results []map[string]string
	for _, a := range agents {
		results = append(results, map[string]string{
			"id":          a.Id,
			"name":        a.Name,
			"type":        a.Type,
			"status":      a.Status,
			"description": a.Description,
		})
	}
	return results, nil
}

func listAgentsProto(ctx context.Context, svc pb.AgentServiceClient) ([]*pb.AgentInfo, error) {
	resp, err := svc.ListAgents(ctx, &pb.ListAgentsRequest{})
	if err != nil {
		return nil, err
	}
	return resp.Agents, nil
}
