package main

import (
	"context"
	"fmt"
	"log"
	"github.com/openconfig/gnmi/proto/gnmi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Hardcoded values for simplicity, replace with appropriate values
	targetAddr := "192.168.151.7:6030" // Address of your GNMI target
	path := "/network-instances/network-instance[name=default]/protocols/protocol[name=BGP]"  // Replace with the desired GNMI path

	// Hardcoded username and password for demonstration
	username := "admin"
	password := "admin"
  
  // Insecure connection instead of using Certs
	creds := insecure.NewCredentials()

	conn, err := grpc.Dial(targetAddr, grpc.WithTransportCredentials(creds), grpc.WithPerRPCCredentials(&authCreds{
		username: username,
		password: password,
	}))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	client := gnmi.NewGNMIClient(conn)

	// Create a GetRequest
	req := &gnmi.GetRequest{
		Path: []*gnmi.Path{
			{
				Elem: parseGNMIElements(path),
			},
		},
	}

	resp, err := client.Get(context.Background(), req)
	if err != nil {
		log.Fatalf("Failed to retrieve data: %v", err)
	}

	fmt.Printf("GNMI Response: %+v\n", resp)
}

func parseGNMIElements(path string) []*gnmi.PathElem {
	var elems []*gnmi.PathElem
	parts := splitPath(path)

	for _, part := range parts {
		elems = append(elems, &gnmi.PathElem{Name: part})
	}

	return elems
}

func splitPath(path string) []string {
	var result []string
	parts := split(path, "/")
	for _, p := range parts {
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func split(s, sep string) []string {
	return []string{}
}

// authCreds is an implementation of grpc.PerRPCCredentials
// for sending username and password
type authCreds struct {
	username string
	password string
}

func (a *authCreds) GetRequestMetadata(context.Context, ...string) (map[string]string, error) {
	return map[string]string{
		"username": a.username,
		"password": a.password,
	}, nil
}

func (a *authCreds) RequireTransportSecurity() bool {
	return false // set to true if you want to ensure the connection is secure
}
