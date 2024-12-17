package utils

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	coreMenuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	coreStoreModels "github.com/kwaaka-team/orders-core/core/storecore/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"strconv"
	"strings"
	"time"

	"github.com/kwaaka-team/orders-core/pkg/menu/dto"
)

func ActiveMenu(menus []coreStoreModels.StoreDSMenu, deliveryService dto.DeliveryService) (string, error) {
	for _, menu := range menus {
		if menu.IsActive && menu.Delivery == deliveryService.String() {
			return menu.ID, nil
		}
	}

	return "", errors.New("the restaurant does not have an active menu")
}

func ProductsMap(menu coreMenuModels.Menu) map[string]coreMenuModels.Product {
	m := map[string]coreMenuModels.Product{}

	for _, product := range menu.Products {
		m[product.ExtID] = product
	}

	return m
}

func AtributesMap(menu coreMenuModels.Menu) map[string]coreMenuModels.Attribute {
	m := map[string]coreMenuModels.Attribute{}

	for _, attribute := range menu.Attributes {
		m[attribute.ExtID] = attribute
	}

	return m
}

func AtributeGroupsMap(menu coreMenuModels.Menu) map[string]coreMenuModels.AttributeGroup {
	m := map[string]coreMenuModels.AttributeGroup{}

	for _, attributeGroup := range menu.AttributesGroups {
		m[attributeGroup.ExtID] = attributeGroup
	}

	return m
}

func ComboMap(menu coreMenuModels.Menu) map[string]coreMenuModels.Combo {
	m := map[string]coreMenuModels.Combo{}

	for _, combo := range menu.Combos {
		m[combo.ID] = combo
	}

	return m
}

func FindAttributeGroupID(
	productID string,
	attributeID string,
	menuProducts map[string]coreMenuModels.Product,
	menuAttributesGroups map[string]coreMenuModels.AttributeGroup,
	menuAttributes map[string]coreMenuModels.Attribute,
	isExternalMenu bool, parentAttributeGroupID string,
) string {

	menuProduct, exist := menuProducts[productID]

	if exist {
		if isExternalMenu {
			for _, menuAttributesGroup := range menuAttributesGroups {
				if _, ok := menuAttributesGroups[menuAttributesGroup.ExtID]; ok && len(menuAttributesGroup.Attributes) > 0 {
					for _, menuAttributeID := range menuAttributesGroup.Attributes {
						if attributeID == menuAttributes[menuAttributeID].PosID && menuAttributeID == menuAttributes[menuAttributeID].ExtID {
							return menuAttributes[menuAttributeID].ParentAttributeGroup
						} else {
							return parentAttributeGroupID
						}
					}
				}
			}
		} else {
			for _, menuAttributeGroupID := range menuProduct.AttributesGroups {
				menuAttributesGroup, ok := menuAttributesGroups[menuAttributeGroupID]

				if ok && len(menuAttributesGroup.Attributes) > 0 {
					for _, menuAttributeID := range menuAttributesGroup.Attributes {
						if attributeID == menuAttributeID {
							return menuAttributeGroupID
						}
					}
				}
			}
		}
	}

	return ""
}

func UploadMenuToS3(storeID, bucketName, deliveryService, shareMenuUrl string, menu interface{}, sv3 *s3.S3) (string, error) {
	trId := strconv.Itoa(int(time.Now().Unix()))
	link := strings.TrimSpace(fmt.Sprintf("publications/%s/%s/%s.json", deliveryService, storeID, trId))
	fileBody, err := json.Marshal(menu)
	if err != nil {
		log.Err(err).Msgf("s3 %v menu publication request load error: %v", deliveryService, err)
		return "", err
	}

	_, err = sv3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(link),
		Body:        strings.NewReader(string(fileBody)),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		log.Err(err).Msgf("s3 %v menu publication request load error", deliveryService)
	}
	objectUrl := fmt.Sprintf("%s/%s", shareMenuUrl, link)

	return objectUrl, nil
}

func UploadMenuDBVersionToS3(storeID, bucketName, deliveryService, shareMenuURL string, menu interface{}, sv3 *s3.S3) (string, error) {
	trId := deliveryService + strconv.Itoa(int(time.Now().Unix()))
	menuDBVersions := "menu_DB_versions"
	link := strings.TrimSpace(fmt.Sprintf("publications/%s/%s/%s.json", menuDBVersions, storeID, trId))
	fileBody, err := json.Marshal(menu)
	if err != nil {
		return "", err
	}

	_, err = sv3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(link),
		Body:        strings.NewReader(string(fileBody)),
		ContentType: aws.String("application/json"),
	})
	if err != nil {
		log.Err(err).Msgf("could not save menu DB version in S3 in store_id %s", storeID)
	}

	objectUrl := fmt.Sprintf("%s/%s", shareMenuURL, link)

	return objectUrl, nil
}
