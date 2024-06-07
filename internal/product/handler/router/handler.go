package router

import (
	"context"
	"fmt"

	grpcpb "github.com/namnv2496/go-coffee-shop-demo/api/grpcpb/gen"

	"github.com/namnv2496/go-coffee-shop-demo/internal/product/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/repo"
)

type Handler interface {
	GetProducts(ctx context.Context, request *grpcpb.GetProductsRequest) (*grpcpb.GetProductsResponse, error)
}
type handler struct {
	grpcpb.UnimplementedProductServiceServer
	itemRepo repo.ItemRepo
}

func NewHandler(
	itemRepo repo.ItemRepo,
) grpcpb.ProductServiceServer {

	return &handler{
		itemRepo: itemRepo,
	}
}

func (s handler) GetProducts(
	ctx context.Context,
	request *grpcpb.GetProductsRequest,
) (*grpcpb.GetProductsResponse, error) {

	page := request.Page
	size := request.Size
	itemType := request.ItemType

	var itemList []domain.Item
	var err error
	if itemType == 0 && request.Id == 0 && request.Name == "" {
		itemList, err = s.itemRepo.GetAll(context.Background(), page, size)
	} else if itemType == 0 && (request.Id != 0 || request.Name != "") {
		itemList, err = s.itemRepo.GetByIdOrName(context.Background(), request.Id, request.Name, page, size)
	} else {
		itemList, err = s.itemRepo.GetByIdOrNameOrType(context.Background(), request.Id, request.Name, int32(itemType), page, size)
	}
	if err != nil {
		fmt.Println("Error when get items: " + string(err.Error()))
		return &grpcpb.GetProductsResponse{}, err
	}

	items := make([]*grpcpb.Item, 0)
	for _, item := range itemList {
		items = append(items, &grpcpb.Item{
			Id:    item.Id,
			Name:  item.Name,
			Price: item.Price,
			Type:  grpcpb.ItemType(item.Type),
			Image: item.Img,
		})
	}
	return &grpcpb.GetProductsResponse{
		Items: items,
	}, nil
}
