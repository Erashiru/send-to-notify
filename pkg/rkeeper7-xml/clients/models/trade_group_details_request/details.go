package trade_group_details_request

type RK7CMD struct {
	CMD         string      `xml:"CMD,attr"`
	RefName     string      `xml:"RefName,attr"`
	Parent      string      `xml:"Parent,attr"`
	PROPFILTERS PropFilters `xml:"PROPFILTERS"`
}

type PropFilters struct {
	PROPFILTER PropFilter `xml:"PROPFILTER"`
}

type PropFilter struct {
	Name  string `xml:"Name,attr"`
	Value string `xml:"Value,attr"`
}

type RK7Query struct {
	RK7CMD RK7CMD `xml:"RK7CMD"`
}
