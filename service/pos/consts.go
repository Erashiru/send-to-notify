package pos

import (
	"github.com/kwaaka-team/orders-core/service/error_solutions/models"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"regexp"
	"strings"
)

const (
	PAYMENT_TYPE_MISSED string = "payment type not found"
	ATTRIBUTE_MISSED    string = "attribute not found in pos menu: "
	PRODUCT_MISSED      string = "product not found in pos menu: "
	OFFERS_MISSED       string = "couldn't get BK offers"
	INTEGRATION_OFF     string = "Integration is off"
	ORDER_ALREADY_EXIST string = "Order already exist in db"
)

const (
	CREATION_TIMEOUT_CODE    string = "2"
	PRODUCT_MISSED_CODE      string = "4"
	ATTRIBUTE_MISSED_CODE    string = "5"
	NEED_TO_HEAL_ORDER_CODE  string = "6"
	PAYMENT_TYPE_MISSED_CODE string = "10"
	BK_ORDER_FAILED          string = "17"
	OFFERS_MISSED_CODE       string = "18"
	INTEGRATION_OFF_CODE     string = "19"
	WSA_TIMEOUT              string = "38"
	ORDER_ALREADY_EXIST_CODE string = "20"

	OTHER_FAIL_REASON_CODE string = "66"
)

const (
	palomaIntTrue              = 1
	palomaDeliveryTypeCourier  = 1
	palomaDeliveryTypeCustomer = 2
)

var ErrUnsupportedMethod = errors.New("unsupported method")

func MatchingCodes(message string, errorSolutions []models.ErrorSolution) string {

	for _, errorSolution := range errorSolutions {
		if errorSolution.ContainsText != "" {
			if errorSolution.Code == "9" {
				re := regexp.MustCompile(`Payment type.*?is deleted\.`)
				match := re.FindString(message)
				if match != "" {
					log.Info().Msgf("message: %s, errorSolution.Code: %s", message, errorSolution.Code)
					return errorSolution.Code
				}
			}
			if strings.Contains(message, errorSolution.ContainsText) {
				log.Info().Msgf("message: %s, errorSolution.Code: %s, containsText: %s", message, errorSolution.Code, errorSolution.ContainsText)
				return errorSolution.Code
			}
		}
	}

	log.Info().Msgf("there was no match for message: %s", message)
	return OTHER_FAIL_REASON_CODE
}

func GetProductIDFromRegexp(message string, solution models.ErrorSolution) string {

	regex := regexp.MustCompile(solution.RegexpToFindProduct)
	match := regex.FindStringSubmatch(message)

	if len(match) > 1 {
		log.Info().Msgf("productID after error solution regexp: %s", match[1])
		return match[1]
	}
	return ""
}
