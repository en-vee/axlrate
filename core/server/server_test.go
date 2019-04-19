package server_test

import (
	"testing"

	"github.com/en-vee/axlrate/core/server"

	"google.golang.org/grpc"
)

type testServer struct {
	server.NetworkComponent
}

/*
func (s *testServer) TestRpc(ctx context.Context, in *TestRequest) (*TestResponse, error) {
	return &TestResponse{Msg: "Hello, Client"}, nil
}

func (s *testServer) RegisterServer(gs *grpc.Server) {
	RegisterTestServiceServer(gs, s)
}
*/

/////////////////////////////////////////
// AddressValidator Tests
/////////////////////////////////////////
func TestInvalidTcpIpAddress(t *testing.T) {
	var err error
	srv := testServer{NetworkComponent: server.NetworkComponent{NetworkType: server.TcpAddressType, Address: "127.5.0.1", PortNumber: 2345}}
	if err = srv.ValidateAddress(); err == nil {
		t.Errorf("Test Case Failed. Want : %v, Got : %v", &server.InvalidAddressErr{}, err)
		//t.Fatalf("Test Case Failed. Want : %v, Got : %v", &InvalidAddressErr{}, err)
	}
	// Validate if error is of type InvalidIpAddressErr
	//t.Logf("Error Value = %v", err)
	if _, ok := err.(*server.InvalidAddressErr); !ok {
		t.Errorf("Test Case Failed. Want : %v, Got : %v", &server.InvalidAddressErr{}, err)
	}
}

func TestInvalidTcpHostAddress(t *testing.T) {
	var err error
	srv := testServer{NetworkComponent: server.NetworkComponent{NetworkType: server.TcpAddressType, Address: "localhost1", PortNumber: 2345}}
	if err = srv.ValidateAddress(); err == nil {
		t.Errorf("Test Case Failed")
	}
	// Validate if error is of type InvalidIpAddressErr
	//t.Logf("Error Value = %v", err)
	if _, ok := err.(*server.InvalidAddressErr); !ok {
		t.Errorf("Test Case Failed")
	}
}

func TestValidTcpHostAddress(t *testing.T) {
	var err error
	srv := testServer{NetworkComponent: server.NetworkComponent{NetworkType: server.TcpAddressType, Address: "localhost", PortNumber: 2345}}

	if err = srv.ValidateAddress(); err != nil {
		t.Errorf("Test Case Failed")
	}
}

func TestCreateNetworkListener(t *testing.T) {

	srvComp := server.NetworkComponent{NetworkType: server.TcpAddressType, Address: "127.0.0.1", PortNumber: 2345}
	if err := srvComp.CreateNetworkListener(); err != nil {
		t.Errorf("Error Creating Network Listener. Error : %v", err)
	}
}

///////////////////////////////////////////////////////////
// StartGrpcServerTests
///////////////////////////////////////////////////////////
func TestStartUnregisteredServerThrowsError(t *testing.T) {
	var err error
	srv := testServer{NetworkComponent: server.NetworkComponent{NetworkType: server.TcpAddressType, Address: "localhost", PortNumber: 2345}}
	if err = srv.StartGrpcServer(grpc.NewServer()); err == nil {
		t.Errorf("test failed. want : %v , got : %v", server.ServiceNotRegisteredErr{}, err)
	}
}
