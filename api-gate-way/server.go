package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/joho/godotenv"
	"github.com/samObot19/shopverse/api-gate-way/authenticate"
	"github.com/samObot19/shopverse/api-gate-way/graph"
	"github.com/samObot19/shopverse/api-gate-way/product-client"
	userclient "github.com/samObot19/shopverse/api-gate-way/user-client"
	orderclient "github.com/samObot19/shopverse/api-gate-way/order-client"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Printf("No .env file found or error loading it: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	userServiceAddress := os.Getenv("USER_SERVICE_ADDRESS")
	if userServiceAddress == "" {
		log.Fatal("USER_SERVICE_ADDRESS not set in .env")
	}

	productServiceAddress := os.Getenv("PRODUCT_SERVICE_ADDRESS")
	if productServiceAddress == "" {
		log.Fatal("PRODUCT_SERVICE_ADDRESS not set in .env")
	}

	orderServiceAddress := os.Getenv("ORDER_SERVICE_ADDRESS")
	if orderServiceAddress == "" {
		log.Fatal("ORDER_SERVICE_ADDRESS not set in .env")
	}

	userConn, err := userclient.ConnectToUserService(userServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to user service: %v", err)
	}
	defer userConn.Close()
	userClient := userclient.NewUserClient(userConn)

	productConn, err := productclient.ConnectToProductService(productServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to product service: %v", err)
	}
	defer productConn.Close()
	productClient := productclient.NewProductClient(productConn)

	orderConn, err := orderclient.ConnectToOrderService(orderServiceAddress)
	if err != nil {
		log.Fatalf("Failed to connect to order service: %v", err)
	}
	defer orderConn.Close()
	orderClient := orderclient.NewOrderClient(orderConn)

	resolver := &graph.Resolver{
		ProductClient: productClient,
		UserClient:    userClient,
		OrderClient:   orderClient,
	}

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: resolver}))

	http.HandleFunc("/signup", authenticate.HandleGoogleAuth)
	http.HandleFunc("/login", authenticate.HandleGoogleAuth)
	http.HandleFunc("/auth/google/callback", authenticate.HandleGoogleCallback)
	http.HandleFunc("/refresh", authenticate.HandleRefreshToken)
	http.HandleFunc("/logout", authenticate.HandleLogout) // Add logout handler

	http.Handle("/query", authenticate.JWTMiddleware(srv))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}