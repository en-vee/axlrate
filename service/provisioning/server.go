package provisioning

import (
	context "context"

	"github.com/en-vee/alog"
	"github.com/en-vee/axlrate/core/server"
	grpc "google.golang.org/grpc"
)

type Server struct {
	// Every type of server should have it's own IMDG and it's own Communications/Discovery ports
	server.NetworkComponent
}

func (s *Server) CreateCustomer(ctx context.Context, req *CreateCustomerRequest) (*CreateCustomerResponse, error) {
	// Extract CustomerId
	alog.Debug("CreateCustomerRequest = %v", req)
	return &CreateCustomerResponse{ObjectId: 1}, nil
}

func (s *Server) RegisterServer(g *grpc.Server) {
	RegisterProvisioningServer(g, s)
}

func init() {

}
