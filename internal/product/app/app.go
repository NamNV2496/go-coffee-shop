package app

import (
	"context"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/domain"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/router"
	"github.com/namnv2496/go-coffee-shop-demo/internal/product/service"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/s3"
	"github.com/namnv2496/go-coffee-shop-demo/pkg/utils"
)

type AppInterface interface {
	Start() error
	AddNewProduct(ctx context.Context, req *gin.Context)
	GetImageInMinio(ctx context.Context, req *gin.Context)
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

func (a App) AddNewProduct(ctx context.Context, req *gin.Context) {

	// Parse the multipart form, 10 << 20 specifies a max upload size of 10MB
	err := req.Request.ParseMultipartForm(10 << 20)
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, "File too large")
		return
	}

	// Retrieve file from the request
	file, header, err := req.Request.FormFile("file")
	if err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	}
	defer file.Close()

	// Check the file type
	fileType := header.Header.Get("Content-Type")
	if fileType != "image/png" && fileType != "image/jpeg" {
		utils.WrapperResponse(req, http.StatusBadRequest, "Only PNG and JPEG files are allowed")
		return
	}

	name := req.Request.FormValue("name")
	price, _ := strconv.Atoi(req.Request.FormValue("price"))
	foodType, _ := strconv.Atoi(req.Request.FormValue("type"))

	if id, err := a.productService.AddNewProduct(
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
	); err != nil {
		utils.WrapperResponse(req, http.StatusBadRequest, err.Error())
		return
	} else {
		utils.WrapperResponse(req, http.StatusOK, id)
	}
}

func (a App) GetImageInMinio(ctx context.Context, req *gin.Context) {

	link, _ := a.productService.GetImageInMinio(ctx, req.Query("name"))

	utils.WrapperResponse(req, http.StatusOK, link)
}
