package router

import (
	"context"
	"fmt"

	grpcpb "github.com/namnv2496/go-coffee-shop-demo/grpc/grpcpb/gen"
	"github.com/namnv2496/go-coffee-shop-demo/internal/counter/domain"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProductGRPCClient interface {
	GetProductByIdOrName(id int32, name string) ([]domain.Item, error)
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

func (c *productGRPCClient) GetProductByIdOrName(id int32, name string) ([]domain.Item, error) {

	client := grpcpb.NewProductServiceClient(c.conn)

	result, err := client.GetProducts(context.Background(), &grpcpb.GetProductsRequest{
		Id:   id,
		Name: name,
	})

	if err != nil {
		fmt.Println("Error: ", err)
	}
	fmt.Println("result: ", result)
	res := make([]domain.Item, 0)
	for _, item := range result.Items {
		res = append(res, domain.Item{
			Id:    item.Id,
			Name:  item.Name,
			Price: item.Price,
			Type:  item.Type,
		})
	}
	return res, nil
}
