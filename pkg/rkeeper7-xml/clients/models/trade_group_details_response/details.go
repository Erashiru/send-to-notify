package trade_group_details_response

type Items struct {
	Item []Item `xml:"Item"`
}

type Childs struct {
	ClassName string `xml:"ClassName,attr"`
}

type Item struct {
	Ident                string `xml:"Ident,attr"`
	ItemIdent            string `xml:"ItemIdent,attr"`
	SourceIdent          string `xml:"SourceIdent,attr"`
	GUIDString           string `xml:"GUIDString,attr"`
	AssignChildsOnServer string `xml:"AssignChildsOnServer,attr"`
	TradeObject          string `xml:"TradeObject,attr"`
	ObjectSifr           string `xml:"ObjectSifr,attr"`
	RefCollName          string `xml:"refCollName,attr"`
	Parent               string `xml:"Parent,attr"`
	Flag                 string `xml:"Flag,attr"`
	Childs               Childs `xml:"Childs"`
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
