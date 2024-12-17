package aws_s3

import (
	"bytes"
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/rs/zerolog/log"
	"io"
	"strings"
)

type Service interface {
	PutObjectFromStruct(link string, body interface{}, bucketName, contentType string) error
	PutObjectFromBytes(link string, bytes []byte, bucketName, contentType string) error
	GetObject(bucketName string, fileKey string) ([]byte, error)
	RemoveNonAlphaNumericSymbols(s string) string
	ReduceSpaces(s string) string
	PutPDF(link string, body []byte, bucketName, contentType string) error
	PutDocument(link string, body []byte, bucketName, contentType string) error
}

type ServiceImpl struct {
	Sv3 *s3.S3
}

func NewS3Service(session *session.Session) *ServiceImpl {
	return &ServiceImpl{
		Sv3: s3.New(session),
	}
}

func (svc *ServiceImpl) getExtension(contentType string) string {
	fileExtension := map[string]string{
		"application/json":   ".json",
		"text/csv":           ".csv",
		"application/pdf":    ".pdf",
		"application/msword": ".doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
	}

	return fileExtension[contentType]
}

func (svc *ServiceImpl) GetObject(bucketName string, fileKey string) ([]byte, error) {
	res, err := svc.Sv3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(fileKey),
	})
	if err != nil {
		return nil, err
	}

	configBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return configBytes, nil
}

func (svc *ServiceImpl) PutObjectFromBytes(link string, bytes []byte, bucketName, contentType string) error {
	_, err := svc.Sv3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(svc.name(link) + svc.getExtension(contentType)),
		Body:        strings.NewReader(string(bytes)),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Err(err).Msgf("s3 load error")
		return err
	}

	return nil
}

func (svc *ServiceImpl) PutObjectFromStruct(link string, body interface{}, bucketName, contentType string) error {
	bytes, err := json.Marshal(body)
	if err != nil {
		log.Err(err).Msgf("marshal body")
		return err
	}

	_, err = svc.Sv3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(svc.name(link) + svc.getExtension(contentType)),
		Body:        strings.NewReader(string(bytes)),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Err(err).Msgf("s3 load error")
		return err
	}

	return nil
}

func (svc *ServiceImpl) name(img string) string {
	raw := strings.TrimPrefix(string(img), "s3://")
	i := strings.Index(raw, "/")
	if i == -1 {
		return string(img)
	}
	return raw[i+1:]
}

func (svc *ServiceImpl) RemoveNonAlphaNumericSymbols(s string) string {
	var result string
	for _, char := range s {
		if svc.isAlphanumeric(char) || char == ' ' || svc.isCyrillic(char) {
			result += string(char)
		}
	}

	return result
}

func (svc *ServiceImpl) isAlphanumeric(char rune) bool {
	return (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9')
}

func (svc *ServiceImpl) isCyrillic(char rune) bool {
	return (char >= 'А' && char <= 'я') || char == 'Ё' || char == 'ё'
}

func (svc *ServiceImpl) ReduceSpaces(s string) string {
	words := strings.Fields(s)
	result := strings.Join(words, " ")

	return result
}

func (svc *ServiceImpl) PutPDF(link string, body []byte, bucketName, contentType string) error {
	return svc.PutDocument(link, body, bucketName, contentType)
}

func (svc *ServiceImpl) PutDocument(link string, body []byte, bucketName, contentType string) error {
	_, err := svc.Sv3.PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(link + svc.getExtension(contentType)),
		Body:        bytes.NewReader(body),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		log.Err(err).Msgf("s3 load error")
		return err
	}

	return nil
}
