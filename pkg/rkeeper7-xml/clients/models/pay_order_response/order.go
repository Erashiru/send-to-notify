package pay_order_response

type RK7QueryResult struct {
	ServerVersion   string     `xml:"ServerVersion,attr"`
	XmlVersion      string     `xml:"XmlVersion,attr"`
	NetName         string     `xml:"NetName,attr"`
	Status          string     `xml:"Status,attr"`
	CMD             string     `xml:"CMD,attr"`
	ErrorText       string     `xml:"ErrorText,attr"`
	DateTime        string     `xml:"DateTime,attr"`
	WorkTime        string     `xml:"WorkTime,attr"`
	Processed       string     `xml:"Processed,attr"`
	ArrivalDateTime string     `xml:"ArrivalDateTime,attr"`
	PrintCheck      PrintCheck `xml:"PrintCheck"`
}

type PrintCheck struct {
	Uni       string `xml:"uni,attr"`
	LineGuid  string `xml:"line_guid,attr"`
	State     string `xml:"state,attr"`
	Amount    string `xml:"amount,attr"`
	CheckNum  string `xml:"CheckNum,attr"`
	Deleted   string `xml:"deleted,attr"`
	PrintTime string `xml:"printTime,attr"`
	StartTime string `xml:"startTime,attr"`
	Author    Author `xml:"Author"`
	Pay       Pay    `xml:"Pay"`
}

type Author struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
	Role Role   `xml:"Role"`
}

type Pay struct {
	ID       string `xml:"id,attr"`
	Code     string `xml:"code,attr"`
	Name     string `xml:"name,attr"`
	Uni      string `xml:"uni,attr"`
	LineGuid string `xml:"line_guid,attr"`
	State    string `xml:"state,attr"`
	Guid     string `xml:"guid,attr"`
	Amount   string `xml:"amount,attr"`
	BasicSum string `xml:"basicSum,attr"`
}

type Role struct {
	ID   string `xml:"id,attr"`
	Code string `xml:"code,attr"`
	Name string `xml:"name,attr"`
}
