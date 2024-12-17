package http

import (
	"crypto/tls"
	"github.com/go-resty/resty/v2"
	"github.com/kwaaka-team/orders-core/pkg"
	"github.com/kwaaka-team/orders-core/pkg/rkeeper7-xml/clients"
	"github.com/rs/zerolog/log"
	"time"
)

type client struct {
	cli                    *resty.Client
	token                  string
	quit                   chan struct{}
	username               string
	password               string
	ucsUsername            string
	ucsPassword            string
	licenseBaseURL         string
	anchor                 string
	objectID               string
	licenseToken           string
	stationID              string
	stationCode            string
	licenseInstanceGUID    string
	childItems             int
	classificatorItemIdent int
	classificatorPropMask  string
	menuItemsPropMask      string
	propFilter             string
	cashier                string
}

func New(conf *clients.Config) (*client, error) {
	//transport := &http.Transport{
	//	Proxy: http.ProxyURL(&url.URL{
	//		Scheme: "http",
	//		User:   url.UserPassword("max0Q8ga", "mAxsWX6F2Y"),
	//		Host:   "91.147.127.252:50100",
	//	}),
	//}

	cli := resty.New().
		SetBaseURL(conf.BaseURL).
		SetBasicAuth(conf.Username, conf.Password).
		SetRetryCount(0).
		SetTimeout(35 * time.Second).
		SetHeaders(map[string]string{
			pkg.ContentTypeHeader: pkg.XMLType,
		}).SetTLSClientConfig(&tls.Config{InsecureSkipVerify: true, MinVersion: tls.VersionTLS10}).
		SetProxy("http://max0Q8ga:mAxsWX6F2Y@185.120.79.119:50100")

	log.Info().Msgf("client: %+v", cli)
	log.Info().Msgf("client request: %+v", cli.R().RawRequest)

	c := &client{
		cli:                    cli,
		quit:                   make(chan struct{}),
		username:               conf.Username,
		password:               conf.Password,
		ucsUsername:            conf.UCSUsername,
		ucsPassword:            conf.UCSPassword,
		licenseBaseURL:         conf.LicenseBaseURL,
		token:                  conf.Token,
		anchor:                 conf.Anchor,
		objectID:               conf.ObjectID,
		stationID:              conf.StationID,
		stationCode:            conf.StationCode,
		licenseInstanceGUID:    conf.LicenseInstanceGUID,
		childItems:             conf.ChildItems,
		classificatorItemIdent: conf.ClassificatorItemIdent,
		classificatorPropMask:  conf.ClassificatorPropMask,
		menuItemsPropMask:      conf.MenuItemsPropMask,
		propFilter:             conf.PropFilter,
		cashier:                conf.Cashier,
	}

	return c, nil
}
