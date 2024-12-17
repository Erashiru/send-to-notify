package http

import (
	dto2 "github.com/kwaaka-team/orders-core/pkg/wolt/clients/dto"
	"time"
)

const (
	baseURL      = "https://pos-integration-service.wolt.com"
	tokenTimeout = 59 * time.Minute
)

const (
	retriesNumber   = 3
	retriesWaitTime = 1 * time.Second
)

const (
	acceptHeader      = "Accept"
	authHeader        = "Authorization"
	contentTypeHeader = "Content-Type"
	jsonType          = "application/json"
	WOLT_API_KEY      = "WOLT-API-KEY"
)

func getUniqueProducts(products dto2.UpdateProducts) dto2.UpdateProducts {
	uniqueProducts := make(map[string]dto2.UpdateProduct)

	// Iterate over products and add to the map if not already present
	for _, product := range products.Product {
		uniqueProducts[product.ExtID] = product
	}

	// Convert the map back to a slice
	uniqueProductSlice := make([]dto2.UpdateProduct, 0, len(uniqueProducts))
	for _, product := range uniqueProducts {
		uniqueProductSlice = append(uniqueProductSlice, product)
	}

	// Set the unique products back to the products.Product array
	products.Product = uniqueProductSlice
	return products

}
func getUniqueAttributes(attributes dto2.UpdateAttributes) dto2.UpdateAttributes {
	uniqueAttributes := make(map[string]dto2.UpdateAttribute)

	// Iterate over attributes and add to the map if not already present
	for _, attribute := range attributes.Attribute {
		uniqueAttributes[attribute.ExtID] = attribute
	}

	// Convert the map back to a slice
	uniqueAttributeSlice := make([]dto2.UpdateAttribute, 0, len(uniqueAttributes))
	for _, attribute := range uniqueAttributes {
		uniqueAttributeSlice = append(uniqueAttributeSlice, attribute)
	}

	// Set the unique attributes back to the attributes.Attribute array
	attributes.Attribute = uniqueAttributeSlice
	return attributes
}
