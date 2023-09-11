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
	return fmt.Sprintf("/idr/%s/", path.Join(paths...))
}

func (idr *IDR) InvestigationComments(inv *Investigation) (comments *InvestigationComments, err error) {
	req := idr.http.R().
		SetError(&APIError{}).
		SetResult(InvestigationComments{}).
		SetHeader("accept-version", "comments-preview").
		SetQueryParam("target", inv.RRN).
		SetQueryParam("size", "100")
	res, err := req.Get(idr.URL("/v1/comments"))
	if err != nil {
		return
	}
	if res.IsError() {
		e := res.Error().(*APIError)
		err = fmt.Errorf("%s: %s", res.Status(), e.Message)
		return
	}
	comments, ok := res.Result().(*InvestigationComments)
	if !ok {
		err = fmt.Errorf("failed to decode comments data")
		return
	}
	return
}

func (idr *IDR) Investigation(id string) (investigation *Investigation, err error) {
	req := idr.http.R()
	req.SetError(&APIError{})
	res, err := req.Get(idr.URL("/v2/investigations", id))
	if err != nil {
		return
	}
	if res.IsError() {
		e := res.Error().(*APIError)
		err = fmt.Errorf("%s: %s", res.Status(), e.Message)
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
	res, err := req.Get(idr.URL("/v2/investigations"))
	if err != nil {
		return
	}
	if res.IsError() {
		e := res.Error().(*APIError)
		err = fmt.Errorf("%s: %s", res.Status(), e.Message)
		return
	}
	var invRes *InvestigationsResponse
	err = json.Unmarshal(res.Body(), &invRes)
	investigations = invRes.Data
	return
}

func (idr *IDR) UpdateInvestigation(id string, update *InvestigationUpdateRequest) (*Investigation, error) {
	req := idr.http.R()
	req.SetError(&APIError{})
	req.SetBody(&update)
	res, err := req.Patch(idr.URL("/v2/investigations", id))
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		e := res.Error().(*APIError)
		err := fmt.Errorf("%s: %s", res.Status(), e.Message)
		return nil, err
	}
	var inv *Investigation
	err = json.Unmarshal(res.Body(), &inv)
	if err != nil {
		return nil, err
	}
	return inv, nil
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
