package server

import "fmt"

// InvalidAddressErr is encountered when the address cannot be resolved to a valid IP/Host
type InvalidAddressErr struct {
	address string
}

func (err *InvalidAddressErr) Error() string {
	return fmt.Sprintf("address validation error : invalid network address : %s", err.address)
}

type ServiceNotRegisteredErr struct {
	service interface{}
}

func (err *ServiceNotRegisteredErr) Error() string {
	return fmt.Sprintf("grpc: service %T not registered", err.service)
}

type GrpcServerNotStartedError struct {
}

func (err *GrpcServerNotStartedError) Error() string {
	return fmt.Sprintf("grpc: server not started")
}
