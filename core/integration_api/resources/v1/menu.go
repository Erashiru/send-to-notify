package v1

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kwaaka-team/orders-core/core/errors"
	menuModels "github.com/kwaaka-team/orders-core/core/menu/models"
	"net/http"
)

func (s *Server) AddNameInProduct(c *gin.Context) {
	menuID := c.Param("menu_id")
	productID := c.Param("product_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.AddNameInProduct(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  productID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("add name in product error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) AddDescriptionInProduct(c *gin.Context) {
	menuID := c.Param("menu_id")
	productID := c.Param("product_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.AddDescriptionInProduct(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  productID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("add description in product error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) AddNameInSection(c *gin.Context) {
	menuID := c.Param("menu_id")
	sectionID := c.Param("section_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.AddNameInSection(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  sectionID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("add name in section error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) AddDescriptionInSection(c *gin.Context) {
	menuID := c.Param("menu_id")
	sectionID := c.Param("section_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.AddDescriptionInSection(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  sectionID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("add descrition in section error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) AddNameInAttributeGroup(c *gin.Context) {
	menuID := c.Param("menu_id")
	attrGroupID := c.Param("attribute_group_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.AddNameInAttributeGroup(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  attrGroupID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("add name in attribute group error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) AddNameInAttribute(c *gin.Context) {
	menuID := c.Param("menu_id")
	attributeID := c.Param("attribute_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.AddNameInAttribute(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  attributeID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("add name in attribute error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) ChangeNameInProduct(c *gin.Context) {
	menuID := c.Param("menu_id")
	productID := c.Param("product_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.ChangeNameInProduct(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  productID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("change name in product error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) ChangeDescriptionInProduct(c *gin.Context) {
	menuID := c.Param("menu_id")
	productID := c.Param("product_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.ChangeDescriptionInProduct(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  productID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("change description in product error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) ChangeNameInSection(c *gin.Context) {
	menuID := c.Param("menu_id")
	sectionID := c.Param("section_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.ChangeNameInSection(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  sectionID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("change name in section error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) ChangeDescriptionInSection(c *gin.Context) {
	menuID := c.Param("menu_id")
	sectionID := c.Param("section_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.ChangeDescriptionInSection(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  sectionID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("change descrition in section error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) ChangeNameInAttributeGroup(c *gin.Context) {
	menuID := c.Param("menu_id")
	attrGroupID := c.Param("attribute_group_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.ChangeNameInAttributeGroup(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  attrGroupID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("change name in attribute group error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) ChangeNameInAttribute(c *gin.Context) {
	menuID := c.Param("menu_id")
	attributeID := c.Param("attribute_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.LanguageDescription
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Errorf(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.ChangeNameInAttribute(c.Request.Context(), menuModels.AddLanguageDescriptionRequest{
		MenuID:    menuID,
		PosMenuID: posMenuID,
		ObjectID:  attributeID,
		Request:   req,
	}); err != nil {
		s.Logger.Errorf("change name in attribute error: %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) AddRegulatoryInformation(c *gin.Context) {
	menuID := c.Param("menu_id")
	productID := c.Param("product_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(errGetQuery, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.RegulatoryInformationValues
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Error(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.AddRegulatoryInformation(c.Request.Context(), menuModels.RegulatoryInformationRequest{
		MenuID:                menuID,
		ProductID:             productID,
		PosMenuID:             posMenuID,
		RegulatoryInformation: req,
	}); err != nil {
		s.Logger.Errorf("add regulatory information error %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}

func (s *Server) ChangeRegulatoryInformation(c *gin.Context) {
	menuID := c.Param("menu_id")
	productID := c.Param("product_id")
	posMenuID, ok := c.GetQuery("pos_menu_id")
	if !ok {
		err := fmt.Errorf("pos menu id query is missing")
		s.Logger.Error(errGetQuery, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	var req menuModels.RegulatoryInformationValues
	if err := c.BindJSON(&req); err != nil {
		s.Logger.Error(errBindBody, err)
		c.AbortWithStatusJSON(http.StatusBadRequest, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	if err := s.menuService.ChangeRegulatoryInformation(c.Request.Context(), menuModels.RegulatoryInformationRequest{
		MenuID:                menuID,
		ProductID:             productID,
		PosMenuID:             posMenuID,
		RegulatoryInformation: req,
	}); err != nil {
		s.Logger.Errorf("change regulatory information error %s", err.Error())
		c.AbortWithStatusJSON(http.StatusInternalServerError, errors.ErrorResponse{Msg: err.Error()})
		return
	}

	c.AbortWithStatus(http.StatusNoContent)
}
