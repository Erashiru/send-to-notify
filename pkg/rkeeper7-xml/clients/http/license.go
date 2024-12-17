package http

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients/models"
	"github.com/rs/zerolog/log"
	"net/http"
	"strings"
	"time"
)

func licenceToken(username, password, token string) string {
	result := username + ";"

	center := md5.New()
	center.Write([]byte(username + password))

	result += hex.EncodeToString(center.Sum(nil)) + ";"

	last := md5.New()
	last.Write([]byte(token))

	result += hex.EncodeToString(last.Sum(nil))

	sEnc := base64.StdEncoding.EncodeToString([]byte(result))
	return sEnc
}

func (rkeeper *client) SetLicense(ctx context.Context) (models.LicenseResponse, error) {
	//httpClient := httpClient.NewHTTPClient(rkeeper.licenseBaseURL)

	path := fmt.Sprintf("/ls5api/api/License/GetLicenseIdByAnchor?anchor=%s", "6%3A"+rkeeper.anchor+"%23"+rkeeper.objectID+"/17")

	cli := resty.New().
		SetBaseURL(rkeeper.licenseBaseURL).
		SetRetryCount(1).
		SetTimeout(15 * time.Second).
		SetHeaders(map[string]string{"usr": licenceToken(rkeeper.ucsUsername, rkeeper.ucsPassword, rkeeper.token)}).
		SetProxy("http://max0Q8ga:mAxsWX6F2Y@185.120.79.119:50100")

	// log.Infof(ctx, "client: %+v", cli)
	// log.Infof(ctx, "client request: %+v", cli.R().RawRequest)

	var result models.LicenseResponse

	response, err := cli.R().
		SetContext(ctx).
		EnableTrace().
		SetResult(&result).
		Get(path)
	if err != nil {
		return models.LicenseResponse{}, err
	}

	log.Info().Msgf("response: %v", response)

	if response.IsError() {
		return models.LicenseResponse{}, fmt.Errorf("error response: %s", response.Error())
	}
	//_, body, err := httpClient.Get(fmt.Sprintf("/ls5api/api/License/GetLicenseIdByAnchor?anchor=%s", "6%3A"+rkeeper.anchor+"%23"+rkeeper.objectID+"/17"), nil, map[string]string{
	//	"usr": licenceToken(rkeeper.ucsUsername, rkeeper.ucsPassword, rkeeper.token),
	//})
	//if err != nil {
	//	return models.LicenseResponse{}, err
	//}
	//
	//var result models.LicenseResponse
	//
	//if err = json.Unmarshal(body, &result); err != nil {
	//	return models.LicenseResponse{}, err
	//}

	rkeeper.licenseToken = result.Id

	return result, nil
}

func (rkeeper *client) GetSeqNumber(ctx context.Context) (models.GetSeqNumberRK7QueryResult, error) {

	path := "/rk7api/v0/xmlinterface.xml"

	var (
		result  models.GetSeqNumberRK7QueryResult
		request = fmt.Sprintf(`
			<RK7Query>
			<RK7CMD CMD="GetXMLLicenseInstanceSeqNumber">
			<LicenseInfo anchor="6:%s#%s/17" licenseToken="%s">
			<LicenseInstance guid="%s"/>
        	</LicenseInfo>
			</RK7CMD>
			</RK7Query>
		`, rkeeper.anchor, rkeeper.objectID, rkeeper.licenseToken, rkeeper.licenseInstanceGUID)
	)

	response, err := rkeeper.cli.R().
		SetContext(ctx).
		EnableTrace().
		SetBody(strings.NewReader(request)).
		SetResult(&result).
		Post(path)
	if err != nil {
		return models.GetSeqNumberRK7QueryResult{}, err
	}

	log.Info().Msgf("get seq number request url: %s, body: %+v", response.Request.URL, response.Request.Body)

	if response.IsError() || response.StatusCode() >= http.StatusBadRequest {
		return models.GetSeqNumberRK7QueryResult{}, fmt.Errorf("get seq number response error: %v", response.Error())
	}

	log.Info().Msgf("get seq number response status: %d, body: %+v", response.RawResponse.StatusCode, string(response.Body()))

	return result, nil
}
