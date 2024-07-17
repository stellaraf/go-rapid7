package rapid7

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/go-resty/resty/v2"
)

type VM struct {
	BaseURL *url.URL
	http    *resty.Client
}

func (vm *VM) URL(paths ...string) string {
	p := []string{"/vm"}
	p = append(p, paths...)
	return path.Join(p...)
}

func (vm *VM) getQuery(search []VMAssetSearchQuery) VMAssetSearchQuery {
	var qp VMAssetSearchQuery
	if len(search) == 0 {
		qp = VMAssetSearchQuery{}
	} else {
		qp = search[0]
	}
	return qp
}

func (vm *VM) AssetSearch(search ...VMAssetSearchQuery) (*Rapid7VMPagedResponse[VMAsset], error) {
	qp := vm.getQuery(search)
	qp.Size = int(VM_ASSET_SEARCH_PAGE_SIZE)
	req := vm.http.R()
	req.SetError(&APIError{})
	req.SetQueryParams(qp.Map())
	res, err := req.Post(vm.URL("/v4/integration/assets"))
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
	var data *Rapid7VMPagedResponse[VMAsset]
	err = json.Unmarshal(res.Body(), &data)
	if err != nil {
		err = errors.Join(fmt.Errorf("failed to parse response"), err)
		return nil, err
	}
	return data, nil
}

func (vm *VM) Assets(search ...VMAssetSearchQuery) ([]VMAsset, error) {
	qp := vm.getQuery(search)
	qp.Size = int(VM_ASSET_SEARCH_PAGE_SIZE)
	res, err := vm.AssetSearch(qp)
	if err != nil {
		return nil, err
	}
	assets := res.Data
	pages := res.Metadata.TotalPages
	page := res.Metadata.Number
	lastPage := pages - 1
	if pages > 1 {
		for {
			page++
			if page == lastPage {
				return assets, nil
			}
			qqp := qp
			qqp.Page = int(page)
			pageRes, err := vm.AssetSearch(qqp)
			if err != nil {
				return nil, err
			}
			assets = append(assets, pageRes.Data...)
		}
	}
	return assets, nil
}

func (vm *VM) AssetCount() (uint64, error) {
	req := vm.http.R().SetQueryParam("size", "1").SetResult(&Rapid7VMPagedResponse[VMAsset]{})
	res, err := req.Post(vm.URL("/v4/integration/assets"))
	if err != nil {
		return 0, err
	}
	if res.IsError() {
		e, ok := res.Error().(*APIError)
		if !ok {
			err = fmt.Errorf(res.Status())
			return 0, err
		}
		err = fmt.Errorf("%s: %s", res.Status(), e.Message)
		return 0, err
	}
	data, ok := res.Result().(*Rapid7VMPagedResponse[VMAsset])
	if !ok {
		return 0, fmt.Errorf("failed to parse response")
	}
	return uint64(data.Metadata.TotalResources), nil
}

func newVM(region, apiKey string) (idr *VM, err error) {
	h := resty.New()
	urlS := fmt.Sprintf("https://%s.api.insight.rapid7.com", strings.ToLower(region))
	u, err := url.Parse(urlS)
	if err != nil {
		return
	}
	h.SetTimeout(RequestTimeout)
	h.SetHeader("X-API-Key", apiKey)
	validate, err := h.R().Get(fmt.Sprintf("https://%s/validate", u.Hostname()))
	if err != nil {
		return
	}
	if validate.StatusCode() > 200 {
		err = fmt.Errorf("failed to authenticate with Rapid7 API")
		return
	}
	h.SetBaseURL(u.String())
	idr = &VM{
		BaseURL: u,
		http:    h,
	}
	return
}
