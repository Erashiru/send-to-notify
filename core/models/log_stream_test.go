package models

import (
	"fmt"
	"os"
	"testing"
)

func TestLogStream(t *testing.T) {
	os.Setenv("AWS_LAMBDA_LOG_STREAM_NAME", "2023/09/18/[53]fefef3e772aa4ae0967ff54fbd976ad0")
	os.Setenv("AWS_LAMBDA_LOG_GROUP_NAME", "/aws/lambda/wolt-handler")
	os.Setenv("REGION", "eu-west-1")

	var logStream LogStream
	stream := logStream.GetLink()
	logLink := logStream.GetLinkWithPattern("6508db05bc74fe719a055b76")

	if stream != "https://eu-west-1.console.aws.amazon.com/cloudwatch/home?region=eu-west-1#logsV2:log-groups/log-group/$252Faws$252Flambda$252Fwolt-handler/log-events/2023$252F09$252F18$252F$255B53$255Dfefef3e772aa4ae0967ff54fbd976ad0" {
		t.Error(fmt.Errorf("invalid log Stream: %v", stream))
		return
	}

	if logLink != "https://eu-west-1.console.aws.amazon.com/cloudwatch/home?region=eu-west-1#logsV2:log-groups/log-group/$252Faws$252Flambda$252Fwolt-handler/log-events/2023$252F09$252F18$252F$255B53$255Dfefef3e772aa4ae0967ff54fbd976ad0$3FfilterPattern$3D6508db05bc74fe719a055b76" {
		t.Error(fmt.Errorf("invalid log Link: %v", logLink))
		return
	}
}
