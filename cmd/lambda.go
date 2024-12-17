package cmd

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"os"
)

func IsLambda() bool {
	if lambdaTaskRoot := os.Getenv("LAMBDA_TASK_ROOT"); lambdaTaskRoot != "" {
		return true
	}
	return false
}

func GetSession() *session.Session {
	region := os.Getenv("REGION")
	if region == "" {
		region = os.Getenv("AWS_REGION")
	}

	session := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))

	return session
}
