package main

import (
	"context"
	"time"

	"github.com/en-vee/alog"

	"github.com/en-vee/axlrate/service/provisioning"
	"google.golang.org/grpc"
)

func main() {
	conn, err := grpc.Dial("localhost:10080", grpc.WithInsecure())
	if err != nil {
		alog.Critical("Error dialing GRPC server : %v", err)
	}
	defer conn.Close()

	client := provisioning.NewProvisioningClient(conn)

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)

	resp, err := client.CreateCustomer(ctx, &provisioning.CreateCustomerRequest{CustomerId: "C12345"})
	if err != nil {
		alog.Critical("Error sending GRPC request : %v", err)
	}

	alog.Info("Response = %v", resp)
}