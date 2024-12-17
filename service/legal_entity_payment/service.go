package legal_entity_payment

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jung-kurt/gofpdf/v2"
	"github.com/kwaaka-team/orders-core/config/general"
	"github.com/kwaaka-team/orders-core/core/menu/models/selector"
	"github.com/kwaaka-team/orders-core/service/aws_s3"
	"github.com/kwaaka-team/orders-core/service/legal_entity_payment/models"
	"github.com/kwaaka-team/orders-core/service/legalentity"
	legalEntityModels "github.com/kwaaka-team/orders-core/service/legalentity/models"
	"github.com/kwaaka-team/orders-core/service/whatsapp"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Service interface {
	CreatePayment(ctx context.Context, payment models.LegalEntityPayment) (string, error)
	GetPaymentByID(ctx context.Context, paymentID string) (models.LegalEntityPayment, error)
	GetList(ctx context.Context, query models.ListLegalEntityPaymentQuery) ([]models.ListLegalEntityPayment, error)
	Update(ctx context.Context, payment models.UpdateLegalEntityPayment) error
	Delete(ctx context.Context, paymentID string) error
	GetLegalEntityPaymentAnalytics(ctx context.Context, query models.LegalEntityPaymentAnalyticsRequest) (models.LegalEntityPaymentAnalyticsResponse, error)
	UploadPDF(ctx context.Context, query models.LegalEntityPaymentDownloadPDFRequest) (string, error)
	CreateBill(ctx context.Context, query models.LegalEntityPaymentCreateBillRequest) error
	ConfirmPayment(ctx context.Context, query models.LegalEntityPaymentConfirmPaymentRequest) error
	SendPayment(ctx context.Context) error
}

type ServiceImpl struct {
	cfg                 general.Configuration
	repo                Repository
	s3Service           aws_s3.Service
	legalEntityService  *legalentity.ServiceImpl
	googleSheetsService *sheets.Service
	whatsAppService     whatsapp.Service
}

func NewService(cfg general.Configuration, repo Repository, s3Service aws_s3.Service, legalEntityService *legalentity.ServiceImpl, whatsAppService whatsapp.Service) (*ServiceImpl, error) {
	sheetsService, err := sheets.NewService(context.Background(), option.WithCredentialsJSON(googleCredentialsToJSON(cfg.GoogleSheetsEmail, cfg.GoogleSheetsKey, cfg.GoogleSheetsCredsType)), option.WithScopes(sheets.SpreadsheetsScope))
	if err != nil {
		return &ServiceImpl{}, err
	}

	return &ServiceImpl{
		cfg:                 cfg,
		repo:                repo,
		s3Service:           s3Service,
		legalEntityService:  legalEntityService,
		googleSheetsService: sheetsService,
		whatsAppService:     whatsAppService,
	}, nil
}

const (
	paymentStatusSend     = "Выслан"
	paymentStatusDontSend = "Не выслан"
	spreadSheetId         = "1McRkJY61tUNIQRhweosjPrNQib8QtsG5FW2KOuGY6Lo"
)

func googleCredentialsToJSON(email, key, credsType string) []byte {
	return []byte(fmt.Sprintf(`{
		"client_email": "%s",
		"private_key": "%s",
		"type": "%s"
	}`, email, key, credsType))
}

func (s *ServiceImpl) CreatePayment(ctx context.Context, payment models.LegalEntityPayment) (string, error) {
	payment.Status = models.UNBILLED.String()
	return s.repo.Create(ctx, payment)
}

func (s *ServiceImpl) GetPaymentByID(ctx context.Context, paymentID string) (models.LegalEntityPayment, error) {
	return s.repo.GetByID(ctx, paymentID)
}

func (s *ServiceImpl) GetList(ctx context.Context, query models.ListLegalEntityPaymentQuery) ([]models.ListLegalEntityPayment, error) {
	allLegalEntities, err := s.legalEntityService.List(ctx, selector.Pagination{Limit: query.Pagination.Limit, Page: query.Pagination.Page - 1}, legalEntityModels.Filter{})
	if err != nil {
		return nil, err
	}

	var (
		legalEntities []legalEntityModels.GetListOfLegalEntities
		res           []models.ListLegalEntityPayment
	)

	switch len(query.LegalEntityIDs) {
	case 0:
		legalEntities = allLegalEntities
	default:
		legalEntitiesFilterMap := make(map[string]bool, 3)
		for i := range query.LegalEntityIDs {
			legalEntitiesFilterMap[query.LegalEntityIDs[i]] = true
		}

		for i := range allLegalEntities {
			if legalEntitiesFilterMap[legalEntities[i].LegalEntityID] {
				legalEntities = append(legalEntities, allLegalEntities[i])
			}
		}
	}

	query.Limit = 0
	for i := range legalEntities {
		query.LegalEntityIDs = []string{legalEntities[i].LegalEntityID}
		legalEntityPayments, err := s.repo.GetList(ctx, query)
		if err != nil {
			return nil, err
		}

		legalEntity := models.ListLegalEntityPayment{
			LegalEntityId:       legalEntities[i].LegalEntityID,
			Name:                legalEntities[i].Name,
			PaymentType:         legalEntities[i].PaymentType,
			Status:              s.getLegalEntityStatusForList(legalEntityPayments),
			BilledAmount:        s.getLegalEntityAmountForList(legalEntityPayments),
			PaidAmount:          s.getLegalEntityPaidAmountForList(legalEntityPayments),
			Brands:              legalEntities[i].Brands,
			LegalEntityPayments: legalEntityPayments,
		}

		res = append(res, legalEntity)
	}

	return res, nil
}

func (s *ServiceImpl) getLegalEntityStatusForList(legalEntityPayments []models.LegalEntityPayment) string {
	statusDefine := int(models.PAID_CONFIRMED)
	for j := range legalEntityPayments {
		if statusDefine > int(models.StringToLegalEntityPaymentStatus(legalEntityPayments[j].Status)) && int(models.StringToLegalEntityPaymentStatus(legalEntityPayments[j].Status)) > 0 {
			statusDefine = int(models.StringToLegalEntityPaymentStatus(legalEntityPayments[j].Status))
		}
	}
	return models.LegalEntityPaymentStatus(statusDefine).String()
}

func (s *ServiceImpl) getLegalEntityAmountForList(legalEntityPayments []models.LegalEntityPayment) float64 {
	var res float64
	for i := range legalEntityPayments {
		res += legalEntityPayments[i].Amount
	}
	return res
}

func (s *ServiceImpl) getLegalEntityPaidAmountForList(legalEntityPayments []models.LegalEntityPayment) float64 {
	var res float64
	for i := range legalEntityPayments {
		res += legalEntityPayments[i].PaidAmount
	}
	return res
}

func (s *ServiceImpl) Update(ctx context.Context, payment models.UpdateLegalEntityPayment) error {
	return s.repo.Update(ctx, payment)
}

func (s *ServiceImpl) Delete(ctx context.Context, paymentID string) error {
	return s.repo.Delete(ctx, paymentID)
}

func (s *ServiceImpl) GetLegalEntityPaymentAnalytics(ctx context.Context, query models.LegalEntityPaymentAnalyticsRequest) (models.LegalEntityPaymentAnalyticsResponse, error) {
	paid, err := s.repo.GetPaidPaymentsAnalytics(ctx, query)
	if err != nil {
		return models.LegalEntityPaymentAnalyticsResponse{}, err
	}

	unpaid, err := s.repo.GetUnpaidPaymentsAnalytics(ctx, query)
	if err != nil {
		return models.LegalEntityPaymentAnalyticsResponse{}, err
	}

	total := models.LegalEntityPaymentAnalytics{
		Amount:          paid.Amount + unpaid.Amount,
		Quantity:        paid.Quantity + unpaid.Quantity,
		AmountPercent:   100,
		QuantityPercent: 100,
	}

	paid.AmountPercent = s.getPercentForAnalytics(paid.Amount, total.Amount)
	paid.QuantityPercent = s.getPercentForAnalytics(float64(paid.Quantity), float64(total.Quantity))
	unpaid.AmountPercent = s.getPercentForAnalytics(unpaid.Amount, total.Amount)
	unpaid.QuantityPercent = s.getPercentForAnalytics(float64(unpaid.Quantity), float64(total.Quantity))

	return models.LegalEntityPaymentAnalyticsResponse{
		Paid:   paid,
		Unpaid: unpaid,
		Total:  total,
	}, nil
}

func (s *ServiceImpl) getPercentForAnalytics(val, total float64) float64 {
	return val * 100 / total
}

func (s *ServiceImpl) UploadPDF(ctx context.Context, query models.LegalEntityPaymentDownloadPDFRequest) (string, error) {
	if query.File == nil {
		return "", errors.New("missing file in PDF")
	}

	payment, err := s.GetPaymentByID(ctx, query.LegalEntityPaymentID)
	if err != nil {
		return "", err
	}

	link, err := s.saveFilePDFInS3(query.LegalEntityPaymentID, payment.LegalEntityID, query.File)
	if err != nil {
		return "", err
	}

	return link, nil
}

func (s *ServiceImpl) saveFilePDFInS3(legalEntityPaymentID, legalEntityID string, filePDF []byte) (string, error) {
	changeSpaces := strings.Replace(legalEntityPaymentID, " ", "_", -1)
	name := fmt.Sprintf("%s_%d", changeSpaces, time.Now().Unix())

	link := strings.TrimSpace(fmt.Sprintf("legal_entities_payments/%s/%s", legalEntityID, name))
	contentType := "application/pdf"
	if err := s.s3Service.PutPDF(link, filePDF, s.cfg.KwaakaFilesBucket, contentType); err != nil {
		return "", err
	}

	resLink := fmt.Sprintf("%s/%s.pdf", s.cfg.KwaakaFilesBaseUrl, link)

	return resLink, nil
}

func (s *ServiceImpl) CreateBill(ctx context.Context, query models.LegalEntityPaymentCreateBillRequest) error {
	if err := s.createBillValidate(query); err != nil {
		return err
	}

	update := s.setCreateBillQuery(query)

	if err := s.Update(ctx, models.UpdateLegalEntityPayment{
		ID:        update.ID,
		Name:      update.Name,
		StartDate: update.StartDate,
		EndDate:   update.EndDate,
		Amount:    update.Amount,
		Status:    update.Status,
		BillingAt: update.BillingAt,
		Bill:      update.Bill,
	}); err != nil {
		return err
	}

	return nil
}

func (s *ServiceImpl) createBillValidate(query models.LegalEntityPaymentCreateBillRequest) error {
	if query.LegalEntityPaymentID == "" {
		return errors.New("missing legal entity payment id")
	}
	if query.BillLink == "" {
		return errors.New("missing bill, add pdf url")
	}
	return nil
}

func (s *ServiceImpl) setCreateBillQuery(query models.LegalEntityPaymentCreateBillRequest) models.UpdateLegalEntityPayment {
	status := models.BILLED.String()
	timeNow := time.Now().UTC()
	res := models.UpdateLegalEntityPayment{
		ID:        &query.LegalEntityPaymentID,
		Bill:      &query.BillLink,
		Status:    &status,
		BillingAt: &timeNow,
	}

	if query.Name != "" {
		res.Name = &query.Name
	}
	if query.Amount != 0 {
		res.Amount = &query.Amount
	}
	if !query.StartDate.IsZero() {
		res.StartDate = &query.StartDate
	}
	if !query.EndDate.IsZero() {
		res.EndDate = &query.EndDate
	}

	return res
}

func (s *ServiceImpl) ConfirmPayment(ctx context.Context, query models.LegalEntityPaymentConfirmPaymentRequest) error {

	payment, err := s.GetPaymentByID(ctx, query.LegalEntityPaymentID)
	if err != nil {
		return err
	}

	timeNow := time.Now().UTC()
	status := models.PAID_CONFIRMED.String()

	switch {
	case payment.Status == models.PAID.String() && payment.BillPayment != "":
		if err := s.Update(ctx, models.UpdateLegalEntityPayment{
			ID:               &query.LegalEntityPaymentID,
			Status:           &status,
			ConfirmPaymentAt: &timeNow,
		}); err != nil {
			return err
		}

	default:
		if query.BillLink == "" {
			return errors.New("missing bill, add pdf url")
		}

		if err := s.Update(ctx, models.UpdateLegalEntityPayment{
			ID:               &query.LegalEntityPaymentID,
			Status:           &status,
			BillPaymentAt:    &timeNow,
			ConfirmPaymentAt: &timeNow,
			BillPayment:      &query.BillLink,
		}); err != nil {
			return err
		}

	}
	return nil
}

func (s *ServiceImpl) SendPayment(ctx context.Context) error {
	log.Info().Msgf("send payment to whatsapp")

	payments, integrations, err := s.getPaymentsFromSheets()
	if err != nil {
		log.Err(err).Msg("send payment failed get payments from sheets")
		return fmt.Errorf("get payments from sheets error: %s", err)
	}
	log.Info().Msgf("send payment to whatsapp, get payments from sheets success")

	integrationMap := make(map[string][]models.IntegrationXlsx, len(integrations))
	for i := range integrations {
		integrationMap[integrations[i].Number] = append(integrationMap[integrations[i].Number], integrations[i])
	}

	for i := range payments {
		log.Info().Msgf("send payment to whatsapp, work with payment run: %s, status: %s", payments[i].Number, payments[i].Status)

		if payments[i].Status != paymentStatusDontSend && payments[i].Status != paymentStatusSend {
			log.Info().Msgf("send payment to whatsapp, payment: %s skipped with status: %s", payments[i].Number, payments[i].Status)
			continue
		}

		pdf, err := s.createPdfForPayment(payments[i], integrationMap[payments[i].Number])
		if err != nil {
			log.Err(err).Msgf("send payment failed to create pdf payment: %s", payments[i].Number)
			continue
		}

		statusCell := fmt.Sprintf("1!K%d", i+2)

		if err = s.selectAndSendPdf(ctx, pdf, payments[i], statusCell); err != nil {
			log.Err(err).Msgf("send payment failed send in whatsapp, payment: %s", payments[i].Number)
			continue
		}

		log.Info().Msgf("send payment to whatsapp, work with payment successfully %s", payments[i].Number)
	}

	return nil
}

func (s *ServiceImpl) getPaymentsFromSheets() ([]models.PaymentXlsx, []models.IntegrationXlsx, error) {
	var (
		paymentsXlsx     []models.PaymentXlsx
		integrationsXlsx []models.IntegrationXlsx
		paymentPage      = "1!A2:K"
		integrationPage  = "2!A2:K"
	)

	paymentPageResp, err := s.googleSheetsService.Spreadsheets.Values.Get(spreadSheetId, paymentPage).Do()
	if err != nil {
		return nil, nil, err
	}

	if len(paymentPageResp.Values) == 0 {
		return nil, nil, errors.New("no payment page data found")
	}

	for _, row := range paymentPageResp.Values {
		payment := models.PaymentXlsx{
			Number:      row[0].(string),
			Name:        row[1].(string),
			Phone:       row[2].(string),
			Bank:        row[3].(string),
			BIK:         row[4].(string),
			Code:        row[5].(string),
			Month:       row[6].(string),
			BillingDate: row[7].(string),
			Buyer:       row[8].(string),
			Contract:    row[9].(string),
			Status:      row[10].(string),
		}
		paymentsXlsx = append(paymentsXlsx, payment)
	}

	integrationPageResp, err := s.googleSheetsService.Spreadsheets.Values.Get(spreadSheetId, integrationPage).Do()
	if err != nil {
		return nil, nil, err
	}

	if len(integrationPageResp.Values) == 0 {
		return nil, nil, errors.New("no integration page data found")
	}

	for _, row := range integrationPageResp.Values {
		integration := models.IntegrationXlsx{
			Number:   row[0].(string),
			Name:     row[1].(string),
			Amount:   row[2].(string),
			Price:    row[3].(string),
			SumPrice: row[4].(string),
		}
		integrationsXlsx = append(integrationsXlsx, integration)
	}

	return paymentsXlsx, integrationsXlsx, nil
}

func (s *ServiceImpl) createPdfForPayment(payment models.PaymentXlsx, integrations []models.IntegrationXlsx) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")

	fontBytes, err := os.ReadFile("fonts/ArialCyrRegular.ttf")
	if err != nil {
		return nil, err
	}

	boldFontFBytes, err := os.ReadFile("fonts/kzArialBold.ttf")
	if err != nil {
		return nil, err
	}

	pdf.AddUTF8FontFromBytes("Arial", "", fontBytes)
	pdf.AddUTF8FontFromBytes("Arial", "B", boldFontFBytes)

	pdf.AddPage()
	pdf.SetMargins(10, 20, 20)

	pdf.SetFont("Arial", "", 8)
	pdf.MultiCell(0, 3, " Внимание! Оплата данного счета означает согласие с условиями поставки товара. Уведомление об оплате\nобязательно, в противном случае не гарантируется наличие товара на складе. Товар отпускается по факту\nприхода денег на р/с Поставщика, самовывозом, при наличии доверенности и документов удостоверяющих\n", "", "R", false)
	pdf.MultiCell(0, 3, "личность.", "", "C", false)
	pdf.Ln(10)

	pdf.SetFont("Arial", "B", 10)
	pdf.Cell(0, 3, "Образец платежного поручения")
	pdf.Ln(4)

	yPosition := pdf.GetY()
	pdf.SetFont("Arial", "B", 9)
	pdf.MultiCell(100, 4, "Бенефициар:\nТоварищество с ограниченной ответственностью «Kwaaka»", "LRT", "L", false)
	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(100, 4, "БИН: 220640026674", "LRB", "L", false)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(pdf.GetX()+100, yPosition)
	pdf.MultiCell(50, 6, "ИИК\nKZ578562203118112745", "1", "C", false)
	pdf.SetXY(pdf.GetX()+150, yPosition)
	pdf.MultiCell(30, 6, "Кбе\n17", "1", "C", false)

	pdf.SetFont("Arial", "", 9)
	pdf.MultiCell(100, 4, "Банк бенефициара:\nАО \"Банк ЦентрКредит", "1", "L", false)
	pdf.SetFont("Arial", "B", 9)
	pdf.SetXY(pdf.GetX()+100, pdf.GetY()-8)
	pdf.MultiCell(30, 4, "БИК\nKCJBKZKX", "1", "C", false)
	pdf.SetXY(pdf.GetX()+130, pdf.GetY()-8)
	pdf.MultiCell(50, 4, "Код назначения платежа\n851", "1", "C", false)
	pdf.Ln(7.5)

	pdf.SetFont("Arial", "B", 14)
	pdf.Cell(0, 3, fmt.Sprintf("Счет на оплату № %s от %s", payment.Number, payment.BillingDate))
	pdf.Ln(7.5)

	pdf.Line(10, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(4)

	yPosition = pdf.GetY()
	pdf.SetFont("Arial", "", 9.5)
	pdf.MultiCell(25, 4, "Поставщик:", "", "L", false)
	pdf.SetXY(pdf.GetX()+25, yPosition)
	pdf.SetFont("Arial", "B", 9.5)
	pdf.MultiCell(155, 4, "БИН / ИИН 220640026674,Товарищество с ограниченной ответственностью «Kwaaka»,Республика Казахстан, Ауэзовский район, г. Алматы, мкр. Аксай 5, дом № 25, к.183", "", "L", false)
	pdf.Ln(4)

	yPosition = pdf.GetY()
	pdf.SetFont("Arial", "", 9.5)
	pdf.MultiCell(25, 4, "Покупатель:", "", "L", false)
	pdf.SetXY(pdf.GetX()+25, yPosition)
	pdf.SetFont("Arial", "B", 9.5)
	pdf.MultiCell(155, 4, payment.Buyer, "", "L", false)
	pdf.Ln(4)

	yPosition = pdf.GetY()
	pdf.SetFont("Arial", "", 9.5)
	pdf.MultiCell(25, 4, "Договор:", "", "L", false)
	pdf.SetXY(pdf.GetX()+25, yPosition)
	pdf.SetFont("Arial", "B", 9.5)
	pdf.MultiCell(155, 4, payment.Contract, "", "L", false)
	pdf.Ln(4)

	yPosition = pdf.GetY()
	pdf.SetFont("Arial", "B", 9.5)
	pdf.MultiCell(10, 4, "No", "1", "C", false)
	pdf.SetXY(pdf.GetX()+10, yPosition)
	pdf.MultiCell(25, 4, "Код", "1", "C", false)
	pdf.SetXY(pdf.GetX()+35, yPosition)
	pdf.MultiCell(60, 4, "Наименование", "1", "C", false)
	pdf.SetXY(pdf.GetX()+95, yPosition)
	pdf.MultiCell(15, 4, "Кол-во", "1", "C", false)
	pdf.SetXY(pdf.GetX()+110, yPosition)
	pdf.MultiCell(10, 4, "Ед", "1", "C", false)
	pdf.SetXY(pdf.GetX()+120, yPosition)
	pdf.MultiCell(30, 4, "Цена", "1", "C", false)
	pdf.SetXY(pdf.GetX()+150, yPosition)
	pdf.MultiCell(30, 4, "Сумма", "1", "C", false)

	var totalSum float64
	yPosition = pdf.GetY()
	pdf.SetFont("Arial", "", 8)
	for i := range integrations {
		number := strconv.Itoa(i + 1)
		pdf.SetXY(pdf.GetX()+35, yPosition)
		pdf.MultiCell(60, 4, integrations[i].Name, "1", "L", false)
		secondPosY := pdf.GetY()
		maxHeight := secondPosY - yPosition
		pdf.SetXY(pdf.GetX(), yPosition)
		pdf.MultiCell(10, maxHeight, number, "1", "C", false)
		pdf.SetXY(pdf.GetX()+10, yPosition)
		pdf.MultiCell(25, maxHeight, "00000000088", "1", "L", false)
		pdf.SetXY(pdf.GetX()+95, yPosition)
		pdf.MultiCell(15, maxHeight, integrations[i].Amount, "1", "R", false)
		pdf.SetXY(pdf.GetX()+110, yPosition)
		pdf.MultiCell(10, maxHeight, "мес", "1", "R", false)
		pdf.SetXY(pdf.GetX()+120, yPosition)
		pdf.MultiCell(30, maxHeight, integrations[i].Price, "1", "R", false)
		pdf.SetXY(pdf.GetX()+150, yPosition)
		pdf.MultiCell(30, maxHeight, integrations[i].SumPrice, "1", "R", false)
		yPosition = secondPosY

		totalSum, err = s.getTotlaSumInPdf(integrations[i].SumPrice, totalSum)
		if err != nil {
			return nil, err
		}
	}
	pdf.Ln(5)

	pdf.SetFont("Arial", "B", 9.5)
	pdf.MultiCell(0, 4, fmt.Sprintf("Итого: %.2f", totalSum), "", "R", false)
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 9.5)
	pdf.MultiCell(0, 4, fmt.Sprintf("Всего наименований %d, на сумму %.2f KZT", len(integrations), totalSum), "", "L", false)
	pdf.SetFont("Arial", "B", 9.5)
	pdf.MultiCell(0, 4, fmt.Sprintf("Всего к оплате: %s", s.getNumberWords(totalSum)), "", "L", false)
	pdf.Ln(5)

	pdf.Line(10, pdf.GetY(), 190, pdf.GetY())
	pdf.Ln(5)

	yPosition = pdf.GetY()
	pdf.MultiCell(25, 4, "Исполнитель", "", "L", false)
	pdf.SetXY(pdf.GetX()+25, pdf.GetY())
	pdf.Line(pdf.GetX(), pdf.GetY(), pdf.GetX()+80, pdf.GetY())
	pdf.SetXY(pdf.GetX()+80, yPosition)
	pdf.SetFont("Arial", "", 8)
	pdf.MultiCell(0, 4, "/Бухгалтер/", "", "L", false)

	var buf bytes.Buffer
	if err = pdf.Output(&buf); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *ServiceImpl) getTotlaSumInPdf(priceSum string, totalSum float64) (float64, error) {
	sumPriceStr := strings.ReplaceAll(priceSum, " ", "")
	sumPriceStr = strings.ReplaceAll(sumPriceStr, ",", ".")

	sum, err := strconv.ParseFloat(sumPriceStr, 64)
	if err != nil {
		return 0, err
	}
	totalSum += sum

	return math.Round(totalSum*100) / 100, nil
}

func (s *ServiceImpl) getNumberWords(number float64) string {
	res := s.numberToWords(number)

	numStr := strconv.FormatFloat(number, 'f', 2, 64)
	numStr = numStr[len(numStr)-2:] + " тиын"

	return res + " тенге " + numStr
}

func (s *ServiceImpl) numberToWords(number float64) string {
	var res string
	ones := []string{"", "один", "два", "три", "четыре", "пять", "шесть", "семь", "восемь", "девять", "десять", "одиннадцать", "двенадцать", "тринадцать", "четырнадцать", "пятнадцать", "шестнадцать", "семнадцать", "восемнадцать", "девятнадцать"}
	tens := []string{"", "", "двадцать", "тридцать", "сорок", "пятьдесят", "шестьдесят", "семьдесят", "восемьдесят", "девяносто"}
	hundreds := []string{"", "сто", "двести", "триста", "четыреста", "пятьсот", "шестьсот", "семьсот", "восемсот", "девятсот"}

	intNum := int(number)
	switch {
	case number < 20:
		res = ones[intNum]
	case number < 100:
		res = tens[intNum/10] + " " + ones[intNum%10]
	case number < 1000:
		res = hundreds[intNum/100] + " " + s.numberToWords(float64(intNum%100))
	case number < 1000000:
		thousands := intNum / 1000
		var thousandWord string
		switch {
		case thousands == 1:
			thousandWord = "тысяча"
		case thousands >= 2 && thousands <= 4:
			thousandWord = ones[thousands] + " тысячи"
		default:
			thousandWord = s.numberToWords(float64(thousands)) + " тысяч"
		}
		res = thousandWord + " " + s.numberToWords(float64(intNum%1000))
	case number < 1000000000:
		millions := intNum / 1000000
		var millionWord string
		switch {
		case millions == 1:
			millionWord = "миллион"
		case millions >= 2 && millions <= 4:
			millionWord = ones[millions] + " миллиона"
		default:
			millionWord = s.numberToWords(float64(millions)) + " миллионов"
		}
		res = millionWord + " " + s.numberToWords(float64(intNum%1000000))
	}

	return res
}

func (s *ServiceImpl) selectAndSendPdf(ctx context.Context, pdf []byte, payment models.PaymentXlsx, statusCell string) error {
	message := s.getMessageFromStatusForFilePdf(payment.Status, payment.Month)

	payment.Phone = s.checkAndGetPhoneNumber(payment.Phone)

	switch payment.Status {
	case paymentStatusDontSend:
		if err := s.whatsAppService.SendFilePdf(ctx, payment.Phone, "file.pdf", message, pdf); err != nil {
			log.Err(err).Msg(fmt.Sprintf("Failed to send payment pdf to whatsapp, payment: %s", payment.Number))
			return err
		}
		log.Info().Msgf("send payment to whatsapp, send message and pdf success payment %s", payment.Number)

		if err := s.updateStatusInSheet(statusCell); err != nil {
			log.Err(err).Msg(fmt.Sprintf("Failed to update payment status in sheets, payment: %s\n", payment.Number))
			return err
		}
		log.Info().Msgf("send payment to whatsapp, update status in sheet success payment %s", payment.Number)

	case paymentStatusSend:
		weekday := time.Now().Weekday()
		switch {
		case s.isSendReminderDay(weekday):
			if err := s.whatsAppService.SendMessageFromBaseEnvs(ctx, payment.Phone, message); err != nil {
				log.Err(err).Msg(fmt.Sprintf("Failed to send payment message to whatsapp, payment: %s", payment.Number))
				return err
			}
			log.Info().Msgf("send payment to whatsapp, send remind message success, payment %s", payment.Number)

		case s.isSendPdfDay(weekday):
			if err := s.whatsAppService.SendFilePdf(ctx, payment.Phone, "file.pdf", message, pdf); err != nil {
				log.Err(err).Msg(fmt.Sprintf("Failed to send payment pdf to whatsapp, payment: %s", payment.Number))
				return err
			}
			log.Info().Msgf("send payment to whatsapp, send remind message and pdf success, payment %s", payment.Number)
		}
	}
	return nil
}

func (s *ServiceImpl) getMessageFromStatusForFilePdf(status, month string) string {
	var message string
	switch status {
	case paymentStatusDontSend:
		message = fmt.Sprintf("Добрый день! \nНаправляю вам счет за оплату интеграцию за %s. Напоминаем, что оплата должна быть произведена до 5 числа текущего месяца. \nСпасибо за сотрудничество! \nС уважением, Kwaaka \n \n*Это автоматическое сообщение.*", month)
	case paymentStatusSend:
		message = "Добрый день!\nХочу напомнить вам о необходимости произвести оплату по вашему счету. Просим вас произвести оплату в ближайшее время.\nБлагодарим за внимание и своевременность.\nС уважением, Kwaaka \n\n*Это автоматическое сообщение.*"
	}
	return message
}

func (s *ServiceImpl) checkAndGetPhoneNumber(phone string) string {
	if len(phone) == 0 {
		return phone
	}

	switch phone[0] {
	case '8':
		phone = "+7" + phone[1:]
	default:
		phone = "+" + phone
	}

	phone = strings.ReplaceAll(phone, " ", "")

	return phone
}

func (s *ServiceImpl) isSendReminderDay(weekDay time.Weekday) bool {
	return weekDay == time.Wednesday || weekDay == time.Friday
}

func (s *ServiceImpl) isSendPdfDay(weekDay time.Weekday) bool {
	return weekDay == time.Monday
}

func (s *ServiceImpl) updateStatusInSheet(statusCell string) error {
	value := &sheets.ValueRange{
		Values: [][]interface{}{
			{paymentStatusSend},
		},
	}
	_, err := s.googleSheetsService.Spreadsheets.Values.Update(spreadSheetId, statusCell, value).ValueInputOption("RAW").Do()
	if err != nil {
		return err
	}

	return nil
}
