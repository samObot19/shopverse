package graph

import (
	user_client "github.com/samObot19/shopverse/api-gate-way/user-client"
	"github.com/samObot19/shopverse/api-gate-way/product-client"
	"github.com/samObot19/shopverse/api-gate-way/order-client"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	ProductClient *productclient.ProductClient
	UserClient    *user_client.UserClient
	OrderClient   *orderclient.OrderClient // Added OrderClient
}
