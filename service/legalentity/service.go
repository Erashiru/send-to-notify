package legalentity

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/service/aws_s3"
	"github.com/kwaaka-team/orders-core/service/legalentity/models"
	leSel "github.com/kwaaka-team/orders-core/service/legalentity/models/selector"
	"github.com/kwaaka-team/orders-core/service/legalentity/repository"
	"github.com/nguyenthenguyen/docx"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
	"unicode"
)

type ServiceImpl struct {
	legalEntityRepo repository.LegalEntityRepository
	s3Info          models.S3Info
	s3Service       aws_s3.Service
}

func NewLegalEntityService(s3Info models.S3Info, repo repository.LegalEntityRepository, s3Service aws_s3.Service) *ServiceImpl {
	return &ServiceImpl{
		legalEntityRepo: repo,
		s3Info:          s3Info,
		s3Service:       s3Service,
	}
}

func (s *ServiceImpl) Create(ctx context.Context, profile models.LegalEntityForm) (string, error) {
	res, err := s.legalEntityRepo.Insert(ctx, profile)
	if err != nil {
		return "", err
	}

	return res, nil
}

func (s *ServiceImpl) Get(ctx context.Context, id string) (models.LegalEntityView, error) {
	legalEntity, err := s.legalEntityRepo.GetByID(ctx, id)
	if err != nil {
		return models.LegalEntityView{}, err
	}

	uniqueBrands := make(map[string]bool)

	for _, brand := range legalEntity.Brands {
		if _, ok := uniqueBrands[brand.Name]; !ok {
			uniqueBrands[brand.Name] = true
		}
	}

	var brands []string
	for b := range uniqueBrands {
		brands = append(brands, b)
	}

	var stores []string
	var integrations models.Integrations
	var paymentAmount models.Cost
	for _, s := range legalEntity.Stores {
		stores = append(stores, s.Name)

		if s.GlovoExists {
			integrations.Glovo++
		}
		if s.WoltExists {
			integrations.Wolt++
		}
		if s.YandexExists {
			integrations.Yandex++
		}

		paymentAmount.Value += s.Fare.Cost.Value
	}

	// TODO: currency is hardcoded, same reason as in List method below
	paymentAmount.Currency = "KZT"

	integrationsCount := integrations.Glovo + integrations.Wolt + integrations.Yandex

	legalEntity.BusinessInfo.Brands = brands
	legalEntity.BusinessInfo.BrandsCount = len(brands)
	legalEntity.BusinessInfo.Stores = stores
	legalEntity.BusinessInfo.StoresCount = len(stores)
	legalEntity.BusinessInfo.Integrations = integrations
	legalEntity.BusinessInfo.IntegrationsCount = integrationsCount
	legalEntity.BusinessInfo.PaymentAmount = paymentAmount

	return legalEntity, nil
}

func (s *ServiceImpl) List(ctx context.Context, pagination selector.Pagination, filter models.Filter) ([]models.GetListOfLegalEntities, error) {
	legalEntities, err := s.legalEntityRepo.List(ctx, pagination, filter)
	if err != nil {
		return nil, err
	}

	output := make([]models.GetListOfLegalEntities, 0)

	for _, le := range legalEntities {
		var temp models.GetListOfLegalEntities

		var overallCount int
		var fareOutput models.FareOutput
		fareOutput.Fares = make(map[string]models.Fare)
		integrationsCount := make(map[string]int)
		for _, s := range le.Stores {
			integrationsCount[s.Fare.Type]++
			fareOutput.Fares[s.Fare.Type] = s.Fare
			// TODO: this should take into account the currency, but how? (need to clarify)
			overallCount += s.Fare.Cost.Value
		}
		for fareType, count := range integrationsCount {
			fare := fareOutput.Fares[fareType]
			fare.IntegrationsAmount = count
			fareOutput.Fares[fareType] = fare
		}
		// TODO: currency is hard coded, because I still don't know how different currencies would be calculated and what currency to use in payment amount
		fareOutput.OverallCost = models.Cost{Value: overallCount, Currency: "KZT"}

		var brands []string
		for _, b := range le.Brands {
			brands = append(brands, b.Name)
		}

		temp.LegalEntityID = le.LegalEntityID
		temp.Name = le.Name
		temp.Brands = brands
		temp.PaymentType = le.PaymentType
		temp.Contacts = le.Contacts
		temp.PaymentAmount = fareOutput.OverallCost
		temp.Status = le.Status
		temp.Fare = fareOutput

		output = append(output, temp)
	}

	return output, nil
}

func (s *ServiceImpl) Update(ctx context.Context, updatedProfile leSel.LegalEntityForm, id string) error {
	err := s.legalEntityRepo.Update(ctx, updatedProfile, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) Disable(ctx context.Context, id string) error {
	err := s.legalEntityRepo.Disable(ctx, id)
	if err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) GetStores(ctx context.Context, pagination selector.Pagination, id string) (models.GetListOfStores, error) {
	storesInfoDB, err := s.legalEntityRepo.GetStores(ctx, pagination, id)
	if err != nil {
		return models.GetListOfStores{}, err
	}

	var getListOfStores models.GetListOfStores
	getListOfStores.StoresInfo = make([]models.StoreInfo, len(storesInfoDB.Stores))

	var overallPaymentAmount int

	for i := 0; i < len(storesInfoDB.Stores); i++ {
		getListOfStores.StoresInfo[i].Name = storesInfoDB.Stores[i].Name
		getListOfStores.StoresInfo[i].Brand = storesInfoDB.Brands[i].Name
		getListOfStores.StoresInfo[i].Fare = storesInfoDB.Stores[i].Fare
		getListOfStores.StoresInfo[i].GlovoExists = storesInfoDB.Stores[i].GlovoExists
		getListOfStores.StoresInfo[i].YandexExists = storesInfoDB.Stores[i].YandexExists
		getListOfStores.StoresInfo[i].WoltExists = storesInfoDB.Stores[i].WoltExists

		overallPaymentAmount += storesInfoDB.Stores[i].Fare.Cost.Value
	}

	getListOfStores.OverallPayment = models.Cost{Value: overallPaymentAmount, Currency: "KZT"} // TODO: hardcoded currency, as everywhere

	return getListOfStores, nil
}

func (s *ServiceImpl) UploadDocument(ctx context.Context, request models.UploadDocumentRequest) (string, error) {
	if request.Data == nil {
		return "", errors.New("missing file in document upload legal entity")
	}

	documentID := strconv.Itoa(int(s.generateDocumentID()))

	link, err := s.saveDocumentInS3(request.DocName, request.Extension, request.LegalEntityID, documentID, request.Data)
	if err != nil {
		return "", err
	}

	insertDocument := models.Document{
		ID:      documentID,
		DocName: request.DocName,
		Type:    request.DocType,
		S3Link:  link,
		Comment: request.Comment,
	}

	err = s.legalEntityRepo.InsertDocument(ctx, request.LegalEntityID, insertDocument)
	if err != nil {
		return "", err
	}

	return documentID, nil
}

func (s *ServiceImpl) DisableDocument(ctx context.Context, legalEntityID, documentID string) error {
	return s.legalEntityRepo.DisableDocument(ctx, legalEntityID, documentID)
}

func (s *ServiceImpl) GetAllDocumentsByLegalEntityID(ctx context.Context, pagination selector.Pagination, filter models.DocumentFilter, legalEntityID string) ([]models.Document, error) {
	return s.legalEntityRepo.GetAllDocumentsByLegalEntityID(ctx, pagination, filter, legalEntityID)
}

func (s *ServiceImpl) GetDocumentDownloadLink(ctx context.Context, legalEntityID, documentID string) (string, error) {
	return s.legalEntityRepo.GetDocumentDownloadLink(ctx, legalEntityID, documentID)
}

func (s *ServiceImpl) GenerateContract(contract models.ContractRequest) (models.ContractResponse, error) {

	var result models.ContractResponse

	decodedData, err := base64.StdEncoding.DecodeString(models.ContractFormBase64)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn base64.StdEncoding.DecodeString : %w", err)
	}
	docBuffer := bytes.NewReader(decodedData)
	file, err := docx.ReadDocxFromMemory(docBuffer, int64(len(models.ContractFormBase64)))
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn docx.ReadDocxFromMemory : %w", err)
	}
	defer file.Close()

	var (
		doc           = file.Editable()
		priceStr      string
		totalPrice    int
		totalPriceStr string
	)

	switch {
	case len(contract.IntegratedRestaurants) >= 1:
		for idx, store := range contract.IntegratedRestaurants {
			//var aggregators strings.Builder
			//for id, aggregator := range store.Aggregators {
			//	aggregators.WriteString(aggregator)
			//	if id != len(store.Aggregators)-1 {
			//		aggregators.WriteString(", ")
			//	}
			//}
			if idx == 0 {
				//err := doc.Replace("Aggregators", aggregators.String(), -1)
				//if err != nil {
				//	return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Aggregators) : %w", err)
				//}
				err = doc.Replace("RestName", store.RestaurantName, -1)
				if err != nil {
					return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(RestName) : %w", err)
				}
				if store.RestaurantAddress != "" {
					err = doc.Replace("RestAddress", store.RestaurantAddress, -1)
					if err != nil {
						return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(RestAddress) : %w", err)
					}
				} else {
					err = doc.Replace("RestAddress", "", -1)
					if err != nil {
						return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace() : %w", err)
					}
				}

				priceStr = formatPrice(store.Price)
				err = doc.Replace("Price", priceStr, -1)
				if err != nil {
					return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Price) : %w", err)
				}
				totalPrice += store.Price
				continue
			}
			body := doc.GetContent()
			searchText := "Цена"
			startIndex := strings.Index(body, "<w:tbl>")
			for startIndex != -1 {
				endIndex := strings.Index(body[startIndex:], "</w:tbl>") + startIndex + len("</w:tbl>")
				table := body[startIndex:endIndex]
				if strings.Contains(table, searchText) {
					rowStartTag := "<w:tr>"
					tableStartIndex := startIndex
					tableEndIndex := endIndex
					lastRowStartIndex := strings.LastIndex(table, rowStartTag)
					penultimateRowEndIndex := lastRowStartIndex
					newRow := fmt.Sprintf(`<w:tr><w:trPr><w:trHeight w:val="545" w:hRule="atLeast"/></w:trPr><w:tc><w:tcPr><w:tcW w:w="613" w:type="dxa"/><w:tcBorders><w:top w:val="single" w:sz="8" w:space="0" w:color="000000"/><w:left w:val="single" w:sz="8" w:space="0" w:color="000000"/><w:bottom w:val="single" w:sz="8" w:space="0" w:color="000000"/><w:right w:val="single" w:sz="8" w:space="0" w:color="000000"/></w:tcBorders></w:tcPr><w:p><w:pPr><w:pStyle w:val="Normal"/><w:jc w:val="center"/><w:rPr><w:rFonts w:ascii="Times New Roman" w:hAnsi="Times New Roman" w:eastAsia="Times New Roman" w:cs="Times New Roman"/></w:rPr></w:pPr><w:r><w:rPr><w:rFonts w:eastAsia="Times New Roman" w:cs="Times New Roman" w:ascii="Times New Roman" w:hAnsi="Times New Roman"/></w:rPr><w:t>%d</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="6452" w:type="dxa"/><w:tcBorders><w:top w:val="single" w:sz="8" w:space="0" w:color="000000"/><w:left w:val="single" w:sz="6" w:space="0" w:color="000000"/><w:bottom w:val="single" w:sz="8" w:space="0" w:color="000000"/><w:right w:val="single" w:sz="8" w:space="0" w:color="000000"/></w:tcBorders></w:tcPr><w:p><w:pPr><w:pStyle w:val="Normal"/><w:jc w:val="center"/><w:rPr><w:rFonts w:ascii="Times New Roman" w:hAnsi="Times New Roman" w:eastAsia="Times New Roman" w:cs="Times New Roman"/><w:b/></w:rPr></w:pPr><w:r><w:rPr><w:rFonts w:eastAsia="Times New Roman" w:cs="Times New Roman" w:ascii="Times New Roman" w:hAnsi="Times New Roman"/></w:rPr><w:t xml:space="preserve">Интеграция 1 точки «%s», </w:t></w:r><w:r><w:rPr><w:rFonts w:eastAsia="Times New Roman" w:cs="Times New Roman" w:ascii="Times New Roman" w:hAnsi="Times New Roman"/><w:shd w:fill="auto" w:val="clear"/></w:rPr><w:t>расположенного по адресу: %s</w:t></w:r></w:p></w:tc><w:tc><w:tcPr><w:tcW w:w="1950" w:type="dxa"/><w:tcBorders><w:top w:val="single" w:sz="8" w:space="0" w:color="000000"/><w:left w:val="single" w:sz="6" w:space="0" w:color="000000"/><w:bottom w:val="single" w:sz="8" w:space="0" w:color="000000"/><w:right w:val="single" w:sz="8" w:space="0" w:color="000000"/></w:tcBorders></w:tcPr><w:p><w:pPr><w:pStyle w:val="Normal"/><w:jc w:val="center"/><w:rPr><w:rFonts w:ascii="Times New Roman" w:hAnsi="Times New Roman" w:eastAsia="Times New Roman" w:cs="Times New Roman"/></w:rPr></w:pPr><w:r><w:rPr><w:rFonts w:eastAsia="Times New Roman" w:cs="Times New Roman" w:ascii="Times New Roman" w:hAnsi="Times New Roman"/></w:rPr><w:t>%s</w:t></w:r></w:p></w:tc></w:tr>`, idx+1, store.RestaurantName, store.RestaurantAddress, formatPrice(store.Price))
					updatedTable := table[:penultimateRowEndIndex] + newRow + table[penultimateRowEndIndex:]
					updatedContent := body[:tableStartIndex] + updatedTable + body[tableEndIndex:]
					doc.SetContent(updatedContent)
					totalPrice += store.Price
				}
				startIndex = strings.Index(body[endIndex:], "<w:tbl>")
				if startIndex != -1 {
					startIndex += endIndex
				}
			}
		}
	default:
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract : %w", err)
	}

	totalPriceStr = formatPrice(totalPrice)
	err = doc.Replace("Total", totalPriceStr, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Total) : %w", err)
	}
	err = doc.Replace("Number", contract.ContractNum, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Number) : %w", err)
	}
	formattedDate, err := formatSignatureDate(contract.SignatureDate)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn formatSignatureDate : %w", err)
	}
	err = doc.Replace("Date", formattedDate, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Date) :  %w", err)
	}

	formattedName := formatFullName(contract.FullNameHead)
	err = doc.Replace("FormattedName", formattedName, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(FormattedName) : %w", err)
	}

	if contract.OrganizationType == "ТОО" || contract.OrganizationType == "TOO" {
		err := doc.Replace("OrgType", contract.OrganizationType, -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(TOO) : %w", err)
		}
		err = doc.Replace("OrgName", contract.OrganizationName, -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(OrgName) : %w", err)
		}
		err = doc.Replace("CompanyHeadTitle", "Директор", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(CompanyHeadTitle) : %w", err)
		}
		err = doc.Replace("OrgDoc", "Устава", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(OrgTypeDoc) : %w", err)
		}
	} else {
		err = doc.Replace("OrgType", contract.OrganizationType, -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(OrgType) : %w", err)
		}
		err = doc.Replace("OrgName", contract.OrganizationName, -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(OrgName) : %w", err)
		}
		err = doc.Replace("CompanyHeadTitle", "", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(CompanyHeadTitle) : %w", err)
		}
		if contract.TicketNumber == "" || contract.TicketDate == "" {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract: ticket number or ticket date not found")
		}
		err = doc.Replace("OrgDoc", "талона № "+contract.TicketNumber+" от "+contract.TicketDate+" г.", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(OrgTypeDoc) : %w", err)
		}
	}
	err = doc.Replace("FullNameHead", contract.FullNameHead, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(FullNameHead) : %w", err)
	}
	err = doc.Replace("FullNameLegal", contract.FullNameLegal, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(FullNameLegal) : %w", err)
	}
	err = doc.Replace("JobTitle", contract.JobTitle, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(JobTitle) : %w", err)
	}
	if contract.PhoneNumber != "" {
		err = doc.Replace("Phone", "+"+contract.PhoneNumber, -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Phone) : %w", err)
		}
	} else {
		err = doc.Replace("Phone", "", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace() : %w", err)
		}
	}
	if contract.Email != "" {
		emailStr := fmt.Sprintf(", %s", contract.Email)
		err = doc.Replace("Email", emailStr, -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Email) : %w", err)
		}
	} else {
		err = doc.Replace("Email", "", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace() : %w", err)
		}
	}
	err = doc.Replace("BIN", contract.BIN, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(BIN) : %w", err)
	}
	err = doc.Replace("LegalAddress", contract.LegalAddress, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(LegalAddress) : %w", err)
	}
	err = doc.Replace("ActualAddress", contract.ActualAddress, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(ActualAddress) : %w", err)
	}

	switch contract.Tariff {
	case "monthly":
		err = doc.Replace("PeriodType", "месяца", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Period) : %w", err)
		}
		err = doc.Replace("PeriodCalendar", "календарный месяц", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(PeriodCalendar) : %w", err)
		}
		err = doc.Replace("PeriodTable", "месяц", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(PeriodPrice) : %w", err)
		}
		err = doc.Replace("PeriodAdverb", "ежемесячно", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(PeriodAdverb) : %w", err)
		}
	case "quarterly":
		err = doc.Replace("PeriodType", "Квартала", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Period) : %w", err)
		}
		err = doc.Replace("PeriodCalendar", "Квартал", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(PeriodCalendar) : %w", err)
		}
		err = doc.Replace("PeriodTable", "Квартал", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(PeriodPrice) : %w", err)
		}
		err = doc.Replace("PeriodAdverb", "ежеквартально", -1)
		if err != nil {
			return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(PeriodAdverb) : %w", err)
		}
	default:
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract: invalid tariff input value")
	}

	err = doc.Replace("Bank", contract.Bank, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(Bank) : %w", err)
	}
	err = doc.Replace("BIK", contract.BIK, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(BIK) : %w", err)
	}
	err = doc.Replace("IIK", contract.IIK, -1)
	if err != nil {
		return models.ContractResponse{}, fmt.Errorf("service/legalentity/service - fn GenerateContract - fn Replace(IIK) : %w", err)
	}

	fileName := contract.OrganizationType + " " + contract.OrganizationName + " " + contract.SignatureDate + ".docx"
	result.FileName = fileName

	modifiedBuffer := new(bytes.Buffer)
	err = doc.Write(modifiedBuffer)
	if err != nil {
		return models.ContractResponse{}, err
	}
	modifiedBase64String := base64.StdEncoding.EncodeToString(modifiedBuffer.Bytes())
	result.ContractBase64 = modifiedBase64String

	return result, nil
}

func formatSignatureDate(dateStr string) (string, error) {
	dateStr = strings.TrimFunc(dateStr, func(r rune) bool {
		return !unicode.IsDigit(r) && r != '/' && r != '-' && r != '.'
	})
	pattern := `^\d{1,2}[-/\.]\d{1,2}[-/\.]\d{2,4}$`
	matched, err := regexp.MatchString(pattern, dateStr)
	if err != nil {
		return "", fmt.Errorf("service/legalentity/service - formatSignatureDate - fn regexp.MatchString : %w", err)
	}
	if !matched {
		return "", fmt.Errorf("wrong data format: %s", dateStr)
	}
	parts := strings.FieldsFunc(dateStr, func(r rune) bool {
		return r == '/' || r == '-' || r == '.'
	})
	if len(parts) == 1 {
		parts = strings.Split(dateStr, "-")
		if len(parts) == 1 {
			parts = strings.Split(dateStr, ".")
		}
	}
	month, _ := strconv.Atoi(parts[0])
	day, _ := strconv.Atoi(parts[1])
	year, _ := strconv.Atoi(parts[2])

	if year < 100 {
		year += 2000
	}
	formattedDateStr := fmt.Sprintf("%02d/%02d/%04d", month, day, year)
	date, err := time.Parse("01/02/2006", formattedDateStr)
	if err != nil {
		return "", fmt.Errorf("service/legalentity/service - formatSignatureDate - fn time.Parse : %w", err)
	}
	months := [...]string{
		"января", "февраля", "марта", "апреля", "мая", "июня",
		"июля", "августа", "сентября", "октября", "ноября", "декабря",
	}
	formattedDate := fmt.Sprintf("«%d» %s %d г.", date.Day(), months[date.Month()-1], date.Year())
	return formattedDate, nil
}

func formatFullName(name string) string {
	parts := strings.Split(name, " ")
	if len(parts) == 3 {
		lastName := parts[0]
		firstName := parts[1]
		patronymic := parts[2]

		for ending, newEnding := range models.LastNameAndPatronymicEndings {
			if strings.HasSuffix(lastName, ending) {
				lastName = strings.TrimSuffix(lastName, ending) + newEnding
				break
			}
		}

		for ending, newEnding := range models.FirstNameEndings {
			if strings.HasSuffix(firstName, ending) {
				firstName = strings.TrimSuffix(firstName, ending) + newEnding
				break
			}
		}

		for ending, newEnding := range models.LastNameAndPatronymicEndings {
			if strings.HasSuffix(patronymic, ending) {
				patronymic = strings.TrimSuffix(patronymic, ending) + newEnding
				break
			}
		}

		return fmt.Sprintf("%s %s %s", lastName, firstName, patronymic)
	}
	return name
}

func formatPrice(price int) string {
	priceStr := strconv.Itoa(price)
	var parts []string
	for i := len(priceStr); i > 0; i -= 3 {
		start := i - 3
		if start < 0 {
			start = 0
		}
		parts = append([]string{priceStr[start:i]}, parts...)
	}
	formattedPrice := strings.Join(parts, " ")
	return formattedPrice
}
