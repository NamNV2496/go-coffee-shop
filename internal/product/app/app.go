package app

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/router"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/s3"
)

type AppInterface interface {
	Start() error
	AddNewProduct(ctx context.Context, req *gin.Context) (int32, error)
	GetImageInMinio(ctx context.Context, req *gin.Context) error
}
type App struct {
	grpcServer     router.ProductServer
	productService service.ProductService
}

func NewApp(
	server router.ProductServer,
	productService service.ProductService,
) *App {
	return &App{
		grpcServer:     server,
		productService: productService,
	}
}

func (a App) Start() error {
	go func() {
		a.grpcServer.StartServerGRPC()
	}()
	return nil
}

func (a App) AddNewProduct(ctx context.Context, req *gin.Context) (int32, error) {

	// Parse the multipart form, 10 << 20 specifies a max upload size of 10MB
	err := req.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		req.String(http.StatusBadRequest, "File too large")
		return 0, errors.New("File too large")
	}

	// Retrieve file from the request
	file, header, err := req.Request.FormFile("file")
	if err != nil {
		req.String(http.StatusBadRequest, "Failed to retrieve file")
		return 0, errors.New("Failed to retrieve file")
	}
	defer file.Close()

	// Check the file type
	fileType := header.Header.Get("Content-Type")
	if fileType != "image/png" && fileType != "image/jpeg" {
		req.String(http.StatusBadRequest, "Only PNG and JPEG files are allowed")
	}

	name := req.Request.FormValue("name")
	price, _ := strconv.Atoi(req.Request.FormValue("price"))
	foodType, _ := strconv.Atoi(req.Request.FormValue("type"))

	// Respond to the client
	req.String(http.StatusOK, fmt.Sprintf("File uploaded successfully!"))

	return a.productService.AddNewProduct(
		ctx,
		s3.BUCKETNAME,
		domain.Item{
			Name:  name,
			Price: int32(price),
			Type:  int32(foodType),
		},
		file,
		header.Size,
		header.Header.Get("Content-Type"),
	)
}

func (a App) GetImageInMinio(ctx context.Context, req *gin.Context) (string, error) {

	link, _ := a.productService.GetImageInMinio(ctx, req.Query("name"))

	req.String(http.StatusBadRequest, link)
	return "", nil
}
