package router

import (
	"context"
	"fmt"

	grpcpb "github.com/namnv2496/go-coffee-shop-demo/api/grpcpb/gen"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductGRPCClient interface {
	GetProductByIdOrNameOrType(id int32, name string, itemType int32, page int32, size int32) ([]*grpcpb.Item, error)
}

type productGRPCClient struct {
	conn *grpc.ClientConn
}

func NewGRPCProductClient(
	config configs.Config,
) (ProductGRPCClient, error) {
	conn, err := grpc.NewClient(config.GRPC.Address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &productGRPCClient{
		conn: conn,
	}, nil
}

func (c *productGRPCClient) GetProductByIdOrNameOrType(
	id int32,
	name string,
	itemType int32,
	page int32,
	size int32,
) ([]*grpcpb.Item, error) {

	client := grpcpb.NewProductServiceClient(c.conn)

	result, err := client.GetProducts(context.Background(), &grpcpb.GetProductsRequest{
		Id:       id,
		Name:     name,
		ItemType: grpcpb.ItemType(itemType),
		Page:     0,
		Size:     50,
	})

	if err != nil {
		fmt.Println("Error: ", err)
	}
	// fmt.Println("result: ", result)
	return result.Items, nil
}
