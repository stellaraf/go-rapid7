package rapid7

import (
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/google/go-querystring/query"
)

type Client struct {
	IDR *IDR
}

type IDR struct {
	BaseURL *url.URL
	http    *resty.Client
}

func (idr *IDR) URL(paths ...string) string {
	base := fmt.Sprintf("/idr/%s/", IDR_VERSION)
	p := path.Join(paths...)
	return base + p
}

func (idr *IDR) Investigation(id string) (investigation *Investigation, err error) {
	req := idr.http.R()
	req.SetError(&APIError{})
	res, err := req.Get(idr.URL("investigations", id))
	if err != nil {
		return
	}
	if res.StatusCode() > 200 {
		e := res.Error().(*APIError)
		err = fmt.Errorf(e.Message)
		return
	}
	err = json.Unmarshal(res.Body(), &investigation)
	return
}

func (idr *IDR) Investigations(q ...*InvestigationsQuery) (investigations []*Investigation, err error) {
	req := idr.http.R()
	req.SetError(&APIError{})
	if len(q) > 0 {
		qe, err := query.Values(q[0])
		if err != nil {
			return nil, err
		}
		req.SetQueryParamsFromValues(qe)
	}
	res, err := req.Get(idr.URL("investigations"))
	if err != nil {
		return
	}
	if res.StatusCode() > 200 {
		e := res.Error().(*APIError)
		err = fmt.Errorf(e.Message)
		return
	}
	var invRes *InvestigationsResponse
	err = json.Unmarshal(res.Body(), &invRes)
	investigations = invRes.Data
	return
}

func newIDR(region, apiKey string) (idr *IDR, err error) {
	h := resty.New()
	urlS := fmt.Sprintf("https://%s.api.insight.rapid7.com", strings.ToLower(region))
	u, err := url.Parse(urlS)
	if err != nil {
		return
	}
	h.SetHeader("X-API-Key", apiKey)
	h.SetHeader("Accept-version", "investigations-preview")
	validate, err := h.R().Get(fmt.Sprintf("https://%s/validate", u.Hostname()))
	if err != nil {
		return
	}
	if validate.StatusCode() > 200 {
		err = fmt.Errorf("failed to authenticate with Rapid7 API")
		return
	}
	h.SetBaseURL(u.String())
	idr = &IDR{
		BaseURL: u,
		http:    h,
	}
	return
}

func New(region, apiKey string) (client *Client, err error) {
	idr, err := newIDR(region, apiKey)
	if err != nil {
		return
	}
	client = &Client{
		IDR: idr,
	}
	return
}
