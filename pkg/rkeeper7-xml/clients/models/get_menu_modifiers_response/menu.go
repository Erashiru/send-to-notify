package get_menu_modifiers_response

type Item struct {
	Ident                string `xml:"Ident,attr"`
	ItemIdent            string `xml:"ItemIdent,attr"`
	SourceIdent          string `xml:"SourceIdent,attr"`
	GUIDString           string `xml:"GUIDString,attr"`
	AssignChildsOnServer string `xml:"AssignChildsOnServer,attr"`
	ActiveHierarchy      string `xml:"ActiveHierarchy,attr"`
	Code                 string `xml:"Code,attr"`
	Name                 string `xml:"Name,attr"`
	AltName              string `xml:"AltName,attr"`
	MainParentIdent      string `xml:"MainParentIdent,attr"`
	Status               string `xml:"Status,attr"`
	VisualTypeImage      string `xml:"VisualType_Image,attr"`
	VisualTypeBColor     string `xml:"VisualType_BColor,attr"`
	VisualTypeTextColor  string `xml:"VisualType_TextColor,attr"`
	VisualTypeFlags      string `xml:"VisualType_Flags,attr"`
	SalesTermsFlag       string `xml:"SalesTerms_Flag,attr"`
	SalesTermsStartSale  string `xml:"SalesTerms_StartSale,attr"`
	SalesTermsStopSale   string `xml:"SalesTerms_StopSale,attr"`
	RightLvl             string `xml:"RightLvl,attr"`
	AvailabilitySchedule string `xml:"AvailabilitySchedule,attr"`
	UseStartSale         string `xml:"UseStartSale,attr"`
	UseStopSale          string `xml:"UseStopSale,attr"`
	SortNum              string `xml:"SortNum,attr"`
	Parent               string `xml:"Parent,attr"`
	ExtCode              string `xml:"ExtCode,attr"`
	ShortName            string `xml:"ShortName,attr"`
	AltShortName         string `xml:"AltShortName,attr"`
	MaxOneDish           string `xml:"MaxOneDish,attr"`
	Flags                string `xml:"Flags,attr"`
	Comment              string `xml:"Comment,attr"`
	Format               string `xml:"Format,attr"`
	Weight               string `xml:"Weight,attr"`
	Kurs                 string `xml:"Kurs,attr"`
	Dish                 string `xml:"Dish,attr"`
	MInterface           string `xml:"MInterface,attr"`
	LargeImagePath       string `xml:"LargeImagePath,attr"`
	SaveInCheck          string `xml:"SaveInCheck,attr"`
	AddToName            string `xml:"AddToName,attr"`
	ReplaceName          string `xml:"ReplaceName,attr"`
	InputName            string `xml:"InputName,attr"`
	AddMenuItemPrice     string `xml:"AddMenuItemPrice,attr"`
	UseLimitedQnt        string `xml:"UseLimitedQnt,attr"`
	UseFormatInput       string `xml:"UseFormatInput,attr"`
	UseKurs              string `xml:"UseKurs,attr"`
	PrintInfoName        string `xml:"PrintInfoName,attr"`
	Childs               Childs `xml:"Childs"`
}

type Childs struct {
	ClassName string `xml:"ClassName,attr"`
}

type Items struct {
	Item []Item `xml:"Item"`
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
