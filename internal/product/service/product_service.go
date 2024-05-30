package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"strings"

	"github.com/namnv2496/go-coffee-shop-demo/internal/product/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/repo"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/s3"
)

type ProductService interface {
	GetAllItems(ctx context.Context, offset int32, limit int32) ([]domain.Item, error)
	GetItemByIdOrName(ctx context.Context, id int32, name string, offset int32, limit int32) ([]domain.Item, error)
	AddNewProduct(ctx context.Context, bucket string, item domain.Item, img multipart.File, size int64, contentType string) (int32, error)
	GetImageInMinio(ctx context.Context, name string) (string, error)
}

type productService struct {
	itemRepo repo.ItemRepo
	s3client s3.S3Client
}

func NewProductService(
	itemRepo repo.ItemRepo,
	s3client s3.S3Client,
) ProductService {

	return &productService{
		itemRepo: itemRepo,
		s3client: s3client,
	}
}

func (s productService) GetAllItems(
	ctx context.Context,
	offset int32,
	limit int32,
) ([]domain.Item, error) {

	return s.itemRepo.GetAll(ctx, offset, limit)
}

func (s productService) GetItemByIdOrName(
	ctx context.Context,
	id int32,
	name string,
	offset int32,
	limit int32,
) ([]domain.Item, error) {
	return s.itemRepo.GetByIdOrName(ctx, id, name, offset, limit)
}

func (s productService) AddNewProduct(
	ctx context.Context,
	bucket string,
	item domain.Item,
	img multipart.File,
	size int64,
	contentType string,
) (int32, error) {

	fileName := s.convertString(item.Name)
	fileName = fmt.Sprintf("image_file_%v.png", fileName)

	// save to DB
	id, err := s.itemRepo.AddNewProduct(ctx, item, fileName)
	if err != nil {
		panic("Fail to save to DB")
	}
	// save to MinIO
	_, err = s.s3client.Write(
		ctx,
		fileName,
		bucket,
		img,
		size,
		contentType,
	)
	if err != nil {
		panic("Fail to upload to minio")
	}
	return id, nil
}

func (s productService) removeDiacritics(str string) string {
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

func (s productService) convertString(str string) string {
	// Step 1: Remove diacritics
	str = s.removeDiacritics(str)

	// Step 2: Replace spaces with underscores
	str = strings.ReplaceAll(str, " ", "_")

	return str
}

func (s productService) GetImageInMinio(ctx context.Context, name string) (string, error) {

	return s.s3client.PreviewImage(ctx, name, s3.BUCKETNAME)
}
