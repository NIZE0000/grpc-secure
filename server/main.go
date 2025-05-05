package main

import (
	"context"
	"crypto/tls"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	pb "grpc-secure/migrator/proto"
)

const apiKey = "secret-api-key"

type server struct {
	pb.UnimplementedMigratorServer
}

func (s *server) MigrateData(ctx context.Context, req *pb.MigrateRequest) (*pb.MigrateResponse, error) {
	log.Printf("Received: ID=%s URL=%s", req.Id, req.SourceUrl)
	return &pb.MigrateResponse{Status: "ok"}, nil
}

// API key interceptor
func authInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["x-api-key"]) == 0 || md["x-api-key"][0] != apiKey {
		return nil, status.Errorf(codes.Unauthenticated, "invalid API key")
	}
	return handler(ctx, req)
}

func main() {
	cert, err := tls.LoadX509KeyPair("cert/server.crt", "cert/server.key")
	if err != nil {
		log.Fatalf("failed to load certs: %v", err)
	}

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer(
		grpc.Creds(credentials.NewServerTLSFromCert(&cert)),
		grpc.UnaryInterceptor(authInterceptor),
	)

	pb.RegisterMigratorServer(s, &server{})
	log.Println("gRPC server running on port 50051 (TLS)")
	s.Serve(lis)
}
