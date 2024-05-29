package router

import (
	"fmt"
	"net"

	grpcpb "github.com/namnv2496/go-coffee-shop-demo/grpc/grpcpb/gen"

	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"google.golang.org/grpc"
)

type ProductServer interface {
	StartServerGRPC() error
}

type productServer struct {
	grpcpb.UnimplementedProductServiceServer
	config     configs.Config
	grpcServer *grpc.Server
	handler    grpcpb.ProductServiceServer
}

func NewGrpcRouterServer(
	config configs.Config,
	grpcServer *grpc.Server,
	handler grpcpb.ProductServiceServer,
) ProductServer {

	return &productServer{
		config:     config,
		grpcServer: grpcServer,
		handler:    handler,
	}
}

func (s *productServer) StartServerGRPC() error {

	grpcpb.RegisterProductServiceServer(s.grpcServer, s.handler)
	fmt.Printf("serve type %v, address: %v", s.config.GRPC.Type, s.config.GRPC.Address)
	fmt.Println()
	lis, _ := net.Listen(s.config.GRPC.Type, s.config.GRPC.Address)
	if err := s.grpcServer.Serve(lis); err != nil {
		panic("Failed to serve: " + err.Error())
	}
	return nil
}
