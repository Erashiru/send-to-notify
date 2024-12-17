package models

type GetSeqNumberRK7QueryResult struct {
	ServerVersion string      `xml:"ServerVersion,attr"`
	XmlVersion    string      `xml:"XmlVersion,attr"`
	NetName       string      `xml:"NetName,attr"`
	Status        string      `xml:"Status,attr"`
	CMD           string      `xml:"CMD,attr"`
	ErrorText     string      `xml:"ErrorText,attr"`
	DateTime      string      `xml:"DateTime,attr"`
	WorkTime      string      `xml:"WorkTime,attr"`
	Processed     string      `xml:"Processed,attr"`
	LicenseInfo   LicenseInfo `xml:"LicenseInfo"`
}

type LicenseInfo struct {
	Anchor          string          `xml:"anchor,attr"`
	LicenseToken    string          `xml:"licenseToken,attrs"`
	LicenseInstance LicenseInstance `xml:"LicenseInstance"`
}

type LicenseInstance struct {
	Guid      string `xml:"guid,attr"`
	SeqNumber string `xml:"seqNumber,attr"`
}
