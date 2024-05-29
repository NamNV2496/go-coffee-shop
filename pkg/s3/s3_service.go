package s3

import (
	"context"
	"fmt"
	"mime/multipart"
	"time"

	"github.com/minio/minio-go"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/configs"
)

type S3Client interface {
	Write(ctx context.Context, fileName string, bucket string, img multipart.File, size int64, contentType string) (int64, error)
	PreviewImage(ctx context.Context, filename string, bucketName string) (string, error)
}

type s3Client struct {
	minioClient *minio.Client
	bucket      string
}

func NewS3Client(
	s3Config configs.S3,
) S3Client {
	minioClient, err := minio.New(s3Config.Address, s3Config.Username, s3Config.Password, false)
	if err != nil {
		fmt.Println("failed to create minio client")
		return nil
	}

	return &s3Client{
		minioClient: minioClient,
		bucket:      BUCKETNAME,
	}
}

func (s s3Client) Write(
	ctx context.Context,
	fileName string,
	bucket string,
	img multipart.File,
	size int64,
	contentType string,
) (int64, error) {
	// Upload the file to MinIO
	return s.minioClient.PutObject(
		bucket,
		fileName,
		img,
		size,
		minio.PutObjectOptions{ContentType: contentType})
}

// func (s s3Client) Read(
// 	ctx context.Context,
// 	filePath string,
// ) (io.ReadCloser, error) {

// 	object, err := s.minioClient.GetObjectWithContext(ctx, s.bucket, filePath, minio.GetObjectOptions{})
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, "failed to get s3 object")
// 	}

// 	return object, nil
// }

// func (s s3Client) MiniClient(ctx context.Context) *minio.Client {
// 	return s.minioClient
// }

func (s s3Client) PreviewImage(
	ctx context.Context,
	filename string,
	bucketName string,
) (string, error) {

	// Get the file URL
	presignedURL, err := s.minioClient.PresignedGetObject(
		bucketName,
		filename,
		10000*time.Second,
		nil)
	if err != nil {
		return "", nil
	}
	// Redirect to the presigned URL for previewing
	return presignedURL.String(), nil
}
