package legalentity

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"strings"
	"time"
)

func (s *ServiceImpl) saveDocumentInS3(docName, extension, legalEntityID, documentID string, fileContents []byte) (string, error) {
	formatted := strings.Replace(docName, " ", "_", -1)
	name := fmt.Sprintf("%s_%s", formatted, documentID)

	link := strings.TrimSpace(fmt.Sprintf("legal_entity_documents/%s/%s", legalEntityID, name))
	var contentType string
	switch extension {
	case ".pdf":
		contentType = "application/pdf"
	case ".doc":
		contentType = "application/msword"
	case ".docx":
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	default:
		log.Err(errors.New("invalid document type in saveDocumentInS3: " + extension))
	}

	if err := s.s3Service.PutDocument(link, fileContents, s.s3Info.KwaakaFilesBucket, contentType); err != nil {
		return "", err
	}

	resLink := fmt.Sprintf("%s/%s%s", s.s3Info.KwaakaFilesBaseUrl, link, extension)
	return resLink, nil
}

func (s *ServiceImpl) generateDocumentID() int64 {
	currentTimestamp := time.Now().UnixNano() / int64(time.Microsecond)
	uniqueID := uuid.New().ID()

	id := currentTimestamp + int64(uniqueID)

	return id
}
