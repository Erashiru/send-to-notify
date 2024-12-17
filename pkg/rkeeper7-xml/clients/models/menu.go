package models

import "github.com/aws/smithy-go/encoding/xml"

type GetMenuRequest struct {
	XMLName xml.Name `xml:"RK7Query"`
	RK7CMD  RK7CMD   `xml:"RK7CMD"`
}

type RK7CMD struct {
	CMD        string `xml:"CMD,attr"`
	RefName    string `xml:"RefName,attr"`
	OnlyActive string `xml:"onlyActive,attr"`
}

type Items struct {
	Item []Item `xml:"Item"`
}

type Item struct {
	Ident                string               `xml:"Ident,attr"`
	ItemIdent            string               `xml:"ItemIdent,attr"`
	SourceIdent          string               `xml:"SourceIdent,attr"`
	GUIDString           string               `xml:"GUIDString,attr"`
	AssignChildsOnServer string               `xml:"AssignChildsOnServer,attr"`
	ActiveHierarchy      string               `xml:"ActiveHierarchy,attr"`
	Code                 string               `xml:"Code,attr"`
	Name                 string               `xml:"Name,attr"`
	AltName              string               `xml:"AltName,attr"`
	MainParentIdent      string               `xml:"MainParentIdent,attr"`
	Status               string               `xml:"Status,attr"`
	VisualTypeImage      string               `xml:"VisualType_Image,attr"`
	VisualTypeBColor     string               `xml:"VisualType_BColor,attr"`
	VisualTypeTextColor  string               `xml:"VisualType_TextColor,attr"`
	VisualTypeFlags      string               `xml:"VisualType_Flags,attr"`
	SalesTermsFlag       string               `xml:"SalesTerms_Flag,attr"`
	SalesTermsStartSale  string               `xml:"SalesTerms_StartSale,attr"`
	SalesTermsStopSale   string               `xml:"SalesTerms_StopSale,attr"`
	RightLvl             string               `xml:"RightLvl,attr"`
	AvailabilitySchedule string               `xml:"AvailabilitySchedule,attr"`
	UseStartSale         string               `xml:"UseStartSale,attr"`
	UseStopSale          string               `xml:"UseStopSale,attr"`
	TaxDishType          string               `xml:"TaxDishType,attr"`
	FutureTaxDishType    string               `xml:"FutureTaxDishType,attr"`
	Parent               string               `xml:"Parent,attr"`
	ExtCode              string               `xml:"ExtCode,attr"`
	ShortName            string               `xml:"ShortName,attr"`
	AltShortName         string               `xml:"AltShortName,attr"`
	PortionWeight        string               `xml:"PortionWeight,attr"`
	PortionName          string               `xml:"PortionName,attr"`
	AltPortion           string               `xml:"AltPortion,attr"`
	Kurs                 string               `xml:"Kurs,attr"`
	QntDecDigits         string               `xml:"QntDecDigits,attr"`
	ModiScheme           string               `xml:"ModiScheme,attr"`
	ComboScheme          string               `xml:"ComboScheme,attr"`
	ModiWeight           string               `xml:"ModiWeight,attr"`
	CookMins             string               `xml:"CookMins,attr"`
	Comment              string               `xml:"Comment,attr"`
	Instruct             string               `xml:"Instruct,attr"`
	Flags                string               `xml:"Flags,attr"`
	TaraWeight           string               `xml:"TaraWeight,attr"`
	ConfirmQnt           string               `xml:"ConfirmQnt,attr"`
	MInterface           string               `xml:"MInterface,attr"`
	MinRestQnt           string               `xml:"MinRestQnt,attr"`
	BarCodes             string               `xml:"BarCodes,attr"`
	PriceMode            string               `xml:"PriceMode,attr"`
	PartsPerPackage      string               `xml:"PartsPerPackage,attr"`
	OpenPrice            string               `xml:"OpenPrice,attr"`
	DontPack             string               `xml:"DontPack,attr"`
	ChangeQntOnce        string               `xml:"ChangeQntOnce,attr"`
	AllowPurchasing      string               `xml:"AllowPurchasing,attr"`
	UseRestControl       string               `xml:"UseRestControl,attr"`
	UseConfirmQnt        string               `xml:"UseConfirmQnt,attr"`
	LabeledProduct       string               `xml:"LabeledProduct,attr"`
	CategPath            string               `xml:"CategPath,attr"`
	SaleObjectType       string               `xml:"SaleObjectType,attr"`
	ComboJoinMode        string               `xml:"ComboJoinMode,attr"`
	ComboSplitMode       string               `xml:"ComboSplitMode,attr"`
	AddLineMode          string               `xml:"AddLineMode,attr"`
	ChangeToCombo        string               `xml:"ChangeToCombo,attr"`
	GuestsDishRating     string               `xml:"GuestsDishRating,attr"`
	RateType             string               `xml:"RateType,attr"`
	MinimumTarifTime     string               `xml:"MinimumTarifTime,attr"`
	MaximumTarifTime     string               `xml:"MaximumTarifTime,attr"`
	IgnoredTarifTime     string               `xml:"IgnoredTarifTime,attr"`
	MinTarifAmount       string               `xml:"MinTarifAmount,attr"`
	MaxTarifAmount       string               `xml:"MaxTarifAmount,attr"`
	RoundTime            string               `xml:"RoundTime,attr"`
	TariffRoundRule      string               `xml:"TariffRoundRule,attr"`
	MoneyRoundRule       string               `xml:"MoneyRoundRule,attr"`
	DefTarifTimeLimit    string               `xml:"DefTarifTimeLimit,attr"`
	ComboDiscount        string               `xml:"ComboDiscount,attr"`
	LargeImagePath       string               `xml:"LargeImagePath,attr"`
	HighLevelGroup1      string               `xml:"HighLevelGroup1,attr"`
	HighLevelGroup2      string               `xml:"HighLevelGroup2,attr"`
	HighLevelGroup3      string               `xml:"HighLevelGroup3,attr"`
	HighLevelGroup4      string               `xml:"HighLevelGroup4,attr"`
	BarCodesText         string               `xml:"BarCodesText,attr"`
	BarcodesFullInfo     string               `xml:"BarcodesFullInfo,attr"`
	ItemKind             string               `xml:"ItemKind,attr"`
	RecommendedMenuItems RecommendedMenuItems `xml:"RecommendedMenuItems"`
	Childs               Childs               `xml:"Childs"`
}

type RecommendedMenuItems struct {
	ClassName string `xml:"ClassName,attr"`
	Items     string `xml:"Items"`
}

type Childs struct {
	ClassName string `xml:"ClassName,attr"`
}

type RK7Reference struct {
	DataVersion         string `xml:"DataVersion,attr"`
	ClassName           string `xml:"ClassName,attr"`
	Name                string `xml:"Name,attr"`
	MinIdent            string `xml:"MinIdent,attr"`
	MaxIdent            string `xml:"MaxIdent,attr"`
	ViewRight           string `xml:"ViewRight,attr"`
	UpdateRight         string `xml:"UpdateRight,attr"`
	ChildRight          string `xml:"ChildRight,attr"`
	DeleteRight         string `xml:"DeleteRight,attr"`
	XMLExport           string `xml:"XMLExport,attr"`
	XMLMask             string `xml:"XMLMask,attr"`
	LeafCollectionCount string `xml:"LeafCollectionCount,attr"`
	TotalItemCount      string `xml:"TotalItemCount,attr"`
	Items               Items  `xml:"Items"`
}

type MenuRK7QueryResult struct {
	ServerVersion   string       `xml:"ServerVersion,attr"`
	XmlVersion      string       `xml:"XmlVersion,attr"`
	NetName         string       `xml:"NetName,attr"`
	Status          string       `xml:"Status,attr"`
	CMD             string       `xml:"CMD,attr"`
	ErrorText       string       `xml:"ErrorText,attr"`
	DateTime        string       `xml:"DateTime,attr"`
	WorkTime        string       `xml:"WorkTime,attr"`
	Processed       string       `xml:"Processed,attr"`
	ArrivalDateTime string       `xml:"ArrivalDateTime,attr"`
	RK7Reference    RK7Reference `xml:"RK7Reference"`
}
