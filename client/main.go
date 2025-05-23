package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"

	pb "grpc-secure/migrator/proto"
)

func main() {
	// sckip tls layer
	// creds := credentials.NewTLS(&tls.Config{
	// 	InsecureSkipVerify: true, // for testing only
	// })

	// creds := credentials.NewTLS(&tls.Config{
	// 	InsecureSkipVerify: true, // for testing only
	// })

	// use tls layer
	certPool := x509.NewCertPool()
	cert, err := os.ReadFile("cert/ca-cert.pem")
	if err != nil {
		log.Fatalf("Failed to read CA cert: %v", err)
	}
	certPool.AppendCertsFromPEM(cert)

	creds := credentials.NewTLS(&tls.Config{
		RootCAs: certPool, // Validate the server's cert properly
	})

	conn, err := grpc.Dial("localhost:50051", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	defer conn.Close()

	client := pb.NewMigratorClient(conn)

	ctx := metadata.AppendToOutgoingContext(context.Background(), "x-api-key", "secret-api-key")

	res, err := client.MigrateData(ctx, &pb.MigrateRequest{
		Id:        "123",
		SourceUrl: "https://example.com/file",
	})
	if err != nil {
		log.Fatalf("gRPC call failed: %v", err)
	}

	log.Println("Response:", res.Status)
}
