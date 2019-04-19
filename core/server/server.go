/*
Package server exposes a function for launching axlrate grpc servers via Launch(...)
*/
package server

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"
)

type NetworkAddressValidator interface {
	ValidateAddress() error
}

type Server interface {
	NetworkAddressValidator
	CreateNetworkListener() error
	RegisterServer(*grpc.Server)
	StartGrpcServer(*grpc.Server) error
	StopGrpcServer() error
}

// NetworkAddressType represents TCP/UDP address type
type NetworkAddressType uint8

const (
	// TcpAddressType represents a TCP Address Type
	TcpAddressType NetworkAddressType = iota
	UdpAddressType
)

// NetworkComponent is a type which encapsulates the definition of network connection properties/parameters for a server i.e. the host,port, etc.
// It satisfies NetworkAddressValidator
type NetworkComponent struct {
	NetworkType NetworkAddressType
	Address     string
	PortNumber  int64
	listener    net.Listener
	grpcServer  *grpc.Server
}

// Launch is the entry point to create and launch a server as a goroutine.
// Callers should arrange to receive on the errorChannel of the server for any errors being communicated by the goroutine.
func Launch(srv Server) <-chan error {

	c := make(chan error, 1)

	go func() {
		// Validate the Network Component address
		if err := srv.ValidateAddress(); err != nil {
			c <- err
		}

		// Create listener
		if err := srv.CreateNetworkListener(); err != nil {
			c <- err
		}

		// Create a GRPC server
		grpcServer := grpc.NewServer()

		// Register the GRPC Server with the axlrate Server instance which also implements the relevant service
		srv.RegisterServer(grpcServer)

		if err := srv.StartGrpcServer(grpcServer); err != nil {
			c <- err
		}
	}()
	return c
}

var networkTypeToStringMap = map[NetworkAddressType]string{
	TcpAddressType: "tcp",
	UdpAddressType: "udp",
}

func (nc *NetworkComponent) ValidateAddress() error {
	var isIPAddress = false
	// Check if the address represents an IPv4/6 address
	// If it does, then check if the address provided is assigned to any of the network interfaces.
	// 		If not, then error out
	if ip := net.ParseIP(nc.Address); ip != nil {
		isIPAddress = true
		if addrs, err := net.LookupHost("localhost"); err == nil {
			ok := false
			for _, addr := range addrs {
				if ok = (addr == nc.Address); ok {
					break
				}
			}

			if !ok {
				// Provided IP address did not match any of the network interface addresses
				return &InvalidAddressErr{nc.Address}
			}

			if isIPAddress {
				return nil
			}
		}
	}
	// If the address is not an IPv4/6 address, then check if it is a valid host name i.e. some name that can be resolved and reached via DNS lookup
	if _, err := net.ResolveIPAddr("ip", nc.Address); err != nil {
		return &InvalidAddressErr{nc.Address}
	}
	return nil
}

func (nc *NetworkComponent) CreateNetworkListener() error {

	var err error

	networkAddress := fmt.Sprintf("%s:%d", nc.Address, nc.PortNumber)
	if nc.listener, err = net.Listen(networkTypeToStringMap[nc.NetworkType], networkAddress); err != nil {
		return translateError(err)
	}
	return nil
}

func (nc *NetworkComponent) StartGrpcServer(g *grpc.Server) error {
	// At best, try to see if the RegisterServer method actually registers the service implementation with the GRPC server
	if sInfo := g.GetServiceInfo(); len(sInfo) == 0 {
		return &ServiceNotRegisteredErr{nc}
	}

	nc.grpcServer = g

	if err := g.Serve(nc.listener); err != nil {
		return err
	}

	return nil
}

func (s *NetworkComponent) StopGrpcServer() error {
	var err error
	if s.listener != nil {
		s.listener.Close()
	}

	if s.grpcServer != nil {
		s.grpcServer.Stop()
	} else {
		err = &GrpcServerNotStartedError{}
	}
	return err
}

func translateError(inputError error) error {
	var err error
	switch inputError := inputError.(type) {
	case *net.OpError:
		err = &InvalidAddressErr{inputError.Addr.String()}
	}
	return err
}

type ContextHandler struct {
	ctx context.Context
}

func (c *ContextHandler) WatchForMessages() error {
	select {
	case <-c.ctx.Done():
	default:
	}
	return nil
}
