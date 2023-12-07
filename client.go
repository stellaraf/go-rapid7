package rapid7

import (
	"encoding/json"
	"errors"
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
	return fmt.Sprintf("/idr/%s", path.Join(paths...))
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
		e, ok := res.Error().(*APIError)
		if !ok {
			return nil, fmt.Errorf(res.Status())
		}
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
		e, ok := res.Error().(*APIError)
		if !ok {
			return nil, fmt.Errorf(res.Status())
		}
		err = fmt.Errorf("%s: %s", res.Status(), e.Message)
		return
	}
	err = json.Unmarshal(res.Body(), &investigation)
	return
}

func (idr *IDR) InvestigationsResponse(q ...*InvestigationsQuery) (*InvestigationsResponse, error) {
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
		return nil, err
	}
	if res.IsError() {
		e, ok := res.Error().(*APIError)
		if !ok {
			err = fmt.Errorf(res.Status())
			return nil, err
		}
		err = fmt.Errorf("%s: %s", res.Status(), e.Message)
		return nil, err
	}
	var invRes *InvestigationsResponse
	err = json.Unmarshal(res.Body(), &invRes)
	if err != nil {
		err = errors.Join(fmt.Errorf("failed to parse response"), err)
		return nil, err
	}
	return invRes, nil
}

func (idr *IDR) Investigations(q ...*InvestigationsQuery) ([]*Investigation, error) {
	res, err := idr.InvestigationsResponse(q...)
	if err != nil {
		return nil, err
	}
	return res.Data, nil
}

func (idr *IDR) InvestigationsAll(q ...*InvestigationsQuery) ([]*Investigation, error) {
	res, err := idr.InvestigationsResponse(q...)
	if err != nil {
		return nil, err
	}
	investigations := res.Data
	pages := res.Metadata.TotalPages
	if pages > 1 {
		page := res.Metadata.Index
		lastPage := pages - 1
		for {
			if page == lastPage {
				return investigations, nil
			}
			var qq *InvestigationsQuery
			if len(q) > 0 {
				qq = q[0]
			} else {
				qq = &InvestigationsQuery{}
			}
			qq.Index = page
			pageRes, err := idr.InvestigationsResponse(qq)
			if err != nil {
				return nil, err
			}
			investigations = append(investigations, pageRes.Data...)
			page = pageRes.Metadata.Index
		}
	}
	return investigations, nil
}

func (idr *IDR) UpdateInvestigation(id string, update *InvestigationUpdateRequest) (*Investigation, error) {
	req := idr.http.R()
	req.SetError(&APIError{})
	req.SetBody(update)
	res, err := req.Patch(idr.URL("/v2/investigations", id))
	if err != nil {
		return nil, err
	}
	if res.IsError() {
		e, ok := res.Error().(*APIError)
		if !ok {
			return nil, fmt.Errorf(res.Status())
		}
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
