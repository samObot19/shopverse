package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/samObot19/shopverse/api-gate-way/graph"
	//"github.com/samObot19/shopverse/api-gate-way/graph/generated"
	"github.com/samObot19/shopverse/api-gate-way/product-client"
	"github.com/samObot19/shopverse/api-gate-way/user-client"
	orderclient "github.com/samObot19/shopverse/api-gate-way/order-client"
	"google.golang.org/grpc"
)

const (
	defaultPort          = "8080"
	productServiceAddress = "localhost:50055" // Replace with actual product service address
	orderServiceAddress   = "localhost:50051" // Replace with actual order service address
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	// Initialize gRPC connection
	conn, err := grpc.Dial("localhost:50500", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to gRPC server: %v", err)
	}
	defer conn.Close()

	// Initialize UserClient
	userClient := user_client.NewUserClient(conn)

	// Connect to the product service
	productConn, err := productclient.ConnectToProductService(productServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productConn.Close()

	// Initialize the ProductClient
	productClient := productclient.NewProductClient(productConn)

	// Connect to the order service
	orderConn, err := orderclient.ConnectToOrderService(orderServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer orderConn.Close()

	// Initialize the OrderClient
	orderClient := orderclient.NewOrderClient(orderConn)

	// Initialize the Resolver with the OrderClient
	resolver := &graph.Resolver{
		ProductClient: productClient,
		UserClient:    userClient,
		OrderClient:   orderClient, // Added OrderClient
	}

	// Create the GraphQL server
	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	// Set up HTTP handlers
	http.Handle("/query", srv)
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
