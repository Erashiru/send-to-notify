package get_menu_by_category

type CommandResult struct {
	CMD          string       `xml:"CMD,attr"`
	Status       string       `xml:"Status,attr"`
	ErrorText    string       `xml:"ErrorText,attr"`
	DateTime     string       `xml:"DateTime,attr"`
	WorkTime     string       `xml:"WorkTime,attr"`
	RK7Reference RK7Reference `xml:"RK7Reference"`
}

type Item struct {
	Ident                    string       `xml:"Ident,attr"`
	GUIDString               string       `xml:"GUIDString,attr"`
	Code                     string       `xml:"Code,attr"`
	Name                     string       `xml:"Name,attr"`
	Status                   string       `xml:"Status,attr"`
	Parent                   string       `xml:"Parent,attr"`
	ModiScheme               string       `xml:"ModiScheme,attr"`
	CLASSIFICATORGROUPS12288 string       `xml:"CLASSIFICATORGROUPS-12288,attr"`
	GenNutritionalValue      string       `xml:"genNutritionalValue,attr"`
	GenitemParams            string       `xml:"genitemParams,attr"`
	GenName0419              string       `xml:"genName0419,attr"`
	GenDescription0419       string       `xml:"genDescription0419,attr"`
	GenName0409              string       `xml:"genName0409,attr"`
	GenName043f              string       `xml:"genName043f,attr"`
	GenIKPU                  string       `xml:"genIKPU,attr"`
	Genphotolink             string       `xml:"genphotolink,attr"`
	GenhideDish              string       `xml:"genhideDish,attr"`
	RIChildItems             RIChildItems `xml:"RIChildItems"`
}

type TClassificatorGroup struct {
	Text         string `xml:",chardata"`
	GUIDString   string `xml:"GUIDString,attr"`
	Code         string `xml:"Code,attr"`
	Name         string `xml:"Name,attr"`
	Ident        string `xml:"Ident,attr"`
	RIChildItems string `xml:"RIChildItems"`
}

type RIChildItems struct {
	TClassificatorGroup []TClassificatorGroup `xml:"TClassificatorGroup"`
}

type Items struct {
	Item []Item `xml:"Item"`
}

type RK7Reference struct {
	DataVersion    string `xml:"DataVersion,attr"`
	TotalItemCount string `xml:"TotalItemCount,attr"`
	ClassName      string `xml:"ClassName,attr"`
	Items          Items  `xml:"Items"`
	RIChildItems   string `xml:"RIChildItems"`
}

type RK7QueryResult struct {
	ServerVersion   string          `xml:"ServerVersion,attr"`
	XmlVersion      string          `xml:"XmlVersion,attr"`
	NetName         string          `xml:"NetName,attr"`
	Status          string          `xml:"Status,attr"`
	Processed       string          `xml:"Processed,attr"`
	ArrivalDateTime string          `xml:"ArrivalDateTime,attr"`
	CommandResult   []CommandResult `xml:"CommandResult"`
}
