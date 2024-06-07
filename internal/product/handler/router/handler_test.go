package router_test

import (
	"context"
	"errors"
	"testing"

	grpcpb "github.com/namnv2496/go-coffee-shop-demo/api/grpcpb/gen"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/handler/router"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/repo"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type HandlerTestSuite struct {
	ctrl         *gomock.Controller
	mockItemRepo *repo.MockItemRepo
	mockHandler  router.Handler
}

func setup(t *testing.T) *HandlerTestSuite {
	ctrl := gomock.NewController(t)
	mockItemRepo := repo.NewMockItemRepo(ctrl)
	mockHandler := router.NewHandler(mockItemRepo)

	return &HandlerTestSuite{
		ctrl:         ctrl,
		mockItemRepo: mockItemRepo,
		mockHandler:  mockHandler,
	}
}

func teardown(suite *HandlerTestSuite) {
	suite.ctrl.Finish()
}

func TestHandlergRPC_GetProducts_FindAll(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	itemList := []domain.Item{
		{
			Id:    int32(1),
			Name:  "caffe",
			Price: int32(50),
			Type:  int32(0),
			Img:   "caffe_muoi",
		},
	}
	ctx := context.Background()
	page := int32(0)
	size := int32(50)
	suite.mockItemRepo.EXPECT().
		GetAll(
			ctx,
			page,
			size,
		).
		Return(itemList, nil)
	result, err := suite.mockHandler.GetProducts(ctx, &grpcpb.GetProductsRequest{
		Id:       int32(0),
		Name:     "",
		ItemType: grpcpb.ItemType(0),
		Page:     int32(0),
		Size:     int32(50),
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result.Items[0].Name, "caffe")
}

func TestHandlergRPC_GetProducts_ById(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	itemList := []domain.Item{
		{
			Id:    int32(1),
			Name:  "caffe",
			Price: int32(50),
			Type:  int32(0),
			Img:   "caffe_muoi",
		},
	}
	ctx := context.Background()
	page := int32(0)
	size := int32(50)
	suite.mockItemRepo.EXPECT().
		GetByIdOrName(
			ctx,
			int32(1),
			"",
			page,
			size,
		).
		Return(itemList, nil)
	result, err := suite.mockHandler.GetProducts(ctx, &grpcpb.GetProductsRequest{
		Id:       int32(1),
		Name:     "",
		ItemType: grpcpb.ItemType(0),
		Page:     int32(0),
		Size:     int32(50),
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result.Items[0].Name, "caffe")
}

func TestHandlergRPC_GetProducts_ByName(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	itemList := []domain.Item{
		{
			Id:    int32(1),
			Name:  "caffe",
			Price: int32(50),
			Type:  int32(0),
			Img:   "caffe_muoi",
		},
	}
	ctx := context.Background()
	page := int32(0)
	size := int32(50)
	suite.mockItemRepo.EXPECT().
		GetByIdOrName(
			ctx,
			int32(0),
			"caffe",
			page,
			size,
		).
		Return(itemList, nil)
	result, err := suite.mockHandler.GetProducts(ctx, &grpcpb.GetProductsRequest{
		Id:       int32(0),
		Name:     "caffe",
		ItemType: grpcpb.ItemType(0),
		Page:     int32(0),
		Size:     int32(50),
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result.Items[0].Name, "caffe")
}

func TestHandlergRPC_GetProducts_ByItemType(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	itemList := []domain.Item{
		{
			Id:    int32(1),
			Name:  "caffe",
			Price: int32(50),
			Type:  int32(0),
			Img:   "caffe_muoi",
		},
	}
	ctx := context.Background()
	page := int32(0)
	size := int32(50)
	suite.mockItemRepo.EXPECT().
		GetByIdOrNameOrType(
			ctx,
			int32(0),
			"",
			int32(1),
			page,
			size,
		).
		Return(itemList, nil)
	result, err := suite.mockHandler.GetProducts(ctx, &grpcpb.GetProductsRequest{
		Id:       int32(0),
		Name:     "",
		ItemType: grpcpb.ItemType(1),
		Page:     int32(0),
		Size:     int32(50),
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result.Items[0].Name, "caffe")
}

func TestHandlergRPC_GetProducts_ByItemTypeAndId(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	itemList := []domain.Item{
		{
			Id:    int32(1),
			Name:  "caffe",
			Price: int32(50),
			Type:  int32(0),
			Img:   "caffe_muoi",
		},
	}
	ctx := context.Background()
	page := int32(0)
	size := int32(50)
	suite.mockItemRepo.EXPECT().
		GetByIdOrNameOrType(
			ctx,
			int32(1),
			"",
			int32(1),
			page,
			size,
		).
		Return(itemList, nil)
	result, err := suite.mockHandler.GetProducts(ctx, &grpcpb.GetProductsRequest{
		Id:       int32(1),
		Name:     "",
		ItemType: grpcpb.ItemType(1),
		Page:     int32(0),
		Size:     int32(50),
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result.Items[0].Name, "caffe")
}

func TestHandlergRPC_GetProducts_ByItemTypeAndName(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	itemList := []domain.Item{
		{
			Id:    int32(1),
			Name:  "caffe",
			Price: int32(50),
			Type:  int32(0),
			Img:   "caffe_muoi",
		},
	}
	ctx := context.Background()
	page := int32(0)
	size := int32(50)
	suite.mockItemRepo.EXPECT().
		GetByIdOrNameOrType(
			ctx,
			int32(0),
			"caffe",
			int32(1),
			page,
			size,
		).
		Return(itemList, nil)
	result, err := suite.mockHandler.GetProducts(ctx, &grpcpb.GetProductsRequest{
		Id:       int32(0),
		Name:     "caffe",
		ItemType: grpcpb.ItemType(1),
		Page:     int32(0),
		Size:     int32(50),
	})
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result.Items[0].Name, "caffe")
}

func TestHandlergRPC_GetProducts_ReturnError(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	itemList := []domain.Item{
		{
			Id:    int32(1),
			Name:  "caffe",
			Price: int32(50),
			Type:  int32(0),
			Img:   "caffe_muoi",
		},
	}
	ctx := context.Background()
	page := int32(0)
	size := int32(50)
	suite.mockItemRepo.EXPECT().
		GetByIdOrNameOrType(
			ctx,
			int32(0),
			"caffe",
			int32(1),
			page,
			size,
		).
		Return(itemList, errors.New("Error"))
	result, err := suite.mockHandler.GetProducts(ctx, &grpcpb.GetProductsRequest{
		Id:       int32(0),
		Name:     "caffe",
		ItemType: grpcpb.ItemType(1),
		Page:     int32(0),
		Size:     int32(50),
	})
	require.Error(t, err)
	require.NotNil(t, result)
}
