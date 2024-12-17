package models

import (
	"fmt"
	"os"
)

const (
	link         = "https://%s.console.aws.amazon.com/cloudwatch/home?region=%s#logsV2:log-groups/log-group/%s/log-events/"
	logStreamEnv = "AWS_LAMBDA_LOG_STREAM_NAME"
	logGroupEnv  = "AWS_LAMBDA_LOG_GROUP_NAME"
	filterSuffix = "$3FfilterPattern$3D"
)

var (
	match = map[rune]string{
		'/': "$252F",
		'[': "$255B",
		']': "$255D",
	}
)

type LogStream struct{}

func (ls LogStream) GetValue() string {
	return os.Getenv(logStreamEnv)
}

func (ls LogStream) GetLogGroup() string {
	return os.Getenv(logGroupEnv)
}

func (ls LogStream) convertToLogFormat(text string) string {
	var result string

	for _, ch := range text {
		if val, ok := match[ch]; ok {
			result += val
			continue
		}
		result += string(ch)
	}

	return result
}

func (ls LogStream) GetLink() string {
	region := os.Getenv(REGION)
	return fmt.Sprintf(link, region, region, ls.convertToLogFormat(ls.GetLogGroup())) + ls.convertToLogFormat(ls.GetValue())
}

func (ls LogStream) GetLinkWithPattern(patterns ...string) string {
	var filter string

	for index, pattern := range patterns {
		filter += pattern

		if index != len(patterns)-1 {
			filter += "+"
		}
	}

	return ls.GetLink() + filterSuffix + filter
}
