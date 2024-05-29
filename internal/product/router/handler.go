package router

import (
	"context"
	"fmt"

	grpcpb "github.com/namnv2496/go-coffee-shop-demo/grpc/grpcpb/gen"

	"github.com/namnv2496/go-coffee-shop-demo/internal/product/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/service"
)

type Handler interface {
	GetProductByIdOrName(ctx context.Context, request *grpcpb.GetProductsRequest) (*grpcpb.GetProductsResponse, error)
	GetProducts(ctx context.Context, request *grpcpb.GetProductsRequest) (*grpcpb.GetProductsResponse, error)
}
type handler struct {
	grpcpb.UnimplementedProductServiceServer
	itemService service.ProductService
}

func NewHandler(
	itemService service.ProductService,
) grpcpb.ProductServiceServer {

	return &handler{
		itemService: itemService,
	}
}

func (s handler) GetProductByIdOrName(
	ctx context.Context,
	request *grpcpb.GetProductsRequest,
) (*grpcpb.GetProductsResponse, error) {

	var page int32
	if request.Page == 0 {
		page = 0
	} else {
		page = request.Page
	}
	var size int32
	if request.Size == 0 {
		size = 50
	} else {
		size = request.Size
	}
	itemList, err := s.itemService.GetItemByIdOrName(context.Background(), request.Id, request.Name, page, size)
	if err != nil {
		panic("Error when get items: " + string(err.Error()))
	}
	fmt.Println(itemList)
	items := make([]*grpcpb.Item, 0)
	for _, item := range itemList {
		items = append(items, &grpcpb.Item{
			Id:    item.Id,
			Name:  item.Name,
			Price: item.Price,
			Type:  item.Type,
		})
	}
	return &grpcpb.GetProductsResponse{
		Items: items,
	}, nil
}

func (s handler) GetProducts(
	ctx context.Context,
	request *grpcpb.GetProductsRequest,
) (*grpcpb.GetProductsResponse, error) {

	var page int32
	if request.Page == 0 {
		page = 0
	} else {
		page = request.Page
	}
	var size int32
	if request.Size == 0 {
		size = 50
	} else {
		size = request.Size
	}

	var itemList []domain.Item
	var err error
	if request.Id != 0 || request.Name != "" {
		itemList, err = s.itemService.GetItemByIdOrName(context.Background(), request.Id, request.Name, page, size)
	} else {
		itemList, err = s.itemService.GetAllItems(context.Background(), page, size)
	}
	if err != nil {
		panic("Error when get items: " + string(err.Error()))
	}
	fmt.Println(itemList)
	items := make([]*grpcpb.Item, 0)
	for _, item := range itemList {
		items = append(items, &grpcpb.Item{
			Id:    item.Id,
			Name:  item.Name,
			Price: item.Price,
			Type:  item.Type,
		})
	}
	return &grpcpb.GetProductsResponse{
		Items: items,
	}, nil
}
