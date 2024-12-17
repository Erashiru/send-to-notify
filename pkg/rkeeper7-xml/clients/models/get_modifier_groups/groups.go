package get_modifier_groups

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
	Code                 string `xml:"Code,attr"`
	Name                 string `xml:"Name,attr"`
	AltName              string `xml:"AltName,attr"`
	MainParentIdent      string `xml:"MainParentIdent,attr"`
	Status               string `xml:"Status,attr"`
	VisualTypeImage      string `xml:"VisualType_Image,attr"`
	VisualTypeBColor     string `xml:"VisualType_BColor,attr"`
	VisualTypeTextColor  string `xml:"VisualType_TextColor,attr"`
	VisualTypeFlags      string `xml:"VisualType_Flags,attr"`
	ExtCode              string `xml:"ExtCode,attr"`
	Prompt               string `xml:"Prompt,attr"`
	AltPrompt            string `xml:"AltPrompt,attr"`
	Flags                string `xml:"Flags,attr"`
	ModiGroupType        string `xml:"ModiGroupType,attr"`
	CommonModifier       string `xml:"CommonModifier,attr"`
	EditRight            string `xml:"EditRight,attr"`
	Parent               string `xml:"Parent,attr"`
	Childs               Childs `xml:"Childs"`
}

type Childs struct {
	ClassName string  `xml:"ClassName,attr"`
	Child     []Child `xml:"Child"`
}

type Child struct {
	ChildIdent string `xml:"ChildIdent,attr"`
	IsTerminal string `xml:"IsTerminal,attr"`
}
