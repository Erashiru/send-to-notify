package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	"github.com/kwaaka-team/orders-core/service/legal_entity_payment/models"
	"io"
	"net/http"
)

// CreateLegalEntityPayment docs
//
//	@Tags		legal entity payment
//	@Title		Method for create legal entity payment
//	@Security	ApiKeyAuth
//	@Summary	Method for create legal entity payment
//	@Param		legal_entity_payment	body		models.LegalEntityPayment	true	"legal_entity_payment"
//	@Success	200						{object}	string
//	@Failure	400						{object}	errors.ErrorResponse
//	@Failure	401						{object}	errors.ErrorResponse
//	@Failure	500						{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/legal-entity-payment/create [post]
func (s *Server) CreateLegalEntityPayment(c *gin.Context) {
	var req models.LegalEntityPayment

	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err.Error())
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	id, err := s.LegalEntityPaymentService.CreatePayment(c.Request.Context(), req)
	if err != nil {
		s.Logger.Errorf("create legal entity payment error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, id)
}

// GetLegalEntityPaymentByID docs
//
//	@Tags		legal entity payment
//	@Title		Method for get legal entity payment by id
//	@Security	ApiKeyAuth
//	@Summary	Method for get legal entity payment by id
//	@Param		legal_entity_payment_id	path		string	true	"legal_entity_payment_id"
//	@Success	200						{object}	models.LegalEntityPayment
//	@Failure	401						{object}	errors.ErrorResponse
//	@Failure	500						{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/legal-entity-payment/{legal_entity_payment_id} [get]
func (s *Server) GetLegalEntityPaymentByID(c *gin.Context) {
	id := c.Param("legal_entity_payment_id")

	legalEntityPayment, err := s.LegalEntityPaymentService.GetPaymentByID(c.Request.Context(), id)
	if err != nil {
		s.Logger.Errorf("get legal entity payment error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, legalEntityPayment)
}

// GetListLegalEntityPayment docs
//
//	@Tags		legal entity payment
//	@Title		Method for get list legal entity payments
//	@Security	ApiKeyAuth
//	@Summary	Method for get list legal entity payments
//	@Param		query	body		models.ListQuery	true	"query"
//	@Param		page	query		string				true	"page"
//	@Param		limit	query		string				true	"limit"
//	@Success	200		{object}	[]models.ListLegalEntityPayment
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/legal-entity-payment/list [post]
func (s *Server) GetListLegalEntityPayment(c *gin.Context) {
	var req models.ListLegalEntityPaymentQuery
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err.Error())
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	page, limit, err := parsePaging(c)
	if err != nil {
		s.Logger.Errorf("page parsing error: %s", err.Error())
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	req.Pagination.Page = page
	req.Pagination.Limit = limit

	res, err := s.LegalEntityPaymentService.GetList(c.Request.Context(), req)
	if err != nil {
		s.Logger.Errorf("get list legal entity payment error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UpdateLegalEntityPayment docs
//
//	@Tags		legal entity payment
//	@Title		Method for update legal entity payment
//	@Security	ApiKeyAuth
//	@Summary	Method for update legal entity payment
//	@Param		legal_entity_payment	body	models.LegalEntityPayment	true	"legal_entity_payment"
//	@Success	204
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/legal-entity-payment/update [put]
func (s *Server) UpdateLegalEntityPayment(c *gin.Context) {
	var req models.LegalEntityPayment

	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err.Error())
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.LegalEntityPaymentService.Update(c.Request.Context(), req.ToUpdateModel()); err != nil {
		s.Logger.Errorf("update legal entity payment error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// DeleteLegalEntityPayment docs
//
//	@Tags		legal entity payment
//	@Title		Method for get legal entity payment by id
//	@Security	ApiKeyAuth
//	@Summary	Method for create legal entity payment by id
//	@Param		legal_entity_payment_id	path	string	true	"legal_entity_payment_id"
//	@Success	204
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/legal-entity-payment/{legal_entity_payment_id} [delete]
func (s *Server) DeleteLegalEntityPayment(c *gin.Context) {
	id := c.Param("legal_entity_payment_id")

	if err := s.LegalEntityPaymentService.Delete(c.Request.Context(), id); err != nil {
		s.Logger.Errorf("delete legal entity payment error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// GetLegalEntityPaymentAnalytics docs
//
//	@Tags		legal entity payment
//	@Title		Method for get legal entity payments analytics
//	@Security	ApiKeyAuth
//	@Summary	Method for get legal entity payments analytics
//	@Param		query	body		models.LegalEntityPaymentAnalyticsRequest	true	"query"
//	@Success	200		{object}	models.LegalEntityPaymentAnalyticsResponse
//	@Failure	400		{object}	errors.ErrorResponse
//	@Failure	401		{object}	errors.ErrorResponse
//	@Failure	500		{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/legal-entity-payment/list [post]
func (s *Server) GetLegalEntityPaymentAnalytics(c *gin.Context) {
	var req models.LegalEntityPaymentAnalyticsRequest

	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err.Error())
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	res, err := s.LegalEntityPaymentService.GetLegalEntityPaymentAnalytics(c.Request.Context(), req)
	if err != nil {
		s.Logger.Errorf("get legal entity payment analytics error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// UploadPDF  docs
//
//	@Tags		legal entity payment
//	@Title		Method for Upload PDF  file for legal entity payment
//	@Security	ApiKeyAuth
//	@Summary	Method Upload PDF  file for legal entity payment
//	@Param		legal_entity_payment_id	query		string	true	"legal_entity_payment_id"
//	@Param		file					formData	file	true	"PDF file"
//	@Success	200						string
//	@Failure	400						{object}	errors.ErrorResponse
//	@Failure	401						{object}	errors.ErrorResponse
//	@Failure	500						{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/legal-entity-payment/upload-pdf [post]
func (s *Server) UploadPDF(c *gin.Context) {

	legalEntityPaymentID, ok := c.GetQuery("legal_entity_payment_id")
	if !ok {
		err := fmt.Errorf("missing legal entity payment id")
		s.Logger.Errorf(errGetQuery, err)
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := c.Request.ParseMultipartForm(0); err != nil {
		s.Logger.Errorf("parse multipart form error: %s", err.Error())
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		s.Logger.Errorf("reading PDF file error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	res, err := s.LegalEntityPaymentService.UploadPDF(c.Request.Context(), models.LegalEntityPaymentDownloadPDFRequest{
		LegalEntityPaymentID: legalEntityPaymentID,
		File:                 data,
	})
	if err != nil {
		s.Logger.Errorf("download PDF file error: %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// CreateBill docs
//
//	@Tags		legal entity payment
//	@Title		Method for create bill for legal entity payment
//	@Security	ApiKeyAuth
//	@Summary	Method for create bill for legal entity payment
//	@Param		query	body	models.LegalEntityPaymentCreateBillRequest	true	"query"
//	@Success	204
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/legal-entity-payment/create-bill [post]
func (s *Server) CreateBill(c *gin.Context) {
	var req models.LegalEntityPaymentCreateBillRequest

	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.LegalEntityPaymentService.CreateBill(c.Request.Context(), req); err != nil {
		s.Logger.Errorf("create bill error %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// ConfirmPayment docs
//
//	@Tags		legal entity payment
//	@Title		Method for confirm payment bill for legal entity payment
//	@Security	ApiKeyAuth
//	@Summary	Method confirm payment bill for legal entity payment
//	@Param		request	body	models.LegalEntityPaymentConfirmPaymentRequest	true	"request"
//	@Success	204
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/v1/kwaaka-admin/legal-entity-payment/confirm-payment [post]
func (s *Server) ConfirmPayment(c *gin.Context) {
	var req models.LegalEntityPaymentConfirmPaymentRequest

	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.JSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.LegalEntityPaymentService.ConfirmPayment(c.Request.Context(), req); err != nil {
		s.Logger.Errorf("confirm payment error %s", err.Error())
		c.JSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

// SendPayment docs
//
//	@Tags		legal entity payment
//	@Title		Method for generate and send payment in whatsapp
//	@Security	ApiKeyAuth
//	@Summary	Method generate and send payment in whatsapp
//	@Success	204
//	@Failure	400	{object}	errors.ErrorResponse
//	@Failure	401	{object}	errors.ErrorResponse
//	@Failure	500	{object}	errors.ErrorResponse
//	@Router		/api/payment/send-in-whatsapp [post]
func (s *Server) SendPayment(c *gin.Context) {
	if err := s.LegalEntityPaymentService.SendPayment(c.Request.Context()); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}
	c.AbortWithStatus(http.StatusNoContent)
}
