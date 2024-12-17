package models

type Document struct {
	ID      string `bson:"id" json:"id"`
	Status  string `bson:"status" json:"status"`
	Type    string `bson:"type" json:"type"`
	DocName string `bson:"name" json:"name"`
	S3Link  string `bson:"s3_link" json:"s3_link"`
	Comment string `bson:"comment,omitempty" json:"comment,omitempty"`
}

type UploadDocumentRequest struct {
	LegalEntityID string
	DocName       string
	DocType       string
	Extension     string
	Data          []byte
	Comment       string
}

type DocumentFilter struct {
	DocType string
}

type S3Info struct {
	KwaakaFilesBucket  string
	KwaakaFilesBaseUrl string
}

type GetDocumentByLegalEntityIDRequest struct {
	Documents []Document `bson:"documents" json:"documents"`
}

type UnwindDocumentsField struct {
	Document Document `bson:"documents"`
}
