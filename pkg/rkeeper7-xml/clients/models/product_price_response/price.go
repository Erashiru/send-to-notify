package product_price_response

// RK7QueryResult was generated 2024-03-02 00:45:05 by https://xml-to-go.github.io/ in Ukraine.
type RK7QueryResult struct {
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

type Items struct {
	Item []Item `xml:"Item"`
}

type Item struct {
	Ident       string `xml:"Ident,attr"`
	ItemIdent   string `xml:"ItemIdent,attr"`
	SourceIdent string `xml:"SourceIdent,attr"`
	GUIDString  string `xml:"GUIDString,attr"`
	ObjectID    string `xml:"ObjectID,attr"`
	Species     string `xml:"Species,attr"`
	PriceType   string `xml:"PriceType,attr"`
	Modifiers   string `xml:"Modifiers,attr"`
	Value       string `xml:"Value,attr"`
}

type RK7Reference struct {
	DataVersion    string `xml:"DataVersion,attr"`
	ClassName      string `xml:"ClassName,attr"`
	Name           string `xml:"Name,attr"`
	MinIdent       string `xml:"MinIdent,attr"`
	MaxIdent       string `xml:"MaxIdent,attr"`
	ViewRight      string `xml:"ViewRight,attr"`
	UpdateRight    string `xml:"UpdateRight,attr"`
	ChildRight     string `xml:"ChildRight,attr"`
	DeleteRight    string `xml:"DeleteRight,attr"`
	XMLExport      string `xml:"XMLExport,attr"`
	XMLMask        string `xml:"XMLMask,attr"`
	TotalItemCount string `xml:"TotalItemCount,attr"`
	Items          Items  `xml:"Items"`
}
