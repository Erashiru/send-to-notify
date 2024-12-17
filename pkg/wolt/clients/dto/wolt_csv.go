package dto

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"
)

type WoltCsvProducts []WoltCsvProduct

func (products WoltCsvProducts) ToCSVBytes() ([]byte, error) {
	csvData := new(bytes.Buffer)
	writer := csv.NewWriter(csvData)

	header := []string{"category_id", "gtin", "merchant_sku", "name", "description", "price",
		"discounted_price", "vat_percentage", "pos_id", "enabled", "max_quantity_per_purchase",
		"min_quantity_per_purchase", "is_age_restricted_item", "age_limit", "item_balance", "use_limited_quantity",
		"weighted_item_sell_by", "weighted_item_min_weight_in_grams", "weighted_item_approx_weight_in_grams",
		"caffeine_content_value", "caffeine_content_units", "use_inventory_management", "is_in_stock", "alcohol_percentage"}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	// Write data
	for _, p := range products {
		discountPrice := ""
		if p.DiscountedPrice != 0 {
			discountPrice = strconv.FormatFloat(p.DiscountedPrice, 'g', -1, 64)
		}

		record := []string{p.CategoryId, "", strconv.Itoa(p.MerchantSku), p.Name, p.Description, strconv.FormatFloat(p.Price, 'g', -1, 64), discountPrice,
			strconv.FormatFloat(p.VatPercentage, 'g', -1, 64), p.PosId, p.Enabled, strconv.Itoa(p.MaxQuantityPerPurchase),
			strconv.Itoa(p.MinQuantityPerPurchase), p.IsAgeRestrictedItem, strconv.Itoa(p.AgeLimit), strconv.Itoa(p.ItemBalance), p.UseLimitedQuantity,
			p.WeightedItemSellBy, p.WeightedItemMinWeightInGrams, p.WeightedItemApproxWeightInGrams, strconv.Itoa(p.CaffeineContentValue), p.CaffeineContentUnits,
			p.UseInventoryManagement, p.IsInStock, fmt.Sprintf("%f", p.AlcoholPercentage)}
		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()

	return csvData.Bytes(), nil
}

type WoltCsvProduct struct {
	CategoryId                      string  `json:"category_id"`
	Gtin                            int64   `json:"gtin,omitempty"`
	MerchantSku                     int     `json:"merchant_sku"`
	Name                            string  `json:"name"`
	Description                     string  `json:"description"`
	Price                           float64 `json:"price"`
	DiscountedPrice                 float64 `json:"discounted_price"`
	VatPercentage                   float64 `json:"vat_percentage"`
	PosId                           string  `json:"pos_id"`
	Enabled                         string  `json:"enabled"`
	MaxQuantityPerPurchase          int     `json:"max_quantity_per_purchase"`
	MinQuantityPerPurchase          int     `json:"min_quantity_per_purchase"`
	IsAgeRestrictedItem             string  `json:"is_age_restricted_item"`
	AgeLimit                        int     `json:"age_limit"`
	ItemBalance                     int     `json:"item_balance"`
	UseLimitedQuantity              string  `json:"use_limited_quantity"`
	WeightedItemSellBy              string  `json:"weighted_item_sell_by"`
	WeightedItemMinWeightInGrams    string  `json:"weighted_item_min_weight_in_grams"`
	WeightedItemApproxWeightInGrams string  `json:"weighted_item_approx_weight_in_grams"`
	CaffeineContentValue            int     `json:"caffeine_content_value"`
	CaffeineContentUnits            string  `json:"caffeine_content_units"`
	UseInventoryManagement          string  `json:"use_inventory_management"`
	IsInStock                       string  `json:"is_in_stock"`
	ProducerInformation             string  `json:"producer_information"`
	AlcoholPercentage               float32 `json:"alcohol_percentage"`
}
