package utils

import (
	"fmt"
	"github.com/kwaaka-team/orders-core/core/externalapi/models"
	coreOrderModels "github.com/kwaaka-team/orders-core/core/models"
	"github.com/rs/zerolog/log"
	"regexp"
	"time"
)

func ParseOrders(req []coreOrderModels.Order) []models.Order {
	orders := make([]models.Order, 0, len(req))

	for _, order := range req {
		orders = append(orders, ParseOrder(order))
	}

	return orders
}

func ParseOrder(req coreOrderModels.Order) models.Order {
	var zone, _ = time.LoadLocation("Asia/Almaty") // TODO From Config

	order := models.Order{
		Discriminator: req.Discriminator,
		EatsId:        req.OrderID,
		RestaurantId:  req.StoreID,
		DeliveryInfo: models.DeliveryInfo{
			ClientName:   req.Customer.Name,
			PhoneNumber:  req.Customer.PhoneNumber,
			DeliveryDate: req.OrderTime.Value.In(zone).Format(timeFormat),
		},
		PaymentInfo: models.PaymentInfo{
			ItemsCost: int(req.EstimatedTotalPrice.Value),
		},
		Persons: req.Persons,
		Comment: req.SpecialRequirements,
		Promos:  []models.OrderPromo{},
	}

	for _, promo := range req.Promos {
		order.Promos = append(order.Promos, models.OrderPromo{
			Type:     promo.Type,
			Discount: promo.Discount,
		})
	}

	switch req.PaymentMethod {
	case "DELAYED":
		order.PaymentInfo.PaymentType = "CARD"
	default:
		order.PaymentInfo.PaymentType = "CASH"
	}

	for _, product := range req.Products {
		productWithoutDiscount := int(product.PriceWithoutDiscount.Value)
		item := models.OrderItem{
			Id:                   product.ID,
			Name:                 product.Name,
			Quantity:             float64(product.Quantity),
			Price:                int(product.Price.Value),
			PriceWithoutDiscount: &productWithoutDiscount,
			Modifications:        []models.OrderModification{},
			Promos:               []models.OrderPromo{},
		}
		if req.DeliveryService == "yandex" {
			item.PriceWithoutDiscount = nil
		}

		for _, promo := range product.Promos {
			item.Promos = append(item.Promos, models.OrderPromo{
				Type:     promo.Type,
				Discount: promo.Discount,
			})
		}

		for _, attribute := range product.Attributes {
			item.Modifications = append(item.Modifications, models.OrderModification{
				Id:       attribute.ID,
				Name:     attribute.Name,
				Quantity: attribute.Quantity,
				Price:    int(attribute.Price.Value),
			})
		}

		order.Items = append(order.Items, item)
	}

	return order
}

func ParseOrderStatus(req coreOrderModels.Order) models.OrderStatusResponse {
	var zone, _ = time.LoadLocation("Asia/Almaty") // TODO From Config

	status := models.OrderStatusResponse{
		Comment:   req.SpecialRequirements,
		UpdatedAt: req.UpdatedAt.Time.In(zone).Format(timeFormat),
	}

	if req.PosType == "paloma" && req.Status == "ON_WAY" {
		req.Status = "COOKING_STARTED"
	}

	switch req.Status {
	case "NEW":
		status.Status = "NEW"
	case "ACCEPTED":
		status.Status = "ACCEPTED_BY_RESTAURANT"
	case "COOKING_STARTED":
		status.Status = "COOKING"
	case "READY_FOR_PICKUP", "COOKING_COMPLETE", "ON_WAY":
		status.Status = "READY"
	case "OUT_FOR_DELIVERY", "CLOSED", "PICKED_UP_BY_CUSTOMER":
		status.Status = "TAKEN_BY_COURIER"
	//case "CANCELLED", "CANCELLED_BY_POS_SYSTEM":
	//	status.Status = "CANCELLED"
	//	status.Comment = req.CancelReason.Reason
	//case "SKIPPED":
	//	status.Status = "CANCELLED"
	//	status.Comment = "Интеграция отключена"
	case "FAILED":
		status.Status = "ACCEPTED_BY_RESTAURANT"

		log.Printf("time range: %v", time.Now().Minute()-req.OrderTime.Value.Minute())
		if time.Now().Minute()-req.OrderTime.Value.Minute() > 15 {
			status.Status = "READY"
		}
	}

	return status
}

// YandexAddressInfo содержит информацию о частях адреса, полученного от Яндекса
type YandexAddressInfo struct {
	City        string
	Street      string
	HouseNumber string
	Entrance    string
	Intercom    string
	Floor       string
	Office      string
}

// Функция для извлечения данных из строки адреса
func ExtractYandexAddressInfo(address string) (YandexAddressInfo, error) {
	var info YandexAddressInfo

	// Регулярные выражения для поиска отдельных частей адреса
	cityRe := regexp.MustCompile(`^(.*?),\s*`)
	streetRe := regexp.MustCompile(`(?i)(?:^|,\s*)([^,]*?(?:улица|проспект|переулок|бульвар|площадь)[^,]*)`)
	houseRe := regexp.MustCompile(`,\s*(\d+)`)
	entranceRe := regexp.MustCompile(`под\.?\s*(\d+)`)
	intercomRe := regexp.MustCompile(`домофон\s*(\d+)`)
	floorRe := regexp.MustCompile(`этаж\s*(\d+)`)
	officeRe := regexp.MustCompile(`(?:кв\.|офис)\s*(\d+)`)

	// Извлечение города
	cityMatch := cityRe.FindStringSubmatch(address)
	if len(cityMatch) > 1 {
		info.City = cityMatch[1]
	}

	// Извлечение названия улицы вместе с типом
	streetMatch := streetRe.FindStringSubmatch(address)
	if len(streetMatch) > 1 {
		// Убираем возможные ведущие и завершающие пробелы
		info.Street = regexp.MustCompile(`^\s*|\s*$`).ReplaceAllString(streetMatch[1], "")
	}

	// Извлечение номера дома
	houseMatch := houseRe.FindStringSubmatch(address)
	if len(houseMatch) > 1 {
		info.HouseNumber = houseMatch[1]
	}

	// Извлечение номера подъезда
	entranceMatch := entranceRe.FindStringSubmatch(address)
	if len(entranceMatch) > 1 {
		info.Entrance = entranceMatch[1]
	}

	// Извлечение номера домофона
	intercomMatch := intercomRe.FindStringSubmatch(address)
	if len(intercomMatch) > 1 {
		info.Intercom = intercomMatch[1]
	}

	// Извлечение номера этажа
	floorMatch := floorRe.FindStringSubmatch(address)
	if len(floorMatch) > 1 {
		info.Floor = floorMatch[1]
	}

	// Извлечение номера офиса или квартиры
	officeMatch := officeRe.FindStringSubmatch(address)
	if len(officeMatch) > 1 {
		info.Office = officeMatch[1]
	}

	// Проверка на наличие хотя бы одного элемента адреса
	if info.City == "" && info.Street == "" && info.HouseNumber == "" && info.Entrance == "" &&
		info.Intercom == "" && info.Floor == "" && info.Office == "" {
		return info, fmt.Errorf("не удалось извлечь информацию из адреса")
	}

	return info, nil
}
