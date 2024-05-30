package router

import (
	"context"
	"fmt"

	grpcpb "github.com/namnv2496/go-coffee-shop-demo/grpc/grpcpb/gen"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductGRPCClient interface {
	GetProductByIdOrName(id int32, name string, page int32, size int32) ([]*grpcpb.Item, error)
}

type productGRPCClient struct {
	conn *grpc.ClientConn
}

func NewGRPCProductClient(
	config configs.Config,
) (ProductGRPCClient, error) {
	conn, err := grpc.Dial(config.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &productGRPCClient{
		conn: conn,
	}, nil
}

func (c *productGRPCClient) GetProductByIdOrName(
	id int32,
	name string,
	page int32,
	size int32,
) ([]*grpcpb.Item, error) {

	client := grpcpb.NewProductServiceClient(c.conn)

	result, err := client.GetProducts(context.Background(), &grpcpb.GetProductsRequest{
		Id:   id,
		Name: name,
		Page: 0,
		Size: 50,
	})

	if err != nil {
		fmt.Println("Error: ", err)
	}
	// fmt.Println("result: ", result)
	return result.Items, nil
}
