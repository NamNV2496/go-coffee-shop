package main

// import (
// 	"context"
// 	"fmt"

// 	grpcpb "github.com/namnv2496/go-coffee-shop-demo/api/grpcpb/gen"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// 	"google.golang.org/grpc/grpclog"
// )

// func main() {
// 	conn, _ := grpc.Dial("localhost:5600", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	c := grpcpb.NewProductServiceClient(conn)

// 	if err := triggerFunction(c); err != nil {
// 		grpclog.Fatal(err)
// 	}
// }

// func triggerFunction(c grpcpb.ProductServiceClient) error {
// 	result, err := c.GetProducts(context.Background(), &grpcpb.GetProductsRequest{
// 		Id:   1,
// 		Name: "thịt",
// 	})

// 	if err != nil {
// 		fmt.Println("Error: ", err)
// 	}
// 	fmt.Println("result: ", result)
// 	return nil
// }

// Client with mux but it wasn't work. It can't forward request body

// package main

// import (
// 	"context"
// 	"flag"
// 	"fmt"
// 	"log"
// 	"net/http"

// 	grpcpb "github.com/namnv2496/go-coffee-shop-demo/api/grpcpb/gen"

// 	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// )

// var grpcServerEndpoint = flag.String("grpc-server-endpoint", "localhost:5600", "gRPC server endpoint")

// func run() error {
// 	ctx := context.Background()
// 	ctx, cancel := context.WithCancel(ctx)
// 	defer cancel()

// 	// Register gRPC server endpoint
// 	// Note: Make sure the gRPC server is running properly and accessible
// 	mux := runtime.NewServeMux()
// 	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

// 	// Log the gRPC endpoint being used
// 	log.Printf("Registering gRPC service handler from endpoint: %s", grpcServerEndpoint)

// 	err := grpcpb.RegisterProductServiceHandlerFromEndpoint(ctx, mux, *grpcServerEndpoint, opts)
// 	if err != nil {
// 		log.Fatalf("Failed to register gRPC service handler: %v", err)
// 		return err
// 	}

// 	fmt.Println("Listening on port 8081 ...")
// 	// Start HTTP server (and proxy calls to gRPC server endpoint)
// 	err = http.ListenAndServe(":8081", mux)
// 	if err != nil {
// 		log.Fatalf("Failed to start HTTP server: %v", err)
// 		return err
// 	}
// 	return nil
// }

// func main() {
// 	if err := run(); err != nil {
// 		log.Fatalf("Failed to run the server: %v", err)
// 	}
// }
