package models

type PriceScale struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type TradeGroup struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}

type Dishes struct {
	Item []DishItem `xml:"Item"`
}

type DishItem struct {
	Ident    string `xml:"Ident,attr"`
	Price    string `xml:"Price,attr"`
	Quantity string `xml:"quantity,attr"`
}

type Modifiers struct {
	Item []ModifierItem `xml:"Item"`
}

type ModifierItem struct {
	Ident string `xml:"Ident,attr"`
	ID    string `xml:"ID,attr"`
	Price string `xml:"Price,attr"`
}

type OrderMenuRK7QueryResult struct {
	ServerVersion string     `xml:"ServerVersion,attr"`
	XmlVersion    string     `xml:"XmlVersion,attr"`
	NetName       string     `xml:"NetName,attr"`
	Status        string     `xml:"Status,attr"`
	CMD           string     `xml:"CMD,attr"`
	ErrorText     string     `xml:"ErrorText,attr"`
	DateTime      string     `xml:"DateTime,attr"`
	WorkTime      string     `xml:"WorkTime,attr"`
	Processed     string     `xml:"Processed,attr"`
	PriceScale    PriceScale `xml:"PriceScale"`
	TradeGroup    TradeGroup `xml:"TradeGroup"`
	Dishes        Dishes     `xml:"Dishes"`
	Modifiers     Modifiers  `xml:"Modifiers"`
}
