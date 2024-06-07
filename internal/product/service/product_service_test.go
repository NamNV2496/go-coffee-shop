package service_test

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"strings"
	"testing"
	"time"

	"github.com/namnv2496/go-coffee-shop-demo/internal/product/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/repo"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/s3"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
)

type ProductServiceTestSuite struct {
	ctrl               *gomock.Controller
	mockItemRepo       *repo.MockItemRepo
	mockS3Client       *s3.MockS3Client
	mockProductService service.ProductService
}

func setup(t *testing.T) *ProductServiceTestSuite {
	ctrl := gomock.NewController(t)

	mockItemRepo := repo.NewMockItemRepo(ctrl)
	mockS3Client := s3.NewMockS3Client(ctrl)
	mockProductService := service.NewProductService(
		mockItemRepo,
		mockS3Client,
	)

	return &ProductServiceTestSuite{
		ctrl:               ctrl,
		mockS3Client:       mockS3Client,
		mockItemRepo:       mockItemRepo,
		mockProductService: mockProductService,
	}
}

func teardown(suite *ProductServiceTestSuite) {
	suite.ctrl.Finish()
}

func TestProductService_GetAllItems(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	const layout = "2006-01-02 15:04:05"
	str := "2024-01-06 18:36:00"
	createdDate, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println(err)
	}

	retArr := []domain.Item{
		{
			Id:          int32(1),
			Name:        "Caffe cốt dừa",
			Price:       int32(50),
			Type:        int32(0),
			Img:         "image_file_ca_phe_cot_dua",
			CreatedDate: createdDate,
		},
		{
			Id:          int32(2),
			Name:        "Caffe bạc sỉu",
			Price:       int32(50),
			Type:        int32(0),
			Img:         "image_file_ca_phe_bac_sui",
			CreatedDate: createdDate,
		},
	}

	ctx := context.Background()
	suite.mockItemRepo.EXPECT().
		GetAll(
			ctx,
			int32(0),
			int32(50),
		).
		Return(retArr, nil)

	result, err := suite.mockProductService.GetAllItems(ctx, int32(0), int32(50))
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result, retArr)
}

func TestProductService_GetItemByIdOrName(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)
	// fake := faker.New()

	const layout = "2006-01-02 15:04:05"
	str := "2024-01-06 18:36:00"
	createdDate, err := time.Parse(layout, str)
	if err != nil {
		fmt.Println(err)
	}

	retArr := []domain.Item{
		{
			Id:          int32(1),
			Name:        "Caffe cốt dừa",
			Price:       int32(50),
			Type:        int32(0),
			Img:         "image_file_ca_phe_cot_dua",
			CreatedDate: createdDate,
		},
		{
			Id:          int32(2),
			Name:        "Caffe bạc sỉu",
			Price:       int32(50),
			Type:        int32(0),
			Img:         "image_file_ca_phe_bac_sui",
			CreatedDate: createdDate,
		},
	}

	ctx := context.Background()
	suite.mockItemRepo.EXPECT().
		GetByIdOrName(
			ctx,
			int32(1),
			"test",
			int32(0),
			int32(50),
		).
		Return([]domain.Item{retArr[0]}, nil)

	result, err := suite.mockProductService.GetItemByIdOrName(ctx, 1, "test", int32(0), int32(50))
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result[0].Id, retArr[0].Id)
}

func TestProductService_AddNewProduct(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	fileName := convertString("caffe muối")
	fileName = fmt.Sprintf("image_file_%v.png", fileName)
	ctx := context.Background()
	item := domain.Item{
		Id:   int32(1),
		Name: "caffe muối",
		Img:  fileName,
	}
	buf := multipart.FileHeader{
		Filename: "image.png",
		Header:   nil,
		Size:     2,
	}
	img, _ := buf.Open()
	bucket_name := "coffee"
	contentType := "image/png"

	suite.mockItemRepo.EXPECT().
		AddNewProduct(
			ctx,
			item,
			fileName,
		).
		Return(int32(1), nil)
	suite.mockS3Client.EXPECT().
		Write(
			ctx,
			fileName,
			bucket_name,
			img,
			int64(2),
			contentType,
		)
	result, err := suite.mockProductService.AddNewProduct(
		ctx,
		bucket_name,
		item,
		img,
		int64(2),
		contentType,
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result, int32(1))
}

func TestProductService_AddNewProduct_failedSaveToDB(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	fileName := convertString("caffe muối")
	fileName = fmt.Sprintf("image_file_%v.png", fileName)
	ctx := context.Background()
	item := domain.Item{
		Id:   int32(1),
		Name: "caffe muối",
		Img:  fileName,
	}
	buf := multipart.FileHeader{
		Filename: "image.png",
		Header:   nil,
		Size:     2,
	}
	img, _ := buf.Open()
	bucket_name := "coffee"
	contentType := "image/png"

	suite.mockItemRepo.EXPECT().
		AddNewProduct(
			ctx,
			item,
			fileName,
		).
		Return(int32(1), errors.New("Fail to save to DB"))

	result, err := suite.mockProductService.AddNewProduct(
		ctx,
		bucket_name,
		item,
		img,
		int64(2),
		contentType,
	)
	// require.Error(t, err)
	require.NotNil(t, result)
	require.Equal(t, err.Error(), "Fail to save to DB")
}

func TestProductService_AddNewProduct_failedSaveToMinio(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	fileName := convertString("caffe muối")
	fileName = fmt.Sprintf("image_file_%v.png", fileName)
	ctx := context.Background()
	item := domain.Item{
		Id:   int32(1),
		Name: "caffe muối",
		Img:  fileName,
	}
	buf := multipart.FileHeader{
		Filename: "image.png",
		Header:   nil,
		Size:     2,
	}
	img, _ := buf.Open()
	bucket_name := "coffee"
	contentType := "image/png"

	suite.mockItemRepo.EXPECT().
		AddNewProduct(
			ctx,
			item,
			fileName,
		).
		Return(int32(1), nil)

	suite.mockS3Client.EXPECT().
		Write(
			ctx,
			fileName,
			bucket_name,
			img,
			int64(2),
			contentType,
		).Return(int64(0), errors.New("Fail to upload to minio"))
	result, err := suite.mockProductService.AddNewProduct(
		ctx,
		bucket_name,
		item,
		img,
		int64(2),
		contentType,
	)
	// require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, err.Error(), "Fail to upload to minio")
}

func removeDiacritics(str string) string {
	replacements := map[rune]rune{
		'đ': 'd', 'Đ': 'D',
		'á': 'a', 'à': 'a', 'ả': 'a', 'ã': 'a', 'ạ': 'a',
		'ă': 'a', 'ắ': 'a', 'ằ': 'a', 'ẳ': 'a', 'ẵ': 'a', 'ặ': 'a',
		'â': 'a', 'ấ': 'a', 'ầ': 'a', 'ẩ': 'a', 'ẫ': 'a', 'ậ': 'a',
		'é': 'e', 'è': 'e', 'ẻ': 'e', 'ẽ': 'e', 'ẹ': 'e',
		'ê': 'e', 'ế': 'e', 'ề': 'e', 'ể': 'e', 'ễ': 'e', 'ệ': 'e',
		'í': 'i', 'ì': 'i', 'ỉ': 'i', 'ĩ': 'i', 'ị': 'i',
		'ó': 'o', 'ò': 'o', 'ỏ': 'o', 'õ': 'o', 'ọ': 'o',
		'ô': 'o', 'ố': 'o', 'ồ': 'o', 'ổ': 'o', 'ỗ': 'o', 'ộ': 'o',
		'ơ': 'o', 'ớ': 'o', 'ờ': 'o', 'ở': 'o', 'ỡ': 'o', 'ợ': 'o',
		'ú': 'u', 'ù': 'u', 'ủ': 'u', 'ũ': 'u', 'ụ': 'u',
		'ư': 'u', 'ứ': 'u', 'ừ': 'u', 'ử': 'u', 'ữ': 'u', 'ự': 'u',
		'ý': 'y', 'ỳ': 'y', 'ỷ': 'y', 'ỹ': 'y', 'ỵ': 'y',
	}

	var sb strings.Builder
	for _, r := range str {
		if repl, found := replacements[r]; found {
			sb.WriteRune(repl)
		} else {
			sb.WriteRune(r)
		}
	}
	return sb.String()
}

func convertString(str string) string {
	str = removeDiacritics(str)
	str = strings.ReplaceAll(str, " ", "_")
	return strings.ToLower(str)
}

func TestProductService_GetImageInMinio(t *testing.T) {
	suite := setup(t)
	defer teardown(suite)

	ctx := context.Background()
	fileName := convertString("caffe muối")
	fileName = fmt.Sprintf("image_file_%v.png", fileName)
	bucket_name := "coffee"

	suite.mockS3Client.EXPECT().
		PreviewImage(
			ctx,
			fileName,
			bucket_name,
		).
		Return("http://minio.io", nil)
	result, err := suite.mockProductService.GetImageInMinio(
		ctx,
		fileName,
	)
	require.NoError(t, err)
	require.NotNil(t, result)
	require.Equal(t, result, "http://minio.io")
}
