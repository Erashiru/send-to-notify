package get_modifier_schema_details

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

type Items struct {
	Item []Item `xml:"Item"`
}

type Item struct {
	Ident                string `xml:"Ident,attr"`
	ItemIdent            string `xml:"ItemIdent,attr"`
	SourceIdent          string `xml:"SourceIdent,attr"`
	GUIDString           string `xml:"GUIDString,attr"`
	AssignChildsOnServer string `xml:"AssignChildsOnServer,attr"`
	ActiveHierarchy      string `xml:"ActiveHierarchy,attr"`
	Name                 string `xml:"Name,attr"`
	ReadOnlyName         string `xml:"ReadOnlyName,attr"`
	ModiScheme           string `xml:"ModiScheme,attr"`
	ModiGroup            string `xml:"ModiGroup,attr"`
	Flags                string `xml:"Flags,attr"`
	UpLimit              string `xml:"UpLimit,attr"`
	DownLimit            string `xml:"DownLimit,attr"`
	SortNum              string `xml:"SortNum,attr"`
	DefaultModifier      string `xml:"DefaultModifier,attr"`
	SHQuantity           string `xml:"SHQuantity,attr"`
	FreeCount            string `xml:"FreeCount,attr"`
	AddUntilUpperLimit   string `xml:"AddUntilUpperLimit,attr"`
	ReplaceDefModifier   string `xml:"ReplaceDefModifier,attr"`
	ChangesPrice         string `xml:"ChangesPrice,attr"`
	UseUpLimit           string `xml:"UseUpLimit,attr"`
	UseDownLimit         string `xml:"UseDownLimit,attr"`
	Childs               Childs `xml:"Childs"`
}

type Childs struct {
	ClassName string `xml:"ClassName,attr"`
}
